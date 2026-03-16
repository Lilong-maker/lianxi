package api

import (
	"fmt"
	"lianxi/srv/dasic/config"
	"lianxi/srv/handler/model"

	"github.com/gin-gonic/gin"
)

func NotifyPay(c *gin.Context) {
	c.Request.ParseForm()
	fmt.Println("1111", c.Request.PostForm)
	TradeStatus := c.PostForm("trade_status")
	if TradeStatus != "TRADE_SUCCESS" {
		return
	}
	outTradeNo := c.PostForm("out_trade_no")
	if outTradeNo == "" {
		return
	}
	// 事务
	tx := config.DB.Begin()
	// 查询订单
	var order model.Order
	err := tx.Where("order_no = ?", outTradeNo).First(&order).Error
	if err != nil {
		tx.Rollback()
		return
	}
	// 幂等：已经支付就不再处理
	if order.PayStatus == 2 {
		tx.Commit()
		return
	}
	// 修改订单状态
	order.PayStatus = 2
	err = tx.Save(&order).Error
	if err != nil {
		tx.Rollback()
		return
	}
	// 查询订单明细
	var orderItems []model.OrderItem
	err = tx.Where("id=?", order.ID).Find(&orderItems).Error
	if err != nil {
		tx.Rollback()
		return
	}
	// 扣库存
	for _, items := range orderItems {
		var goods model.Goods
		err = tx.Where("id=?", items.GoodsID).First(&goods).Error
		if err != nil {
			tx.Rollback()
			return
		}
		goods.Num -= items.Num
		err = tx.Save(&goods).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
}
