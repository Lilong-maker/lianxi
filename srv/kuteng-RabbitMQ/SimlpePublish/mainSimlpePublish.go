package main

import (
	"fmt"
	"lianxi/srv/kuteng-RabbitMQ/RabbitMQ"

	"github.com/go-redis/redis/v8"
)

func main() {
	// 初始化 Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "115.190.43.83:6379",
		Password: "4ay1nkal3u8ed77y",
		DB:       0,
	})

	sendMsg("test_queue_003", "Hello kuteng222!", redisClient)
}

func sendMsg(queue string, msg string, redisClient *redis.Client) {
	rabbitmq := RabbitMQ.NewRabbitMQSimple(queue)
	defer rabbitmq.Destory()
	rabbitmq.SetRedisClient(redisClient)

	err := rabbitmq.PublishSimple(msg)
	if err != nil {
		fmt.Printf("发送失败：%v\n", err)
		return
	}
	fmt.Println("发送成功")
}
