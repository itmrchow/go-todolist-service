package v1

import "github.com/gin-gonic/gin"

type TodoHandler interface {
	CreateTodo(c *gin.Context)
	FindTodo(c *gin.Context)
	UpdateTodo(c *gin.Context)
}
