package cmd

import (
	"github.com/spf13/cobra"
)

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all account aliases",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	accountsCmd.AddCommand(accountsListCmd)
}
