// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: AuthHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/hash"
	"github.com/sky-xhsoft/sky-gin-server/pkg/response"
	"github.com/sky-xhsoft/sky-gin-server/pkg/token"
	"gorm.io/gorm"
)

func init() {
	Register("AuthHandler", &AuthHandler{})
}

type AuthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func (h *AuthHandler) HandlerName() string {
	return "AuthHandler"
}

func (h *AuthHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
	h.redis = ctx.Redis
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数格式错误")
		return
	}

	var user models.SysUser
	if err := h.db.Where("USERNAME = ? AND IS_ACTIVE = 'Y'", req.Username).First(&user).Error; err != nil {
		response.WithCode(c, 401, "用户不存在或被禁用", nil)
		return
	}

	if !hash.CheckPassword(user.Password, req.Password) {
		response.WithCode(c, 401, "密码错误", nil)
		return
	}

	tk := token.GenerateToken()
	if err := token.Save(h.redis, tk, user); err != nil {
		response.Fail(c, "Token 生成失败")
		return
	}

	response.Ok(c, gin.H{
		"token": tk,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"trueName": user.TrueName,
		},
	})
}
