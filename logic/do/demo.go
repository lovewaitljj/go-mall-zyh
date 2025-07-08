package do

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

// Demo demo用户分单表
type Demo struct {
	ID         uint32                `json:"id"`
	Account    string                `json:"account"`
	UpperLimit uint32                `json:"upper_limit0"` // 接单上限
	Lang       string                `json:"lang"`         // 语言
	Game       string                `json:"game"`
	Operator   string                `json:"operator"`   // op
	IsDel      soft_delete.DeletedAt `json:"is_del"`     // 软删除，查询会自动带上is_del = 0
	CreatedAt  time.Time             `json:"created_at"` // 创建时间
	UpdatedAt  time.Time             `json:"updated_at"` // 更新时间
}
