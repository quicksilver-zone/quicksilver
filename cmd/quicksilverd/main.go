package main

import (
	"os"
	"path/filepath"

	"cosmossdk.io/log"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/quicksilver-zone/quicksilver/v7/cmd/quicksilverd/cmd"

	"github.com/quicksilver-zone/quicksilver/v7/app"
	cmdcfg "github.com/quicksilver-zone/quicksilver/v7/cmd/quicksilverd/config"
)

func main() {
	cmdcfg.SetupConfig()
	cmdcfg.RegisterDenoms()

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	app.DefaultNodeHome = filepath.Join(userHomeDir, ".quicksilverd")

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		log.NewLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
