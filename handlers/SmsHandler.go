// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: SmsHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/5/13
// Project Description:
// ----------------------------------------------------------------------------

package handlers

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/config"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/pkg/consts"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"log"
)

func init() {
	Register("SmsHandler", &SmsHandler{})
}

type SmsHandler struct {
	redis  *redis.Client
	config *config.Config
}

func (h *SmsHandler) HandlerName() string {
	return "SmsHandler"
}

func (h *SmsHandler) SetOption(ctx *core.AppContext) {
	h.redis = ctx.Redis
	h.config = ctx.Config
}

// SendSms 通过阿里云短信 SDK 发送短信并将验证码存入 Redis
func (s *SmsHandler) SendSms(c *gin.Context) {
	// 获取传入的 phone 参数
	var requestBody struct {
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println("Invalid request parameters:", err)
		// 使用统一的 ErrorResp 返回错误
		ecode.ErrorResp(c, 400)
		return
	}

	// 生成 4 位验证码
	code := utils.RandDigit(4)

	// 使用阿里云的 SDK 初始化短信客户端
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou",
		s.config.Sms.AccessKeyId,
		s.config.Sms.AccessKeySecret)

	if err != nil {
		log.Println("Failed to create client: %v", err)
		// 使用统一的 ErrorResp 返回错误
		ecode.ErrorResp(c, 400)
		return
	}

	// 构建发送短信的请求
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = s.config.Sms.SignName                       // 短信签名
	request.TemplateCode = s.config.Sms.TemplateCode               // 短信模板Code
	request.PhoneNumbers = requestBody.Phone                       // 发送的电话号码
	request.TemplateParam = fmt.Sprintf("{\"code\":\"%s\"}", code) // 短信内容，替换模板中的变量

	// 发送短信
	responseResult, err := client.SendSms(request)
	if err != nil {
		log.Println("Failed to send SMS: %v", err)
		// 使用统一的 ErrorResp 返回错误
		ecode.ErrorResp(c, 500)
		return
	}

	// 将验证码存储到 Redis，并设置过期时间（例如 5 分钟）
	ctx := context.Background()
	redisKey := consts.RedisSmsCodePrefix + requestBody.Phone // 键：sms:code:手机号
	err = s.redis.SetEX(ctx, redisKey, code, consts.SmsCodeExpire).Err()
	if err != nil {
		log.Printf("Failed to store verification code in Redis: %v", err)
		// 使用统一的 ErrorResp 返回错误
		ecode.ErrorResp(c, 500)
		return
	}
	log.Println(requestBody.Phone, responseResult)
	// 短信发送成功
	if responseResult.Code == "OK" {
		// 使用统一的 SuccessResp 返回成功结果
		ecode.SuccessResp(c, gin.H{"message": "SMS sent successfully", "verification_code": code})
	} else {
		// 使用统一的 ErrorResp 返回失败错误
		ecode.ErrorRespWithData(c, 500, responseResult.Message)
	}
}
