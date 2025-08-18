package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	v1 "itmrchow/go-todolist-service/internal/delivery/http/dto/v1"
	"itmrchow/go-todolist-service/internal/domain/usecase"
)

var _ TodoHandler = &TodoHandlerImpl{}

type TodoHandlerImpl struct {
	logger zerolog.Logger
	todoUc usecase.TodoUseCase // 假設有一個 TodoUserService
}

func NewTodoHandlerImpl(logger zerolog.Logger, todoUc usecase.TodoUseCase) *TodoHandlerImpl {
	return &TodoHandlerImpl{
		logger: logger,
		todoUc: todoUc,
	}
}

func (t *TodoHandlerImpl) CreateTodo(c *gin.Context) {
	// Parse HTTP request body into HTTP DTO
	var httpReq v1.CreateTodoRequest

	if err := c.ShouldBindJSON(&httpReq); err != nil {
		t.logger.Error().Err(err).Msg("invalid request format")
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

func (t *TodoHandlerImpl) FindTodo(c *gin.Context) {
	// Parse HTTP request body into HTTP DTO
	var httpReq v1.FindTodoRequest
	if err := c.ShouldBindJSON(&httpReq); err != nil {
		t.logger.Error().Err(err).Msg("invalid request format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request format",
		})
		return
	}

	// Convert HTTP DTO to UseCase DTO
	ucReq := usecase.FindTodoRequest{
		Keyword:     httpReq.Keyword,
		Status:      httpReq.Status,
		CreatedFrom: httpReq.CreatedFrom,
		CreatedTo:   httpReq.CreatedTo,
		DueFrom:     httpReq.DueFrom,
		DueTo:       httpReq.DueTo,
		Pagination:  httpReq.Pagination,
	}

	// Call Usecase
	ucResp, err := t.todoUc.FindTodo(c, ucReq)

	if err != nil {
		// Default error handling
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Convert UseCase response to HTTP DTO
	todos := make([]v1.TodoItem, len(ucResp.Todos))
	for i, todo := range ucResp.Todos {
		todos[i] = v1.TodoItem{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status,
			DueDate:     todo.DueDate,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}

	// Return success response
	httpResp := v1.FindTodoResponse{
		Todos:      todos,
		Pagination: ucResp.Pagination,
	}

	c.JSON(http.StatusOK, httpResp)
}

func (t *TodoHandlerImpl) UpdateTodo(c *gin.Context) {
	// Parse HTTP request body into HTTP DTO
	var httpReq v1.UpdateTodoRequest
	if err := c.ShouldBindJSON(&httpReq); err != nil {
		t.logger.Error().Err(err).Msg("invalid request format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request format",
		})
		return
	}

	// Convert HTTP DTO to UseCase DTO
	ucReq := usecase.UpdateTodoRequest{
		ID:          httpReq.ID,
		Title:       httpReq.Title,
		Description: httpReq.Description,
		Status:      httpReq.Status,
		DueDate:     httpReq.DueDate,
	}

	// Call usecase
	err := t.todoUc.UpdateTodo(c, ucReq)
	if err != nil {
		// Handle different error types
		if strings.Contains(err.Error(), "validation fail") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
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

	// Return HTTP 204 No Content for successful update
	c.AbortWithStatus(http.StatusNoContent)
}

func (t *TodoHandlerImpl) DeleteTodo(c *gin.Context) {
	// Parse HTTP request body into HTTP DTO
	var httpReq v1.DeleteTodoRequest
	if err := c.ShouldBindJSON(&httpReq); err != nil {
		t.logger.Error().Err(err).Msg("invalid request format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request format",
		})
		return
	}

	// Call usecase
	err := t.todoUc.DeleteTodo(c, httpReq.ID)
	if err != nil {
		// Handle different error types
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
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

	c.AbortWithStatus(http.StatusNoContent)
}
