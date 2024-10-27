package cmd

import (
	"github.com/spf13/cobra"
)

var accountsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an account alias",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	accountsCmd.AddCommand(accountsRemoveCmd)
}
