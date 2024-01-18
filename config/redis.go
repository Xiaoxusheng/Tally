package config

import (
	"Tally/global"
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var rdb *redis.Client
var x sync.Once

func InitRedis() {
	x.Do(
		func() {
			rdb = redis.NewClient(&redis.Options{
				Addr:            Config.Redis.Addr,
				Password:        Config.Redis.Password, // 没有密码，默认值
				DB:              Config.Redis.Db,
				PoolSize:        Config.Redis.PoolSize,
				MinIdleConns:    Config.Redis.MinIdleConns,
				MaxIdleConns:    Config.Redis.MaxIdleConns,
				ConnMaxIdleTime: Config.Redis.ConnMaxIdleTime * time.Second,
			})
			ctx := context.Background()
			res, err := rdb.Ping(ctx).Result()
			if err != nil {
				global.Global.Log.Error(err)
				panic(err)
			}
			if res == "PONG" {
				global.Global.Log.Info("redis连接成功!")
			}
			global.Global.Redis = rdb
			global.Global.Ctx = ctx
		})
}
