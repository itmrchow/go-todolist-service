package router

import "github.com/gin-gonic/gin"

// Router defines the interface for HTTP routing management.
type Router interface {
	SetupRoutes() *gin.Engine
	RegisterV1Routes(routerGroup *gin.RouterGroup)
}
