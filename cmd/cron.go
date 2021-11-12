package cmd

import (
	"log"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/spf13/cobra"
	"github.com/ataklychev/gitlab_registry_cleaner/config"
	"github.com/ataklychev/gitlab_registry_cleaner/logger"
	"github.com/ataklychev/gitlab_registry_cleaner/repository"
	"github.com/ataklychev/gitlab_registry_cleaner/service"
)

// cronCmd represents the cron command
// nolint:gochecknoglobals
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Clean gitlab registry by cron, every day at specific time",
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		config, err := config.LoadConfig(".env")
		if err != nil {
			log.Fatal("cannot load config:", err)
		}

		// init logger
		_logger, _loggerSync := logger.NewLogger(config.Production)
		defer _loggerSync()

		startImmediately := false
		if len(config.CronTime) == 0 {
			startImmediately = true
			config.CronTime = time.Now().Format("15:04")
		}

		// every day at specific time start garbage collection
		err = gocron.Every(1).Day().At(config.CronTime).Do(func() {
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
