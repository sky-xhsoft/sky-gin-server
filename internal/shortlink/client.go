// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: client.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/20
// Project Description:
// ----------------------------------------------------------------------------

package shortlink

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const shortLinkAPI = "https://api.xiaomark.com/v2/sl/link/create"

// CreateShortLinkRequest 请求结构体
type CreateShortLinkRequest struct {
	Apikey               string `json:"apikey"`
	GroupID              string `json:"group_id"`
	TargetURL            string `json:"target_url"`
	Name                 string `json:"name,omitempty"`
	Domain               string `json:"domain,omitempty"`
	Key                  string `json:"key,omitempty"`
	KeyLength            int    `json:"key_length,omitempty"`
	EscapeFromWechat     bool   `json:"escape_from_wechat,omitempty"`
	AdvancedBotDetection bool   `json:"advanced_bot_detection,omitempty"`
	Webhook              bool   `json:"webhook,omitempty"`
	WebhookScene         string `json:"webhook_scene,omitempty"`
}

// CreateShortLinkResponse 响应结构体
type CreateShortLinkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		LinkURL string `json:"link_url"`
	} `json:"data"`
}

// CallShortLinkAPI 发起请求
func CallShortLinkAPI(req *CreateShortLinkRequest) (string, error) {
	body, _ := json.Marshal(req)
	resp, err := http.Post(shortLinkAPI, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result CreateShortLinkResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if result.Code != 0 {
		return "", errors.New(result.Message)
	}

	return result.Data.LinkURL, nil
}
