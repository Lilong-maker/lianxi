package RabbitMQ

import (
	"testing"
)

// 测试发送消息
func TestSendMsg(t *testing.T) {
	SendMsg("test_queue", "hello mq")
}

// 测试消费消息
func TestSubscribeMsg(t *testing.T) {
	SubscribeMsg("test_queue", func(msg string) {
		t.Log("消费：", msg)
	})
}

// 测试库存扣减
func TestStockDeduct(t *testing.T) {
	SendStockDeductMsg("1001", 1)
}
