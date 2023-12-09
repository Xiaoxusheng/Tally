package utils

import (
	"Tally/dao"
	"Tally/global"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Set 这个脚本异步写入数据库
func Set() {
	t := time.NewTicker(time.Minute * 10)
	for {
		list := Get("set")
		select {
		case <-t.C:
			fmt.Println("list", list)
			if list == nil {
				continue
			}
			go func() {
				for i := 0; i < len(list); i++ {
					if strings.Contains(list[i], global.BlogSetLikesKey) {
						fmt.Println(list[i][len("blogLikesSet"):])
						val := global.Global.Redis.Get(global.Global.Ctx, global.BlogLikesKey+list[i][len("blogLikesSet"):]).Val()
						err := dao.UpDateLikes(list[i][len(global.BlogSetLikesKey):], val)
						if err != nil {
							log.Println("插入出差", err)
							return
						}
					}
				}
			}()
		}
	}
}

func Get(t string) []string {
	list := global.Global.Redis.Keys(global.Global.Ctx, "*").Val()
	valList := make([]string, 0, len(list))

	for i := 0; i < len(list); i++ {
		if t == global.Global.Redis.Type(global.Global.Ctx, list[i]).Val() {
			valList = append(valList, list[i])
		}
		fmt.Println(global.Global.Redis.Type(global.Global.Ctx, list[i]).Val())
	}
	return valList
}
