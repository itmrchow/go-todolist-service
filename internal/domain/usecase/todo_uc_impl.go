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
