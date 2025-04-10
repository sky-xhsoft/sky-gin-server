// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: registry.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package routers

import (
	"github.com/gin-gonic/gin"
	redis2 "github.com/go-redis/redis/v8"
)

var routeList []func(*gin.Engine)
var redis *redis2.Client

// 外部设置 redis 客户端（供中间件使用）
func SetRedis(r *redis2.Client) {
	redis = r
}

// 模块注册路由方法（模块内部调用）
func Register(f func(*gin.Engine)) {
	routeList = append(routeList, f)
}

// 系统启动统一加载所有路由
func Load(engine *gin.Engine) {
	for _, f := range routeList {
		f(engine)
	}
}
