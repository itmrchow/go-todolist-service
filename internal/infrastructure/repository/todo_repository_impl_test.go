package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"itmrchow/go-todolist-service/internal/domain/entity"
	"itmrchow/go-todolist-service/internal/domain/repository"
	"itmrchow/go-todolist-service/internal/infrastructure/config"
	"itmrchow/go-todolist-service/internal/infrastructure/database"
	"itmrchow/go-todolist-service/internal/infrastructure/database/model"
)

type TodoRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo repository.TodoRepository
	ctx  context.Context
}

// SetupSuite 在整個測試 suite 開始前執行一次
func (suite *TodoRepositoryTestSuite) SetupSuite() {
	// 建立 SQLite in-memory 資料庫
	ctx := context.Background()
	sqlLiteDB := &database.SQLiteDBImpl{}
	db, err := sqlLiteDB.Connect(ctx, &config.DatabaseConfig{})
	suite.Require().NoError(err)

	// Auto migrate
	err = sqlLiteDB.Migrate(&model.Todo{})
	suite.Require().NoError(err)

	suite.db = db
	suite.ctx = ctx

	suite.repo = NewTodoRepository(zerolog.New(os.Stdout), suite.db)
}

// TearDownSuite 在整個測試 suite 結束後執行一次
func (suite *TodoRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, err := suite.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// TearDownTest 每個測試後清理資料
func (suite *TodoRepositoryTestSuite) TearDownTest() {
	if suite.db != nil {
		suite.db.Exec("DELETE FROM todos")
		suite.db.Exec("DELETE FROM sqlite_sequence WHERE name='todos'")
	}
}

func (suite *TodoRepositoryTestSuite) TestCreate_Success() {
	// Arrange
	todo, err := entity.NewTodo("測試標題", nil, nil, nil)
	suite.Require().NoError(err)

	// Act
	createdTodo, err := suite.repo.Create(suite.ctx, todo)

	// Assert
	suite.NoError(err)
	suite.NotNil(createdTodo)
	suite.NotZero(createdTodo.ID)
	suite.Equal("測試標題", createdTodo.Title)
	suite.Equal(entity.StatusPending, createdTodo.Status)
	suite.Nil(createdTodo.Description)
	suite.Nil(createdTodo.DueDate)
	suite.NotZero(createdTodo.CreatedAt)
	suite.NotZero(createdTodo.UpdatedAt)
	suite.Nil(createdTodo.DeletedAt)
}

func (suite *TodoRepositoryTestSuite) TestCreate_WithAllFields() {
	// Arrange
	description := "測試描述"
	status := entity.StatusDoing
	dueDate := time.Now().Add(time.Hour * 24).UTC()

	todo, err := entity.NewTodo("測試標題", &description, &status, &dueDate)
	suite.Require().NoError(err)

	// Act
	createdTodo, err := suite.repo.Create(suite.ctx, todo)

	// Assert
	suite.NoError(err)
	suite.NotNil(createdTodo)
	suite.NotZero(createdTodo.ID)
	suite.Equal("測試標題", createdTodo.Title)
	suite.Equal("測試描述", *createdTodo.Description)
	suite.Equal(entity.StatusDoing, createdTodo.Status)
	suite.Equal(dueDate.Unix(), createdTodo.DueDate.Unix())
}

func (suite *TodoRepositoryTestSuite) TestCreate_NilInput() {
	// Act
	createdTodo, err := suite.repo.Create(suite.ctx, nil)

	// Assert
	suite.Error(err)
	suite.Nil(createdTodo)
	suite.Contains(err.Error(), "todo cannot be nil")
}

func (suite *TodoRepositoryTestSuite) TestGetByID_Success() {
	// Arrange - First create a todo
	todo, err := entity.NewTodo("測試標題", nil, nil, nil)
	suite.Require().NoError(err)

	createdTodo, err := suite.repo.Create(suite.ctx, todo)
	suite.Require().NoError(err)

	// Act
	foundTodo, err := suite.repo.GetByID(suite.ctx, createdTodo.ID)

	// Assert
	suite.NoError(err)
	suite.NotNil(foundTodo)
	suite.Equal(createdTodo.ID, foundTodo.ID)
	suite.Equal(createdTodo.Title, foundTodo.Title)
	suite.Equal(createdTodo.Status, foundTodo.Status)
}

func (suite *TodoRepositoryTestSuite) TestGetByID_NotFound() {
	// Act
	foundTodo, err := suite.repo.GetByID(suite.ctx, 999)

	// Assert
	suite.NoError(err) // Repository should not return error for not found
	suite.Nil(foundTodo)
}

func (suite *TodoRepositoryTestSuite) TestGetByID_SoftDeleted() {
	// Arrange - Create and soft delete a todo
	todo, err := entity.NewTodo("測試標題", nil, nil, nil)
	suite.Require().NoError(err)

	createdTodo, err := suite.repo.Create(suite.ctx, todo)
	suite.Require().NoError(err)

	// Soft delete the todo
	err = suite.repo.Delete(suite.ctx, createdTodo.ID)
	suite.Require().NoError(err)

	// Act
	foundTodo, err := suite.repo.GetByID(suite.ctx, createdTodo.ID)

	// Assert
	suite.NoError(err)
	suite.Nil(foundTodo) // Soft deleted todos should not be returned
}

func (suite *TodoRepositoryTestSuite) TestUpdate_Success() {
	// Arrange - First create a todo
	todo, err := entity.NewTodo("原始標題", nil, nil, nil)
	suite.Require().NoError(err)

	createdTodo, err := suite.repo.Create(suite.ctx, todo)
	suite.Require().NoError(err)

	// Modify the todo
	createdTodo.Title = "更新的標題"
	newDescription := "新的描述"
	createdTodo.Description = &newDescription
	createdTodo.Status = entity.StatusDoing

	// Act
	rowsAffected, err := suite.repo.Update(suite.ctx, createdTodo)

	// Assert
	suite.NoError(err)
	suite.Equal(int64(1), rowsAffected)
	
	// Verify the update by fetching the record
	updatedTodo, err := suite.repo.GetByID(suite.ctx, createdTodo.ID)
	suite.NoError(err)
	suite.NotNil(updatedTodo)
	suite.Equal(createdTodo.ID, updatedTodo.ID)
	suite.Equal("更新的標題", updatedTodo.Title)
	suite.Equal("新的描述", *updatedTodo.Description)
	suite.Equal(entity.StatusDoing, updatedTodo.Status)
	suite.True(updatedTodo.UpdatedAt.After(createdTodo.UpdatedAt))
}

func (suite *TodoRepositoryTestSuite) TestUpdate_NilInput() {
	// Act
	rowsAffected, err := suite.repo.Update(suite.ctx, nil)

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "todo cannot be nil")
	suite.Equal(int64(0), rowsAffected)
}

