// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: code.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package ecode

const (
	Success          = 0
	ErrInvalidParam  = 10001
	ErrUnauthorized  = 10002
	ErrUserNotFound  = 10003
	ErrPasswordWrong = 10004

	ErrTokenExpired = 20005
	ErrTokenCreate  = 20006
	ErrTokenEmpty   = 20007
	ErrTokenError   = 20008

	ErrRequest = 30001

	ErrServer = 50000
)
