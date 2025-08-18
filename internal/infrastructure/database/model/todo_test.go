package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"itmrchow/go-todolist-service/internal/domain/entity"
)

func TestTodo_TableName(t *testing.T) {
	todo := Todo{}
	assert.Equal(t, "todos", todo.TableName())
}

func TestTodo_BeforeCreate(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  string
		expectedStatus string
	}{
		{
			name:           "Valid status pending",
			initialStatus:  "pending",
			expectedStatus: "pending",
		},
		{
			name:           "Valid status doing",
			initialStatus:  "doing",
			expectedStatus: "doing",
		},
		{
			name:           "Valid status done",
			initialStatus:  "done",
			expectedStatus: "done",
		},
		{
			name:           "Invalid status should default to pending",
			initialStatus:  "invalid_status",
			expectedStatus: "pending",
		},
		{
			name:           "Empty status should default to pending",
			initialStatus:  "",
			expectedStatus: "pending",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := &Todo{Status: tt.initialStatus}
			err := todo.BeforeCreate(nil) // tx is not used in our implementation
			
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, todo.Status)
		})
	}
}

func TestEntityToModel(t *testing.T) {
	dueDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	description := "測試描述"
	deletedAt := time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC)
	
	entityTodo := &entity.Todo{
		ID:          1,
		Title:       "測試標題",
		Description: &description,
		Status:      entity.StatusDoing,
		DueDate:     &dueDate,
		CreatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		DeletedAt:   &deletedAt,
	}

	modelTodo := EntityToModel(entityTodo)

	assert.NotNil(t, modelTodo)
	assert.Equal(t, uint(1), modelTodo.ID)
	assert.Equal(t, "測試標題", modelTodo.Title)
	assert.Equal(t, "測試描述", *modelTodo.Description)
	assert.Equal(t, "doing", modelTodo.Status)
	assert.Equal(t, dueDate.Unix(), modelTodo.DueDate.Unix())
	assert.Equal(t, entityTodo.CreatedAt.Unix(), modelTodo.CreatedAt.Unix())
	assert.Equal(t, entityTodo.UpdatedAt.Unix(), modelTodo.UpdatedAt.Unix())
	assert.True(t, modelTodo.DeletedAt.Valid)
	assert.Equal(t, deletedAt.Unix(), modelTodo.DeletedAt.Time.Unix())
}

func TestEntityToModel_WithNilFields(t *testing.T) {
	entityTodo := &entity.Todo{
		ID:          1,
		Title:       "測試標題",
		Description: nil,
		Status:      entity.StatusPending,
		DueDate:     nil,
		CreatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		DeletedAt:   nil,
	}

	modelTodo := EntityToModel(entityTodo)

	assert.NotNil(t, modelTodo)
	assert.Equal(t, uint(1), modelTodo.ID)
	assert.Equal(t, "測試標題", modelTodo.Title)
	assert.Nil(t, modelTodo.Description)
	assert.Equal(t, "pending", modelTodo.Status)
	assert.Nil(t, modelTodo.DueDate)
	assert.False(t, modelTodo.DeletedAt.Valid)
}

func TestEntityToModel_NilInput(t *testing.T) {
	modelTodo := EntityToModel(nil)
	assert.Nil(t, modelTodo)
}

func TestModelToEntity(t *testing.T) {
	dueDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	description := "測試描述"
	deletedAt := time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC)
	
	modelTodo := &Todo{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		Title:       "測試標題",
		Description: &description,
		Status:      "doing",
		DueDate:     &dueDate,
	}

	entityTodo := ModelToEntity(modelTodo)

	assert.NotNil(t, entityTodo)
	assert.Equal(t, uint(1), entityTodo.ID)
	assert.Equal(t, "測試標題", entityTodo.Title)
	assert.Equal(t, "測試描述", *entityTodo.Description)
	assert.Equal(t, entity.StatusDoing, entityTodo.Status)
	assert.Equal(t, dueDate.Unix(), entityTodo.DueDate.Unix())
	assert.Equal(t, modelTodo.CreatedAt.Unix(), entityTodo.CreatedAt.Unix())
	assert.Equal(t, modelTodo.UpdatedAt.Unix(), entityTodo.UpdatedAt.Unix())
	assert.NotNil(t, entityTodo.DeletedAt)
	assert.Equal(t, deletedAt.Unix(), entityTodo.DeletedAt.Unix())
}

func TestModelToEntity_WithNilFields(t *testing.T) {
	modelTodo := &Todo{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			DeletedAt: gorm.DeletedAt{Valid: false},
		},
		Title:       "測試標題",
		Description: nil,
		Status:      "pending",
		DueDate:     nil,
	}

	entityTodo := ModelToEntity(modelTodo)

	assert.NotNil(t, entityTodo)
	assert.Equal(t, uint(1), entityTodo.ID)
	assert.Equal(t, "測試標題", entityTodo.Title)
	assert.Nil(t, entityTodo.Description)
	assert.Equal(t, entity.StatusPending, entityTodo.Status)
	assert.Nil(t, entityTodo.DueDate)
	assert.Nil(t, entityTodo.DeletedAt)
}

func TestModelToEntity_InvalidStatus(t *testing.T) {
	modelTodo := &Todo{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		Title:  "測試標題",
		Status: "invalid_status",
	}

	entityTodo := ModelToEntity(modelTodo)

	assert.NotNil(t, entityTodo)
	// Should default to pending for invalid status
	assert.Equal(t, entity.StatusPending, entityTodo.Status)
}

func TestModelToEntity_NilInput(t *testing.T) {
	entityTodo := ModelToEntity(nil)
	assert.Nil(t, entityTodo)
}

func TestModelsToEntities(t *testing.T) {
	models := []*Todo{
		{
			Model: gorm.Model{ID: 1},
			Title: "第一個 Todo",
			Status: "pending",
		},
		{
			Model: gorm.Model{ID: 2},
			Title: "第二個 Todo",
			Status: "doing",
		},
	}

	entities := ModelsToEntities(models)

	assert.Len(t, entities, 2)
	assert.Equal(t, uint(1), entities[0].ID)
	assert.Equal(t, "第一個 Todo", entities[0].Title)
	assert.Equal(t, entity.StatusPending, entities[0].Status)
	assert.Equal(t, uint(2), entities[1].ID)
	assert.Equal(t, "第二個 Todo", entities[1].Title)
	assert.Equal(t, entity.StatusDoing, entities[1].Status)
}

func TestModelsToEntities_NilInput(t *testing.T) {
	entities := ModelsToEntities(nil)
	assert.Nil(t, entities)
}

func TestModelsToEntities_EmptySlice(t *testing.T) {
	models := []*Todo{}
	entities := ModelsToEntities(models)
	assert.NotNil(t, entities)
	assert.Len(t, entities, 0)
}