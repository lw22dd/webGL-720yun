package config

type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint" yaml:"endpoint"`
	AccessKey string `mapstructure:"access_key" yaml:"access_key"`
	SecretKey string `mapstructure:"secret_key" yaml:"secret_key"`
	Bucket    string `mapstructure:"bucket" yaml:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl" yaml:"use_ssl"`
	Region    string `mapstructure:"region" yaml:"region"`
}
