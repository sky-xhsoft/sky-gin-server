// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ResourceItemHandler.go
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

type ResourceItemHandler struct {
	db *gorm.DB
}

func (h *ResourceItemHandler) HandlerName() string {
	return "ResourceItemHandler"
}

func init() {
	Register("ResourceItemHandler", &ResourceItemHandler{})
}

func (h *ResourceItemHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
}

// 资源明细新增
func (h *ResourceItemHandler) Create(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req models.ChrResourceItem
	if err := c.ShouldBindJSON(&req); err != nil {
		ecode.ErrorResp(c, ecode.ErrInvalidParam)
		return
	}

	if req.ChrResourceId == nil || *req.ChrResourceId == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}

	if req.Type == "RTMP" {
		if req.RtmpUrl == "" || req.CutTimes == nil {
			ecode.Resp(c, ecode.ErrInvalidParam, "RTMP 类型需提供 RTMP_URL 和 CUT_TIMES")
			return
		}

		// 如果存在 RTMP 类型记录，更新而不是新建
		var existing models.ChrResourceItem
		err := tx.Where("CHR_RESOURCE_ID = ? AND TYPE = ? AND IS_ACTIVE = 'Y'", req.ChrResourceId, "RTMP").First(&existing).Error
		if err == nil {
			utils.FillUpdateMeta(c, &req)
			if err := tx.Model(&models.ChrResourceItem{}).Where("ID = ?", existing.ID).Updates(&req).Error; err != nil {
				c.Error(err)
				ecode.Resp(c, ecode.ErrServer, err.Error())
				return
			}
			ecode.SuccessResp(c, existing.ID)
			return
		}
	} else if req.Type == "VIDEO" {
		if req.VideoUrl == "" {
			ecode.Resp(c, ecode.ErrInvalidParam, "VIDEO 类型需提供 VIDEO_URL")
			return
		}
	}

	utils.FillCreateMeta(c, &req)

	if err := tx.Create(&req).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, req.ID)
}

// 更新资源明细
func (h *ResourceItemHandler) Update(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req models.ChrResourceItem
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源明细ID")
		return
	}

	if req.ChrResourceId == nil || *req.ChrResourceId == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}

	if req.Type == "RTMP" {
		if req.RtmpUrl == "" || req.CutTimes == nil {
			ecode.Resp(c, ecode.ErrInvalidParam, "RTMP 类型需提供 RTMP_URL 和 CUT_TIMES")
			return
		}
	} else if req.Type == "VIDEO" {
		if req.VideoUrl == "" {
			ecode.Resp(c, ecode.ErrInvalidParam, "VIDEO 类型需提供 VIDEO_URL")
			return
		}
	}

	utils.FillUpdateMeta(c, &req)

	if err := tx.Model(&models.ChrResourceItem{}).Where("ID = ?", req.ID).Updates(&req).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, "更新成功")
}

// 删除明细资源
func (h *ResourceItemHandler) Delete(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	id := c.Query("id")
	if id == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源明细ID")
		return
	}
	if err := tx.Model(&models.ChrResourceItem{}).Where("ID = ?", id).Update("IS_ACTIVE", "N").Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, "已删除")
}

// 获取资源组资源明细
func (h *ResourceItemHandler) ListByResource(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	resourceID := c.Query("resourceId")
	if resourceID == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少或项目ID")
		return
	}

	var list []models.ChrResourceItem
	if err := tx.Where("CHR_RESOURCE_ID = ? AND IS_ACTIVE = 'Y'", resourceID).Find(&list).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, list)
}
