package cmd

import (
	"context"

	"github.com/charmbracelet/huh"
	gh "github.com/cli/go-gh/v2/pkg/api"
	"github.com/google/go-github/v66/github"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nint8835/repo-archiver/pkg/config"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer repos between accounts",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var sourceAccount string
		var destinationAccount string
		var repository string

		var accountOptions []huh.Option[string]
		for displayName := range config.Instance.Accounts {
			accountOptions = append(accountOptions, huh.NewOption(displayName, displayName))
		}

		authenticatedClient, err := gh.DefaultHTTPClient()
		checkError(err, "error creating authenticated client")

		ghClient := github.NewClient(authenticatedClient)

		repositoryOptionsFunc := func() []huh.Option[string] {
			var repositoryOptions []huh.Option[string]

			repos, _, err := ghClient.Repositories.ListByUser(
				context.Background(),
				config.Instance.Accounts[sourceAccount].Name,
				nil,
			)
			if err != nil {
				log.Warn().Err(err).Msg("error listing repositories")
				return repositoryOptions
			}

			for _, repo := range repos {
				repositoryOptions = append(repositoryOptions, huh.NewOption(*repo.Name, *repo.Name))
			}

			return repositoryOptions
		}

		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Options(accountOptions...).
					Title("Source Account").
					Value(&sourceAccount),
				huh.NewSelect[string]().
					Options(accountOptions...).
					Title("Destination Account").
					Value(&destinationAccount),
				huh.NewSelect[string]().
					OptionsFunc(repositoryOptionsFunc, &sourceAccount).
					Title("Repository").
					Value(&repository),
			),
		).WithTheme(huh.ThemeCatppuccin()).Run()
		checkError(err, "error prompting for transfer information")
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
}
