package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

// GetTxCmd returns the cli transaction commands for the airdrop module.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Transaction subcommands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetClaimTxCmd(),
	)

	return txCmd
}

func GetClaimTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [chainID] [action]",
		Short: "claim airdrop for the given action in the given zone",
		Example: strings.TrimSpace(
			fmt.Sprintf("$ %s tx %s claim %s %s",
				version.AppName,
				types.ModuleName,
				exampleChainID,
				exampleAction,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chainID := args[0]
			action, err := strconv.ParseInt(args[1], 10, 32)

			msg := &types.MsgClaim{
				ChainId: chainID,
				Action:  int32(action),
				Address: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
