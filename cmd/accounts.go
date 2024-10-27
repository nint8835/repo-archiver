package cmd

import (
	"github.com/spf13/cobra"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage account aliases",
}

func init() {
	rootCmd.AddCommand(accountsCmd)
}
