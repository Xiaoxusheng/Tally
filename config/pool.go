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
func InitPool() {
	once := sync.Once{}
	once.Do(
		func() {
			global.Global.Pool = workerpool.New(Config.Pool.Num)
			global.Global.Pool.StopWait()
		},
	)
}
