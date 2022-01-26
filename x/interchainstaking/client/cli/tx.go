package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetTxCmd returns a root CLI command handler for all x/bank transaction commands.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Interchain staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(GetRegisterZoneTxCmd())

	return txCmd
}

// GetRegisterZoneTxCmd returns a CLI command handler for creating a MsgSend transaction.
func GetRegisterZoneTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [identifier] [chain_id] [local_denom] [remote_denom]",
		Short: `Send funds from one account to another.`,
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			if err != nil {
				return err
			}
			identifier := args[0]
			chain_id := args[1]
			local_denom := args[2]
			remote_denom := args[3]

			msg := types.NewMsgRegisterZone(identifier, chain_id, local_denom, remote_denom, clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
