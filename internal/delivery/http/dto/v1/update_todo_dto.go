package v1

import (
	"time"
)

// UpdateTodoRequest represents the HTTP request body for updating a todo
type UpdateTodoRequest struct {
	ID          uint       `json:"id" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Description *string    `json:"description"`
	Status      *string    `json:"status" binding:"omitempty,oneof=pending doing done"`
	DueDate     *time.Time `json:"due_date"`
}

// No UpdateTodoResponse needed - using HTTP 204 No Content