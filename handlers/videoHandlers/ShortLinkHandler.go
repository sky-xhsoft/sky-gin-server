// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ShortLinkHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/20
// Project Description:
// ----------------------------------------------------------------------------

package videoHandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/handlers"
	"github.com/sky-xhsoft/sky-gin-server/internal/shortlink"
	"net/http"
)

// ShortLinkHandler 短链处理器
type ShortLinkHandler struct {
}

// HandlerName 返回注册名
func (h *ShortLinkHandler) HandlerName() string {
	return "ShortLinkHandler"
}

func init() {
	handlers.Register("ShortLinkHandler", &ShortLinkHandler{})
}

// SetOption 注入上下文
func (h *ShortLinkHandler) SetOption(ctx *core.AppContext) {
}

// 创建短链请求结构
type createReq struct {
	URL  string `json:"url" binding:"required"`
	Name string `json:"name"`
}

// CreateShortLink 创建短链接口
func (h *ShortLinkHandler) CreateShortLink(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	shortURL, err := shortlink.CreateShortLink(req.URL, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_link": shortURL})
}
