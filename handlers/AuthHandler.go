// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: AuthHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/consts"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
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

// 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ecode.ErrorResp(c, ecode.ErrInvalidParam)
		return
	}

	// 检查手机号是否已存在
	var exist models.SysUser
	if err := h.db.Where("PHONE = ?", req.Phone).First(&exist).Error; err == nil {
		ecode.ErrorResp(c, ecode.ErrUserAlreadyExists)
		return
	}

	// 校验验证码
	key := "sms:code:" + req.Phone
	code, err := h.redis.Get(context.Background(), key).Result()
	if err != nil || code != req.Code {
		ecode.ErrorResp(c, ecode.ErrInvalidCode)
		return
	}

	passwordHash, _ := hash.HashPassword(req.Password)
	// 创建用户
	user := models.SysUser{
		Username: req.Username,
		Phone:    req.Phone,
		Password: passwordHash,
	}
	if err := h.db.Create(&user).Error; err != nil {
		ecode.ErrorResp(c, ecode.ErrServer)
		return
	}

	// 生成 token
	tk := token.GenerateToken()
	if err := token.SaveUser(h.redis, tk, &user); err != nil {
		ecode.ErrorResp(c, ecode.ErrTokenCreate)
		return
	}

	response.Ok(c, gin.H{
		"token": tk,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"trueName": user.TrueName,
			"phone":    user.Phone,
		},
	})
}

// 手机验证码登录（支持自动注册）
func (h *AuthHandler) PhoneLogin(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ecode.ErrorResp(c, ecode.ErrInvalidParam)
		return
	}

	// 校验验证码
	key := "sms:code:" + req.Phone
	code, err := h.redis.Get(context.Background(), key).Result()
	if err != nil || code != req.Code {
		ecode.ErrorResp(c, ecode.ErrInvalidCode)
		return
	}

	// 查询用户
	var user models.SysUser
	tx := h.db.Where("PHONE = ?", req.Phone).First(&user)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		//设置默认密码 abc123
		passwordHash, _ := hash.HashPassword(consts.DefaultPassword)
		// 自动注册新用户
		user = models.SysUser{
			Username: "u" + req.Phone,
			Phone:    req.Phone,
			Password: passwordHash, // 可选：设置空密码
		}
		if err := h.db.Create(&user).Error; err != nil {
			ecode.ErrorResp(c, ecode.ErrServer)
			return
		}
	} else if tx.Error != nil {
		ecode.ErrorResp(c, ecode.ErrServer)
		return
	}

	// 生成 token
	tk := token.GenerateToken()
	if err := token.SaveUser(h.redis, tk, &user); err != nil {
		ecode.ErrorResp(c, ecode.ErrTokenCreate)
		return
	}

	response.Ok(c, gin.H{
		"token": tk,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"trueName": user.TrueName,
			"phone":    user.Phone,
		},
	})
}

// 用户账号密码登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ecode.ErrorResp(c, ecode.ErrInvalidParam)
		return
	}

	var user models.SysUser
	if err := h.db.Where("USERNAME = ? AND IS_ACTIVE = 'Y'", req.Username).First(&user).Error; err != nil {
		ecode.ErrorResp(c, ecode.ErrUserNotFound)
		return
	}

	if !hash.CheckPassword(user.Password, req.Password) {
		ecode.ErrorResp(c, ecode.ErrPasswordWrong)
		return
	}

	tk := token.GenerateToken()
	if err := token.SaveUser(h.redis, tk, &user); err != nil {
		ecode.ErrorResp(c, ecode.ErrTokenCreate)
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

// 用户退出登录(若query中不存在token，则登出当前账户)
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenStr := c.Query("token")

	//如果可以从url中获取token 则删除对应token
	//否则当前登录用户推出
	if tokenStr == "" {
		tokenStr = c.GetHeader("token")
	}

	if tokenStr == "" {
		ecode.ErrorResp(c, ecode.ErrTokenEmpty)
		return
	}

	if err := token.DeleteToken(h.redis, tokenStr); err != nil {
		response.WithCode(c, ecode.ErrTokenError, err.Error(), nil)
		return
	}

	ecode.SuccessResp(c, nil)
}
