// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: permission.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/middleware"
	"strings"
)

func WithPermission(r *gin.Engine, method, path string, handler gin.HandlerFunc, perm string) {
	perm = strings.ToUpper(perm)
	var mws []gin.HandlerFunc

	switch perm {
	case "S", "D":
		if redis != nil {
			mws = append(mws, middleware.TokenAuth(redis))
		}
		// if perm == "D" { mws = append(mws, middleware.PermissionCheck()) } // TODO
	case "P":
		// no middleware
	default:
		// 默认使用 TokenAuth
		if redis != nil {
			mws = append(mws, middleware.TokenAuth(redis))
		}
	}

	// 按方法注册
	r.Handle(method, path, append(mws, handler)...)
}
