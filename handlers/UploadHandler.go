// handlers/disk_handler.go
package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ossUtil"
	"net/http"
	"strconv"
)

// UploadHandler 兼容动态注册
type UploadHandler struct {
	ctx *core.AppContext
}

func (h *UploadHandler) HandlerName() string {
	return "UploadHandler"
}

func (h *UploadHandler) SetOption(ctx *core.AppContext) {
	h.ctx = ctx
}

// 注册 upload handler
func init() {
	Register("UploadHandler", &UploadHandler{})
}

// 上传处理逻辑
func (h *UploadHandler) UploadFile(c *gin.Context) {
	// 解析参数
	projectID, _ := strconv.ParseUint(c.PostForm("project_id"), 10, 64)
	//parentID, _ := strconv.ParseUint(c.PostForm("parent_id"), 10, 64)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

	// 保存本地临时文件
	tempPath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 生成 OSS Key 并上传
	ossKey := fmt.Sprintf("uploads/%s", file.Filename)
	ossUrl, err := ossUtil.GetClient().UploadFile(c, file, ossKey, strconv.FormatUint(projectID, 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OSS upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":     "Upload success",
		"oss_url": ossUrl,
	})
}
