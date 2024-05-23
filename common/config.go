package common

import (
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"log"
)

var Conf = new(config)

type config struct {
	Server *ServerConfig `mapstructure:"server" json:"server"`
	Logs   *LogsConfig   `mapstructure:"logs" json:"logs"`
}

type ServerConfig struct {
	Port      *Port    `mapstructure:"port" json:"port"`
	Path      string   `mapstructure:"path" json:"path"`
	Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
	Auth      *Auth    `mapstructure:"auth" json:"auth"`
	Limit     *Limit   `mapstructure:"limit" json:"limit"`
}

type Port struct {
	Http int `mapstructure:"http" json:"http"`
	Tcp  int `mapstructure:"tcp" json:"tcp"`
}

type Auth struct {
	Enable   bool   `mapstructure:"enable" json:"enable"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
}

type Limit struct {
	Enable   bool `mapstructure:"enable" json:"enable"`
	Duration int  `mapstructure:"duration" json:"duration"`
	Count    int  `mapstructure:"count" json:"count"`
}

type LogsConfig struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

// InitConfig 设置读取配置信息
func InitConfig() {
	viper.SetConfigFile("config.yml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s\n", err)
		return
	}
	err := viper.Unmarshal(&Conf)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return
	}
}
