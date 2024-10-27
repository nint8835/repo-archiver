package cmd

import (
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/nint8835/repo-archiver/pkg/config"
)

var accountsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new account alias",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		displayName := ""
		accountName := ""
		isArchive := true

		err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Display Name").
					Description("A friendly name for the account").
					Value(&displayName),
				huh.NewInput().
					Title("Account Name").
					Description("The actual name of the GitHub user or organization").
					Value(&accountName),
				huh.NewConfirm().
					Title("Is this account an archive?").
					Description("Archive accounts will have repos archived on transfer").
					Value(&isArchive),
			),
		).WithTheme(huh.ThemeCatppuccin()).Run()
		checkError(err, "error prompting for account information")

		config.Instance.Accounts[displayName] = config.Account{
			Name:      accountName,
			IsArchive: isArchive,
		}

		err = config.Save()
		checkError(err, "error saving config")
	},
}

func init() {
	accountsCmd.AddCommand(accountsAddCmd)
}
