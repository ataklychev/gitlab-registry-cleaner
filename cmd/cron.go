package cmd

import (
	"time"

	"github.com/ataklychev/gitlab-registry-cleaner/config"
	"github.com/ataklychev/gitlab-registry-cleaner/logger"
	"github.com/ataklychev/gitlab-registry-cleaner/repository"
	"github.com/ataklychev/gitlab-registry-cleaner/service"
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/cobra"
)

// cronCmd represents the cron command
// nolint:gochecknoglobals
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Clean gitlab registry by cron, every day at specific time",
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		config := config.LoadConfig()

		// init logger
		_logger, _loggerSync := logger.NewLogger(config.Debug)
		defer _loggerSync()

		startImmediately := false
		if len(config.CronTime) == 0 {
			startImmediately = true
			config.CronTime = time.Now().Format("15:04")
		}

		// every day at specific time start garbage collection
		err := gocron.Every(1).Day().At(config.CronTime).Do(func() {
			gitlabRepo := repository.NewGitlabRepository(repository.NewGitlabClient(config.AccessToken, config.BaseAPIURL), _logger)
			gc := service.NewGarbageCollectionService(config.Threshold, gitlabRepo, _logger)
			gc.Run()
		})
		if nil != err {
			_logger.Error(err)
		}

		// run job immediately
		if startImmediately {
			gocron.RunAll()
		}

		// start all the pending jobs
		<-gocron.Start()
	},
}

// nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(cronCmd)
}
