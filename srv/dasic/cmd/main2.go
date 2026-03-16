package main

import (
	"lianxi/RabbitMQ"
	_ "lianxi/srv/dasic/inits"
)

func main() {
	RabbitMQ.SendStockDeductMsg("111", 1)
}
