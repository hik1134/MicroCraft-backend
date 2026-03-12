package config

import (
	"os"
	perr "MicroCraft/pkg/errors"
	"gopkg.in/yaml.v3"
)

//全局配置对象
var Conf *Config

type Config struct {
	Mysql  MysqlConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
	Server ServerConfig `yaml:"server"`
	Email  EmailConfig  `yaml:"email" json:"email"`
	Jwt    JwtConfig    `yaml:"jwt"`
}

type JwtConfig struct {
	Secret        string `yaml:"secret"`
	ExpireSeconds int64  `yaml:"expire_seconds"`
}

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

//初始化配置
func InitConfig() error {
	file, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return perr.Wrap(perr.CONFIG_LOAD_FAIL, err)
	}
	var c Config
	if err := yaml.Unmarshal(file, &c); err != nil {
		return perr.Wrap(perr.CONFIG_PARSE_FAIL, err)
	}
	Conf = &c
	return nil
}