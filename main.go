package main

import (
	"errors"
	"github.com/go-study-lab/go-mall/common/app"
	"github.com/go-study-lab/go-mall/common/errcode"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/config"
	"github.com/go-study-lab/go-mall/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()
	// TODO: 后面会把应用日志统一收集到文件， 这里根据运行环境判断, 只在dev环境下才使用gin.Logger()输出信息到控制台
	g.Use(gin.Logger(), middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	g.GET("/config-read", func(c *gin.Context) {
		database := config.Database
		c.JSON(http.StatusOK, gin.H{
			"type":     database.Type,
			"max_life": database.MaxLifeTime,
		})
	})
	// logger门面的测试
	g.GET("/logger-test", func(c *gin.Context) {
		logger.Info(c, "logger test", "key", "keyName", "val", 2)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	g.GET("/panic-log-test", func(c *gin.Context) {
		var a map[string]string
		a["k"] = "v"
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   a,
		})
	})

	// 测试error包装返回
	g.GET("/customized-error-test", func(ctx *gin.Context) {
		// 使用wrap包装原因error生成 项目error
		err := errors.New("a dao error")
		appErr1 := errcode.Wrap("包装错误", err)
		logger.Error(ctx, "记录错误", "err", appErr1)

		// 预定义的ErrServer, 给其追加错误原因的error
		err = errors.New("a domain error")
		apiErr2 := errcode.ErrServer.WithCause(err)
		logger.Error(ctx, "API执行中出现错误", "err", apiErr2)
		ctx.JSON(http.StatusOK, gin.H{
			"code1": appErr1.Code(),
			"msg1":  appErr1.Msg(),
			"code2": apiErr2.Code(),
			"msg2":  apiErr2.Msg(),
		})
	})
	// 测试返回响应封装
	g.GET("/response-obj", func(c *gin.Context) {
		data := map[string]int{
			"a": 1,
			"b": 2,
		}
		app.NewResponse(c).Success(data)
		return
	})
	// 测试返回响应封装err
	g.GET("/response-err", func(c *gin.Context) {
		baseErr := errors.New("a base error")
		wrapErr := errcode.Wrap("a base error 错误", baseErr)
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(wrapErr))
		return
	})

	// 测试返回响应封装page
	g.GET("/response-page", func(c *gin.Context) {
		pagination := app.NewPagination(c)
		data := []struct {
			Name string
			Age  int
		}{
			{Name: "a", Age: 1},
			{Name: "b", Age: 2},
		}
		pagination.SetTotalRows(2)
		app.NewResponse(c).SetPagination(pagination).Success(data)
		return
	})

	g.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
