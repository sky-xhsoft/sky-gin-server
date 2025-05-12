package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/middleware"
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
	redis    *redis.Client
}

func NewRouteBinder(db *gorm.DB, s *Server) *RouteBinder {
	return &RouteBinder{
		db:     db,
		router: s.Engine,
		log:    s.Log,
		redis:  s.RedisClient,
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

	if err := rb.db.Table("sys_api").Where("IS_ACTIVE = 'Y'").Find(&routes).Error; err != nil {
		rb.log.Error("加载路由失败:", err)
		return err
	}

	// 路由冲突检测表
	existingRoutes := map[string]map[string]bool{}
	for _, route := range rb.router.Routes() {
		if _, ok := existingRoutes[route.Method]; !ok {
			existingRoutes[route.Method] = map[string]bool{}
		}
		existingRoutes[route.Method][route.Path] = true
	}

	// 权限类型对应的中间件
	middlewareMap := map[string][]gin.HandlerFunc{
		"P": {middleware.WithTransaction(rb.db)}, // Public
		"S": {middleware.TokenAuth(rb.redis),
			middleware.WithTransaction(rb.db)},
		"D": {middleware.TokenAuth(rb.redis),
			middleware.WithTransaction(rb.db)}, // TODO: future: 加数据权限验证中间件
	}

	// 绑定新路由
	for _, r := range routes {
		// 检查静态路由冲突
		if existingRoutes[strings.ToUpper(r.Method)][r.Path] {
			rb.log.Warnf("路由冲突，静态路由已注册: [%s] %s，跳过数据库加载", r.Method, r.Path)
			continue
		}

		parts := strings.Split(r.Handle, ".")
		if len(parts) != 2 {
			rb.log.Warnf("Handle 格式不合法: %s", r.Handle)
			continue
		}
		handlerName, methodName := parts[0], parts[1]

		instance, ok := rb.handlers.Load(handlerName)
		if !ok {
			rb.log.Warnf("未注册 handler: %s", handlerName)
			continue
		}

		method := reflect.ValueOf(instance).MethodByName(methodName)
		if !method.IsValid() {
			rb.log.Warnf("Handler %s 不包含方法: %s", handlerName, methodName)
			continue
		}

		fn, ok := method.Interface().(func(*gin.Context))
		if !ok {
			rb.log.Warnf("方法签名不匹配: %s.%s", handlerName, methodName)
			continue
		}

		perm := strings.ToUpper(r.Permission)
		mws := middlewareMap[perm]
		rb.router.Handle(r.Method, r.Path, append(mws, fn)...)

		rb.log.Infof("[路由加载] %s %s → %s.%s", r.Method, r.Path, handlerName, methodName)
	}

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
