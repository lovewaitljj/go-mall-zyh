package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/util"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// infrastructure 中存放项目运行需要的基础中间件

func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("traceid")
		pSpanId := c.Request.Header.Get("spanid")
		spanId := util.GenerateSpanID(c.Request.RemoteAddr)
		if traceId == "" { // 如果traceId 为空，证明是链路的发端，把它设置成此次的spanId，发端的spanId是root spanId,一般为网关
			traceId = spanId // trace 标识整个请求的链路, span则标识链路中的不同服务
		}
		c.Set("traceid", traceId)
		c.Set("spanid", spanId)
		c.Set("pspanid", pSpanId)
		c.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 包装一下 gin.ResponseWriter，通过这种方式拦截写响应
// 让gin写响应的时候先写到 bodyLogWriter 再写gin.ResponseWriter ，
// 这样利用中间件里输出访问日志时就能拿到响应了
// https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		// HTTP 请求的 body 是一个流（io.ReadCloser），数据只能被读取一次,这里需要重置
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		start := time.Now()
		//封装一个writer，往buf和http响应分别写一份
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		accessLog(c, "access_start", time.Since(start), reqBody, nil)
		defer func() {
			accessLog(c, "access_end", time.Since(start), reqBody, blw.body.String())
		}()
		c.Next()
		return
	}
}

func accessLog(c *gin.Context, accessType string, dur time.Duration, reqBody []byte, ResBody interface{}) {
	req := c.Request
	// TODO: 实现Token认证后再把访问日志里也加上token记录
	// token := c.Request.Header.Get("token")
	logger.Info(c, "AccessLog",
		"type", accessType,
		"ip", c.ClientIP(),
		//"token", token,
		"method", req.Method,
		"path", req.URL.Path,
		"query", req.URL.RawQuery, // ?后面的值，例如："q=golang&page=2"
		"req", string(reqBody),
		"res", ResBody,
		"time(ms)", int64(dur/time.Millisecond))
}

func GinPanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				//broken pipe错误  例子：用户在下载文件时突然关闭浏览器，或者长时间没有操作时浏览器自动断开了连接。
				// 不会打印堆栈信息，因为这只是网络错误。
				if brokenPipe {
					logger.Error(c, "http request broken pipe", "path", c.Request.URL.Path, "errcode", err, "request", string(httpRequest))
					c.Error(err.(error))
					c.Abort()
					return
				}
				// 程序panic，打印堆栈信息
				logger.Error(c, "http_request_panic", "path", c.Request.URL.Path, "errcode", err, "request", string(httpRequest), "stack", string(debug.Stack()))
				c.AbortWithError(http.StatusInternalServerError, err.(error))
			}
		}()
		c.Next()
	}
}
