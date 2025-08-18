package usecase

import (
	"context"
	"time"
)

//go:generate mockgen -source=todo_uc.go -destination=todo_uc_mock.go -package=usecase
type TodoUseCase interface {

	// CreateTodo creates a new todo and returns the created todo with assigned ID
	// Error:
	// - validation fail
	// - internal fail
	CreateTodo(ctx context.Context, req CreateTodoRequest) (*CreateTodoResponse, error)

	// GetTodo(ctx context.Context, id uint) (*GetTodoResponse, error)
	// UpdateTodo(ctx context.Context, id uint, req UpdateTodoRequest) (*UpdateTodoResponse, error)
	// DeleteTodo(ctx context.Context, id uint) error
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
