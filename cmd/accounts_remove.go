package cmd

import (
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/nint8835/repo-archiver/pkg/config"
)

var accountsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an account alias",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var options []huh.Option[string]
		for displayName := range config.Instance.Accounts {
			options = append(options, huh.NewOption(displayName, displayName))
		}

		var selectedOption string
		err := huh.NewSelect[string]().
			Options(options...).
			Title("Select an account to remove").
			Value(&selectedOption).
			WithTheme(huh.ThemeCatppuccin()).
			Run()
		checkError(err, "error prompting for account selection")

		delete(config.Instance.Accounts, selectedOption)

		err = config.Save()
		checkError(err, "error saving config")
	},
}

func init() {
	accountsCmd.AddCommand(accountsRemoveCmd)
}
