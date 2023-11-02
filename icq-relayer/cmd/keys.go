package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"golang.org/x/term"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagCoinType           = "coin-type"
	defaultCoinType uint32 = sdk.CoinType
)

var (
	// FlagAccountPrefix allows the user to override the prefix for a given account
	FlagAccountPrefix = ""
)

// keysCmd represents the keys command
func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "manage keys held by the relayer for each chain",
	}

	cmd.AddCommand(keysAddCmd())
	cmd.AddCommand(keysRestoreCmd())
	cmd.AddCommand(keysDeleteCmd())
	cmd.AddCommand(keysListCmd())
	cmd.AddCommand(keysShowCmd())
	cmd.AddCommand(keysEnumerateCmd())
	cmd.AddCommand(keysExportCmd())

	return cmd
}

// keysAddCmd respresents the `keys add` command
func keysAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [name]",
		Aliases: []string{"a"},
		Short:   "adds a key to the keychain associated with a particular chain",
		Long:    "if no name is passed, 'default' is used",
		Args:    cobra.RangeArgs(0, 1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys add
$ %s keys add test_key
$ %s k a osmo_key --chain osmosis`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := cfg.GetDefaultClient()
			var keyName string
			if len(args) == 0 {
				keyName = cl.Config.Key
			} else {
				keyName = args[0]
			}
			if cl.KeyExists(keyName) {
				return errKeyExists(keyName)
			}

			ko, err := cl.AddKey(keyName, 118)
			if err != nil {
				return err
			}

			out, err := json.Marshal(&ko)
			if err != nil {
				return err
			}

			fmt.Println(string(out))
			return nil
		},
	}
	// TODO: wire this up
	cmd.Flags().Uint32(flagCoinType, defaultCoinType, "coin type number for HD derivation")

	return cmd
}

// keysRestoreCmd respresents the `keys add` command
func keysRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "restore [name]",
		Aliases: []string{"r"},
		Short:   "restores a mnemonic to the keychain associated with a particular chain",
		Args:    cobra.ExactArgs(1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys restore --chain ibc-0 testkey
$ %s k r --chain ibc-1 faucet-key`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := cfg.GetDefaultClient()
			keyName := args[0]
			if cl.KeyExists(keyName) {
				return errKeyExists(keyName)
			}

			fmt.Print("Enter mnemonic 🔑: ")
			mnemonic, _ := term.ReadPassword(0)
			fmt.Println()

			address, err := cl.RestoreKey(keyName, string(mnemonic), 118)
			if err != nil {
				return err
			}

			fmt.Println(address)
			return nil
		},
	}
	// TODO: wire this up
	cmd.Flags().Uint32(flagCoinType, defaultCoinType, "coin type number for HD derivation")
	return cmd
}

// keysDeleteCmd respresents the `keys delete` command
func keysDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [name]",
		Aliases: []string{"d"},
		Short:   "deletes a key from the keychain associated with a particular chain",
		Args:    cobra.ExactArgs(1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys delete ibc-0 -y
$ %s keys delete ibc-1 key2 -y
$ %s k d ibc-2 testkey`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := cfg.GetDefaultClient()
			chainName := cl.Config.ChainID
			keyName := args[0]
			if !cl.KeyExists(keyName) {
				return errKeyDoesntExist(keyName)
			}

			if skip, _ := cmd.Flags().GetBool("skip"); !skip {
				fmt.Printf("Are you sure you want to delete key(%s) from chain(%s)? (Y/n)\n", keyName, chainName)
				if !askForConfirmation() {
					return nil
				}
			}

			if err := cl.DeleteKey(keyName); err != nil {
				panic(err)
			}

			fmt.Printf("key %s deleted\n", keyName)
			return nil
		},
	}

	return skipConfirm(cmd)
}

func askForConfirmation() bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("please type (y)es or (n)o and then press enter")
		return askForConfirmation()
	}
}

// keysListCmd respresents the `keys list` command
func keysListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "lists keys from the keychain associated with a particular chain",
		Args:    cobra.NoArgs,
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys list ibc-0
$ %s k l ibc-1`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := cfg.GetDefaultClient()
			info, err := cl.ListAddresses()
			if err != nil {
				return err
			}

			for key, val := range info {
				fmt.Printf("key(%s) -> %s\n", key, val)
			}

			return nil
		},
	}

	return cmd
}

