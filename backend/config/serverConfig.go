package config

type ServerConfig struct {
	Port         int    `mapstructure:"port" yaml:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"`
	Mode         string `mapstructure:"mode" yaml:"mode"`
}
