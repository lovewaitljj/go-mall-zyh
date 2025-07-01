package model

import (
	"time"
)

// Demo demo用户分单表
type Demo struct {
	ID         uint32    `gorm:"primaryKey;column:id;type:int(11) unsigned;not null"`
	Account    string    `gorm:"column:account;type:varchar(64);not null;default:''"`  // 用户名称
	UpperLimit uint32    `gorm:"column:upper_limit;type:int(4);not null;default:0"`    // 接单上限
	Game       string    `gorm:"column:game;type:varchar(1024);not null"`              // 游戏（废弃）
	Lang       string    `gorm:"column:lang;type:varchar(512);not null"`               // 语言
	GameCat    string    `gorm:"column:game_cat;type:text;not null"`                   // 游戏、问题分类和系统标签
	Operator   string    `gorm:"column:operator;type:varchar(64);not null;default:''"` // op
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null"`             // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;default:null"`         // 更新时间
}

func (Demo) TableName() string {
	return "demo"
}
