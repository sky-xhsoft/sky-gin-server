// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ecode.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package ecode

import (
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/pkg/response"
)

func Resp(c *gin.Context, code int, data interface{}) {
	response.WithCode(c, code, Msg(code), data)
}

func SuccessResp(c *gin.Context, data interface{}) {
	Resp(c, Success, data)
}

func ErrorResp(c *gin.Context, code int) {
	Resp(c, code, nil)
}
