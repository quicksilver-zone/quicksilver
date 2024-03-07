package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// GetQueryCmd returns the cli query commands for interchainstaking module.
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
		GetCmdZones(),
		GetDelegatorIntentCmd(),
		GetDepositAccountCmd(),
		GetMappedAccountsCmd(),
		GetWithdrawalRecordsCmd(),
		GetUserWithdrawalRecordsCmd(),
		GetZoneWithdrawalRecords(),
		GetUnbondingRecordsCmd(),
		GetReceiptsCmd(),
		GetTxStatusCmd(),
		GetZoneRedelegationRecordsCmd(),
		GetZoneValidatorsCmd(),
		GetZoneCmd(),
	)

	return cmd
}

// GetCmdZonesInfos provide running epochInfos.
func GetCmdZones() *cobra.Command {
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

			req := &types.QueryZonesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Zones(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetDelegatorIntentCmd returns the intents of the user for the given chainID.
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

// GetMappedAccountsCmd returns the mapped account for the given address.
func GetMappedAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mapped-accounts [address]",
		Short: "Query mapped accounts for a given address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// args
			address := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryMappedAccountsRequest{
				Address: address,
			}

			res, err := queryClient.MappedAccounts(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetWithdrawalRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdrawal-records",
		Short: "Query all withdrawal records",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query interchainstaking withdrawal-records`,
				version.AppName,
			)),
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

			req := &types.QueryWithdrawalRecordsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.WithdrawalRecords(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetUserWithdrawalRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-withdrawal-record [user-address]",
		Short: "Query withdrawal record for a given address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// args
			address := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryUserWithdrawalRecordsRequest{
				UserAddress: address,
			}

			res, err := queryClient.UserWithdrawalRecords(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetZoneWithdrawalRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone-withdrawal-records [chain-id]",
		Short: "Query withdrawal records for a given zone.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// args
			chainID := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryWithdrawalRecordsRequest{
				ChainId: chainID,
			}

			res, err := queryClient.ZoneWithdrawalRecords(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetUnbondingRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-records",
		Short: "Query all unbonding records",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query interchainstaking unbonding-records`,
				version.AppName,
			)),
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

			req := &types.QueryUnbondingRecordsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.UnbondingRecords(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetReceiptsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recepts",
		Short: "Query all receipts",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query interchainstaking receipts`,
				version.AppName,
			)),
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

			req := &types.QueryReceiptsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Receipts(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetTxStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx-status [chain-id] [tx-hash]",
		Short: "Query the status of a transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			chainID := args[0]
			txHash := args[1]

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryTxStatusRequest{
				ChainId: chainID,
				TxHash:  txHash,
			}

			res, err := queryClient.TxStatus(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetZoneRedelegationRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone-redelegation-records",
		Short: "Query re-delegation records for a given zone.",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query interchainstaking zone-redelegation-records`,
				version.AppName,
			)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			chainID := args[0]
			req := &types.QueryRedelegationRecordsRequest{
				ChainId: chainID,
			}

			// TODO: refactor this. Should be RedelegationRecords
			res, err := queryClient.RedelegationRecords(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetZoneValidatorCmd returns the validators for the given zone.
func GetZoneValidatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone-validators",
		Short: "Query validators for a given zone.",
		Example: strings.TrimSpace(
			fmt.Sprintf(`$ %s query interchainstaking zone-validators`,
				version.AppName,
			)),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryZoneValidatorsRequest{}

			res, err := queryClient.ZoneValidators(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetZoneCmd returns the information about the zone.
func GetZoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone [chain-id]",
		Short: "Query zone information for a given chain.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			chainID := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryZoneRequest{
				ChainId: chainID,
			}

			res, err := queryClient.Zone(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
