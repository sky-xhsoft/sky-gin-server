// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ResourceItemHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package videoHandlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/handlers"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ossUtil"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type ResourceItemHandler struct {
	db *gorm.DB
}

func (h *ResourceItemHandler) HandlerName() string {
	return "ResourceItemHandler"
}

func init() {
	handlers.Register("ResourceItemHandler", &ResourceItemHandler{})
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

	var resouce models.ChrResource
	if err := tx.Where(" id = ?", *req.ChrResourceId).First(&resouce).Error; err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}
	req.ProjectId = resouce.ProjectId

	if req.Type == "RTMP" {
		if req.RtmpUrl == "" || req.CutTimes == nil {
			ecode.Resp(c, ecode.ErrInvalidParam, "RTMP 类型需提供 RTMP_URL 和 CUT_TIMES")
			return
		}

		// 如果存在 RTMP 类型记录，更新而不是新建
		var existing models.ChrResourceItem
		err := tx.Where("CHR_RESOURCE_ID = ? AND TYPE = ? AND IS_ACTIVE = 'Y'", req.ChrResourceId, "RTMP").First(&existing).Error
		if err == nil {
			models.FillUpdateMeta(c, &req)
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

	models.FillCreateMeta(c, &req)

	if err := tx.Create(&req).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}

	h.updateResourceStats(tx, *req.ChrResourceId)

	ecode.SuccessResp(c, req)
}

// 更新资源明细
func (h *ResourceItemHandler) Update(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源明细ID")
		return
	}

	if req["ID"] == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源明细ID")
		return
	}

	models.FillUpdateMetaMap(c, req)

	if err := tx.Model(&models.ChrResourceItem{}).Where("ID = ?", req["ID"]).Updates(&req).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}

	var item models.ChrResourceItem
	if err := tx.Model(&models.ChrResourceItem{}).Where("ID =?", req["ID"]).First(&item); err == nil {
		return
	}

	h.updateResourceStats(tx, *item.ChrResourceId)

	ecode.SuccessResp(c, "更新成功")
}

// 删除明细资源
func (h *ResourceItemHandler) Delete(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	id := c.Query("ID")
	if id == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源明细ID")
		return
	}
	var item models.ChrResourceItem
	if err := tx.First(&item, id).Error; err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "资源明细不存在")
		return
	}

	if err := tx.Model(&models.ChrResourceItem{}).Where("ID = ?", id).Update("IS_ACTIVE", "N").Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	h.updateResourceStats(tx, *item.ChrResourceId)
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

// UploadSingleVideoFile 上传视频文件至 OSS 并自动写入 ChrResourceItem
func (h *ResourceItemHandler) UploadSingleVideoFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "文件读取失败")
		return
	}
	defer file.Close()

	resourceID := c.PostForm("resourceId")
	if resourceID == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}

	rid, err := utils.ParseUint(resourceID)
	if err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}

	var resouce models.ChrResource
	if err := h.db.Where(" id = ?", resourceID).First(&resouce).Error; err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}

	ossClient := ossUtil.GetClient()
	fileRecord, err := ossClient.UploadSingleFile(c, header, strconv.Itoa(int(rid)), header.Filename)
	if err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, "文件上传失败: "+err.Error())
		return
	}

	param, _ := json.Marshal(fileRecord)

	// 转为 MB，保留 1 位小数
	sizeInMB := float64(fileRecord.FileSize) / 1024.0 / 1024.0
	sizeRounded := math.Round(sizeInMB*10) / 10 // 保留 1 位小数

	item := models.ChrResourceItem{
		ChrResourceId: &rid,
		ProjectId:     resouce.ProjectId,
		Name:          header.Filename,
		Type:          "VIDEO",
		VideoUrl:      fileRecord.OSSURL,
		VideoFileSize: sizeRounded,
		VideoFileType: fileRecord.FileType,
		VideoParam:    string(param),
	}
	models.FillCreateMeta(c, &item)

	tx := utils.GetTx(c, h.db)
	if err := tx.Create(&item).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, "资源记录保存失败: "+err.Error())
		return
	}
	h.updateResourceStats(tx, rid)

	ecode.SuccessResp(c, item)
}

func (h *ResourceItemHandler) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "文件读取失败")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "至少上传一个文件")
		return
	}

	resourceID := c.PostForm("resourceId")
	if resourceID == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少项目ID或资源组ID")
		return
	}

	rid, err := utils.ParseUint(resourceID)
	if err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "ID格式错误")
		return
	}

	var resouce models.ChrResource
	if err := h.db.Where(" id = ?", resourceID).First(&resouce).Error; err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}

	ossClient := ossUtil.GetClient()
	tx := utils.GetTx(c, h.db)
	var result []models.ChrResourceItem

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		fileRecord, err := ossClient.UploadSingleFile(c, fileHeader, strconv.Itoa(int(rid)), fileHeader.Filename)
		file.Close()
		if err != nil {
			c.Error(err)
			continue
		}

		param, _ := json.Marshal(fileRecord)

		// 转为 MB，保留 1 位小数
		sizeInMB := float64(fileRecord.FileSize) / 1024.0 / 1024.0
		sizeRounded := math.Round(sizeInMB*10) / 10 // 保留 1 位小数

		item := models.ChrResourceItem{
			ChrResourceId: &rid,
			ProjectId:     resouce.ProjectId,
			Name:          fileHeader.Filename,
			Type:          "VIDEO",
			VideoUrl:      fileRecord.OSSURL,
			VideoFileSize: sizeRounded,
			VideoFileType: fileRecord.FileType,
			VideoParam:    string(param),
		}
		models.FillCreateMeta(c, &item)

		if err := tx.Create(&item).Error; err != nil {
			c.Error(err)
			continue
		}

		result = append(result, item)
	}
	h.updateResourceStats(tx, rid)

	if len(result) == 0 {
		ecode.Resp(c, ecode.ErrServer, "上传失败，请重试")
		return
	}
	ecode.SuccessResp(c, result)
}

func (h *ResourceItemHandler) updateResourceStats(tx *gorm.DB, rid uint) {
	var totalSize float64
	var totalQty int64
	var resource models.ChrResource

	if err := tx.Model(&models.ChrResource{}).Where("ID =?", rid).First(&resource); err == nil {
		return
	}

	tx.Model(&models.ChrResourceItem{}).
		Where("CHR_PROJECT_ID = ? AND TYPE = ? AND IS_ACTIVE = 'Y'", resource.ProjectId, "VIDEO").
		Count(&totalQty).
		Select("COALESCE(SUM(VIDEO_FILE_SIZE), 0)").
		Scan(&totalSize)

	tx.Model(&models.ChrProject{}).Where("ID = ?", resource.ProjectId).Updates(map[string]interface{}{
		"SIZE": totalSize,
		"QTY":  totalQty,
	})
}
