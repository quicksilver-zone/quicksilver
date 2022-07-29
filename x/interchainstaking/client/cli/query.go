package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group epochs queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		Aliases:                    []string{"ics"},
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdZonesInfos(),
		GetDelegatorIntentCmd(),
		GetDepositAccountCmd(),
	)

	return cmd
}

// GetCmdRegisteredZonesInfos provide running epochInfos
func GetCmdZonesInfos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zones",
		Short: "Query registered zones ",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query interchainstaking zones`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryRegisteredZonesInfoRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.RegisteredZoneInfos(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetDelegatorIntentCmd returns the intents of the user for the given chainID
// (zone).
func GetDelegatorIntentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intent [chain_id] [delegator_addr]",
		Short: "Query delegation intent for a given chain.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// args
			chainID := args[0]
			delegatorAddr := args[1]

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryDelegatorIntentRequest{
				ChainId:          chainID,
				DelegatorAddress: delegatorAddr,
			}

			res, err := queryClient.DelegatorIntent(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetDepositAccountCmd returns the deposit account for the given chainID
// (zone).
func GetDepositAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-account [chain_id]",
		Short: "Query deposit account address for a given chain.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// args
			chainID := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryDepositAccountForChainRequest{
				ChainId: chainID,
			}

			res, err := queryClient.DepositAccount(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
