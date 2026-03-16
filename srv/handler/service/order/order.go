package order

import (
	"context"
	"errors"
	"fmt"
	"lianxi/RabbitMQ"
	"lianxi/pkg"
	"lianxi/srv/dasic/config"
	"lianxi/srv/handler/model"
	order2 "lianxi/srv/proto/order"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	order2.UnimplementedOrderServer
}

func (s *Server) OrderAdd(_ context.Context, in *order2.OrderAddReq) (*order2.OrderAddResp, error) {
	timeStr := time.Now().Format("20060102150405")
	uuidStr := uuid.New().String()
	orderSn := fmt.Sprintf("%v%v", timeStr, uuidStr[:8])
	total := 0.0
	var orderItems []*model.OrderItem
	for _, item := range in.List {
		var goods model.Goods
		err := goods.FindGoodsById(config.DB, item.GoodsId)
		if err != nil {
			return nil, errors.New("商品不存在")
		}
		subTotal := goods.Price * float64(item.Quantity)
		total += float64(subTotal)
		orderItem := model.OrderItem{
			OrderNo:    orderSn,
			GoodsID:    goods.ID,
			GoodsName:  goods.Name,
			GoodsPrice: float64(goods.Price),
			Num:        int(item.Quantity),
		}
		orderItems = append(orderItems, &orderItem)
		go func() {
			RabbitMQ.SendStockDeductMsg(strconv.Itoa(int(goods.ID)), int(item.Quantity))
		}()

	}
	order := model.Order{
		OrderNo:    orderSn,
		UserID:     int(in.UserID),
		TotalPrice: total,
		PayStatus:  0,
	}
	err := order.OrderAdd(config.DB)
	if err != nil {
		return nil, errors.New("订单创建失败")
	}
	for i, _ := range orderItems {
		orderItems[i].ID = uint(int(order.ID))
	}
	err = order.OrderItemAdd(config.DB, orderItems)
	if err != nil {
		return nil, errors.New("明细添加失败")
	}
	pay := pkg.Alipay(orderSn, total)

	return &order2.OrderAddResp{
		OrderSn: orderSn,
		Total:   float32(total),
		PayUrl:  pay,
	}, nil
}
