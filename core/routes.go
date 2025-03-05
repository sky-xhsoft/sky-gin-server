package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"sync"
)

type RouteBinder struct {
	db       *gorm.DB
	router   *gin.Engine
	handlers sync.Map // 存储已注册的处理函数实例
	mu       sync.RWMutex
	log      *zap.SugaredLogger
}

func NewRouteBinder(db *gorm.DB, s *Server) *RouteBinder {
	return &RouteBinder{
		db:     db,
		router: s.Engine,
		log:    s.Log,
	}
}

// RegisterHandler 注册处理函数实例
func (rb *RouteBinder) RegisterHandler(name string, instance interface{}) {
	rb.handlers.Store(name, instance)
}

// 动态加载路由
func (rb *RouteBinder) LoadRoutes() error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	var routes = []models.SysRoutes{}

	if err := rb.db.Table("sys_routes").Where("IS_ACTIVE = 'Y'").Find(&routes).Error; err != nil {
		return err
	}

	// 绑定新路由
	for _, route := range routes {
		parts := strings.Split(route.Handle, ".")
		if len(parts) != 2 {
			continue
		}
		handlerName, methodName := parts[0], parts[1]

		instance, ok := rb.handlers.Load(handlerName)
		if !ok {
			continue
		}

		handlerValue := reflect.ValueOf(instance)
		method := handlerValue.MethodByName(methodName)
		if !method.IsValid() {
			continue
		}

		handler, ok := method.Interface().(func(c *gin.Context))
		if !ok {
			continue
		}

		rb.router.Handle(route.Method, route.Path, handler)
	}
	rb.log.Info(rb.router.Routes())
	return nil
}

var RoutesModule = fx.Module("routes",
	fx.Provide(NewRouteBinder),
	fx.Invoke(func(lc fx.Lifecycle, rb *RouteBinder) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return rb.LoadRoutes()
			},
		})
	}),
)
