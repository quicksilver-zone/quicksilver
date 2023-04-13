package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// AddZonedropCmd returns add-zonedrop cobra Command.
func AddZonedropCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-zonedrop [chain_id] [start_time] [duration] [decay] [actions]",
		Short: "Add an zonedrop to genesis.json",
		Long:  `Add an zonedrop to genesis.json.`,

		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			chainID := args[0]

			layout := "2006-01-02T15:04:05Z"
			startTime, err := time.Parse(layout, args[1])
			if err != nil {
				panic(err)
			}

			duration, err := time.ParseDuration(args[2])
			if err != nil {
				panic(err)
			}
			decay, err := time.ParseDuration(args[3])
			if err != nil {
				panic(err)
			}
			actionString := args[4]

			airdrop := types.ZoneDrop{
				ChainId:     chainID,
				StartTime:   startTime,
				Duration:    duration,
				Decay:       decay,
				Allocation:  0,
				Actions:     []sdk.Dec{},
				IsConcluded: false,
			}

			actions := strings.Split(actionString, ",")
			for _, action := range actions {
				weight := sdk.MustNewDecFromStr(action)
				airdrop.Actions = append(airdrop.Actions, weight)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			airdropGenState := types.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			// assert zonedrop exists
			for _, zd := range airdropGenState.ZoneDrops {
				if zd.ChainId == airdrop.ChainId {
					panic("ZoneDrop for this chainId already exists")
				}
			}

			airdropGenState.ZoneDrops = append(airdropGenState.ZoneDrops, &airdrop)

			airdropGenStateBz, err := clientCtx.Codec.MarshalJSON(airdropGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal airdrop genesis state: %w", err)
			}

			appState[types.ModuleName] = airdropGenStateBz
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// BulkGenesisAirdropCmd returns add-genesis-airdrop cobra Command.
func BulkGenesisAirdropCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk-genesis-airdrop [file.csv] [chain_id]",
		Short: "Add an airdrop claim to genesis.json, from csv",
		Long:  `Add an airdrop claim to genesis.json, from csv. The zone drop record must already exist.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			csvfile, err := os.Open(args[0])
			if err != nil {
				log.Fatalln("Couldn't open the csv file", err)
			}

			r := csv.NewReader(bufio.NewReader(csvfile))

			claimRecords := make([]*types.ClaimRecord, 0)
			// Iterate through the records
			for {
				// Read each record from csv
				record, err := r.Read()
				if err == io.EOF {
					break
				}

				addr, err := sdk.AccAddressFromBech32(record[0])
				if err != nil {
					return err
				}

				allocation, err := strconv.ParseUint(record[2], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse allocation: %w", err)
				}

				baseValue, err := strconv.ParseUint(record[1], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse base_value: %w", err)
				}

				claimRecord := types.ClaimRecord{
					Address:       addr.String(),
					ChainId:       args[1],
					MaxAllocation: allocation,
					BaseValue:     baseValue,
				}

				if err := claimRecord.ValidateBasic(); err != nil {
					return err
				}

				claimRecords = append(claimRecords, &claimRecord)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			airdropGenState := types.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			var zoneDrop *types.ZoneDrop
			// assert zonedrop exists
			for _, zd := range airdropGenState.ZoneDrops {
				if zd.ChainId == claimRecords[0].ChainId {
					zoneDrop = zd
				}
			}

			if zoneDrop == nil {
				return fmt.Errorf("zoneDrop doesn't exist for chain ID: %s", claimRecords[0].ChainId)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}
			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			zoneclaims := map[string]bool{}
			existing := airdropGenState.ClaimRecords

			for _, i := range existing {
				if i.ChainId == claimRecords[0].ChainId {
					zoneclaims[i.Address] = true
				}
			}

		OUTER:
			for idx, claimRecord := range claimRecords {
				if idx%100 == 0 {
					fmt.Printf("(%d/%d)...\n", idx, len(claimRecords))
				}

				if _, exists := zoneclaims[claimRecord.Address]; exists {
					return fmt.Errorf("airdrop claimRecord already exists for user %s on chain ID: %s", claimRecord.Address, claimRecord.ChainId)
				}

				// Add the new account to the set of genesis accounts and sanitize the
				// accounts afterwards.
				zoneDrop.Allocation += claimRecord.MaxAllocation

				// add base account for airdrop recipient, containing 1uqck
				balances := banktypes.Balance{Address: claimRecord.Address, Coins: sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))}
				addr, _ := sdk.AccAddressFromBech32(claimRecord.Address)
				genAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

				if err := genAccount.Validate(); err != nil {
					return fmt.Errorf("failed to validate new genesis account: %w", err)
				}

				if !accs.Contains(addr) {
					// If this account does not exist in accs, create it.
					accs = append(accs, genAccount)
					accs = authtypes.SanitizeGenesisAccounts(accs)

					bankGenState.Balances = append(bankGenState.Balances, balances)
					bankGenState.Supply = bankGenState.Supply.Add(balances.Coins...)
				}
				continue OUTER

			}
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

			airdropGenState.ClaimRecords = append(airdropGenState.ClaimRecords, claimRecords...)

			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs

			authGenStateBz, err := clientCtx.Codec.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[authtypes.ModuleName] = authGenStateBz

			airdropGenStateBz, err := clientCtx.Codec.MarshalJSON(airdropGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal airdrop genesis state: %w", err)
			}

			appState[types.ModuleName] = airdropGenStateBz

			bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}

			appState[banktypes.ModuleName] = bankGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
