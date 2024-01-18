package config

import (
	"Tally/global"
	"sync"
	"time"
)

type Mutex struct {
}

var once sync.Once

func InitMutex() {
	once.Do(
		func() {
			global.Global.Mutex = new(Mutex)
		})
}

/*简单版的互斥锁，后期使用lua脚本保证原子操作*/

// Lock 互斥锁模块
func (m *Mutex) Lock(key, id string) bool {
	if key == "" {
		return false
	}
	//设置5s 如果未释放锁自动释放
	return global.Global.Redis.SetNX(global.Global.Ctx, key, id, time.Microsecond*5000).Val()
}

// Unlock 释放锁
func (m *Mutex) Unlock(key string) bool {
	if key == "" {
		return false
	}
	return global.Global.Redis.Del(global.Global.Ctx, key).Val() == global.Success
}
