package config

import (
	"fmt"
	"os"
	"webGL-720yun/app/models"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Conf 全局配置变量
var Conf = new(models.Config)

// Init 支持从环境变量和配置文件加载
func Init() error {
	// 首先使用viper初始化配置（兼容原有的viper配置方式）
	if err := initViperConfig(); err != nil {
		return fmt.Errorf("初始化viper配置失败: %v", err)
	}

	// 从viper读取配置到Conf结构体
	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("解析viper配置失败: %v", err)
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
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

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
	var filename string
	if Conf.App.Debug {
		filename = "config/config.dev.yaml"
	} else {
		filename = "config/config.prod.yaml"
	}

	// 读取配置文件
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 将yaml内容解析到map中
	var configMap map[string]interface{}
	if err := yaml.Unmarshal(data, &configMap); err != nil {
		return err
	}

	// 使用mapstructure库解析map到结构体
	if err := mapstructure.Decode(configMap, Conf); err != nil {
		return err
	}

	return nil
}

// GetViper 获取viper实例（供其他组件使用）
func GetViper() *viper.Viper {
	return viper.GetViper()
}
