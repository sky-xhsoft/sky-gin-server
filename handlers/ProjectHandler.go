package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	db *gorm.DB
}

func (h *ProjectHandler) HandlerName() string {
	return "ProjectHandler"
}

func init() {
	Register("ProjectHandler", &ProjectHandler{})
}

func (h *ProjectHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
}

// 创建项目
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req models.ChrProject
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

	userVal, _ := c.Get("User")
	u := userVal.(*models.SysUser)

	projectUser := models.ChrProjectUser{
		ProjectId: req.ID,
		UserId:    u.ID,
		Prem:      "A",
		IsOwner:   "Y",
	}
	utils.FillCreateMeta(c, &projectUser)

	if err := tx.Create(&projectUser).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, "项目已建，但成员写入失败: "+err.Error())
		return
	}

	ecode.SuccessResp(c, req.ID)
}

// 更新项目
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req models.ChrProject
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少项目ID")
		return
	}

	utils.FillUpdateMeta(c, &req)

	if err := tx.Model(&models.ChrProject{}).Where("ID = ?", req.ID).Updates(&req).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, "更新成功")
}

// 添加项目成员
func (h *ProjectHandler) AddProjectUser(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req models.ChrProjectUser
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
	ecode.SuccessResp(c, "成员添加成功")
}

// 成员列表
func (h *ProjectHandler) ListProjectUsers(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	projectID := c.Query("projectId")
	if projectID == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少 projectId")
		return
	}

	var list []models.ChrProjectUser
	if err := tx.Where("CHR_PROJECT_ID = ? AND IS_ACTIVE = 'Y'", projectID).Find(&list).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, list)
}

// 移除项目成员（逻辑删除）
func (h *ProjectHandler) RemoveProjectUser(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	id := c.Query("id")
	if id == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少 id")
		return
	}

	if err := tx.Model(&models.ChrProjectUser{}).Where("ID = ?", id).Update("IS_ACTIVE", "N").Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}
	ecode.SuccessResp(c, "已移除")
}

// 获取当前用户可访问的项目列表
func (h *ProjectHandler) ListMyProjects(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	userVal, exists := c.Get("User")
	if !exists {
		ecode.Resp(c, ecode.ErrUnauthorized, "未登录")
		return
	}
	user := userVal.(*models.SysUser)

	var projectIDs []uint
	if err := tx.Model(&models.ChrProjectUser{}).
		Where("SYS_USER_ID = ? AND IS_ACTIVE = 'Y'", user.ID).
		Pluck("CHR_PROJECT_ID", &projectIDs).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}

	var projects []models.ChrProject
	if err := tx.Where("ID IN ? AND IS_ACTIVE = 'Y'", projectIDs).
		Find(&projects).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}

	ecode.SuccessResp(c, projects)
}
