package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	v1 "itmrchow/go-todolist-service/internal/delivery/http/dto/v1"
	"itmrchow/go-todolist-service/internal/domain/usecase"
)

var _ TodoHandler = &TodoHandlerImpl{}

type TodoHandlerImpl struct {
	todoUc usecase.TodoUseCase // 假設有一個 TodoUserService
}

func NewTodoHandlerImpl(todoUc usecase.TodoUseCase) *TodoHandlerImpl {
	return &TodoHandlerImpl{
		todoUc: todoUc,
	}
}

func (t *TodoHandlerImpl) CreateTodo(c *gin.Context) {
	// Parse HTTP request body into HTTP DTO
	var httpReq v1.CreateTodoRequest

	if err := c.ShouldBindJSON(&httpReq); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request format",
		})
		return
	}

	if httpReq.Status == nil {
		httpReq.Status = new(string)
		*httpReq.Status = "pending"
	}

	// Convert HTTP DTO to UseCase DTO
	ucReq := usecase.CreateTodoRequest{
		Title:       httpReq.Title,
		Description: httpReq.Description,
		Status:      *httpReq.Status,
		DueDate:     httpReq.DueDate,
	}

	// Call usecase
	ucResp, err := t.todoUc.CreateTodo(c, ucReq)
	if err != nil {
		// Handle different error types
		if strings.Contains(err.Error(), "validation fail") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Default error handling
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Convert UseCase response to HTTP DTO
	httpResp := v1.CreateTodoResponse{
		ID: ucResp.ID,
	}

	// Return success response
	c.JSON(http.StatusOK, httpResp)
}
