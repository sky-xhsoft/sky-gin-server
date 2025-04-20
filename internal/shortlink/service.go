// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: service.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/20
// Project Description:
// ----------------------------------------------------------------------------

package shortlink

// CreateShortLink 封装短链业务逻辑
func CreateShortLink(targetURL, name string) (string, error) {
	req := &CreateShortLinkRequest{
		Apikey:    "0abd8f02531d8d04e78dd06a54f5ea11",
		GroupID:   "54kkbnvp",
		TargetURL: targetURL,
		Name:      name,
		KeyLength: 4,
	}
	return CallShortLinkAPI(req)
}
