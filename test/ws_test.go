package test

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestWs(t *testing.T) {
	//
	rdb := redis.NewClient(&redis.Options{
		Addr:     "xlei.love:6379",
		Password: "admin123", // 没有密码，默认值
		DB:       0,
		PoolSize: 100,
	})
	ctx := context.Background()
	res, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	if res == "PONG" {
		log.Println("redis连接成功!")
	}
	fmt.Println()

	list := rdb.Keys(context.Background(), "*").Val()
	fmt.Println("list", list)
	for i := 0; i < len(list); i++ {
		if "set" == rdb.Type(context.Background(), list[i]).Val() {
			if strings.Contains(list[i], "blogLikesSet") {
				fmt.Println(list[i][len("blogLikesSet"):])
			}
		}
		fmt.Println(rdb.Type(context.Background(), list[i]).Val() == "set", strings.Contains(list[i], "blogLikesSet"))

	}

	fmt.Println(strings.ReplaceAll(time.Now().Format(time.DateOnly+"-"+time.TimeOnly), ":", "-"))
	dir, err := os.ReadDir("../log")
	if err != nil {
		return
	}
	//strings.Split(res.Name(), ".")
	parse, err := time.Parse(time.DateOnly, strings.Split(dir[0].Name(), ".")[0])
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(parse)

	fmt.Println(parse.Add(-time.Hour).After(parse), len(dir))
}

func Test1(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "xlei.love:6379",
		Password: "admin123", // 没有密码，默认值
		DB:       0,
		PoolSize: 100,
	})
	ctx := context.Background()
	res, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	if res == "PONG" {
		log.Println("redis连接成功!")
	}

	t2 := time.Now()
	t1 := time.Date(t2.Year(), 1, 1, 0, 0, 0, 0, t2.Location())
	//获取redis中的数据
	val := rdb.IncrBy(ctx, "key", 1).Val()
	fmt.Println(val, t2.Unix()-t1.Unix(), t2.Format(time.DateOnly))

	fmt.Println(t2.Unix() - t1.Unix() | val)

}
