// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: transaction.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WithTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()
		defer func() {
			if len(c.Errors) > 0 {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
		c.Set("TX", tx)
		c.Next()
	}
}
