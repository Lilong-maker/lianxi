package model

import "gorm.io/gorm"

// Member 注册会员表
type Member struct {
	gorm.Model
	MemberNo        string `gorm:"type:varchar(50);uniqueIndex;comment:会员编号"`
	Username        string `gorm:"type:varchar(50);uniqueIndex;comment:用户名"`
	Password        string `gorm:"type:varchar(100);comment:密码"`
	RealName        string `gorm:"type:varchar(50);comment:真实姓名"`
	Mobile          string `gorm:"type:varchar(20);uniqueIndex;comment:手机号"`
	Email           string `gorm:"type:varchar(100);comment:邮箱"`
	Avatar          string `gorm:"type:varchar(500);comment:头像"`
	Gender          int    `gorm:"type:tinyint;default:0;comment:性别 0未知 1男 2女"`
	Birthday        string `gorm:"type:varchar(20);comment:生日"`
	RegisterIP      string `gorm:"type:varchar(50);comment:注册IP"`
	LastLoginIP     string `gorm:"type:varchar(50);comment:最后登录IP"`
	LastLoginTime   int64  `gorm:"type:bigint;comment:最后登录时间"`
	Status          int    `gorm:"type:tinyint;default:1;comment:状态 0禁用 1启用"`
	MemberLevelID   int    `gorm:"type:int;comment:会员等级ID"`
	TotalPoints     int    `gorm:"type:int;default:0;comment:总积分"`
	AvailablePoints int    `gorm:"type:int;default:0;comment:可用积分"`
}
