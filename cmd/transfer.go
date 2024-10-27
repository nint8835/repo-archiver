package cmd

import (
	"context"
	"errors"

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
		var repoName string

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
				// TODO: Paginate properly
				&github.RepositoryListByUserOptions{ListOptions: github.ListOptions{PerPage: 200}},
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
					Value(&repoName),
			),
		).WithTheme(huh.ThemeCatppuccin()).Run()
		checkError(err, "error prompting for transfer information")

		ctx := context.Background()
		ownerUsername := config.Instance.Accounts[sourceAccount].Name
		shouldBeArchived := config.Instance.Accounts[destinationAccount].IsArchive
		targetUsername := config.Instance.Accounts[destinationAccount].Name

		log.Info().Msgf("Transferring %s from %s to %s", repoName, ownerUsername, targetUsername)

		log.Debug().Msg("Getting current repository state")
		repository, _, err := ghClient.Repositories.Get(ctx, ownerUsername, repoName)
		checkError(err, "error getting repository")

		if repository.GetArchived() != shouldBeArchived {
			log.Debug().Msg("Archiving repository")
			_, _, err = ghClient.Repositories.Edit(ctx, ownerUsername, repoName, &github.Repository{Archived: github.Bool(shouldBeArchived)})
			checkError(err, "error archiving repository")
			log.Debug().Msg("Repository archived")
		}

		log.Debug().Msg("Transferring repository to new owner")
		_, _, err = ghClient.Repositories.Transfer(ctx, ownerUsername, repoName, github.TransferRequest{NewOwner: targetUsername})

		var acceptedError *github.AcceptedError
		if !errors.As(err, &acceptedError) {
			checkError(err, "error transferring repository")
		}

		log.Info().Msg("Repository transferred successfully!")
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
}
