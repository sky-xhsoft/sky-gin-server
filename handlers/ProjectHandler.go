package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"gorm.io/gorm"
	"strings"
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
	models.FillCreateMeta(c, &req)

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
	models.FillCreateMeta(c, &projectUser)

	if err := tx.Create(&projectUser).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, "项目已建，但成员写入失败: "+err.Error())
		return
	}

	ecode.SuccessResp(c, req)
}

// 更新项目
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "参数解析失败")
		return
	}

	idVal, exists := req["ID"]
	if !exists {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少项目ID")
		return
	}

	id, ok := idVal.(float64)
	if !ok || uint(id) == 0 {
		ecode.Resp(c, ecode.ErrInvalidParam, "项目ID非法")
		return
	}
	req["ID"] = uint(id) // 显式转为 uint

	// 填充更新元信息
	models.FillUpdateMetaMap(c, req)

	if err := tx.Model(&models.ChrProject{}).Where("ID = ?", uint(id)).Updates(req).Error; err != nil {
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
	models.FillCreateMeta(c, &req)

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

	projectID := c.Query("ID")
	if projectID == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少 ID")
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

	id := c.Query("ID")
	if id == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少 ID")
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

	// 解析前端传入的排序参数，默认降序
	order := strings.ToLower(c.DefaultQuery("order", "desc"))
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	var projects []models.ChrProject
	if err := tx.Where("ID IN ? AND IS_ACTIVE = 'Y'", projectIDs).
		Order("CREATE_TIME " + order).
		Find(&projects).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err.Error())
		return
	}

	ecode.SuccessResp(c, projects)
}

func (h *ProjectHandler) ListResourceItemByProject(c *gin.Context) {
	tx := utils.GetTx(c, h.db)

	projectId := c.Query("projectId")
	var project models.ChrProject
	//获取项目详情
	if err := tx.Model(models.ChrProject{}).Where("ID = ?", projectId).First(&project).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err)
		return
	}
	//根据项目获取所有资源明细
	var resources []models.ChrResource
	if err := tx.Where("CHR_PROJECT_ID =? and IS_ACTIVE='Y'", projectId).
		Order("CREATE_TIME desc").Find(&resources).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrServer, err)
		return
	}

	for k, v := range resources {
		var items []models.ChrResourceItem
		if err := tx.Where("CHR_RESOURCE_ID =? and IS_ACTIVE='Y' ", v.ID).
			Order("CREATE_TIME desc").Find(&items).Error; err != nil {
			continue
		}
		resources[k].Items = items
	}

	project.Resouse = resources

	ecode.SuccessResp(c, project)

}
