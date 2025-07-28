/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/BurntSushi/toml"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/logger"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/runner"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/types"
)

const (
	FlagHomePath = "home"
	FlagLogLevel = "log-level"
)

func init() {
	rootCmd.AddCommand(StartCommand())
	rootCmd.AddCommand(VersionCommand())
	rootCmd.AddCommand(InitConfigCommand())
}

func InitConfigCommand() *cobra.Command {
	initConfigCommand := &cobra.Command{
		Use:   "init",
		Short: "Initialize the config",
		Long:  `Initialize the config`,
		RunE: func(cmd *cobra.Command, args []string) error {
			homepath, err := cmd.Flags().GetString(FlagHomePath)
			if err != nil {
				return err
			}
			config := types.NewConfig()
			configFilePath := filepath.Join(homepath, "config.toml")
			if _, err := os.Stat(configFilePath); err == nil {
				return fmt.Errorf("config file already exists at %s", configFilePath)
			}
			f, err := os.Create(configFilePath)
			if err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}
			defer f.Close()

			encoder := toml.NewEncoder(f)
			if err := encoder.Encode(config); err != nil {
				return fmt.Errorf("failed to encode config to TOML: %w", err)
			}
			fmt.Printf("Config file created at %s\n", configFilePath)
			return nil
		},
	}
	return initConfigCommand
}

func VersionCommand() *cobra.Command {
	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of icq-relayer",
		Long:  `Print the version number of icq-relayer`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", runner.VERSION)
			fmt.Printf("Quicksilver Version: %s\n", runner.QUICKSILVER_VERSION)
			fmt.Printf("Commit: %s\n", runner.COMMIT)
		},
	}
	return versionCommand
}

func StartCommand() *cobra.Command {
	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start the server",
		Long:  `Start the server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			homepath, err := cmd.Flags().GetString(FlagHomePath)
			if err != nil {
				return err
			}
			logLevel, err := cmd.Flags().GetString(FlagLogLevel)
			if err != nil {
				return err
			}

			logger := logger.New(logger.LogLevel(logLevel))

			config := InitConfig(homepath, logger)
			config.FilterHA()

			rpcClient, err := rpchttp.New(config.DefaultChain.RpcUrl, "/websocket")
			if err != nil {
				return err
			}

			encodingCfg := app.MakeEncodingConfig()

			clientContext := client.Context{}.
				WithCodec(encodingCfg.Marshaler).
				WithInterfaceRegistry(encodingCfg.InterfaceRegistry).
				WithTxConfig(encodingCfg.TxConfig).
				WithLegacyAmino(encodingCfg.Amino).
				WithInput(os.Stdin).
				WithAccountRetriever(authtypes.AccountRetriever{}).
				WithHomeDir(homepath).
				WithNodeURI(config.DefaultChain.RpcUrl).
				WithClient(rpcClient).
				WithViper("")

			config.ClientContext = &clientContext
			config.ProtoCodec = codec.NewProtoCodec(clientContext.InterfaceRegistry)
			ctx := context.Background()
			level.Info(logger).Log("msg", "starting the server and listening for epochs")

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)

			go runner.Run(ctx, &config, logger, CreateErrHandler(c, logger))

			for sig := range c {
				level.Info(logger).Log("msg", "Signal Received, shutting down", "signal", sig.String())
				break
			}
			return nil
		},
	}

	startCommand.Flags().String(FlagHomePath, types.DefaultConfigPath, "homedir")
	startCommand.Flags().String(FlagLogLevel, "info", "log level")
	return startCommand
}

func InitConfig(homepath string, logger log.Logger) types.Config {
	cfg := types.InitializeConfigFromToml(homepath, logger)
	return cfg
}

func CreateErrHandler(sigC chan os.Signal, logger log.Logger) func(err error) {
	return func(err error) {
		level.Error(logger).Log("msg", "fatal error received", "err", err)
		sigC <- syscall.SIGABRT
	}
}
