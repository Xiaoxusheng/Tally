package utils

import (
	"Tally/global"
	"time"
)

type Mutex struct {
}

/*简单版的互斥锁，后期使用lua脚本保证原子性操作*/

// Lock 互斥锁模块
func (m *Mutex) lock(key, id string) bool {
	if key == "" {
		return false
	}
	return global.Global.Redis.SetNX(global.Global.Ctx, key, id, time.Microsecond*5000).Val()
}

// Unlock释放锁
func (m *Mutex) unlock(key string) bool {
	if key == "" {
		return false
	}
	return global.Global.Redis.Del(global.Global.Ctx, key).Val() == global.Success
}
