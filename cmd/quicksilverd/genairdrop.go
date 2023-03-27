package main

import (
	"encoding/json"
	"fmt"
	"strconv"

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

// AddGenesisAirdropCmd returns add-genesis-airdrop cobra Command.
func AddGenesisAirdropCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-airdrop [address] [chain_id] [allocation] [base_value]",
		Short: "Add an airdrop claim to genesis.json",
		Long: `Add an airdrop claim to genesis.json. The zone drop record must already exist.
`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			allocation, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse allocation: %w", err)
			}

			baseValue, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse base_value: %w", err)
			}

			claimRecord := types.ClaimRecord{
				Address:       args[0],
				ChainId:       args[1],
				MaxAllocation: allocation,
				BaseValue:     baseValue,
			}

			if err := claimRecord.ValidateBasic(); err != nil {
				return err
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
				if zd.ChainId == claimRecord.ChainId {
					zoneDrop = zd
				}
			}

			if zoneDrop == nil {
				return fmt.Errorf("zoneDrop doesn't exist for chain ID: %s", claimRecord.ChainId)
			}

			for _, cr := range airdropGenState.ClaimRecords {
				if cr.ChainId == claimRecord.ChainId && cr.Address == claimRecord.Address {
					return fmt.Errorf("airdrop claimRecord already exists for user %s on chain ID: %s", claimRecord.Address, claimRecord.ChainId)
				}
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			airdropGenState.ClaimRecords = append(airdropGenState.ClaimRecords, &claimRecord)
			zoneDrop.Allocation += claimRecord.MaxAllocation

			airdropGenStateBz, err := clientCtx.Codec.MarshalJSON(airdropGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal airdrop genesis state: %w", err)
			}

			appState[types.ModuleName] = airdropGenStateBz

			// add base account for airdrop recipient, containing 1uqck
			balances := banktypes.Balance{Address: addr.String(), Coins: sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))}
			genAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			if !accs.Contains(addr) {

				// If this account does not exist in accs, create it.
				accs = append(accs, genAccount)
				accs = authtypes.SanitizeGenesisAccounts(accs)

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

				bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
				bankGenState.Balances = append(bankGenState.Balances, balances)
				bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)
				bankGenState.Supply = bankGenState.Supply.Add(balances.Coins...)

				bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
				if err != nil {
					return fmt.Errorf("failed to marshal bank genesis state: %w", err)
				}

				appState[banktypes.ModuleName] = bankGenStateBz
			}
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
