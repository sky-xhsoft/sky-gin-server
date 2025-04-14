// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ParseUint.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/15
// Project Description:
// ----------------------------------------------------------------------------

package utils

import "strconv"

// ParseUint 将字符串转为 uint 类型，出错返回 error
func ParseUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
