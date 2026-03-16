package RabbitMQ

import (
	"log"

	"github.com/streadway/amqp"
)

const MQURL = "amqp://pengyilong:123456@115.190.43.83:5672/pengyilong"

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	Mqurl     string
}

func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
}

func (r *RabbitMQ) Destroy() {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "连接RabbitMQ失败")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "创建Channel失败")
	return rabbitmq
}

func (r *RabbitMQ) PublishSimple(message string) {
	_, err := r.channel.QueueDeclare(
		r.QueueName, false, false, false, false, nil,
	)
	r.failOnErr(err, "声明队列失败")

	err = r.channel.Publish(
		r.Exchange, r.QueueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	r.failOnErr(err, "发送消息失败")
}
