package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

// GetQueryCmd returns the cli query commands for the eventmanager module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Query subcommands for the %s module", types.ModuleName),
		Aliases:                    []string{"em"},
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetEventsQueryCmd(),
	)

	return cmd
}

func GetEventsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Query the events for given chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			query := &types.QueryEventsRequest{
				ChainId: args[0],
			}

			res, err := queryClient.Events(context.Background(), query)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
