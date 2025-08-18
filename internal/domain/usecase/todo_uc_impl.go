package usecase

import (
	"context"
	"errors"

	"itmrchow/go-todolist-service/internal/domain/entity"
	"itmrchow/go-todolist-service/internal/domain/repository"
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
