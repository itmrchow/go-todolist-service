package v1

import (
	"time"

	"itmrchow/go-todolist-service/internal/utils/dto"
)

// FindTodoRequest represents the HTTP request body for finding todos
type FindTodoRequest struct {
	Keyword     *string           `json:"keyword"`
	Status      *string           `json:"status"`
	CreatedFrom *time.Time        `json:"created_from"`
	CreatedTo   *time.Time        `json:"created_to"`
	DueFrom     *time.Time        `json:"due_from"`
	DueTo       *time.Time        `json:"due_to"`
	Pagination  dto.PaginationReq `json:"pagination"`
}

// FindTodoResponse represents the HTTP response body for finding todos
type FindTodoResponse struct {
	Todos      []TodoItem         `json:"todos"`
	Pagination dto.PaginationResp `json:"pagination"`
}

// TodoItem represents a single todo item in the response
type TodoItem struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
