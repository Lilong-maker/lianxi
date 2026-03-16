package model

import "gorm.io/gorm"

// Inventory 库存表
type Inventory struct {
	gorm.Model
	ProductID      int     `gorm:"type:int;index;comment:商品ID"`
	WarehouseID    int     `gorm:"type:int;comment:仓库ID"`
	SKU            string  `gorm:"type:varchar(100);uniqueIndex;comment:商品SKU"`
	Stock          int     `gorm:"type:int;default:0;comment:库存数量"`
	LockedStock    int     `gorm:"type:int;default:0;comment:锁定库存"`
	AvailableStock int     `gorm:"type:int;default:0;comment:可用库存"`
	WarnStock      int     `gorm:"type:int;default:0;comment:预警库存"`
	CostPrice      float64 `gorm:"type:decimal(10,2);comment:成本价"`
	Location       string  `gorm:"type:varchar(100);comment:库位"`
	Status         int     `gorm:"type:tinyint;default:1;comment:状态 0禁用 1启用"`
}
