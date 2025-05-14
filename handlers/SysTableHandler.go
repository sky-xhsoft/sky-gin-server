// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: SysTableHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/5/15
// Project Description:
// ----------------------------------------------------------------------------

package handlers

import (
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

type SysTableHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func (h *SysTableHandler) HandlerName() string {
	return "SysTableHandler"
}

func (h *SysTableHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
	h.redis = ctx.Redis
}

// EditPage 渲染 /api/sys_table/edit 页面
func (h *SysTableHandler) EditPage(c *gin.Context) {
	var table models.SysTable
	idStr := c.Query("id")
	if idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			_ = h.db.First(&table, id).Error
		}
	}

	tmpl, err := template.ParseFiles("static/template/table_add_edit.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "模板加载失败: %v", err)
		return
	}

	data := map[string]interface{}{
		"SYS_TABLE": &table,
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(c.Writer, data)
}

// save 方法 /api/sys_table/save 接口
func (h *SysTableHandler) Save(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	var table models.SysTable
	idRaw := req["ID"]

	if idFloat, ok := idRaw.(float64); ok && int(idFloat) != 0 {
		if err := h.db.First(&table, int(idFloat)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
			return
		}
	}

	if idFloat, ok := idRaw.(float64); ok && int(idFloat) == 0 {
		if ID, err := utils.GetTableId(c, h.db, h.redis, "sys_table"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败", "detail": err.Error()})
			return
		} else {
			req["ID"] = ID
		}

		models.FillCreateMetaMap(c, req)
		if err := h.db.Table("sys_table").Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败", "detail": err.Error()})
			return
		}
	} else {
		models.FillUpdateMetaMap(c, req)
		if err := tx.Model(&models.SysTable{}).Where("ID = ?", req["ID"]).Updates(req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, req)
}
