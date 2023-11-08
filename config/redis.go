package config

import (
	"Tally/global"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:            Config.Redis.Addr,
		Password:        Config.Redis.Password, // 没有密码，默认值
		DB:              Config.Redis.Db,
		PoolSize:        Config.Redis.PoolSize,
		MinIdleConns:    Config.Redis.MinIdleConns,
		MaxIdleConns:    Config.Redis.MaxIdleConns,
		ConnMaxIdleTime: Config.Redis.ConnMaxIdleTime,
	})
	res, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println(err)
		return
	}
	if res == "PONG" {
		log.Println("连接redis 成功!")
	}
	global.Global.Redis = rdb
}
