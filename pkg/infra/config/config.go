package config

import "github.com/spf13/viper"

type MongoCfg struct {
	Uri string `mapstructure:"mongo_uri"`
}

type HTTPServerCfg struct {
}

type Config struct {
	Mongo MongoCfg `mapstructure:",squash"` // squash adds all values from cfgs here
}

func Get() (*Config, error) {
	config := new(Config)

	viper.SetDefault("mongo_uri", "mongodb://localhost:27017")

	viper.AutomaticEnv()

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
