package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetTxCmd returns a root CLI command handler for all x/bank transaction commands.
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

	return txCmd
}

// GetSignalIntentTxCmd returns a CLI command handler for signalling validator
// delegation intent.
func GetSignalIntentTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signal-intent [chain_id] [delegation_intent]",
		Short: `Signal validator delegation intent.`,
		Long: `signal validator delegation intent by providing a comma seperated string
containing a decimal weight and the bech32 validator address,
e.g. "0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz"`,
		Example: `signal-intent [chain_id] 0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chain_id := args[0]
			intents, err := types.IntentsFromString(args[1])
			if err != nil {
				return fmt.Errorf("%v, see example: %v", err, cmd.Example)
			}

			msg := types.NewMsgSignalIntent(chain_id, intents, clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetRegisterZoneTxCmd returns a CLI command handler for creating a MsgSend transaction.
func GetRequestRedemptionTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem [coins] [destination_address]",
		Short: `Redeem tokens.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			if err != nil {
				return err
			}
			coins := args[0]
			destination_address := args[1]

			msg := types.NewMsgRequestRedemption(coins, destination_address, clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSubmitRegisterProposal implements the command to submit a register-zone proposal
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

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewRegisterZoneProposal(proposal.Title, proposal.Description, proposal.ConnectionId, proposal.BaseDenom,
				proposal.LocalDenom, proposal.AccountPrefix, proposal.MultiSend, proposal.LiquidityModule)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
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

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

// GetCmdSubmitRegisterProposal implements the command to submit a register-zone proposal
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

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
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

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
