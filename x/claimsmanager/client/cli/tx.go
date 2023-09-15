package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
)

// GetTxCmd returns a root CLI command handler for all x/bank transaction commands.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "ClaimsManager transaction subcommands",
		Aliases:                    []string{"cm"},
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return txCmd
}
