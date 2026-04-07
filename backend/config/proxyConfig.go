package config

type ProxyConfig struct {
	Host    string   `mapstructure:"host" yaml:"host"`
	Port    int      `mapstructure:"port" yaml:"port"`
	NoProxy []string `mapstructure:"no_proxy" yaml:"no_proxy"`
}
