package response

// Demo demo用户分单表
type Demo struct {
	Account    string `json:"account"`
	UpperLimit uint32 `json:"upper_limit0"` // 接单上限
	Lang       string `json:"lang"`         // 语言
	Game       string `json:"game"`
	Operator   string `json:"operator"`   // op
	CreatedAt  string `json:"created_at"` // 创建时间
	UpdatedAt  string `json:"updated_at"` // 更新时间
}
