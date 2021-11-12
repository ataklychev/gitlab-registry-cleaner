package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Debug       bool   `mapstructure:"DEBUG"`
	Threshold   int    `mapstructure:"THRESHOLD"`
	BaseAPIURL  string `mapstructure:"BASE_API_URL"`
	AccessToken string `mapstructure:"ACCESS_TOKEN"`
	CronTime    string `mapstructure:"CRON_TIME"`
}

func LoadConfig() *Config {
	defaultThreshold := 3

	viper.AutomaticEnv()
	viper.SetDefault("DEBUG", true)
	viper.SetDefault("THRESHOLD", defaultThreshold)
	viper.SetDefault("CRON_TIME", "01:11")

	return &Config{
		Debug:       viper.GetBool("DEBUG"),
		AccessToken: viper.GetString("ACCESS_TOKEN"),
		BaseAPIURL:  viper.GetString("BASE_API_URL"),
		Threshold:   viper.GetInt("THRESHOLD"),
		CronTime:    viper.GetString("CRON_TIME"),
	}
}
