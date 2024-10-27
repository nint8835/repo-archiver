package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nint8835/repo-archiver/pkg/config"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use: "repo-archiver",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Interface("config", config.Instance).Msg("loaded config")
	},
}

func initLogging() {
	parsedLevel, err := zerolog.ParseLevel(logLevel)
	checkError(err, "failed to parse log level")

	zerolog.SetGlobalLevel(parsedLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func loadConfig() {
	err := config.Load()
	checkError(err, "failed to load config")
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level to use for output")
	_ = rootCmd.RegisterFlagCompletionFunc(
		"log-level",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			allLevels := []string{"debug", "info", "warn", "error", "fatal", "disabled"}

			var matchingLevels []string

			for _, level := range allLevels {
				if strings.HasPrefix(level, toComplete) {
					matchingLevels = append(matchingLevels, level)
				}
			}

			return matchingLevels, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(loadConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
