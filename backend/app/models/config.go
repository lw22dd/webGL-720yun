package models

// Config 系统配置
type Config struct {
	App      AppConfig      `mapstructure:"app" yaml:"app"`
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Logger   LoggerConfig   `mapstructure:"logger" yaml:"logger"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Redis    RedisConfig    `mapstructure:"redis" yaml:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt" yaml:"jwt"`
	NoAuth   []string       `mapstructure:"no_auth" yaml:"no_auth"`
	Proxy    ProxyConfig    `mapstructure:"proxy" yaml:"proxy"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name        string `mapstructure:"name" yaml:"name"`
	Version     string `mapstructure:"version" yaml:"version"`
	Debug       bool   `mapstructure:"debug" yaml:"debug"`
	Environment string `mapstructure:"environment" yaml:"environment"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port" yaml:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"`
	Mode         string `mapstructure:"mode" yaml:"mode"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level" yaml:"level"`
	Filename   string `mapstructure:"filename" yaml:"filename"`
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`
	Compress   bool   `mapstructure:"compress" yaml:"compress"`
}

// DatabaseConfig 数据库配置
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

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
	PoolSize int    `mapstructure:"pool_size" yaml:"pool_size"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret         string `mapstructure:"secret" yaml:"secret"`
	AccessTimeout  int    `mapstructure:"access_timeout" yaml:"access_timeout"`   // 访问令牌过期时间(分钟)
	RefreshTimeout int    `mapstructure:"refresh_timeout" yaml:"refresh_timeout"` // 刷新令牌过期时间(小时)
	Issuer         string `mapstructure:"issuer" yaml:"issuer"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Host    string   `mapstructure:"host" yaml:"host"`
	Port    int      `mapstructure:"port" yaml:"port"`
	NoProxy []string `mapstructure:"no_proxy" yaml:"no_proxy"`
}