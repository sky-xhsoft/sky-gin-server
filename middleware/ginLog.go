package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

// 自定义 ResponseWriter 用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // 记录响应体
	return w.ResponseWriter.Write(b)
}

func GinLogger(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// 读取请求体
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 重置请求体，防止后续读取失败

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw // 替换原始 Writer

		// 继续处理请求
		c.Next()

		// 记录日志
		log.Infow("HTTP Request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"body", string(bodyBytes), // 请求体
			"response", blw.body.String(),
			"ip", c.ClientIP(),
			"latency", time.Since(start),
			"userAgent", c.Request.UserAgent(),
		)
	}
}
