// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: tx.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTx(c *gin.Context, fallback *gorm.DB) *gorm.DB {
	if txVal, exists := c.Get("TX"); exists {
		if tx, ok := txVal.(*gorm.DB); ok {
			return tx
		}
	}
	return fallback
}
