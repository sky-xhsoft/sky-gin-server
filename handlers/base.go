// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: base.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package handlers

import "github.com/sky-xhsoft/sky-gin-server/core"

// 注册handlers
var registry = make(map[string]Handler)

func Register(name string, h Handler) {
	registry[name] = h
}

func GetRegistry() map[string]Handler {
	return registry
}

// handler接口
type Handler interface {
	HandlerName() string
	SetOption(ctx *core.AppContext)
}

// 动态加载handlers
func LoadHandlers(rb *core.RouteBinder, ctx *core.AppContext) {
	for name, h := range GetRegistry() {
		h.SetOption(ctx)
		rb.RegisterHandler(name, h)
	}
}
