package config

import (
	"time"

	"github.com/spf13/viper"
)

type MongoCfg struct {
	Uri     string        `mapstructure:"mongo_uri"`
	Timeout time.Duration `mapstructure:"mongo_timeout"`
}

type Config struct {
	Mongo             MongoCfg      `mapstructure:",squash"`
	RecentItemsPeriod time.Duration `mapstructure:"recent_items_period"`
}

func Get() (*Config, error) {
	config := new(Config)

	viper.SetDefault("mongo_uri", "mongodb://localhost:27017")
	viper.SetDefault("mongo_timeout", "10s")

	viper.SetDefault("recent_items_period", "72h")

	viper.AutomaticEnv()

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
