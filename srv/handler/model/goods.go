package model

import (
	"gorm.io/gorm"
)

type Goods struct {
	gorm.Model
	Name  string  `gorm:"type:varchar(30)"`
	Price float64 `gorm:"type:decimal(10,2)"`
	Num   int     `gorm:"type:int(11)"`
}

// ============ Create ============

func (o *Goods) GoodsAdd(db *gorm.DB) error {
	return db.Create(&o).Error
}

// ============ Read ============

// FindGoods 根据名称查找商品（检查是否存在）
func (o *Goods) FindGoods(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).First(&o).Error
}

// GetGoodsByID 根据ID获取商品
func (o *Goods) GetGoodsByID(db *gorm.DB, id uint) error {
	return db.Where("id = ?", id).First(&o).Error
}

// GoodsList 分页查询商品列表
func GoodsList(db *gorm.DB, page, pageSize int) ([]Goods, int64, error) {
	var goods []Goods
	var total int64

	// 查询总数
	err := db.Model(&Goods{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&goods).Error
	if err != nil {
		return nil, 0, err
	}

	return goods, total, nil
}

// ============ Update ============

func (o *Goods) GoodsUpdate(db *gorm.DB) error {
	return db.Model(&o).Updates(map[string]interface{}{
		"name":  o.Name,
		"price": o.Price,
		"num":   o.Num,
	}).Error
}

// ============ Delete ============

func (o *Goods) GoodsDelete(db *gorm.DB) error {
	return db.Delete(&o).Error
}

func (o *Goods) FindGoodsById(db *gorm.DB, id int64) interface{} {
	return db.Debug().Where("id = ?", id).Find(&o).Error
}
