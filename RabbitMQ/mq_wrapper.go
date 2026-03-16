package RabbitMQ

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"lianxi/srv/dasic/inits"
	"log"
	"time"
)

func SendMsg(topic string, msg string) {
	mq := NewRabbitMQSimple(topic)
	defer mq.Destroy()
	mq.PublishSimple(msg)
	fmt.Println("消息入队成功:", topic, " | ", msg)
}

func SubscribeMsg(topic string, handler func(msg string)) {
	mq := NewRabbitMQSimple(topic)
	defer mq.Destroy()
	q, _ := mq.channel.QueueDeclare(
		topic, false, false, false, false, nil,
	)
	msgs, _ := mq.channel.Consume(
		q.Name, "", true, false, false, false, nil,
	)
	go func() {
		for d := range msgs {
			msg := string(d.Body)
			if !isDuplicate(msg) {
				handler(msg)
			} else {
				log.Println("重复消息，已过滤:", msg)
			}
		}
	}()
	fmt.Println("开始监听队列:", topic)
	select {}
}
func isDuplicate(msg string) bool {
	hash := getMsgHash(msg)
	key := "mq:msg:" + hash
	success, _ := inits.Rdb.SetNX(inits.Ctx, key, "1", 86400*time.Second).Result()
	return !success
}
func getMsgHash(msg string) string {
	h := md5.New()
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}
