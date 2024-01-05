package config

import (
	"Tally/global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"time"
)

// Service 服务器
type Service struct {
	Port int    `json:"port" yaml:"port"`
	Name string `json:"name" yaml:"name"`
}

// Mysql 数据库
type Mysql struct {
	Database        string `json:"database" yaml:"database"`
	Username        string `json:"username" yaml:"username"`
	Password        string `json:"password" yaml:"password"`
	MaxIdleCons     int    `json:"maxIdleCons" yaml:"maxIdleCons"`
	MaxOpenCons     int    `json:"max_open_cons" yaml:"maxOpenCons"`
	ConnMaxLifetime int    `json:"connMaxLifetime" yaml:"connMaxLifetime"`
	Url             string `json:"url" yaml:"url"`
}

// Redis
type Redis struct {
	Addr            string        `json:"addr" yaml:"addr"`
	Password        string        `json:"password" yaml:"password"`
	Db              int           `json:"db" yaml:"db"`
	PoolSize        int           `json:"PoolSize" yaml:"poolSize"`
	MinIdleConns    int           `json:"minIdleConns" yaml:"minIdleConns"`
	MaxIdleConns    int           `json:"maxIdleConns" yaml:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime" yaml:"connMaxIdleTime"`
}

type Jwt struct {
	Key  string
	Time time.Duration
}

type Logs struct {
	Leave   string `json:"leave,omitempty" yaml:"leave"`
	Prefix  string `json:"prefix,omitempty" yaml:"prefix"`
	Path    string `json:"path,omitempty" yaml:"path"`
	Maxsize int    `json:"maxsize,omitempty" yaml:"maxsize"`
}

type Oauth2 struct {
	ClientID     string `json:"clientID,omitempty" yaml:"clientID"`
	ClientSecret string `json:"clientSecret,omitempty" yaml:"clientSecret"`
	AuthURL      string `json:"authURL,omitempty" yaml:"authURL"`
	TokenURL     string `json:"tokenURL,omitempty" yaml:"tokenURL"`
	RedirectURL  string `json:"redirectURL,omitempty" yaml:"redirectURL"`
	Scopes       string `json:"scopes,omitempty" yaml:"scopes"`
}

type SparkDesk struct {
	Appid     string `json:"appid,omitempty"  yaml:"appid"`
	ApiSecret string `json:"apiSecret,omitempty" yaml:"apiSecret"`
	ApiKey    string `json:"apiKey,omitempty" yaml:"apiKey"`
	HostUrl   string `json:"hostUrl" yaml:"hostUrl"`
}

type TencentCos struct {
	Url       string `json:"url" yaml:"url"`
	SecretId  string `json:"secretId,omitempty"  yaml:"secretId"`
	SecretKey string `json:"secretKey,omitempty" yaml:"secretKey"`
}

type Pool struct {
	Num int `json:"num" yaml:"num"`
}

type Configs struct {
	Service    Service
	Mysql      Mysql
	Redis      Redis
	Jwt        Jwt
	Logs       Logs
	Oauth2     Oauth2
	SparkDesk  SparkDesk
	TencentCos TencentCos
	Pool       Pool
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
		global.Global.Log.Error("初始化失败")
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
