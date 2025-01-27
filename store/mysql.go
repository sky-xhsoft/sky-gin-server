// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: mysql.go
// Author: xhsoftware-skyzhou
// Created On: 2025/1/24
// Project Description:
// ----------------------------------------------------------------------------

package store

import (
	"github.com/sky-xhsoft/sky-gin-server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func NewMysql(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Mysql.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		}})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(32)
	sqlDB.SetMaxOpenConns(512)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
