// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: random.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/19
// Project Description:
// ----------------------------------------------------------------------------

package utils

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 生成 n 位随机字符串
func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 生成 n 位随机数字字符串
func RandDigit(n int) string {
	digits := []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}
