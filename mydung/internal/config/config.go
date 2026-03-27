package config

import (
	"github.com/spf13/viper"
)
type Config struct {
	// Add your config fields here
	AppName string `mapstructure:"app_name"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
