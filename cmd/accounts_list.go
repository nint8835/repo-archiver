package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/nint8835/repo-archiver/pkg/config"
)

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all account aliases",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		tableData := pterm.TableData{
			{"Display Name", "Account Name", "Is Archive?"},
		}

		for displayName, account := range config.Instance.Accounts {
			tableData = append(tableData, []string{displayName, account.Name, pterm.Sprintf("%t", account.IsArchive)})
		}

		_ = pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	},
}

func init() {
	accountsCmd.AddCommand(accountsListCmd)
}
