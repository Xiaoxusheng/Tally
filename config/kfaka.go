package config

import (
	"Tally/global"
	"github.com/segmentio/kafka-go"
	"strconv"
	"sync"
)

var ks = new(sync.Once)

func InitKafka() {
	// 创建 Kafka 连接
	ks.Do(func() {
		conn, err := kafka.DialContext(global.Global.Ctx, "tcp", Config.Kafka.Address+":"+strconv.Itoa(Config.Kafka.Port))
		if err != nil {
			global.Global.Log.Error("connect kafka  fail err:", err)
			return
		}
		global.Global.Log.Info("kafka 连接成功！")
		global.Global.KafKa = conn
	})
}
