// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: message.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package ecode

var messages = map[int]string{
	Success:          "操作成功",
	ErrInvalidParam:  "参数错误",
	ErrUnauthorized:  "未授权访问",
	ErrUserNotFound:  "用户不存在",
	ErrPasswordWrong: "密码错误",

	ErrTokenExpired: "登录已过期",
	ErrTokenCreate:  "token 生成失败",
	ErrTokenEmpty:   "token 为空",
	ErrTokenError:   "token 注销失败",

	ErrRequest: "请求失败",

	ErrServer: "服务器内部错误",
}

func Msg(code int) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}