// keysShowCmd respresents the `keys show` command
func keysShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [name]",
		Aliases: []string{"s"},
		Short:   "shows a key from the keychain associated with a particular chain",
		Long:    "if no name is passed, name in config is used",
		Args:    cobra.RangeArgs(0, 1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys show ibc-0
$ %s keys show ibc-1 key2
$ %s k s ibc-2 testkey`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := cfg.GetDefaultClient()
			var keyName string
			if len(args) == 0 {
				keyName = cl.Config.Key
			} else {
				keyName = args[0]
			}
			if !cl.KeyExists(keyName) {
				return errKeyDoesntExist(keyName)
			}

			if FlagAccountPrefix != "" {
				cl.Config.AccountPrefix = FlagAccountPrefix
			}

			address, err := cl.ShowAddress(keyName)
			if err != nil {
				return err
			}

			fmt.Println(address)
			return nil
		},
	}

	cmd.Flags().StringVar(&FlagAccountPrefix, "prefix", "", "Encode the key with the user specified prefix")

	return cmd
}

type KeyEnumeration struct {
	KeyName   string            `json:"key_name"`
	Addresses map[string]string `json:"addresses"`
}

// keysEnumerateCmd respresents the `keys enumerate` command
func keysEnumerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "enumerate [name]",
		Aliases: []string{"e"},
		Short:   "enumerates the address for a given key across all configured chains",
		Long:    "if no name is passed, name in config is used",
		Args:    cobra.RangeArgs(0, 1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys enumerate
$ %s keys enumerate key2
$ %s k e key2`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := cfg.GetDefaultClient()
			var keyName string
			if len(args) == 0 {
				keyName = cl.Config.Key
			} else {
				keyName = args[0]
			}
			account, err := cl.AccountFromKeyOrAddress(keyName)
			if err != nil {
				return err
			}

			var chains []string
			for chain := range cfg.Chains {
				chains = append(chains, chain)
			}
			sort.Strings(chains)

			addresses := make(map[string]string)
			for _, chain := range chains {
				client := cfg.GetClient(chain)
				address, err := client.EncodeBech32AccAddr(account)
				if err != nil {
					return err
				}
				addresses[chain] = address
			}

			return cl.PrintObject(addresses)
		},
	}

	// cmd.Flags().StringVar(&FlagAccountPrefix, "prefix", "", "Encode the key with the user specified prefix")

	return cmd
}

// keysExportCmd respresents the `keys export` command
func keysExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export [name]",
		Aliases: []string{"e"},
		Short:   "exports a privkey from the keychain associated with a particular chain",
		Args:    cobra.ExactArgs(1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys export ibc-0 testkey
$ %s k e ibc-2 testkey`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := cfg.GetDefaultClient()
			keyName := args[1]
			if !cl.KeyExists(keyName) {
				return errKeyDoesntExist(keyName)
			}

			info, err := cl.ExportPrivKeyArmor(keyName)
			if err != nil {
				return err
			}

			fmt.Println(info)
			return nil
		},
	}

	return cmd
}

func errKeyExists(name string) error {
	return fmt.Errorf("a key with name %s already exists", name)
}

func errKeyDoesntExist(name string) error {
	return fmt.Errorf("a key with name %s doesn't exist", name)
}

func skipConfirm(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP("skip", "y", false, "output using yaml")
	err := viper.BindPFlag("skip", cmd.Flags().Lookup("skip"))
	if err != nil {
		return nil
	}
	return cmd
}
