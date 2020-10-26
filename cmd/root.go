package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(txtCmd)
}

var rootCmd = &cobra.Command{
	Use:   "novelcli",
	Short: "novel cli tools",
	Long:  `generate novels from web url`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		cmd.UsageString()
	},
}

//Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
