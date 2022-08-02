package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

// GetQueryCmd returns the cli query commands for the airdrop module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Query subcommands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand()

	return cmd
}

func GetParamsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: fmt.Sprintf("Query the current %s parameters", types.ModuleName),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryParamsRequest{}
			res, err := queryClient.Params(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetZoneDropQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone [chain_id]",
		Short: "Query airdrop details of the specified zone",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query %s zone %s`,
				version.AppName,
				types.ModuleName,
				exampleChainID,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// args
			chainID := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryZoneDropRequest{
				ChainId: chainID,
			}

			res, err := queryClient.ZoneDrop(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetAccountBalanceQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-balance [chain_id]",
		Short: "Query airdrop account balance of the specified zone",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query %s account-balance %s`,
				version.AppName,
				types.ModuleName,
				exampleChainID,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// args
			chainID := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryAccountBalanceRequest{
				ChainId: chainID,
			}

			res, err := queryClient.AccountBalance(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetZoneDropsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone-drops [status]",
		Short: "Query all airdrops of the specified status",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query %s zone-drops %s`,
				version.AppName,
				types.ModuleName,
				exampleStatus,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// args
			status, ok := types.Status_value[args[0]]
			if !ok {
				return types.ErrUnknownStatus
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryZoneDropsRequest{
				Status:     types.Status(status),
				Pagination: pageReq,
			}

			res, err := queryClient.ZoneDrops(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetClaimRecordQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-record [chain_id] [address]",
		Short: "Query airdrop claim record details of the given address for the given zone.",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query %s claim-record %s %s`,
				version.AppName,
				types.ModuleName,
				exampleChainID,
				exampleAddress,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// args
			chainID := args[0]
			address, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryClaimRecordRequest{
				ChainId: chainID,
				Address: address.String(),
			}

			res, err := queryClient.ClaimRecord(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
