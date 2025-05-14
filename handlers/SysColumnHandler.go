// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: SysColumnHandler.go
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

type SysColumnHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func (h *SysColumnHandler) HandlerName() string {
	return "SysColumnHandler"
}

func (h *SysColumnHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
	h.redis = ctx.Redis
}

// 编辑页面渲染
func (h *SysColumnHandler) EditPage(c *gin.Context) {
	idStr := c.Query("id")
	tableIdStr := c.Query("tableId")

	var column models.SysColumn
	if idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			h.db.First(&column, id)
		}
	} else if tableIdStr != "" {
		if tid, err := strconv.Atoi(tableIdStr); err == nil {
			column.SysTableId = uint(tid)
		}
	}

	tmpl, err := template.ParseFiles("static/template/sys_column.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "模板加载失败: %v", err)
		return
	}

	data := map[string]interface{}{
		"Data": column,
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(c.Writer, data)
}

// 保存逻辑
func (h *SysColumnHandler) Save(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	idRaw := req["ID"]
	if idFloat, ok := idRaw.(float64); ok && int(idFloat) == 0 {
		if newID, err := utils.GetTableId(c, h.db, h.redis, "sys_column"); err == nil {
			req["ID"] = newID
			models.FillCreateMetaMap(c, req)
			if err := h.db.Table("sys_column").Create(&req).Error; err == nil {
				c.JSON(http.StatusOK, req)
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "新增失败"})
		return
	}

	models.FillUpdateMetaMap(c, req)
	if err := tx.Model(&models.SysColumn{}).Where("ID = ?", req["ID"]).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, req)
}
