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
					//记录点赞
					if strings.Contains(list[i], global.BlogSetLikesKey) {
						fmt.Println(list[i][len("blogLikesSet"):], list[i])
						val := global.Global.Redis.Get(global.Global.Ctx, global.BlogLikesKey+list[i][len(global.BlogSetLikesKey):]).Val()
						fmt.Println("v", val)
						if val == "" {
							continue
						}
						err := dao.UpdateLikes(list[i][len(global.BlogSetLikesKey):], val)
						if err != nil {
							log.Println("插入出差", err)
							return
						}
					}
					//收藏
					if strings.Contains(list[i], global.BlogCollectRem) {
						//拼接
						val := global.Global.Redis.SMembers(global.Global.Ctx, list[i]).Val()
						global.Global.Log.Info(val)
						for j := 0; j < len(val); j++ {
							err := dao.DeleteBlogCollect(val[j])
							if err != nil {
								log.Println("删除出错", err)
								return
							}
						}
					}
					if strings.Contains(list[i], global.BlogCollects) {
						//拼接
						val := global.Global.Redis.SMembers(global.Global.Ctx, list[i]).Val()
						global.Global.Log.Error(val)

						global.Global.Log.Info(val)
						for j := 0; j < len(val); j++ {
							err := dao.UpdateBlogCollect(val[j])
							if err != nil {
								log.Println("更新出错", err)
								return
							}
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
		global.Global.Log.Info(global.Global.Redis.Type(global.Global.Ctx, list[i]).Val())
	}
	return valList
}