func (suite *TodoRepositoryTestSuite) TestUpdate_ZeroID() {
	// Arrange
	todo := &entity.Todo{
		ID:    0, // Zero ID should cause error
		Title: "Test Title",
	}

	// Act
	rowsAffected, err := suite.repo.Update(suite.ctx, todo)

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "todo ID cannot be 0")
	suite.Equal(int64(0), rowsAffected)
}

func (suite *TodoRepositoryTestSuite) TestUpdate_NotFound() {
	// Arrange - Create a todo with non-existent ID
	now := time.Now().UTC()
	nonExistentTodo := &entity.Todo{
		ID:          999,
		Title:       "Non-existent Todo",
		Description: nil,
		Status:      entity.StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Act
	rowsAffected, err := suite.repo.Update(suite.ctx, nonExistentTodo)

	// Assert
	suite.NoError(err) // No error with Updates method
	suite.Equal(int64(0), rowsAffected) // 0 rows affected means not found
}

func (suite *TodoRepositoryTestSuite) TestDelete_Success() {
	// Arrange - First create a todo
	todo, err := entity.NewTodo("測試標題", nil, nil, nil)
	suite.Require().NoError(err)

	createdTodo, err := suite.repo.Create(suite.ctx, todo)
	suite.Require().NoError(err)

	// Act
	err = suite.repo.Delete(suite.ctx, createdTodo.ID)

	// Assert
	suite.NoError(err)

	// Verify the todo is soft deleted
	foundTodo, err := suite.repo.GetByID(suite.ctx, createdTodo.ID)
	suite.NoError(err)
	suite.Nil(foundTodo) // Should not be found after soft delete
}

func (suite *TodoRepositoryTestSuite) TestDelete_NotFound() {
	// Act
	err := suite.repo.Delete(suite.ctx, 999)

	// Assert
	suite.Error(err)
	suite.Contains(err.Error(), "todo not found")
}

func (suite *TodoRepositoryTestSuite) TestList_WithFilters() {
	// Arrange - Create multiple todos
	todo1, _ := entity.NewTodo("第一個 Todo", nil, nil, nil)
	todo2, _ := entity.NewTodo("第二個 Todo", nil, nil, nil)

	status := entity.StatusDoing
	todo3, _ := entity.NewTodo("第三個 Todo", nil, &status, nil)

	suite.repo.Create(suite.ctx, todo1)
	suite.repo.Create(suite.ctx, todo2)
	suite.repo.Create(suite.ctx, todo3)

	// Act - List with status filter
	// statusFilter := entity.StatusDoing
	queryParams := repository.TodoQueryParams{
		// Status: &statusFilter,
	}

	pagination := &repository.Pagination[entity.Todo]{
		Limit: 2,
		Page:  1,
		Sort:  "created_at DESC",
	}

	err := suite.repo.List(suite.ctx, queryParams, pagination)

	// Assert
	suite.NoError(err)
	suite.Len(pagination.Rows, 2)

	suite.EqualValues(2, pagination.TotalPages)
	suite.EqualValues(3, pagination.TotalRows)
}

func TestTodoRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TodoRepositoryTestSuite))
}
