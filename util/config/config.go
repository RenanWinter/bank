package config

import (
	"time"

	"github.com/spf13/viper"
)

type EnvironmentName string

const (
	LOCAL      EnvironmentName = "local"
	STAGING    EnvironmentName = "staging"
	PRODUCTION EnvironmentName = "production"
)

type Config struct {
	DBDriver          string          `mapstructure:"DB_DRIVER"`
	DBSource          string          `mapstructure:"DB_SOURCE"`
	ServerAddress     string          `mapstructure:"SERVER_ADDRESS"`
	Debug             bool            `mapstructure:"DEBUG"`
	Environment       EnvironmentName `mapstructure:"ENVIRONMENT"`
	TokenSymmetricKey string          `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDuration     time.Duration   `mapstructure:"TOKEN_DURATION_MILISECONDS"`
}

var Env Config

func LoadConfig(path string) (configuration Config, err error) {
	if Env.DBDriver != "" {
		return Env, nil
	}
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&configuration)

	if configuration.Environment == "" {
		configuration.Environment = PRODUCTION
	}

	Env = configuration

	return
}
