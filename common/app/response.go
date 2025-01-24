package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/common/errcode"
	"github.com/go-study-lab/go-mall/common/logger"
	"net/http"
)

type response struct {
	ctx        *gin.Context
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	RequestId  string      `json:"request_id"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
}

func NewResponse(c *gin.Context) *response {
	return &response{ctx: c}
}

func (r *response) SetPagination(p *pagination) *response {
	r.Pagination = p
	return r
}

func (r *response) Success(data interface{}) {
	r.Code = errcode.Success.Code()
	r.Msg = errcode.Success.Msg()
	requestId := ""
	if val, exists := r.ctx.Get("traceid"); exists {
		requestId = val.(string)
	}
	r.RequestId = requestId
	r.Data = data
	r.ctx.JSON(http.StatusOK, r)
}

func (r *response) SuccessOk() {
	r.Success("")
}

func (r *response) Error(err *errcode.AppError) {
	r.Code = err.Code()
	r.Msg = err.Msg()
	requestId := ""
	if val, exists := r.ctx.Get("traceid"); exists {
		requestId = val.(string)
	}
	r.RequestId = requestId
	// 兜底记一条错误日志
	logger.Error(r.ctx, "api_response_error", "err", err)
	r.ctx.JSON(http.StatusInternalServerError, r)
}
