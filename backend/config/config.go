package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var Conf = &Config{}

type Config struct {
	App      AppConfig      `mapstructure:"app" yaml:"app"`
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Logger   LoggerConfig   `mapstructure:"logger" yaml:"logger"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Redis    RedisConfig    `mapstructure:"redis" yaml:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt" yaml:"jwt"`
	MinIO    MinIOConfig    `mapstructure:"minio" yaml:"minio"`
	NoAuth   []string       `mapstructure:"no_auth" yaml:"no_auth"`
	Proxy    ProxyConfig    `mapstructure:"proxy" yaml:"proxy"`
}

func Init() error {
	viper.SetDefault("app.debug", true)
	viper.SetDefault("server.port", 7000)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("jwt.access_timeout", 60)
	viper.SetDefault("jwt.refresh_timeout", 168)
	viper.SetDefault("jwt.issuer", "webGL-720yun")

	if err := loadFromConfigFile(); err != nil {
		if err := initViperConfig(); err != nil {
			return fmt.Errorf("加载配置文件失败: %v", err)
		}

		if err := viper.Unmarshal(Conf); err != nil {
			return fmt.Errorf("解析viper配置失败: %v", err)
		}
	}

	if err := loadFromEnv(); err != nil {
		return err
	}

	return nil
}

func initViperConfig() error {
	viper.SetConfigName("config.dev")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../../config")
	viper.AddConfigPath("../..")

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

func loadFromEnv() error {
	if os.Getenv("APP_DEBUG") == "true" {
		Conf.App.Debug = true
	}
	return nil
}

func loadFromConfigFile() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	filename := currentDir + "/config/config.dev.yaml"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = currentDir + "/config.dev.yaml"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return fmt.Errorf("配置文件不存在: %s", filename)
		}
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	if err := yaml.Unmarshal(data, Conf); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

func GetViper() *viper.Viper {
	return viper.GetViper()
}
