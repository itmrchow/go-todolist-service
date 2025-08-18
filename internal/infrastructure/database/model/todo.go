package model

import (
	"time"

	"gorm.io/gorm"
	"itmrchow/go-todolist-service/internal/domain/entity"
)

// Todo represents the GORM model for todo table
type Todo struct {
	gorm.Model
	Title       string     `gorm:"type:varchar(80);not null;comment:Todo標題，最多20個中文字符" json:"title"`
	Description *string    `gorm:"type:text;comment:Todo描述，最多100個中文字符" json:"description"`
	Status      string     `gorm:"type:varchar(20);not null;default:'pending';comment:Todo狀態;index" json:"status"`
	DueDate     *time.Time `gorm:"type:timestamp;null;comment:到期日期，UTC時間;index" json:"due_date"`
}

// TableName specifies the table name for GORM
func (Todo) TableName() string {
	return "todos"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (t *Todo) BeforeCreate(tx *gorm.DB) error {
	// Ensure status is valid, default to pending if invalid
	if t.Status != "pending" && t.Status != "doing" && t.Status != "done" {
		t.Status = "pending"
	}
	return nil
}

// EntityToModel converts domain entity to GORM model
func EntityToModel(entityTodo *entity.Todo) *Todo {
	if entityTodo == nil {
		return nil
	}

	model := &Todo{
		Model: gorm.Model{
			ID:        entityTodo.ID,
			CreatedAt: entityTodo.CreatedAt,
			UpdatedAt: entityTodo.UpdatedAt,
		},
		Title:       entityTodo.Title,
		Description: entityTodo.Description,
		Status:      string(entityTodo.Status),
		DueDate:     entityTodo.DueDate,
	}

	// Handle DeletedAt conversion
	if entityTodo.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{
			Time:  *entityTodo.DeletedAt,
			Valid: true,
		}
	}

	return model
}

// ModelToEntity converts GORM model to domain entity
func ModelToEntity(modelTodo *Todo) *entity.Todo {
	if modelTodo == nil {
		return nil
	}

	// Convert status string to entity.TodoStatus
	var status entity.TodoStatus
	switch modelTodo.Status {
	case "pending":
		status = entity.StatusPending
	case "doing":
		status = entity.StatusDoing
	case "done":
		status = entity.StatusDone
	default:
		// Default to pending for invalid status
		status = entity.StatusPending
	}

	entityTodo := &entity.Todo{
		ID:          modelTodo.ID,
		Title:       modelTodo.Title,
		Description: modelTodo.Description,
		Status:      status,
		DueDate:     modelTodo.DueDate,
		CreatedAt:   modelTodo.CreatedAt,
		UpdatedAt:   modelTodo.UpdatedAt,
	}

	// Handle DeletedAt conversion from gorm.DeletedAt to *time.Time
	if modelTodo.DeletedAt.Valid {
		entityTodo.DeletedAt = &modelTodo.DeletedAt.Time
	} else {
		entityTodo.DeletedAt = nil
	}

	return entityTodo
}

// ModelsToEntities converts slice of GORM models to slice of domain entities
func ModelsToEntities(modelTodos []*Todo) []*entity.Todo {
	if modelTodos == nil {
		return nil
	}

	entities := make([]*entity.Todo, len(modelTodos))
	for i, model := range modelTodos {
		entities[i] = ModelToEntity(model)
	}
	return entities
}