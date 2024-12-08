package errcode

type AppError struct {
	code  int    `json:"code"`
	msg   string `json:"msg"`
	cause string `json:"cause"`
}
