package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Production  bool   `mapstructure:"PRODUCTION"`
	Threshold   int    `mapstructure:"THRESHOLD"`
	BaseAPIURL  string `mapstructure:"BASE_API_URL"`
	AccessToken string `mapstructure:"ACCESS_TOKEN"`
	CronTime    string `mapstructure:"CRON_TIME"`
}

func LoadConfig(path string) (*Config, error) {
	config := new(Config)

	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(config)

	return config, err
}
