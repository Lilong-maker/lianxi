package model

import "gorm.io/gorm"

// Product 电商商品表
type Product struct {
	gorm.Model
	ProductNo     string  `gorm:"type:varchar(50);uniqueIndex;comment:商品编号"`
	ProductName   string  `gorm:"type:varchar(200);not null;comment:商品名称"`
	ProductDesc   string  `gorm:"type:text;comment:商品描述"`
	ProductImg    string  `gorm:"type:varchar(1000);comment:商品图片(多图逗号分隔)"`
	ProductDetail string  `gorm:"type:text;comment:商品详情"`
	CategoryID    int     `gorm:"type:int;comment:分类ID"`
	BrandID       int     `gorm:"type:int;comment:品牌ID"`
	OriginalPrice float64 `gorm:"type:decimal(10,2);comment:原价"`
	SalePrice     float64 `gorm:"type:decimal(10,2);not null;comment:售价"`
	CostPrice     float64 `gorm:"type:decimal(10,2);comment:成本价"`
	Stock         int     `gorm:"type:int;default:0;comment:库存数量"`
	Sales         int     `gorm:"type:int;default:0;comment:销量"`
	Sort          int     `gorm:"type:int;default:0;comment:排序"`
	Status        int     `gorm:"type:tinyint;default:1;comment:状态 0下架 1上架 2删除"`
	IsHot         int     `gorm:"type:tinyint;default:0;comment:是否热销 0否 1是"`
	IsNew         int     `gorm:"type:tinyint;default:0;comment:是否新品 0否 1是"`
	IsRecommend   int     `gorm:"type:tinyint;default:0;comment:是否推荐 0否 1是"`
}
