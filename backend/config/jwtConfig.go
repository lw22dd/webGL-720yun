package config

type JWTConfig struct {
	Secret         string `mapstructure:"secret" yaml:"secret"`
	AccessTimeout  int    `mapstructure:"access_timeout" yaml:"access_timeout"`
	RefreshTimeout int    `mapstructure:"refresh_timeout" yaml:"refresh_timeout"`
	Issuer         string `mapstructure:"issuer" yaml:"issuer"`
}
