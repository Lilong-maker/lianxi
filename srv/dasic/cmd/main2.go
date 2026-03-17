package main

import (
	"lianxi/RabbitMQ"
	_ "lianxi/srv/dasic/inits"
)

func main() {
	RabbitMQ.SendStockDeductMsg("香蕉", 2)
}
