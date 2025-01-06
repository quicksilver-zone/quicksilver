package cmd

import (
	"github.com/spf13/cobra"
)

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	// home, err := cmd.PersistentFlags().GetString(flags.FlagHome)
	// if err != nil {
	// 	return err
	// }
	// cfg = &types.Config{}
	// cfgPath := path.Join(home, "config.yaml")
	// _, err = os.Stat(cfgPath)
	// if err != nil {
	// 	if !os.IsNotExist(err) { // Return immediately
	// 		return err
	// 	}

	// 	if err := types.CreateConfig(home); err != nil {
	// 		return err
	// 	}
	// }

	// viper.SetConfigFile(cfgPath)
	// err = viper.ReadInConfig()
	// if err != nil {
	// 	fmt.Println("Failed to read in config:", err)
	// 	os.Exit(1)
	// }

	// read the config file bytes
	// file, err := os.ReadFile(viper.ConfigFileUsed())
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	os.Exit(1)
	// }

	// // unmarshall them into the struct
	// if err = yaml.Unmarshal(file, cfg); err != nil {
	// 	fmt.Println("Error unmarshalling config:", err)
	// 	os.Exit(1)
	// }

	// validate configuration
	// if err = types.ValidateConfig(cfg); err != nil {
	// 	fmt.Println("Error parsing chain config:", err)
	// 	os.Exit(1)
	// }
	return nil
}
