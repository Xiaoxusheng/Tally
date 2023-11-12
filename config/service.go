package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"time"
)

// 服务器
type Service struct {
	Port int    `json:"port" yaml:"port"`
	Name string `json:"name" yaml:"name"`
}

// 数据库
type Mysql struct {
	Database string `json:"database" yaml:"database"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Url      string `json:"url" yaml:"url"`
}

// redis
type Redis struct {
	Addr            string        `json:"addr" yaml:"addr"`
	Password        string        `json:"password" yaml:"password"`
	Db              int           `json:"db" yaml:"db"`
	PoolSize        int           `json:"PoolSize" yaml:"poolSize"`
	MinIdleConns    int           `json:"minIdleConns" yaml:"minIdleConns"`
	MaxIdleConns    int           `json:"maxIdleConns" yaml:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime" yaml:"connMaxIdleTime"`
}

// Jwt
type Jwt struct {
	Key  string
	Time time.Duration
}

type Configs struct {
	Service Service
	Mysql   Mysql
	Redis   Redis
	Jwt     Jwt
}

var Config Configs

func InitService() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("1", err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Println("初始化失败")
		return
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		err := viper.Unmarshal(&Config)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(Config)
	})
	fmt.Println(Config)
	viper.WatchConfig()
}
