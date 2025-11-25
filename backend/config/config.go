package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Conf 全局配置实例
var Conf = &Config{}

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

// Init 支持从环境变量和配置文件加载
func Init() error {
	// 设置默认值
	viper.SetDefault("app.debug", true)
	viper.SetDefault("server.port", 7000)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("jwt.access_timeout", 60)
	viper.SetDefault("jwt.refresh_timeout", 168)
	viper.SetDefault("jwt.issuer", "webGL-720yun")

	// 从配置文件加载
	if err := loadFromConfigFile(); err != nil {
		// 尝试使用viper加载配置文件作为备选方案
		if err := initViperConfig(); err != nil {
			return fmt.Errorf("加载配置文件失败: %v", err)
		}

		// 从viper读取配置到Conf结构体
		if err := viper.Unmarshal(Conf); err != nil {
			return fmt.Errorf("解析viper配置失败: %v", err)
		}
	}

	// 从环境变量加载
	if err := loadFromEnv(); err != nil {
		return err
	}

	return nil
}

// initViperConfig 初始化viper配置（兼容原有配置方式）
func initViperConfig() error {
	viper.SetConfigName("config.dev")
	viper.SetConfigType("yaml")

	// 添加多个可能的配置文件搜索路径
	// 当前目录的config子目录
	viper.AddConfigPath("./config")
	// 当前目录
	viper.AddConfigPath(".")
	// 上级目录的config子目录
	viper.AddConfigPath("../config")
	// 上级目录
	viper.AddConfigPath("..")
	// 上上级目录的config子目录
	viper.AddConfigPath("../../config")
	// 上上级目录
	viper.AddConfigPath("../..")

	// 设置默认值
	viper.SetDefault("server.port", 7000)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("jwt.access_timeout", 60)
	viper.SetDefault("jwt.refresh_timeout", 168)
	viper.SetDefault("jwt.issuer", "webGL-720yun")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	return nil
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv() error {
	// 读取环境配置
	if os.Getenv("APP_DEBUG") == "true" {
		Conf.App.Debug = true
	}
	return nil
}

// loadFromConfigFile 从配置文件加载配置
func loadFromConfigFile() error {
	// 获取当前文件所在目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	// 构造配置文件路径
	filename := currentDir + "/config/config.dev.yaml"

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 尝试从当前目录直接加载
		filename = currentDir + "/config.dev.yaml"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return fmt.Errorf("配置文件不存在: %s", filename)
		}
	}

	// 读取配置文件
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 将yaml内容解析到Conf结构体
	if err := yaml.Unmarshal(data, Conf); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

// GetViper 获取viper实例（供其他组件使用）
func GetViper() *viper.Viper {
	return viper.GetViper()
}
