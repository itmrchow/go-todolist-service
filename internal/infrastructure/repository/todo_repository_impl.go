package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"itmrchow/go-todolist-service/internal/domain/entity"
	"itmrchow/go-todolist-service/internal/domain/repository"
	"itmrchow/go-todolist-service/internal/infrastructure/database/model"
)

var _ repository.TodoRepository = &TodoRepositoryImpl{}

// TodoRepositoryImpl implements the TodoRepository interface using GORM
type TodoRepositoryImpl struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewTodoRepository creates a new TodoRepository instance
func NewTodoRepository(logger zerolog.Logger, db *gorm.DB) repository.TodoRepository {
	return &TodoRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// Create creates a new todo and returns the created todo with assigned ID
func (r *TodoRepositoryImpl) Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	if todo == nil {
		return nil, errors.New("todo cannot be nil")
	}

	// Convert entity to model
	todoModel := model.EntityToModel(todo)
	if todoModel == nil {
		return nil, errors.New("failed to convert entity to model")
	}

	// Create in database
	if err := r.db.WithContext(ctx).Create(todoModel).Error; err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	// Convert back to entity and return
	createdEntity := model.ModelToEntity(todoModel)
	return createdEntity, nil
}

// GetByID retrieves a todo by its ID
// Returns nil if todo is not found or is soft deleted
func (r *TodoRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Todo, error) {
	var todoModel model.Todo

	// Query with soft delete scope (GORM automatically adds WHERE deleted_at IS NULL)
	err := r.db.WithContext(ctx).First(&todoModel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found, not an error
		}
		return nil, fmt.Errorf("failed to get todo by id %d: %w", id, err)
	}

	// Convert model to entity
	entity := model.ModelToEntity(&todoModel)
	return entity, nil
}

// Update updates an existing todo and returns the number of affected rows
func (r *TodoRepositoryImpl) Update(ctx context.Context, todo *entity.Todo) (int64, error) {
	if todo == nil {
		return 0, errors.New("todo cannot be nil")
	}

	if todo.ID == 0 {
		return 0, errors.New("todo ID cannot be 0")
	}

	// Convert entity to model
	todoModel := model.EntityToModel(todo)
	if todoModel == nil {
		return 0, errors.New("failed to convert entity to model")
	}

	// Use Updates to only update existing records (not insert new ones)
	result := r.db.WithContext(ctx).Model(&model.Todo{}).
		Where("id = ?", todo.ID).
		Updates(todoModel)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to update todo: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// Delete soft deletes a todo (sets DeletedAt timestamp)
func (r *TodoRepositoryImpl) Delete(ctx context.Context, id uint) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&model.Todo{}, id)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete todo: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// List retrieves todos with pagination and filtering options
func (r *TodoRepositoryImpl) List(
	ctx context.Context,
	queryParams repository.TodoQueryParams,
	pagination *repository.Pagination[entity.Todo],
) error {
	var todoModels []*model.Todo

	query := r.db.WithContext(ctx)
	// Apply filters
	query = r.applyFilters(query, queryParams)

	// Execute query
	if err := query.Scopes(Paginate(model.Todo{}, pagination, query)).Find(&todoModels).Error; err != nil {
		return fmt.Errorf("failed to list todos: %w", err)
	}

	// Convert models to entities
	entities := model.ModelsToEntities(todoModels)
	pagination.Rows = entities

	return nil
}

// Count returns the total count of todos (excluding soft deleted ones)
func (r *TodoRepositoryImpl) Count(ctx context.Context, filters repository.TodoQueryParams) (int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Todo{})

	// Apply filters
	query = r.applyFilters(query, filters)

	// Count
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count todos: %w", err)
	}

	return count, nil
}

// applyFilters applies filtering conditions to the query
func (r *TodoRepositoryImpl) applyFilters(query *gorm.DB, qP repository.TodoQueryParams) *gorm.DB {
	// Filter by status
	if qP.Status != nil {
		query = query.Where("status = ?", string(*qP.Status))
	}

	// Filter by due date range
	if qP.DueTo != nil {
		query = query.Where("due_date < ?", *qP.DueTo)
	}
	if qP.DueFrom != nil {
		query = query.Where("due_date > ?", *qP.DueFrom)
	}

	// Search in title and description
	if qP.Keyword != nil && *qP.Keyword != "" {
		search := "%" + *qP.Keyword + "%"
		query = query.Where(
			r.db.Where("title LIKE ?", search).
				Or("description LIKE ?", search),
		)
	}

	return query
}
