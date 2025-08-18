package v1

import (
	"time"
)

// CreateTodoRequest represents the HTTP request body for creating a todo
type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}

// CreateTodoResponse represents the HTTP response body after creating a todo
type CreateTodoResponse struct {
	ID uint `json:"id"`
}
