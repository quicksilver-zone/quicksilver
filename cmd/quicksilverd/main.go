package main

import (
	"fmt"
	"os"
	"path/filepath"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/quicksilver-zone/quicksilver/app"
	cmdcfg "github.com/quicksilver-zone/quicksilver/cmd/config"
)

func main() {
	cmdcfg.SetupConfig()
	cmdcfg.RegisterDenoms()

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	app.DefaultNodeHome = filepath.Join(userHomeDir, ".quicksilverd")

	rootCmd, _ := NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "QUICKSILVERD", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}

	os.Exit(1)
}
