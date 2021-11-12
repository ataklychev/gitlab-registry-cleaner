package cmd

import (
	"github.com/ataklychev/gitlab-registry-cleaner/config"
	"github.com/ataklychev/gitlab-registry-cleaner/logger"
	"github.com/ataklychev/gitlab-registry-cleaner/repository"
	"github.com/ataklychev/gitlab-registry-cleaner/service"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
// nolint:gochecknoglobals
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean gitlab registry",
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		config := config.LoadConfig()

		// init logger
		_logger, _loggerSync := logger.NewLogger(config.Debug)
		defer _loggerSync()

		// start gc
		gitlabRepo := repository.NewGitlabRepository(repository.NewGitlabClient(config.AccessToken, config.BaseAPIURL), _logger)
		gc := service.NewGarbageCollectionService(config.Threshold, gitlabRepo, _logger)
		gc.Run()
	},
}

// nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(cleanCmd)
}
