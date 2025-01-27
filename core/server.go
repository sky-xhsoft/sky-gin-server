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
	"gorm.io/gorm"
	"net/http"
)

var logger = log.GetLogger()

type Server struct {
	Debug       bool
	Config      *config.Config
	Engine      *gin.Engine
	Db          *gorm.DB
	RedisClient *redis.Client
}

func NewServer(config *config.Config) *Server {
	return &Server{
		Debug:  false,
		Config: config,
		Engine: gin.Default(),
	}
}

func (s *Server) Run(c *config.Config) error {
	r := gin.New()
	r.Use(gin.Logger())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return r.Run(c.System.Port)
}
