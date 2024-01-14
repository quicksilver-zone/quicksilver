package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/server"
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

	if err := svrcmd.Execute(rootCmd, "QUICKSILVERD", app.DefaultNodeHome); err != nil {
		var exitError *server.ErrorCode
		if errors.As(err, &exitError) {
			os.Exit(exitError.Code)
		}

		os.Exit(1)
	}

	os.Exit(0)
}
