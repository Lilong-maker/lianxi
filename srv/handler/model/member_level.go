package model

import "gorm.io/gorm"

// MemberLevel 会员等级表
type MemberLevel struct {
	gorm.Model
	LevelName      string  `gorm:"type:varchar(50);uniqueIndex;comment:等级名称"`
	LevelNo        int     `gorm:"type:int;uniqueIndex;comment:等级编号"`
	RequiredPoints int     `gorm:"type:int;default:0;comment:所需积分"`
	DiscountRate   float64 `gorm:"type:decimal(5,2);default:1.00;comment:折扣率"`
	PointsRate     float64 `gorm:"type:decimal(5,2);default:1.00;comment:积分倍率"`
	Description    string  `gorm:"type:varchar(500);comment:等级描述"`
	Sort           int     `gorm:"type:int;default:0;comment:排序"`
	Status         int     `gorm:"type:tinyint;default:1;comment:状态 0禁用 1启用"`
}
