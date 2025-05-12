// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: server.go
// Author: xhsoftware-skyzhou
// Created On: 2025/1/23
// Project Description:
// ----------------------------------------------------------------------------

package core

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/config"
	"github.com/sky-xhsoft/sky-gin-server/middleware"
	"github.com/sky-xhsoft/sky-gin-server/pkg/log"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ossUtil"
	"github.com/sky-xhsoft/sky-gin-server/routers"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var logger = log.GetLogger()

type Server struct {
	Debug       bool
	Config      *config.Config
	Engine      *gin.Engine
	Db          *gorm.DB
	Log         *zap.SugaredLogger
	RedisClient *redis.Client
}

func NewServer(config *config.Config, db *gorm.DB, redis *redis.Client, log *zap.SugaredLogger) *Server {
	return &Server{
		Debug:       false,
		Config:      config,
		Db:          db,
		RedisClient: redis,
		Log:         log,
		Engine:      gin.New(),
	}
}

var ServerModule = fx.Module("Server",
	fx.Provide(NewServer),

	fx.Invoke(func(s *Server) {
		s.Engine.Static("/static/", "./static/")
		s.Engine.Use(middleware.CORSMiddleware(), middleware.GinLogger(s.Log), gin.Recovery())
	}),

	fx.Invoke(func(cfg *config.Config) {
		if err := ossUtil.Init(cfg); err != nil {
			logger.Info("OSS 初始化失败: " + err.Error())
		}
	}),

	//支持用户手动routers注入
	fx.Invoke(func(s *Server) {
		routers.SetRedis(s.RedisClient) // 提供给权限中间件用
		routers.Load(s.Engine)          // 加载所有注册模块路由
	}),
)
