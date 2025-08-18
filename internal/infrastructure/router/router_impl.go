package router

import (
	"github.com/gin-gonic/gin"

	"itmrchow/go-todolist-service/internal/delivery/http/handler"
	v1 "itmrchow/go-todolist-service/internal/delivery/http/handler/v1"
	"itmrchow/go-todolist-service/internal/delivery/http/middleware"
)

var _ Router = &RouterImpl{}

// RouterImpl implements the Router interface.
type RouterImpl struct {
	healthHandler *handler.HealthHandler
	todoV1Handler v1.TodoHandler
}

// NewRouter creates a new router instance.
func NewRouter(
	healthHandler *handler.HealthHandler,
	todoV1Handler v1.TodoHandler,
) *RouterImpl {
	return &RouterImpl{
		healthHandler: healthHandler,
		todoV1Handler: todoV1Handler,
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

	routerGroup.POST("/create-todo", r.todoV1Handler.CreateTodo) // 新增todo
	routerGroup.POST("/find-todo", r.todoV1Handler.FindTodo)     // 查詢todo
	routerGroup.POST("/update-todo", r.todoV1Handler.UpdateTodo) // 更新todo

	// 目前 v1 路由群組為空，未來將在此新增業務邏輯路由
	// 例如：
	// routerGroup.GET("/todos", todoHandler.GetTodos)

}
