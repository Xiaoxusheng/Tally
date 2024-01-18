package config

import (
	"Tally/global"
	"github.com/segmentio/kafka-go"
	"log"
)

func InitKafka() {
	//broker := "KAFKA_BROKER_HOST:KAFKA_BROKER_PORT" // Kafka broker 地址和端口号
	//topic := "your-topic"                           // Kafka 主题名称

	// 创建 Kafka 连接
	conn, err := kafka.DialContext(global.Global.Ctx, "tcp", "xlei.love:9092")
	if err != nil {
		log.Fatal("Failed to connect to Kafka:", err)
	}
	defer conn.Close()
}
