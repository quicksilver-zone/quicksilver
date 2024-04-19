/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/runner"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runner.Run(cfg, cmd.Flag("home").Value.String())
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
