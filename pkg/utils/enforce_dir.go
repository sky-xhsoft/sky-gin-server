// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: enforce_dir.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/15
// Project Description:
// ----------------------------------------------------------------------------

package utils

import "os"

// EnsureDir 确保路径存在
func EnsureDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
