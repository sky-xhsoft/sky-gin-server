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
	"github.com/sky-xhsoft/sky-gin-server/config"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/pkg/log"
	"github.com/sky-xhsoft/sky-gin-server/store"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var logger = log.GetLogger()

// AppLifecycle 应用程序生命周期
type Lifecycle struct {
}

// OnStart 应用程序启动时执行
func (l *Lifecycle) OnStart(context.Context) error {
	logger.Info("AppLifecycle OnStart")
	return nil
}

// OnStop 应用程序停止时执行
func (l *Lifecycle) OnStop(context.Context) error {
	logger.Info("AppLifecycle OnStop")
	os.Exit(0)
	return nil
}

func NewAppLifeCycle() *Lifecycle {
	return &Lifecycle{}
}

func main() {
	// 使用 WaitGroup 等待应用所有 goroutine 完成
	var wg sync.WaitGroup

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

		// 创建应用服务
		fx.Provide(core.NewServer),

		fx.Invoke(func(s *core.Server, c *config.Config) {
			wg.Add(1)
			go func() {
				defer wg.Done() // 确保 goroutine 执行结束后完成通知

				logger.Info("Starting application server...")
				if err := s.Run(c); err != nil {
					logger.Error("Server failed:", err)
					os.Exit(0) // 服务失败则退出
				}
			}()
		}),

		fx.Provide(NewAppLifeCycle),

		// 注册生命周期回调函数
		fx.Invoke(func(lifecycle fx.Lifecycle, lc *Lifecycle) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return lc.OnStart(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return lc.OnStop(ctx)
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
