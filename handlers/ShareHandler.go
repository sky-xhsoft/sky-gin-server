// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ShareHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/19
// Project Description:
// ----------------------------------------------------------------------------

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/config"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ShareHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func (h *ShareHandler) HandlerName() string {
	return "ShareHandler"
}

func init() {
	Register("ShareHandler", &ShareHandler{})
}

func (h *ShareHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
	h.cfg = ctx.Config
}

func (h *ShareHandler) CreateShare(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req models.ChrShare
	if err := c.ShouldBindJSON(&req); err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "无效请求参数")
		return
	}

	// 生成唯一 6位字符串 key
	var genKey string
	for {
		genKey = utils.RandString(6)
		var count int64
		h.db.Model(&models.ChrShare{}).Where("`key` = ?", genKey).Count(&count)
		if count == 0 {
			break
		}
	}
	req.Key = genKey

	// 生成 4位数字密码
	genPwd := utils.RandDigit(4)

	// 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(genPwd), bcrypt.DefaultCost)
	if err != nil {
		ecode.Resp(c, ecode.ErrServer, "密码加密失败")
		return
	}
	req.Password = string(hashed)

	models.FillCreateMeta(c, &req)

	if err := tx.Create(&req).Error; err != nil {
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}

	// 响应明文 password 供展示
	ecode.SuccessResp(c, gin.H{
		"id":       req.ID,
		"key":      genKey,
		"url":      h.cfg.System.ShareUrl + genKey,
		"password": genPwd,
	})
}

func (h *ShareHandler) GetShare(c *gin.Context) {
	key := c.Param("key")
	inputPwd := c.Query("password")

	if key == "" || inputPwd == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少 key 或 password")
		return
	}

	var share models.ChrShare
	err := h.db.Where("`key` = ? AND is_active = 'Y'", key).First(&share).Error
	if err != nil {
		ecode.Resp(c, ecode.ErrRequest, "分享链接无效或已失效")
		return
	}

	// 验证密码
	if share.Password != "" {
		if bcrypt.CompareHashAndPassword([]byte(share.Password), []byte(inputPwd)) != nil {
			ecode.Resp(c, ecode.ErrUnauthorized, "密码错误")
			return
		}
	}
	var project models.ChrProject
	if share.ChrProjectID != nil {
		if err := h.db.Where(" ID = ?", share.ChrProjectID).First(&project).Error; err != nil {
			ecode.Resp(c, ecode.ErrRequest, "分享链接无效或资源不存在")
			return
		}

		var resource []models.ChrResource
		if share.ChrResourceID != nil {

			if err := h.db.Where(" ID = ?", share.ChrResourceID).First(&resource).Error; err != nil {
				ecode.Resp(c, ecode.ErrRequest, "分享链接无效或资源不存在")
				return
			}
			project.Resouse = resource

			var items []models.ChrResourceItem
			if share.ChrResourceItemID != nil {
				if err := h.db.Where(" ID = ? and type = 'VIDEO' ", share.ChrResourceItemID).First(&items).Error; err != nil {
					ecode.Resp(c, ecode.ErrRequest, "分享链接无效或资源不存在")
					return
				}
			} else {
				if err := h.db.Where(" CHR_RESOURCE_ID = ?", share.ChrResourceID).Find(&items).Error; err != nil {
					ecode.Resp(c, ecode.ErrRequest, "分享链接无效或资源不存在")
					return
				}
			}
			resource[0].Items = items
		} else {
			if share.SysDiskFileID == nil {
				if err := h.db.Where(" CHR_PROJECT_ID = ?", share.ChrProjectID).Find(&resource).Error; err != nil {
					ecode.Resp(c, ecode.ErrRequest, "分享链接无效或资源不存在")
					return
				}

				for k, v := range resource {
					var items []models.ChrResourceItem
					if err := h.db.Where(" CHR_RESOURCE_ID = ? and Type='VIDEO' and is_active='Y' ", v.ID).Find(&items).Error; err != nil {
						ecode.Resp(c, ecode.ErrRequest, "分享链接无效或资源不存在")
						return
					}
					resource[k].Items = items
				}
				project.Resouse = resource
			} else {
				var files []models.SysDiskFile
				if err := h.db.Where(" ID =? ", share.SysDiskFileID).Find(&files).Error; err != nil {
					ecode.Resp(c, ecode.ErrRequest, "分享链接无效或资源不存在")
					return
				}
				project.Files = files
			}
		}

	}

	// 可选：只返回关键数据（根据需要调整）
	ecode.SuccessResp(c, project)
}
