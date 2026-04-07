package config

type AppConfig struct {
	Name        string `mapstructure:"name" yaml:"name"`
	Version     string `mapstructure:"version" yaml:"version"`
	Debug       bool   `mapstructure:"debug" yaml:"debug"`
	Environment string `mapstructure:"environment" yaml:"environment"`
}
