package config

type DatabaseConfig struct {
	Host         string `mapstructure:"host" yaml:"host"`
	Port         int    `mapstructure:"port" yaml:"port"`
	Username     string `mapstructure:"username" yaml:"username"`
	Password     string `mapstructure:"password" yaml:"password"`
	Database     string `mapstructure:"database" yaml:"database"`
	Charset      string `mapstructure:"charset" yaml:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
}
