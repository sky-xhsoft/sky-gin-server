package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"gorm.io/gorm"
)

type SysUserHandler struct {
	Db *gorm.DB
}

func (h *SysUserHandler) GetUserByID(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		ecode.ErrorResp(c, ecode.ErrInvalidParam)
		return
	}

	var result map[string]interface{}
	if err := h.Db.Raw("SELECT * FROM sys_user WHERE id = ?", id).Scan(&result).Error; err != nil {
		ecode.Resp(c, ecode.ErrRequest, err.Error())
		return
	}

	ecode.SuccessResp(c, result)
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
