package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		body := c.Request.Body

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw // 替换原始 Writer

		c.Next()

		log.Infow("HTTP Request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"body", body,
			"response", blw.body.String(),
			"ip", c.ClientIP(),
			"latency", time.Since(start),
			"userAgent", c.Request.UserAgent(),
		)
	}
}
