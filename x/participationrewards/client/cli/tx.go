package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// GetTxCmd returns a root CLI command handler for all x/bank transaction commands.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Participation rewards transaction subcommands",
		Aliases:                    []string{"pr"},
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(GetSubmitClaimTxCmd())

	return txCmd
}

func GetSubmitClaimTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [zone] [src-zone] [claim-type] [payload-file].json",
		Short: `Submit proof of assets held in the given zone.`,
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			zone := args[0]
			srcZone := args[1]
			claimTypeStr := args[2]
			fileName := args[3]

			contents, err := os.ReadFile(fileName)
			if err != nil {
				return err
			}

			var proofs []*types.Proof

			if err = json.Unmarshal(contents, &proofs); err != nil {
				return err
			}

			claimType, ok := types.ClaimType_value[claimTypeStr]
			if !ok {
				return fmt.Errorf("invalid claim type: %s", claimTypeStr)
			}

			msg := types.NewMsgSubmitClaim(clientCtx.GetFromAddress(), zone, srcZone, types.ClaimType(claimType), proofs)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdAddProtocolDataProposal implements the command to submit a add protocol data proposal
func GetCmdAddProtocolDataProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-protocol-data [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a add protocol data proposal",
		Long: strings.TrimSpace(
			`Submit an add protocol data proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal add-protocol-data <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "Add Osmosis Atom/qAtom Pool",
  "description": "Add Osmosis Atom/qAtom Pool to support participation rewards",
  "protocol": "osmosis",
  "key": "pools/XXX",
  "type": "osmosispool",
  "data": {
	"poolID": "596",
	"ibcToken": "27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
	"localDenom": "uqatom"
  },
  "deposit": "512000000uqck"
}
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := ParseAddProtocolDataProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			// parseData based on protocol

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewAddProtocolDataProposal(proposal.Title, proposal.Description, proposal.Type, proposal.Protocol, proposal.Key,
				proposal.Data)

			msg, err := govv1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func ParseAddProtocolDataProposal(cdc codec.JSONCodec, proposalFile string) (types.AddProtocolDataProposalWithDeposit, error) {
	proposal := types.AddProtocolDataProposalWithDeposit{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
