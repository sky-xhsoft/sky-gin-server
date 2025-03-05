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
	"github.com/sky-xhsoft/sky-gin-server/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
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
		Engine:      gin.Default(),
	}
}

func (s *Server) Run(c *config.Config) error {
	r := gin.New()
	r.Use(gin.Logger())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/query", func(c *gin.Context) {
		var result map[string]interface{}
		_ = s.Db.Raw("select * from crm_vip").Scan(&result).Error
		c.JSON(http.StatusOK, result)
		s.Log.Info(result)
	})

	return r.Run(c.System.Port)
}
