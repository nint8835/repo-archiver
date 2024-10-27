package cmd

import (
	"github.com/spf13/cobra"
)

var accountsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new account alias",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	accountsCmd.AddCommand(accountsAddCmd)
}
