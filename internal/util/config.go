package util

import "github.com/spf13/viper"

// All configuration of the application
type Config struct {
	Listen   string `mapstructure:"LISTEN"`
	CertPath string `mapstructure:"CERT_PATH"`
	KeyPath  string `mapstructure:"KEY_PATH"`
}

// Reads configuration from environment file then environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.SetDefault("LISTEN", ":17777")
	viper.SetDefault("CERT_PATH", "")
	viper.SetDefault("KEY_PATH", "")

	viper.BindEnv("LISTEN", "CERT_PATH", "KEY_PATH")
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	return
}
