package cmd

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ingenuity-build/interchain-queries/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	homePath       string
	overridenChain string
	debug          bool
	cfg            *config.Config
	defaultHome    = os.ExpandEnv("$HOME/.icq")
	appName        = "lens"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "interchain-queries",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.claim-and-delegate.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// reads `homeDir/config.yaml` into `var config *Config` before each command
		if err := initConfig(rootCmd); err != nil {
			return err
		}
		return nil
	}

	// --home flag
	rootCmd.PersistentFlags().StringVar(&homePath, flags.FlagHome, defaultHome, "set home directory")
	if err := viper.BindPFlag(flags.FlagHome, rootCmd.PersistentFlags().Lookup(flags.FlagHome)); err != nil {
		panic(err)
	}

	// --debug flag
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringP("output", "o", "json", "output format (json, indent, yaml)")
	if err := viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVar(&overridenChain, "chain", "", "override default chain")
	if err := viper.BindPFlag("chain", rootCmd.PersistentFlags().Lookup("chain")); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(keysCmd())
}
