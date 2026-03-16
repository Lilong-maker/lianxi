package RabbitMQ

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

const MQURL = "amqp://pengyilong:123456@115.190.43.83:5672/pengyilong"

var (
	redisCtx = context.Background()
	redisKey = "mq:send:idempotent:"
)

type RabbitMQ struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	QueueName   string
	Exchange    string
	Key         string
	Mqurl       string
	redisClient *redis.Client
}

func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

func (r *RabbitMQ) SetRedisClient(client *redis.Client) {
	r.redisClient = client
}

func (r *RabbitMQ) Destory() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// PublishIdempotent 幂等性发送消息
func (r *RabbitMQ) PublishIdempotent(businessId string, message string) error {
	if r.redisClient == nil {
		return fmt.Errorf("Redis client 未设置")
	}

	msgHash := md5.Sum([]byte(message))
	idempotentKey := fmt.Sprintf("%s%s:%x", redisKey, businessId, msgHash)

	setnx, err := r.redisClient.SetNX(redisCtx, idempotentKey, "1", 24*time.Hour).Result()
	if err != nil {
		return fmt.Errorf("幂等性检查失败：%w", err)
	}

	if !setnx {
		log.Printf("消息已发送过，跳过：businessId=%s", businessId)
		return nil
	}

	err = r.publishMessage(message)
	if err != nil {
		r.redisClient.Del(redisCtx, idempotentKey)
		return fmt.Errorf("发送消息失败：%w", err)
	}

	fmt.Println("发送成功！")
	return nil
}

// PublishSimple 简单发送
func (r *RabbitMQ) PublishSimple(message string) error {
	return r.publishMessage(message)
}

func (r *RabbitMQ) publishMessage(message string) error {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("声明队列失败：%w", err)
	}

	err = r.channel.Confirm(false)
	if err != nil {
		return fmt.Errorf("启用 Confirm 模式失败：%w", err)
	}

	ackChan := make(chan uint64, 1)
	nackChan := make(chan uint64, 1)
	r.channel.NotifyConfirm(ackChan, nackChan)

	err = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		return fmt.Errorf("发送消息失败：%w", err)
	}

	select {
	case <-ackChan:
		return nil
	case <-nackChan:
		return fmt.Errorf("消息被 broker 拒绝")
	}
}

// SubsribeMsgIdempotent 带 Redis 幂等性的消费模式
func (r *RabbitMQ) SubsribeMsgIdempotent(handler func(amqp.Delivery)) error {
	if r.redisClient == nil {
		return fmt.Errorf("Redis client 未设置")
	}

	_, err := r.channel.QueueDeclare(
		r.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("声明队列失败：%w", err)
	}

	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("开始消费失败：%w", err)
	}

	go func() {
		for d := range msgs {
			messageID := d.MessageId
			if messageID == "" {
				messageID = fmt.Sprintf("%s:%d", r.QueueName, d.DeliveryTag)
			}

			isDuplicate, err := checkDuplicate(r.redisClient, messageID)
			if err != nil {
				log.Printf("幂等性检查失败：%v，重新入队", err)
				d.Nack(false, true)
				continue
			}
			if isDuplicate {
				log.Printf("重复消息，跳过：%s", messageID)
				d.Ack(false)
				continue
			}

			handler(d)
			markProcessed(r.redisClient, messageID)
			d.Ack(false)
		}
	}()

	return nil
}

// SubsribeMsg 简单消费模式
func (r *RabbitMQ) SubsribeMsg(handler func(amqp.Delivery)) error {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("声明队列失败：%w", err)
	}

	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("开始消费失败：%w", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("收到消息：%s", d.Body)
			handler(d)
			log.Println("消息处理成功")
		}
	}()

	return nil
}

func checkDuplicate(client *redis.Client, messageID string) (bool, error) {
	if messageID == "" {
		return false, nil
	}

	key := "mq:consume:" + messageID
	exists, err := client.SetNX(redisCtx, key, "1", 3*time.Hour).Result()
	if err != nil {
		return false, err
	}

	return !exists, nil
}

func markProcessed(client *redis.Client, messageID string) error {
	if messageID == "" {
		return nil
	}

	key := "mq:consume:" + messageID
	_, err := client.Expire(redisCtx, key, 3*time.Hour).Result()
	if err != nil {
		return err
	}

	return nil
}
