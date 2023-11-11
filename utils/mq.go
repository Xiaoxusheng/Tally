package utils

import (
	"time"
)

// Async 异步更新
func Async() {
	t := time.NewTicker(time.Hour * 4)
	for {
		select {
		case <-t.C:
			//val := global.Global.Redis.Keys(global.Global.Ctx, "*").Val()

		}
	}
}
