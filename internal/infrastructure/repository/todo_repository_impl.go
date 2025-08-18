package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"itmrchow/go-todolist-service/internal/domain/entity"
	"itmrchow/go-todolist-service/internal/domain/repository"
	"itmrchow/go-todolist-service/internal/infrastructure/database/model"
	"itmrchow/go-todolist-service/internal/infrastructure/logger"
)

var _ repository.TodoRepository = &TodoRepositoryImpl{}

// TodoRepositoryImpl implements the TodoRepository interface using GORM
type TodoRepositoryImpl struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewTodoRepository creates a new TodoRepository instance
func NewTodoRepository(db *gorm.DB, logger logger.Logger) repository.TodoRepository {
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

// Update updates an existing todo
func (r *TodoRepositoryImpl) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	if todo == nil {
		return nil, errors.New("todo cannot be nil")
	}

	if todo.ID == 0 {
		return nil, errors.New("todo ID cannot be 0")
	}

	// Convert entity to model
	todoModel := model.EntityToModel(todo)
	if todoModel == nil {
		return nil, errors.New("failed to convert entity to model")
	}

	// Update in database
	result := r.db.WithContext(ctx).Save(todoModel)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update todo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("todo not found")
	}

	// Convert back to entity and return
	updatedEntity := model.ModelToEntity(todoModel)
	return updatedEntity, nil
}

// Delete soft deletes a todo (sets DeletedAt timestamp)
func (r *TodoRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Todo{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete todo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("todo not found")
	}

	return nil
}

// List retrieves todos with pagination and filtering options
func (r *TodoRepositoryImpl) List(ctx context.Context, options repository.ListOptions) ([]*entity.Todo, error) {
	query := r.db.WithContext(ctx).Model(&model.Todo{})

	// Apply filters
	query = r.applyFilters(query, options.Filters)

	// Apply sorting
	if options.SortBy != "" {
		order := options.SortBy
		if options.SortOrder == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		// Default sorting
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}
	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	// Execute query
	var todoModels []*model.Todo
	if err := query.Find(&todoModels).Error; err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	// Convert models to entities
	entities := model.ModelsToEntities(todoModels)
	return entities, nil
}

// Count returns the total count of todos (excluding soft deleted ones)
func (r *TodoRepositoryImpl) Count(ctx context.Context, filters repository.ListFilters) (int64, error) {
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
func (r *TodoRepositoryImpl) applyFilters(query *gorm.DB, filters repository.ListFilters) *gorm.DB {
	// Filter by status
	if filters.Status != nil {
		query = query.Where("status = ?", string(*filters.Status))
	}

	// Filter by due date range
	if filters.DueBefore != nil {
		query = query.Where("due_date < ?", *filters.DueBefore)
	}
	if filters.DueAfter != nil {
		query = query.Where("due_date > ?", *filters.DueAfter)
	}

	// Search in title and description
	if filters.Search != nil && *filters.Search != "" {
		search := "%" + *filters.Search + "%"
		query = query.Where(
			r.db.Where("title LIKE ?", search).
				Or("description LIKE ?", search),
		)
	}

	return query
}
