package main

import (
	"fmt"
	"lianxi/srv/kuteng-RabbitMQ/RabbitMQ"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

func main() {
	// 初始化 Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "115.190.43.83:6379",
		Password: "4ay1nkal3u8ed77y",
		DB:       0,
	})

	err := subsribeMsg("test_queue_003", redisClient)
	if err != nil {
		log.Fatalf("订阅失败：%v", err)
	}

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown consumer...")
}

func subsribeMsg(topic string, redisClient *redis.Client) error {
	rabbitmq := RabbitMQ.NewRabbitMQSimple(topic)
	rabbitmq.SetRedisClient(redisClient)

	err := rabbitmq.SubsribeMsg(func(d amqp.Delivery) {
		fmt.Printf("Message: %s\n", d.Body)
		d.Ack(false)
	})
	if err != nil {
		return fmt.Errorf("订阅失败：%w", err)
	}
	return nil
}
