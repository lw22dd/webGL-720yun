package config

type LoggerConfig struct {
	Level      string `mapstructure:"level" yaml:"level"`
	Filename   string `mapstructure:"filename" yaml:"filename"`
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`
	Compress   bool   `mapstructure:"compress" yaml:"compress"`
}
