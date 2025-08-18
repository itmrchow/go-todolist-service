package usecase

import (
	"context"
	"time"

	"itmrchow/go-todolist-service/internal/utils/dto"
)

//go:generate mockgen -source=todo_uc.go -destination=todo_uc_mock.go -package=usecase
type TodoUseCase interface {

	// CreateTodo creates a new todo and returns the created todo with assigned ID
	// Error:
	// - validation fail
	// - internal fail
	CreateTodo(ctx context.Context, req CreateTodoRequest) (*CreateTodoResponse, error)

	FindTodo(ctx context.Context, req FindTodoRequest) (*FindTodoResponse, error)

	// GetTodo(ctx context.Context, id uint) (*GetTodoResponse, error)
	UpdateTodo(ctx context.Context, req UpdateTodoRequest) error
	DeleteTodo(ctx context.Context, id uint) error
}

type CreateTodoRequest struct {
	Title       string
	Description *string
	Status      string // "pending", "doing", "done"
	DueDate     *time.Time
}

type CreateTodoResponse struct {
	ID uint
}

type FindTodoRequest struct {
	Keyword     *string           `json:"keyword"`
	Status      *string           `json:"status"`
	CreatedFrom *time.Time        `json:"created_from"`
	CreatedTo   *time.Time        `json:"created_to"`
	DueFrom     *time.Time        `json:"due_from"`
	DueTo       *time.Time        `json:"due_to"`
	Pagination  dto.PaginationReq `json:"pagination"`
}

type FindTodoResponse struct {
	Todos      []TodoResponse     `json:"todos"`
	Pagination dto.PaginationResp `json:"pagination"`
}

type TodoResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type UpdateTodoRequest struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`       // always required for validation
	Description *string    `json:"description"` // nil=keep current, ""=clear, "value"=update
	Status      *string    `json:"status"`      // nil=keep current, "value"=update
	DueDate     *time.Time `json:"due_date"`    // nil=keep current, time=update
}
