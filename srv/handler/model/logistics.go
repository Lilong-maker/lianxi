package model

import "gorm.io/gorm"

// Logistics 物流配送表
type Logistics struct {
	gorm.Model
	OrderNo          string  `gorm:"type:varchar(50);index;comment:订单编号"`
	LogisticsNo      string  `gorm:"type:varchar(50);uniqueIndex;comment:物流单号"`
	LogisticsCompany string  `gorm:"type:varchar(100);comment:物流公司"`
	ReceiverName     string  `gorm:"type:varchar(50);comment:收货人姓名"`
	ReceiverMobile   string  `gorm:"type:varchar(20);comment:收货人手机"`
	ReceiverProvince string  `gorm:"type:varchar(50);comment:收货省份"`
	ReceiverCity     string  `gorm:"type:varchar(50);comment:收货城市"`
	ReceiverArea     string  `gorm:"type:varchar(50);comment:收货区县"`
	ReceiverAddress  string  `gorm:"type:varchar(500);comment:收货详细地址"`
	SenderName       string  `gorm:"type:varchar(50);comment:发货人姓名"`
	SenderMobile     string  `gorm:"type:varchar(20);comment:发货人手机"`
	SenderProvince   string  `gorm:"type:varchar(50);comment:发货省份"`
	SenderCity       string  `gorm:"type:varchar(50);comment:发货城市"`
	SenderArea       string  `gorm:"type:varchar(50);comment:发货区县"`
	SenderAddress    string  `gorm:"type:varchar(500);comment:发货详细地址"`
	Weight           float64 `gorm:"type:decimal(10,2);comment:重量(kg)"`
	Freight          float64 `gorm:"type:decimal(10,2);comment:运费"`
	PaymentMethod    int     `gorm:"type:tinyint;comment:支付方式 1现付 2到付 3月结"`
	Status           int     `gorm:"type:tinyint;default:1;comment:状态 1待发货 2已发货 3运输中 4已签收 5已拒收 6已取消"`
	ShipmentTime     int64   `gorm:"type:bigint;comment:发货时间"`
	DeliveryTime     int64   `gorm:"type:bigint;comment:签收时间"`
	Remark           string  `gorm:"type:varchar(500);comment:备注"`
}
