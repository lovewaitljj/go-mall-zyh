package errcode

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

type AppError struct {
	code     int    `json:"code"`
	msg      string `json:"msg"`
	cause    error  `json:"cause"`    // 底层真正的错误：例如gorm.notfound
	occurred string `json:"occurred"` // 保存由底层错误导致AppErr发生时的位置
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	formattedErr := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Cause    string `json:"cause"`
		Occurred string `json:"occurred"`
	}{
		Code:     e.Code(),
		Msg:      e.Msg(),
		Occurred: e.occurred,
	}

	if e.cause != nil {
		formattedErr.Cause = e.cause.Error()
	}
	errByte, _ := json.Marshal(formattedErr)
	return string(errByte)
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Msg() string {
	return e.msg
}

// Wrap wrap用于没有自定义错误外的错误的封装
func Wrap(msg string, err error) *AppError {
	if err == nil {
		return nil
	}
	appErr := &AppError{code: -1, msg: msg, cause: err}
	appErr.occurred = getAppErrOccurredInfo()
	return appErr
}

// WithCause 用于自定义错误中的底层错误封装
func (e *AppError) WithCause(err error) *AppError {
	e.cause = err
	e.occurred = getAppErrOccurredInfo()
	return e
}

func newError(code int, msg string) *AppError {
	if code > -1 {
		if _, duplicated := codes[code]; duplicated {
			panic(fmt.Sprintf("错误码 %d 不能重复, 请检查后更换", code))
		}
	}
	codes[code] = struct{}{}
	return &AppError{code: code, msg: msg}
}

// getAppErrOccurredInfo 获取项目中调用Wrap或者WithCause方法时的程序位置, 方便排查问题
func getAppErrOccurredInfo() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	file = path.Base(file)
	funcName := runtime.FuncForPC(pc).Name()
	triggerInfo := fmt.Sprintf("func: %s, file: %s, line: %d", funcName, file, line)
	return triggerInfo
}
