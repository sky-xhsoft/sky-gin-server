// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ResourceHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"gorm.io/gorm"
)

type ResourceHandler struct {
	db *gorm.DB
}

func (h *ResourceHandler) HandlerName() string {
	return "ResourceHandler"
}

func init() {
	Register("ResourceHandler", &ResourceHandler{})
}

func (h *ResourceHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
}

// 创建资源组
func (h *ResourceHandler) CreateResource(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req models.ChrResource
	if err := c.ShouldBindJSON(&req); err != nil {
		ecode.ErrorResp(c, ecode.ErrInvalidParam)
		return
	}
	utils.FillCreateMeta(c, &req)

	if err := tx.Create(&req).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, req.ID)
}

// 更新资源组
func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req models.ChrResource
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源ID")
		return
	}
	utils.FillUpdateMeta(c, &req)

	if err := tx.Model(&models.ChrResource{}).Where("ID = ?", req.ID).Updates(&req).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, "更新成功")
}

// 删除资源组（逻辑删除）
func (h *ResourceHandler) DeleteResource(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	id := c.Query("id")
	if id == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源ID")
		return
	}
	if err := tx.Model(&models.ChrResource{}).Where("ID = ?", id).Update("IS_ACTIVE", "N").Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, "删除成功")
}

// 查询资源组列表
func (h *ResourceHandler) ListResources(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	projectID := c.Query("projectId")
	if projectID == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少项目ID")
		return
	}

	var list []models.ChrResource
	if err := tx.Where("CHR_PROJECT_ID = ? AND IS_ACTIVE = 'Y'", projectID).Find(&list).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, list)
}
