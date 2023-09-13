package main_test

import (
	"fmt"
	"testing"

	"github.com/quicksilver-zone/quicksilver/app"
	quicksilverd "github.com/quicksilver-zone/quicksilver/cmd/quicksilverd"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
)

func TestInitCmd(t *testing.T) {
	rootCmd, _ := quicksilverd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",             // Test the init cmd
		"quicksilver-test", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
		fmt.Sprintf("--%s=%s", flags.FlagChainID, "quicksilver-1"),
	})

	err := svrcmd.Execute(rootCmd, "QUICKSILVERD", app.DefaultNodeHome)
	require.NoError(t, err)
}
