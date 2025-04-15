// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: DiskHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/15
// Project Description:
// ----------------------------------------------------------------------------

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ossUtil"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"gorm.io/gorm"
	"path/filepath"
	"strconv"
	"strings"
)

type DiskHandler struct {
	db *gorm.DB
}

func (h *DiskHandler) HandlerName() string {
	return "DiskHandler"
}

func init() {
	Register("DiskHandler", &DiskHandler{})
}

func (h *DiskHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
}

// Upload 上传单个或多个文件到 OSS
func (h *DiskHandler) Upload(c *gin.Context) {
	projectIDStr := c.PostForm("projectId")
	parentIDStr := c.PostForm("parentId")
	projectID, _ := strconv.ParseUint(projectIDStr, 10, 64)
	parentID, _ := strconv.ParseUint(parentIDStr, 10, 64)

	form, err := c.MultipartForm()
	if err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "无效的文件上传表单")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "未上传任何文件")
		return
	}

	tx := utils.GetTx(c, h.db)
	ossClient := ossUtil.GetClient()
	var results []models.SysDiskFile

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		ext := filepath.Ext(fileHeader.Filename)
		//key := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		ossFile, err := ossClient.UploadSingleFile(c, fileHeader, strconv.Itoa(int(projectID)))
		if err != nil {
			continue
		}

		diskFile := models.SysDiskFile{
			ProjectId: uint(projectID),
			ParentId:  uint(parentID),
			FileName:  fileHeader.Filename,
			FileType:  "F",
			FilePath:  ossFile.OSSURL,
			FileExt:   strings.TrimPrefix(ext, "."),
			FileSize:  fileHeader.Size,
			MimeType:  fileHeader.Header.Get("Content-Type"),
		}
		utils.FillCreateMeta(c, &diskFile)
		tx.Create(&diskFile)
		results = append(results, diskFile)
	}

	ecode.SuccessResp(c, results)
}

// 创建文件夹
func (h *DiskHandler) CreateFolder(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req models.SysDiskFile
	if err := c.ShouldBindJSON(&req); err != nil {
		ecode.ErrorResp(c, ecode.ErrInvalidParam)
		return
	}
	if req.FileName == "" || req.FileType != "D" {
		ecode.Resp(c, ecode.ErrInvalidParam, "必须指定文件夹名称和类型为D")
		return
	}
	utils.FillCreateMeta(c, &req)

	if err := tx.Create(&req).Error; err != nil {
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, req.ID)
}

// 文件列表
func (h *DiskHandler) ListFiles(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	projectId := c.Query("projectId")
	parentId := c.Query("parentId")

	var list []models.SysDiskFile
	db := tx.Where("IS_ACTIVE = 'Y'")
	if projectId != "" {
		db = db.Where("PROJECT_ID = ?", projectId)
	}
	if parentId != "" {
		db = db.Where("PARENT_ID = ?", parentId)
	} else {
		db = db.Where("PARENT_ID IS NULL")
	}
	if err := db.Find(&list).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, list)
}

// 删除文件（逻辑删除）
func (h *DiskHandler) DeleteFile(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	id := c.Query("id")
	if id == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少文件ID")
		return
	}

	if err := tx.Model(&models.SysDiskFile{}).Where("ID = ?", id).Update("IS_ACTIVE", "N").Error; err != nil {
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, "删除成功")
}

// 搜索文件
func (h *DiskHandler) SearchFile(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	keyword := c.Query("keyword")
	projectId := c.Query("projectId")

	if keyword == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "请输入关键词")
		return
	}

	var list []models.SysDiskFile
	if err := tx.Where("IS_ACTIVE = 'Y' AND PROJECT_ID = ? AND FILE_NAME LIKE ?", projectId, "%"+keyword+"%").
		Find(&list).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, list)
}

// 移动文件
func (h *DiskHandler) MoveFile(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	var req struct {
		ID       uint `json:"id"`
		ParentID uint `json:"parentId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少参数")
		return
	}
	if err := tx.Model(&models.SysDiskFile{}).Where("ID = ?", req.ID).
		Update("PARENT_ID", req.ParentID).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, "移动成功")
}

// 复制文件
func (h *DiskHandler) CopyFile(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	var req struct {
		ID       uint `json:"id"`
		ParentID uint `json:"parentId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少参数")
		return
	}

	var origin models.SysDiskFile
	if err := tx.First(&origin, req.ID).Error; err != nil {
		ecode.Resp(c, ecode.ErrRequest, "源文件不存在")
		return
	}

	newFile := origin
	newFile.ID = 0
	newFile.ParentId = req.ParentID
	newFile.FileName += "_copy"
	utils.FillCreateMeta(c, &newFile)

	if err := tx.Create(&newFile).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}

	ecode.SuccessResp(c, newFile.ID)
}
