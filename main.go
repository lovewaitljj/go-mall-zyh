package main

import (
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/config"
	"github.com/go-study-lab/go-mall/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()
	// TODO: 后面会把应用日志统一收集到文件， 这里根据运行环境判断, 只在dev环境下才使用gin.Logger()输出信息到控制台
	g.Use(gin.Logger(), middleware.StartTrace())
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	g.GET("/config-read", func(c *gin.Context) {
		logger.ZapLoggerTest("123")
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
	g.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
