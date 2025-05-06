// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: main.go
// Author: xhsoftware-skyzhou
// Created On: 2025/1/18
// Project Description:
// ----------------------------------------------------------------------------

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/config"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/handlers"
	_ "github.com/sky-xhsoft/sky-gin-server/handlers/videoHandlers"
	"github.com/sky-xhsoft/sky-gin-server/pkg/log"
	"github.com/sky-xhsoft/sky-gin-server/store"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger = log.GetLogger()

// AppLifecycle 应用程序生命周期
type Lifecycle struct {
}

// OnStart 应用程序启动时执行
func (l *Lifecycle) OnStart(ctx context.Context, srv *http.Server) error {
	logger.Info("AppLifecycle OnStart")
	srv.ListenAndServe()
	return nil
}

// OnStop 应用程序停止时执行
func (l *Lifecycle) OnStop(ctx context.Context, srv *http.Server) error {
	logger.Info("AppLifecycle OnStop")
	srv.Shutdown(ctx)
	os.Exit(0)
	return nil
}

func NewAppLifeCycle() *Lifecycle {
	return &Lifecycle{}
}

func main() {
	// 使用 WaitGroup 等待应用所有 goroutine 完成
	//var wg sync.WaitGroup

	app := fx.New(
		//初始化配置应用配置
		fx.Provide(func() *config.Config {
			config, err := config.LoadConfig(config.GetConfigFile())
			if err != nil {
				logger.Fatal(err)
			}
			return config
		}),
		//提供日志服务
		fx.Provide(log.GetLogger),

		//提供数据库服务
		fx.Provide(store.NewMysql),

		//提供缓存服务
		fx.Provide(store.NewRedis),

		//提供全局gin服务
		fx.Provide(gin.New),

		core.ServerModule,
		core.RoutesModule,

		//支持handlers注入
		fx.Invoke(func(rb *core.RouteBinder, db *gorm.DB, redis *redis.Client, logger *zap.SugaredLogger, cfg *config.Config) {
			appCtx := &core.AppContext{
				DB:     db,
				Redis:  redis,
				Logger: logger,
				Config: cfg,
			}
			handlers.LoadHandlers(rb, appCtx)
		}),

		fx.Provide(NewAppLifeCycle),

		// 注册生命周期回调函数
		fx.Invoke(func(lifecycle fx.Lifecycle, lc *Lifecycle, s *core.Server) {
			srv := &http.Server{
				Addr:    s.Config.System.Port,
				Handler: s.Engine,
			}
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return lc.OnStart(ctx, srv)
				},
				OnStop: func(ctx context.Context) error {
					return lc.OnStop(ctx, srv)
				},
			})
		}),
	)

	// 启动应用程序
	go func() {
		if err := app.Start(context.Background()); err != nil {
			logger.Error(err)
		}
	}()

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭应用程序
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		logger.Fatal(err)
	}

}
