package usecase

import (
	"context"
	"errors"

	"itmrchow/go-todolist-service/internal/domain/entity"
	"itmrchow/go-todolist-service/internal/domain/repository"
	"itmrchow/go-todolist-service/internal/utils/dto"
)

var _ TodoUseCase = &todoUseCaseImpl{}

type todoUseCaseImpl struct {
	todoRepo repository.TodoRepository
}

func NewTodoUseCaseImpl(todoRepo repository.TodoRepository) TodoUseCase {
	return &todoUseCaseImpl{
		todoRepo: todoRepo,
	}
}

// CreateTodo
func (t *todoUseCaseImpl) CreateTodo(ctx context.Context, req CreateTodoRequest) (*CreateTodoResponse, error) {
	// create todo
	status := entity.TodoStatus(req.Status)

	todoEntity, err := entity.NewTodo(req.Title, req.Description, &status, req.DueDate)
	if err != nil {
		return nil, errors.Join(errors.New("validation fail"), err)
	}

	// repository save model
	todoEntity, err = t.todoRepo.Create(ctx, todoEntity)
	if err != nil {
		return nil, errors.Join(errors.New("internal fail"), err)
	}

	// return response
	return &CreateTodoResponse{ID: todoEntity.ID}, nil
}

// FindTodo
func (t *todoUseCaseImpl) FindTodo(ctx context.Context, req FindTodoRequest) (*FindTodoResponse, error) {

	queryParams := repository.TodoQueryParams{
		Keyword:     req.Keyword,
		DueFrom:     req.DueFrom,
		DueTo:       req.DueTo,
		CreatedFrom: req.CreatedFrom,
		CreatedTo:   req.CreatedTo,
	}

	// status
	if req.Status != nil {
		status := entity.TodoStatus(*req.Status)
		if status.IsValid() {
			queryParams.Status = &status
		}
	}

	// pagination
	pagination := &repository.Pagination[entity.Todo]{
		Limit: req.Pagination.PageSize,
		Page:  req.Pagination.Page,
		Sort:  req.Pagination.SortBy + " " + req.Pagination.SortOrder,
	}

	err := t.todoRepo.List(ctx, queryParams, pagination)
	if err != nil {
		return nil, errors.Join(errors.New("internal fail"), err)
	}

	// response
	todos := make([]TodoResponse, len(pagination.Rows))
	for i, todo := range pagination.Rows {
		todos[i] = TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      string(todo.Status),
			DueDate:     todo.DueDate,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}

	resp := &FindTodoResponse{
		Todos: todos,
		Pagination: dto.PaginationResp{
			Page:       pagination.Page,
			PageSize:   pagination.Limit,
			TotalCount: int(pagination.TotalRows),
			TotalPages: pagination.TotalPages,
		},
	}

	return resp, nil
}

// UpdateTodo updates an existing todo with partial update support
func (t *todoUseCaseImpl) UpdateTodo(ctx context.Context, req UpdateTodoRequest) error {
	// Validate request
	if req.ID == 0 {
		return errors.New("validation fail: ID cannot be 0")
	}

	// First, get the existing todo to check if it exists
	existingTodo, err := t.todoRepo.GetByID(ctx, req.ID)
	if err != nil {
		return errors.Join(errors.New("internal fail"), err)
	}
	if existingTodo == nil {
		return errors.New("not found: todo not found")
	}

	// Create updated entity - start with existing values
	updatedTodo := &entity.Todo{
		ID:          req.ID,
		Title:       req.Title,                // Title is always required
		Description: existingTodo.Description, // Default to existing
		Status:      existingTodo.Status,      // Default to existing
		DueDate:     existingTodo.DueDate,     // Default to existing
		CreatedAt:   existingTodo.CreatedAt,
		UpdatedAt:   existingTodo.UpdatedAt,
	}

	// Partial update logic: only update fields that are provided (not nil)

	// Update Description if provided
	if req.Description != nil {
		if *req.Description == "" {
			// Empty string means clear the field
			updatedTodo.Description = nil
		} else {
			// Non-empty string means update to new value
			updatedTodo.Description = req.Description
		}
	}

	// Update Status if provided
	if req.Status != nil {
		status := entity.TodoStatus(*req.Status)
		if !status.IsValid() {
			return errors.New("validation fail: invalid status")
		}
		updatedTodo.Status = status
	}

	// Update DueDate if provided
	if req.DueDate != nil {
		updatedTodo.DueDate = req.DueDate
	}

	// Validate updated todo using entity rules
	if _, err := entity.NewTodo(updatedTodo.Title, updatedTodo.Description, &updatedTodo.Status, updatedTodo.DueDate); err != nil {
		return errors.Join(errors.New("validation fail"), err)
	}

	// Update in repository
	rowsAffected, err := t.todoRepo.Update(ctx, updatedTodo)
	if err != nil {
		return errors.Join(errors.New("internal fail"), err)
	}

	if rowsAffected == 0 {
		return errors.New("not found: todo not found")
	}

	return nil
}

// DeleteTodo deletes a todo by ID
func (t *todoUseCaseImpl) DeleteTodo(ctx context.Context, id uint) error {
	// Validate request
	if id == 0 {
		return errors.New("validation fail: ID cannot be 0")
	}

	// Delete in repository
	rowsAffected, err := t.todoRepo.Delete(ctx, id)
	if err != nil {
		return errors.Join(errors.New("internal fail"), err)
	}

	if rowsAffected == 0 {
		return errors.New("not found: todo not found")
	}

	return nil
}
