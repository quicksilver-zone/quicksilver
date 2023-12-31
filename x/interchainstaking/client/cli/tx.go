package cli

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// GetTxCmd returns a root CLI command handler for all x/interchainstaking transaction commands.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Aliases:                    []string{"ics"},
		Short:                      "Interchain staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(GetSignalIntentTxCmd())
	txCmd.AddCommand(GetRequestRedemptionTxCmd())
	txCmd.AddCommand(GetReopenChannelTxCmd())

	return txCmd
}

// GetSignalIntentTxCmd returns a CLI command handler for signalling validator
// delegation intent.
func GetSignalIntentTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signal-intent [chainID] [delegation_intent]",
		Short: `Signal validator delegation intent.`,
		Long: `signal validator delegation intent by providing a comma separated string
containing a decimal weight and the bech32 validator address,
e.g. "0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz"`,
		Example: `signal-intent [chain_id] 0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chainID := args[0]

			msg := types.NewMsgSignalIntent(chainID, args[1], clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetRequestRedemptionTxCmd returns a CLI command handler for creating a Request transaction.
func GetRequestRedemptionTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem [coin] [destination_address]",
		Short: `Redeem tokens.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			destinationAddress := args[1]
			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return fmt.Errorf("unable to parse coin %s", args[0])
			}

			msg := types.NewMsgRequestRedemption(coin, destinationAddress, clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetReopenChannelTxCmd returns a CLI command handler for creating a Reopen ICA port transaction.
func GetReopenChannelTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reopen [connection] [port]",
		Short: `Reopen closed ICA port.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			connectionID := args[0]
			port := args[1]

			msg := types.NewMsgGovReopenChannel(connectionID, port, clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSubmitRegisterProposal implements the command to submit a register-zone proposal.
func GetCmdSubmitRegisterProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-zone [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a zone registration proposal",
		Long: strings.TrimSpace(
			`Submit a zone registration proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal register-zone <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "Register cosmoshub-4",
  "description": "Onboard the cosmoshub-4 zone to Quicksilver",
  "connection_id": "connection-3",
  "base_denom": "uatom",
  "local_denom": "uqatom",
  "account_prefix": "cosmos",
  "multi_send": true,
  "liquidity_module": false,
  "deposit": "512000000uqck"
}
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := ParseZoneRegistrationProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			err = sdk.ValidateDenom(proposal.LocalDenom)
			if err != nil {
				return err
			}

			err = sdk.ValidateDenom(proposal.BaseDenom)
			if err != nil {
				return err
			}

			if proposal.MessagesPerTx < 1 {
				return errors.New("messages_per_tx must be a positive non-zero integer")
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewRegisterZoneProposal(proposal.Title, proposal.Description, proposal.ConnectionId, proposal.BaseDenom,
				proposal.LocalDenom, proposal.AccountPrefix, proposal.ReturnToSender, proposal.UnbondingEnabled, proposal.DepositsEnabled, proposal.LiquidityModule, proposal.Decimals, proposal.MessagesPerTx, proposal.Is_118)

			msg, err := govv1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func ParseZoneRegistrationProposal(cdc codec.JSONCodec, proposalFile string) (types.RegisterZoneProposalWithDeposit, error) {
	proposal := types.RegisterZoneProposalWithDeposit{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	if reflect.DeepEqual(proposal, types.RegisterZoneProposalWithDeposit{}) {
		return proposal, errors.New("cannot unmarshal empty JSON object")
	}

	return proposal, nil
}

// GetCmdSubmitUpdateProposal implements the command to update a register-zone proposal.
func GetCmdSubmitUpdateProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-zone [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a zone update proposal",
		Long: strings.TrimSpace(
			`Submit a zone update proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal register-zone <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "Enable liquidity module for cosmoshub-4",
  "description": "Update cosmoshub-4 to enable liquidity module",
  "chain_id": "cosmoshub-4",
  "changes": [{
      "key": "liquidity_module",
      "value": "true",
  }],
  "deposit": "512000000uqck"
}
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := ParseZoneUpdateProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewUpdateZoneProposal(proposal.Title, proposal.Description, proposal.ChainId, proposal.Changes)

			msg, err := govv1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func ParseZoneUpdateProposal(cdc codec.JSONCodec, proposalFile string) (types.UpdateZoneProposalWithDeposit, error) {
	proposal := types.UpdateZoneProposalWithDeposit{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	if reflect.DeepEqual(proposal, types.UpdateZoneProposalWithDeposit{}) {
		return proposal, errors.New("cannot unmarshal empty JSON object")
	}

	return proposal, nil
}
