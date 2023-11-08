package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Configs struct {
	Mysql *gorm.DB      `json:"mysql"`
	Redis *redis.Client `json:"redis"`
}

var Global Configs
