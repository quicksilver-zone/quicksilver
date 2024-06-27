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

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/quicksilver-zone/quicksilver/app"
	quicksilverconfig "github.com/quicksilver-zone/quicksilver/cmd/config"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/runner"
	"github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/types"
	servercfg "github.com/quicksilver-zone/quicksilver/server/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	FlagHomePath = "home"
)

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	customAppTemplate, customAppConfig := servercfg.AppConfig(quicksilverconfig.BaseDenom)

	srvCfg, ok := customAppConfig.(servercfg.Config)
	if !ok {
		panic(fmt.Errorf("unknown app config type %T", customAppConfig))
	}

	srvCfg.StateSync.SnapshotInterval = 1500
	srvCfg.StateSync.SnapshotKeepRecent = 2

	return customAppTemplate, srvCfg
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
			log.Info().Msgf("Config file created at %s\n", configFilePath)
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
			log.Info().Msg(runner.VERSION)
		},
	}
	return versionCommand
}

func StartCommand() *cobra.Command {
	startCommand := &cobra.Command{
		Use:     "start",
		Short:   "Start the server [from_relayer_key_or_address]",
		Long:    `Start the server [from_relayer_key_or_address]`,
		Example: fmt.Sprintf("%s start icq1 -y", appName),
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			homepath, err := cmd.Flags().GetString(FlagHomePath)
			if err != nil {
				return err
			}

			config := InitConfig(homepath)

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
			log.Info().Msg("starting the server and listening for epochs")

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)
			go runner.Run(ctx, &config, CreateErrHandler(c), cmd)

			for sig := range c {
				log.Info().Msgf("Signal Received (%s) - gracefully shutting down", sig.String())
				break
			}
			return nil
		},
	}

	startCommand.Flags().String(FlagHomePath, types.DefaultHomePath, "homedir")
	flags.AddTxFlagsToCmd(startCommand)
	return startCommand
}

func InitConfig(homepath string) types.Config {
	cfg := types.InitializeConfigFromToml(homepath)
	return cfg
}

func CreateErrHandler(sigC chan os.Signal) func(err error) {
	return func(err error) {
		log.Err(err).Msg("Aborting")
		sigC <- syscall.SIGABRT
	}
}
