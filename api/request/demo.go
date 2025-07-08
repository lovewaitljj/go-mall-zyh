package request

// DemoCreate 请求
type DemoCreate struct {
	Account    string `json:"account" binding:"required"`
	UpperLimit uint32 `json:"upper_limit" ` // 接单上限
	Lang       string `json:"lang"`         // 语言
	Game       string `json:"game" binding:"required""`
}
