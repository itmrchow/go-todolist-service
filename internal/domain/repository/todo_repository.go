package repository

import (
	"context"
	"time"

	"itmrchow/go-todolist-service/internal/domain/entity"
)

// TodoRepository defines the interface for todo data persistence operations
//
//go:generate mockgen -source=todo_repository.go -destination=todo_repository_mock.go -package=repository
type TodoRepository interface {
	// Create creates a new todo and returns the created todo with assigned ID
	Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)

	// GetByID retrieves a todo by its ID
	// Returns nil if todo is not found or is soft deleted
	GetByID(ctx context.Context, id uint) (*entity.Todo, error)

	// Update updates an existing todo and returns the number of affected rows
	Update(ctx context.Context, todo *entity.Todo) (int64, error)

	// Delete soft deletes a todo (sets DeletedAt timestamp)
	Delete(ctx context.Context, id uint) error

	// List retrieves todos with pagination and filtering options
	List(ctx context.Context, queryParams TodoQueryParams, pagination *Pagination[entity.Todo]) error
}

// Pagination defines options for listing todos
type Pagination[T any] struct {
	Limit      int    `json:"limit,omitempty;query:limit"`
	Page       int    `json:"page,omitempty;query:page"`
	Sort       string `json:"sort,omitempty;query:sort"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Rows       []*T   `json:"rows"`
}

func (p *Pagination[T]) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}
func (p *Pagination[T]) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination[T]) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination[T]) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

// TodoQueryParams defines filters for listing todos
type TodoQueryParams struct {
	Status      *entity.TodoStatus // filter by status
	CreatedFrom *time.Time         `json:"created_from"`
	CreatedTo   *time.Time         `json:"created_to"`
	DueFrom     *time.Time         `json:"due_from"`
	DueTo       *time.Time         `json:"due_to"`
	Keyword     *string            // search in title and description
}
