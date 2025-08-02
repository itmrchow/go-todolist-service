package router

import (
	"github.com/gin-gonic/gin"

	"itmrchow/go-todolist-service/internal/delivery/http/handler"
	"itmrchow/go-todolist-service/internal/delivery/http/middleware"
)

var _ Router = &RouterImpl{}

// RouterImpl implements the Router interface.
type RouterImpl struct {
	healthHandler *handler.HealthHandler
}

// NewRouter creates a new router instance.
func NewRouter() *RouterImpl {
	return &RouterImpl{
		healthHandler: handler.NewHealthHandler(),
	}
}

// SetupRoutes configures and returns the Gin engine with all routes.
func (r *RouterImpl) SetupRoutes() *gin.Engine {
	engine := gin.Default()

	// 註冊全域中間件
	engine.Use(middleware.CORS())
	engine.Use(middleware.ErrorHandler())

	// 註冊基礎路由
	engine.GET("/health", r.healthHandler.Health)
	engine.GET("/version", r.healthHandler.Version)

	// 設定 v1 API 路由群組
	v1Group := engine.Group("/api/v1")
	r.RegisterV1Routes(v1Group)

	return engine
}

// RegisterV1Routes registers all v1 API routes.
func (r *RouterImpl) RegisterV1Routes(routerGroup *gin.RouterGroup) {
	// 目前 v1 路由群組為空，未來將在此新增業務邏輯路由
	// 例如：
	// routerGroup.GET("/todos", todoHandler.GetTodos)
	// routerGroup.POST("/todos", todoHandler.CreateTodo)
}
