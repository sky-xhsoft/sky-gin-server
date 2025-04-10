package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/pkg/response"
	"gorm.io/gorm"
)

type SysUserHandler struct {
	Db *gorm.DB
}

func (h *SysUserHandler) GetUserByID(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.Fail(c, "id 参数不能为空")
		return
	}

	var result map[string]interface{}
	if err := h.Db.Raw("SELECT * FROM sys_user WHERE id = ?", id).Scan(&result).Error; err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Ok(c, result)
}

func init() {
	Register("SysUserHandler", &SysUserHandler{}) // 注册空实例，由框架注入
}

func (h *SysUserHandler) HandlerName() string {
	return "SysUserHandler"
}

func (h *SysUserHandler) SetOption(ctx *core.AppContext) {
	h.Db = ctx.DB
}
