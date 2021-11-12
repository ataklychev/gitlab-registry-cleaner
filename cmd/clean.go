package cmd

import (
	"log"

	"github.com/ataklychev/gitlab_registry_cleaner/config"
	"github.com/ataklychev/gitlab_registry_cleaner/logger"
	"github.com/ataklychev/gitlab_registry_cleaner/repository"
	"github.com/ataklychev/gitlab_registry_cleaner/service"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
// nolint:gochecknoglobals
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean gitlab registry",
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		config, err := config.LoadConfig(".env")
		if err != nil {
			log.Fatal("cannot load config:", err)
		}

		// init logger
		_logger, _loggerSync := logger.NewLogger(config.Production)
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
