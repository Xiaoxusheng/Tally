package config

import (
	"Tally/global"
	"github.com/gammazero/workerpool"
	"sync"
)

/*
InitPool
初始化协程池
*/
var c = new(sync.Once)

func InitPool() {
	c.Do(
		func() {
			global.Global.Pool = workerpool.New(Config.Pool.Num)
		},
	)
}
