package repository

import (
	"context"

	"itmrchow/go-todolist-service/internal/domain/entity"
)

// TodoRepository defines the interface for todo data persistence operations
type TodoRepository interface {
	// Create creates a new todo and returns the created todo with assigned ID
	Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)

	// GetByID retrieves a todo by its ID
	// Returns nil if todo is not found or is soft deleted
	GetByID(ctx context.Context, id uint) (*entity.Todo, error)

	// Update updates an existing todo
	Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)

	// Delete soft deletes a todo (sets DeletedAt timestamp)
	Delete(ctx context.Context, id uint) error

	// List retrieves todos with pagination and filtering options
	List(ctx context.Context, options ListOptions) ([]*entity.Todo, error)

	// Count returns the total count of todos (excluding soft deleted ones)
	Count(ctx context.Context, filters ListFilters) (int64, error)
}

// ListOptions defines options for listing todos
type ListOptions struct {
	Filters ListFilters

	// Pagination
	Limit  int
	Offset int

	// Sorting
	SortBy    string // field to sort by (e.g., "created_at", "title", "due_date")
	SortOrder string // "asc" or "desc"
}

// ListFilters defines filters for listing todos
type ListFilters struct {
	Status    *entity.TodoStatus // filter by status
	DueBefore *string            // filter todos due before this date (RFC3339 format)
	DueAfter  *string            // filter todos due after this date (RFC3339 format)
	Search    *string            // search in title and description
}
