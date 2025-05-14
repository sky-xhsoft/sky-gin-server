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
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
)

func init() {
	Register("ApiDocHandler", &ApiDocHandler{})
}

type ApiDocHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func (h *ApiDocHandler) HandlerName() string {
	return "ApiDocHandler"
}

func (h *ApiDocHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
	h.redis = ctx.Redis
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

// EditPage 渲染接口编辑页面 /api/edit?id=xx
func (h *ApiDocHandler) EditPage(c *gin.Context) {
	var api models.SysApi
	idStr := c.Query("ID")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			_ = h.db.WithContext(context.Background()).First(&api, id).Error
		}
	}

	tmpl, err := template.ParseFiles("static/template/api_edit.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "模板加载失败: %v", err)
		return
	}

	// 传入编辑接口数据和方法/权限枚举
	data := map[string]interface{}{
		"Api":     &api,
		"Methods": []string{"GET", "POST", "PUT", "DELETE"},
		"Permissions": []struct {
			Value string
			Label string
		}{
			{"P", "Public"},
			{"S", "需要登录"},
			{"D", "数据权限"},
		},
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(c.Writer, data)
}

// Save 处理 POST /api/save 请求（JSON 提交）
func (h *ApiDocHandler) Save(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	var api models.SysApi
	idRaw := req["ID"]

	if idFloat, ok := idRaw.(float64); ok && int(idFloat) != 0 {
		if err := h.db.First(&api, req["ID"]).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "数据不存在"})
			return
		}
	}

	if idFloat, ok := idRaw.(float64); ok && int(idFloat) == 0 {

		if ID, err := utils.GetTableId(c, h.db, h.redis, "sys_api"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败", "detail": err.Error()})
			return
		} else {
			req["ID"] = ID
		}

		models.FillCreateMetaMap(c, req)
		if err := h.db.Table("sys_api").Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败", "detail": err.Error()})
			return
		}
	} else {

		models.FillUpdateMetaMap(c, req)
		if err := tx.Model(&models.SysApi{}).Where("ID = ?", req["ID"]).Updates(req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
			return
		}
	}

	// 返回 JSON 以便前端跳转到 /api/edit?id={ID}
	c.JSON(http.StatusOK, req)
}
