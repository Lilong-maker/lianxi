package RabbitMQ

import (
	"fmt"
	"log"
)

func SendStockDeductMsg(goodsId string, num int) {
	msg := fmt.Sprintf(`{"goodsId":"%s","num":%d}`, goodsId, num)
	SendMsg("stock_queue", msg)
}

func ConsumeStockDeduct() {
	SubscribeMsg("stock_queue", func(msg string) {
		log.Println("执行库存扣减，消息：", msg)
		log.Println("库存扣减日志已记录")
	})
}
