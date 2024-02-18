package utils

import (
	"Tally/global"
	_ "embed"
	"errors"
	"time"
)

// 内嵌lua脚本
//
//go:embed  lock.lua
var lock string

type Mutex struct {
	//自己的key
	key string
	// 每个协程独有的uid,避免误删除其他锁
	value string
	//锁自动释放时间
	expire time.Duration
}

// Lock 枷锁
func (m *Mutex) Lock() error {
	if m.expire == 0 {
		m.expire = time.Second * 5
	}
	return global.Global.Redis.SetNX(global.Global.Ctx, m.key, m.value, m.expire).Err()
}

// Unlock 解锁
func (m *Mutex) Unlock() error {
	/*	script := `
			-- KEYS[1] 就是分布式锁的key
		    -- ARGV[1] 就是预期的存在redis的value
		    if redis.call('get', KEYS[1]) == ARGV[1] then
		    --是自己的锁
		      return redis.call('del', KEYS[1])
		    else
		    --不是自己的锁
		      return  0
		    end
		`*/
	//这里用lua脚本来保证原子性操作
	res, err := global.Global.Redis.Eval(global.Global.Ctx, lock, []string{m.key}, m.value).Int()
	if err != nil || res != 1 {
		return errors.New(global.RedisLockErr)
	}
	return nil

}
