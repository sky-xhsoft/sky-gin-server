// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ApiDocHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/5/14
// Project Description:
// ----------------------------------------------------------------------------

package handlers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

func init() {
	Register("ApiDocHandler", &ApiDocHandler{})
}

type ApiDocHandler struct {
	db *gorm.DB
}

func (h *ApiDocHandler) HandlerName() string {
	return "ApiDocHandler"
}

func (h *ApiDocHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
}

func (h *ApiDocHandler) DocPage(c *gin.Context) {
	var apis []models.SysApi
	err := h.db.WithContext(context.Background()).Where("IS_ACTIVE = ?", "Y").Find(&apis).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "接口数据获取失败"})
		return
	}

	tmpl, err := template.New("api_docs.html").Funcs(template.FuncMap{
		"parseJson": func(s string) []map[string]string {
			var arr []map[string]string
			_ = json.Unmarshal([]byte(s), &arr)
			return arr
		},
	}).ParseFiles("static/template/api_docs.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "模板加载失败: %v", err)
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(c.Writer, apis)
}
