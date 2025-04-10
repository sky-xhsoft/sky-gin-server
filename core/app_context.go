// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: app_context.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package core

import (
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.SugaredLogger
	Config *config.Config
}
