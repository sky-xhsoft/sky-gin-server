package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"net/http"
)

type SysUserHandler struct {
	Db *gorm.DB
}

func NewSysUserHandler(db *gorm.DB) *SysUserHandler {
	return &SysUserHandler{
		Db: db,
	}
}

func (h *SysUserHandler) GetUserByID(c *gin.Context) {
	id := c.Query("id")
	var result map[string]interface{}
	h.Db.Raw("select * from sys_user t where t.id= ?", id).Scan(&result)
	c.JSON(http.StatusOK, result)
}

var SysUserHandlerModule = fx.Module("SysUserHandler",
	fx.Provide(NewSysUserHandler),
	fx.Invoke(func(rb *core.RouteBinder, h *SysUserHandler) {
		rb.RegisterHandler("SysUserHandler", h)
	}),
)
