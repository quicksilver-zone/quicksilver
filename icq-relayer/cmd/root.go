package cmd

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/spf13/viper"
	tmcfg "github.com/tendermint/tendermint/config"
	"os"

	"github.com/spf13/cobra"
)

var (
	homePath       string
	overridenChain string
	defaultHome    = os.ExpandEnv("$HOME/.icq")
	appName        = "icq-relayer"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.EnableCommandSorting = false
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithBroadcastMode(flags.BroadcastBlock).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithHomeDir(defaultHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "A relayer for the Quicksilver interchain queries module",
		Long:  `A relayer for Quicksilver interchain-queries, allowing cryptographically verifiable cross-chain KV lookups.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := tmcfg.DefaultConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
		},
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.claim-and-delegate.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("output", "o", "json", "output format (json, indent, yaml)")
	if err := viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVar(&overridenChain, "chain", "", "override default chain")
	if err := viper.BindPFlag("chain", rootCmd.PersistentFlags().Lookup("chain")); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(StartCommand())
	rootCmd.AddCommand(VersionCommand())
	rootCmd.AddCommand(InitConfigCommand())
	rootCmd.AddCommand(keys.Commands(defaultHome))

	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := svrcmd.Execute(rootCmd, "ICQRELAYER", defaultHome); err != nil {
		var exitError *server.ErrorCode
		if errors.As(err, &exitError) {
			os.Exit(exitError.Code)
		}

		os.Exit(1)
	}
}

// withUsage wraps a PositionalArgs to display usage only when the PositionalArgs
// variant is violated.
func withUsage(inner cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := inner(cmd, args); err != nil {
			cmd.Root().SilenceUsage = false
			cmd.SilenceUsage = false
			return err
		}

		return nil
	}
}
