package v1

import "github.com/gin-gonic/gin"

type TodoHandler interface {
	CreateTodo(c *gin.Context)
}
