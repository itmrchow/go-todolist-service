package entity

import (
	"errors"
	"time"
)

// TodoStatus represents the status of a todo item
type TodoStatus string

const (
	StatusPending TodoStatus = "pending"
	StatusDoing   TodoStatus = "doing"
	StatusDone    TodoStatus = "done"
)

// IsValid checks if the TodoStatus is one of the valid values
func (s TodoStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusDoing, StatusDone:
		return true
	default:
		return false
	}
}

// Todo represents a todo item in the domain layer
type Todo struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      TodoStatus `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// NewTodo creates a new Todo with validation
func NewTodo(title string, description *string, status *TodoStatus, dueDate *time.Time) (*Todo, error) {
	// Validate title
	if len(title) == 0 {
		return nil, errors.New("title cannot be empty")
	}

	titleRunes := []rune(title)
	if len(titleRunes) > 20 {
		return nil, errors.New("title cannot exceed 20 characters")
	}

	// Validate description
	if description != nil {
		descriptionRunes := []rune(*description)
		if len(descriptionRunes) > 100 {
			return nil, errors.New("description cannot exceed 100 characters")
		}
	}

	// Set default status if not provided
	todoStatus := StatusPending
	if status != nil {
		if !status.IsValid() {
			return nil, errors.New("invalid status")
		}
		todoStatus = *status
	}

	now := time.Now().UTC()

	// Validate due date
	if dueDate != nil && !dueDate.After(now) {
		return nil, errors.New("due date must be in the future")
	}

	return &Todo{
		Title:       title,
		Description: description,
		Status:      todoStatus,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil, // New todos are not deleted
	}, nil
}

// IsDeleted checks if the todo is soft deleted
func (t *Todo) IsDeleted() bool {
	return t.DeletedAt != nil
}

// Delete soft deletes the todo by setting DeletedAt timestamp
func (t *Todo) Delete() {
	now := time.Now().UTC()
	t.DeletedAt = &now
	t.UpdatedAt = now
}

// Restore restores a soft deleted todo
func (t *Todo) Restore() {
	t.DeletedAt = nil
	t.UpdatedAt = time.Now().UTC()
}
