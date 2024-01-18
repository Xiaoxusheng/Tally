package global

import (
	"context"
	"github.com/gammazero/workerpool"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Configs struct {
	Mysql *gorm.DB               `json:"mysql"`
	Redis *redis.Client          `json:"redis"`
	Log   *log.Logger            `json:"log"`
	Ctx   context.Context        `json:"ctx"`
	Pool  *workerpool.WorkerPool `json:"pool"`
	KafKa *kafka.Conn            `json:"kafKa"`
	Mutex Mutex                  `json:"mutex"`
}

var (
	Global Configs
)
