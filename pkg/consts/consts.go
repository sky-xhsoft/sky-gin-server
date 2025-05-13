// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: consts.go
// Author: xhsoftware-skyzhou
// Created On: 2025/5/13
// Project Description:
// ----------------------------------------------------------------------------

package consts

import "time"

const (
	// Redis sms key 常量前缀
	RedisSmsCodePrefix = "sms:code:"

	//sms 验证码 有效是新建
	SmsCodeExpire = 5 * time.Minute

	DefaultPassword = "abc123"
)
