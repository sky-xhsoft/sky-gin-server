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
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/token"
)

func TokenAuth(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.GetHeader("Token")
		if t == "" {
			ecode.ErrorResp(c, ecode.ErrUnauthorized)
			c.Abort()
			return
		}

		user, err := token.GetUser(redisClient, t)
		if err != nil {
			ecode.ErrorResp(c, ecode.ErrTokenExpired)
			c.Abort()
			return
		}

		c.Set("User", user) // 注入上下文
		c.Set("token", t)
		c.Next()
	}
}
