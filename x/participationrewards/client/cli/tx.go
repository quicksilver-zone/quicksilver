package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// GetTxCmd returns a root CLI command handler for all x/bank transaction commands.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Participation rewards transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(GetSubmitClaimTxCmd())

	return txCmd
}

func GetSubmitClaimTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [zone] [asset_type]",
		Short: `Submit proof of assets held in the given zone.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			zone := args[0]
			assetType := args[1]

			msg := types.NewMsgSubmitClaim(clientCtx.GetFromAddress(), zone, assetType)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
