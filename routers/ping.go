// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ping.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
)

func init() {
	Register(func(r *gin.Engine) {
		WithPermission(r, "GET", "/ping", demo, "P")

	})
}

func demo(c *gin.Context) {
	ecode.SuccessResp(c, nil)
}
