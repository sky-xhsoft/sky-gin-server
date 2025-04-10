// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: token_auth.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/pkg/response"
	"github.com/sky-xhsoft/sky-gin-server/pkg/token"
)

func TokenAuth(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.GetHeader("Token")
		if t == "" {
			response.WithCode(c, 401, "缺少 Token", nil)
			c.Abort()
			return
		}

		userID, err := token.Get(redisClient, t)
		if err != nil {
			response.WithCode(c, 401, "Token 无效或已过期", nil)
			c.Abort()
			return
		}

		c.Set("User", userID) // 注入上下文
		c.Next()
	}
}
