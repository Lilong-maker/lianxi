package model

import "gorm.io/gorm"

// Points 积分表
type Points struct {
	gorm.Model
	MemberID    int    `gorm:"type:int;index;comment:会员ID"`
	Points      int    `gorm:"type:int;comment:积分变化数量"`
	Balance     int    `gorm:"type:int;comment:变化后积分余额"`
	Type        int    `gorm:"type:tinyint;comment:积分类型 1获得 2消费 3过期 4管理员调整"`
	Reason      string `gorm:"type:varchar(200);comment:积分说明"`
	RelatedType string `gorm:"type:varchar(50);comment:关联类型(订单/签到/活动等)"`
	RelatedID   int    `gorm:"type:int;comment:关联ID"`
	ExpireTime  int64  `gorm:"type:bigint;comment:过期时间"`
}
