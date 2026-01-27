package config

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type MongoCfg struct {
	Uri     string        `mapstructure:"uri"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type AuthCfg struct {
	PrivateKey string `mapstructure:"private_key"`
}

type Config struct {
	Mongo             MongoCfg      `mapstructure:"mongo"`
	Auth              AuthCfg       `mapstructure:"auth"`
	RecentItemsPeriod time.Duration `mapstructure:"recent_items_period"`
	RecentUsersCount  int64         `mapstructure:"recent_users_count"`
}

func Read() {
	viper.SetConfigType("hcl")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("mongo.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongo.timeout", "10s")

	viper.SetDefault("auth.private_key", "")

	viper.SetDefault("recent_items_period", "72h")
	viper.SetDefault("recent_users_count", 2)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Warn().Err(err).Send()
		return
	}

	log.Info().Msgf("Using config file: %s", viper.ConfigFileUsed())
}

func Get() (*Config, error) {
	config := new(Config)

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
