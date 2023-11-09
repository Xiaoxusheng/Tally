package global

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Configs struct {
	Mysql *gorm.DB      `json:"mysql"`
	Redis *redis.Client `json:"redis"`
	Ctx   context.Context
}

var Global Configs
