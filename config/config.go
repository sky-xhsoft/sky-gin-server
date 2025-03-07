// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: config.go
// Author: xhsoftware-skyzhou
// Created On: 2025/1/19
// Project Description:
// ----------------------------------------------------------------------------

package config

import (
	"encoding/json"
	"github.com/sky-xhsoft/sky-gin-server/pkg/log"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Redis Redis `mapstructure:"Redis" json:"Redis" yaml:"Redis"` //redis配置文件

	Mysql Mysql `mapstructure:"Mysql" json:"Mysql" yaml:"Mysql"` //Mysql配置文件

	System System `mapstructure:"System" json:"System" yaml:"System"` //system配置文件
}

type Mysql struct {
	DSN string `mapstructure:"dsn"`
}

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
}

type System struct {
	Port    string `mapstructure:"port" json:"port" yaml:"port"`          //服务开启端口
	LogPath string `mapstructure:"logPath" json:"logPath" yaml:"logPath"` //服务开启端口
}

var logger = log.GetLogger()

// 配置文件初始化
func GetConfigFile() string {
	//默认config 配置文件名称
	configFile := "config"
	//获取环境变量中配置的 环境
	env := os.Getenv("SKY_ENV")
	if env != "" {
		configFile = "config" + "_" + env
	}

	logger.Info("Loading config file: ", configFile)
	return configFile
}

// 加载config 配置文件
func LoadConfig(configFile string) (*Config, error) {
	var config *Config
	// 初始化 Viper
	viper.SetConfigName(configFile) // 配置文件名 (不需要扩展名)
	viper.AddConfigPath(".")        // 配置文件路径
	viper.SetConfigType("yml")      // 配置文件格式

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("read config error", err)
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		logger.Error("Unable to decode into struct, %v", err)
	}

	// 将 Config 结构体转化为 JSON
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		logger.Error("Error marshalling config:", err)
	}

	logger.Debugln("config file:", string(jsonData))

	return config, nil

}
