package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"itmrchow/go-todolist-service/internal/domain/entity"
	"itmrchow/go-todolist-service/internal/domain/repository"
	"itmrchow/go-todolist-service/internal/utils/dto"
)

type TodoUseCaseTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	mockRepo *repository.MockTodoRepository
	uc       TodoUseCase
}

// 執行測試套件
func TestTodoUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TodoUseCaseTestSuite))
}

// SetupTest 在每個測試前執行
func (suite *TodoUseCaseTestSuite) SetupTest() {
	// 每個測試前創建新的 mock controller
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockRepo = repository.NewMockTodoRepository(suite.ctrl)
	suite.uc = NewTodoUseCaseImpl(suite.mockRepo)
}

// TearDownTest 在每個測試後執行
func (suite *TodoUseCaseTestSuite) TearDownTest() {
	if suite.ctrl != nil {
		suite.ctrl.Finish()
	}
}

func (suite *TodoUseCaseTestSuite) TestUpdateTodo() {
	ctx := context.Background()

	tests := []struct {
		name         string
		req          UpdateTodoRequest
		setupMock    func()
		expectErrMsg string
	}{
		{
			name: "validation_fail_zero_id",
			req: UpdateTodoRequest{
				ID:    0,
				Title: "Updated Title",
			},
			setupMock: func() {
				// No mock setup needed for validation error
			},
			expectErrMsg: "validation fail",
		},
		{
			name: "get_existing_todo_fail",
			req: UpdateTodoRequest{
				ID:    1,
				Title: "Updated Title",
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					GetByID(ctx, uint(1)).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectErrMsg: "internal fail",
		},
		{
			name: "todo_not_found",
			req: UpdateTodoRequest{
				ID:    999,
				Title: "Updated Title",
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					GetByID(ctx, uint(999)).
					Return(nil, nil). // not found
					Times(1)
			},
			expectErrMsg: "not found",
		},
		{
			name: "invalid_status",
			req: UpdateTodoRequest{
				ID:     1,
				Title:  "Updated Title",
				Status: stringPtr("invalid_status"),
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					GetByID(ctx, uint(1)).
					Return(&entity.Todo{
						ID:     1,
						Title:  "Original Title",
						Status: entity.StatusPending,
					}, nil).
					Times(1)
			},
			expectErrMsg: "validation fail",
		},
		{
			name: "entity_validation_fail",
			req: UpdateTodoRequest{
				ID:    1,
				Title: "", // Empty title should fail validation
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					GetByID(ctx, uint(1)).
					Return(&entity.Todo{
						ID:     1,
						Title:  "Original Title",
						Status: entity.StatusPending,
					}, nil).
					Times(1)
			},
			expectErrMsg: "validation fail",
		},
		{
			name: "repository_update_fail",
			req: UpdateTodoRequest{
				ID:    1,
				Title: "Updated Title",
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					GetByID(ctx, uint(1)).
					Return(&entity.Todo{
						ID:     1,
						Title:  "Original Title",
						Status: entity.StatusPending,
					}, nil).
					Times(1)

				suite.mockRepo.EXPECT().
					Update(ctx, gomock.Any()).
					Return(int64(0), errors.New("database error")).
					Times(1)
			},
			expectErrMsg: "internal fail",
		},
		{
			name: "no_rows_affected",
			req: UpdateTodoRequest{
				ID:    1,
				Title: "Updated Title",
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					GetByID(ctx, uint(1)).
					Return(&entity.Todo{
						ID:     1,
						Title:  "Original Title",
						Status: entity.StatusPending,
					}, nil).
					Times(1)

				suite.mockRepo.EXPECT().
					Update(ctx, gomock.Any()).
					Return(int64(0), nil). // 0 rows affected
					Times(1)
			},
			expectErrMsg: "not found",
		},
		{
			name: "success_partial_update",
			req: UpdateTodoRequest{
				ID:          1,
				Title:       "Updated Title",
				Description: stringPtr("Updated Description"),
				Status:      stringPtr("doing"),
				// DueDate not provided - should keep existing
			},
			setupMock: func() {
				existingDueDate := time.Now().Add(time.Second)
				suite.mockRepo.EXPECT().
					GetByID(ctx, uint(1)).
					Return(&entity.Todo{
						ID:          1,
						Title:       "Original Title",
						Description: stringPtr("Original Description"),
						Status:      entity.StatusPending,
						DueDate:     &existingDueDate,
						CreatedAt:   timeNow(),
						UpdatedAt:   timeNow(),
					}, nil).
					Times(1)

				suite.mockRepo.EXPECT().
					Update(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, todo *entity.Todo) (int64, error) {
						// Verify partial update logic
						assert.Equal(suite.T(), "Updated Title", todo.Title)
						assert.NotNil(suite.T(), todo.Description)
						assert.Equal(suite.T(), "Updated Description", *todo.Description)
						assert.Equal(suite.T(), entity.StatusDoing, todo.Status)
						assert.NotNil(suite.T(), todo.DueDate) // Should keep existing
						return int64(1), nil
					}).
					Times(1)
			},
			expectErrMsg: "",
		},
		{
			name: "success_clear_description",
			req: UpdateTodoRequest{
				ID:          1,
				Title:       "Updated Title",
				Description: stringPtr(""), // Clear description
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					GetByID(ctx, uint(1)).
					Return(&entity.Todo{
						ID:          1,
						Title:       "Original Title",
						Description: stringPtr("Original Description"),
						Status:      entity.StatusPending,
						CreatedAt:   timeNow(),
						UpdatedAt:   timeNow(),
					}, nil).
					Times(1)

				suite.mockRepo.EXPECT().
					Update(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, todo *entity.Todo) (int64, error) {
						// Verify description was cleared
						assert.Nil(suite.T(), todo.Description)
						return int64(1), nil
					}).
					Times(1)
			},
			expectErrMsg: "",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.setupMock()

			// Execute
			err := suite.uc.UpdateTodo(ctx, tt.req)

			// Verify
			if tt.expectErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErrMsg)
			}
		})
	}
}

// helper functions for test
func stringPtr(s string) *string {
	return &s
}

func timeNow() time.Time {
	return time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
}

func (suite *TodoUseCaseTestSuite) TestCreateTodo() {
	ctx := context.Background()

	tests := []struct {
		name         string
		req          CreateTodoRequest
		setupMock    func()
		expectResp   *CreateTodoResponse
		expectErrMsg string
	}{
		{
			name: "create_todo_entity_fail",
			req:  CreateTodoRequest{}, // 空 title 會讓 NewTodo 返回錯誤
			setupMock: func() {
				// NewTodo 會失敗，所以不會調用 repository
			},
			expectResp:   nil,
			expectErrMsg: "validation fail",
		},
		{
			name: "create_todo_db_fail",
			req: CreateTodoRequest{
				Title:       "測試標題",
				Description: nil,
				Status:      "pending",
				DueDate:     nil,
			},
			setupMock: func() {
				// NewTodo 會成功，但 repository.Create 失敗
				suite.mockRepo.EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectResp:   nil,
			expectErrMsg: "internal fail",
		},
		{
			name: "create_todo_success",
			req: CreateTodoRequest{
				Title:       "測試標題",
				Description: nil,
				Status:      "pending",
				DueDate:     nil,
			},
			setupMock: func() {
				// NewTodo 成功，repository.Create 也成功
				suite.mockRepo.EXPECT().
					Create(ctx, gomock.Any()).
					Return(&entity.Todo{ID: 1}, nil).
					Times(1)
			},
			expectResp:   &CreateTodoResponse{ID: 1},
			expectErrMsg: "",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			// 設置 mock 期望
			tt.setupMock()

			// 執行測試
			resp, err := suite.uc.CreateTodo(ctx, tt.req)

			// 驗證結果
			assert.Equal(t, tt.expectResp, resp)

			if tt.expectErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErrMsg)
			}
		})
	}
}

func (suite *TodoUseCaseTestSuite) TestFindTodo() {
	ctx := context.Background()

	tests := []struct {
		name         string
		req          FindTodoRequest
		setupMock    func()
		expectResp   *FindTodoResponse
		expectErrMsg string
	}{
		{
			name: "todoRepo.List err",
			req: FindTodoRequest{
				Keyword: stringPtr("test"),
				Status:  stringPtr("pending"),
				Pagination: dto.PaginationReq{
					Page:      1,
					PageSize:  10,
					SortBy:    "created_at",
					SortOrder: "asc",
				},
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					List(ctx, gomock.Any(), gomock.Any()).
					Return(errors.New("database error")).
					Times(1)
			},
			expectResp:   nil,
			expectErrMsg: "internal fail",
		},
		{
			name: "invalid status ignored - success",
			req: FindTodoRequest{
				Keyword: stringPtr("test"),
				Status:  stringPtr("invalid_status"), // 無效的 status 會被忽略
				Pagination: dto.PaginationReq{
					Page:      1,
					PageSize:  10,
					SortBy:    "created_at",
					SortOrder: "desc",
				},
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					List(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, queryParams repository.TodoQueryParams, pagination *repository.Pagination[entity.Todo]) error {
						// 驗證無效的 status 沒有被設置
						assert.Nil(suite.T(), queryParams.Status)
						
						// 模擬成功的查詢結果
						pagination.Rows = []*entity.Todo{
							{
								ID:          1,
								Title:       "測試標題",
								Description: stringPtr("測試描述"),
								Status:      entity.StatusPending,
								CreatedAt:   timeNow(),
								UpdatedAt:   timeNow(),
							},
						}
						pagination.TotalRows = 1
						pagination.TotalPages = 1
						return nil
					}).
					Times(1)
			},
			expectResp: &FindTodoResponse{
				Todos: []TodoResponse{
					{
						ID:          1,
						Title:       "測試標題",
						Description: stringPtr("測試描述"),
						Status:      "pending",
						CreatedAt:   timeNow(),
						UpdatedAt:   timeNow(),
					},
				},
				Pagination: dto.PaginationResp{
					Page:       1,
					PageSize:   10,
					TotalCount: 1,
					TotalPages: 1,
				},
			},
			expectErrMsg: "",
		},
		{
			name: "success with valid status",
			req: FindTodoRequest{
				Keyword: stringPtr("測試"),
				Status:  stringPtr("doing"),
				Pagination: dto.PaginationReq{
					Page:      1,
					PageSize:  5,
					SortBy:    "title",
					SortOrder: "asc",
				},
			},
			setupMock: func() {
				suite.mockRepo.EXPECT().
					List(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, queryParams repository.TodoQueryParams, pagination *repository.Pagination[entity.Todo]) error {
						// 驗證有效的 status 有被正確設置
						assert.NotNil(suite.T(), queryParams.Status)
						assert.Equal(suite.T(), entity.StatusDoing, *queryParams.Status)
						
						// 模擬查詢結果
						pagination.Rows = []*entity.Todo{
							{
								ID:          2,
								Title:       "測試標題2",
								Description: nil,
								Status:      entity.StatusDoing,
								CreatedAt:   timeNow(),
								UpdatedAt:   timeNow(),
							},
							{
								ID:          3,
								Title:       "測試標題3",
								Description: stringPtr("詳細描述"),
								Status:      entity.StatusDoing,
								CreatedAt:   timeNow(),
								UpdatedAt:   timeNow(),
							},
						}
						pagination.TotalRows = 2
						pagination.TotalPages = 1
						return nil
					}).
					Times(1)
			},
			expectResp: &FindTodoResponse{
				Todos: []TodoResponse{
					{
						ID:          2,
						Title:       "測試標題2",
						Description: nil,
						Status:      "doing",
						CreatedAt:   timeNow(),
						UpdatedAt:   timeNow(),
					},
					{
						ID:          3,
						Title:       "測試標題3",
						Description: stringPtr("詳細描述"),
						Status:      "doing",
						CreatedAt:   timeNow(),
						UpdatedAt:   timeNow(),
					},
				},
				Pagination: dto.PaginationResp{
					Page:       1,
					PageSize:   5,
					TotalCount: 2,
					TotalPages: 1,
				},
			},
			expectErrMsg: "",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			// 設置 mock 期望
			tt.setupMock()

			// 執行測試
			resp, err := suite.uc.FindTodo(ctx, tt.req)

			// 驗證結果
			if tt.expectErrMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectResp, resp)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErrMsg)
				assert.Nil(t, resp)
			}
		})
	}
}

func (suite *TodoUseCaseTestSuite) TestDeleteTodo() {
	ctx := context.Background()

	tests := []struct {
		name         string
		id           uint
		setupMock    func()
		expectErrMsg string
	}{
		{
			name: "validation_fail_zero_id",
			id:   0,
			setupMock: func() {
				// No mock setup needed for validation error
			},
			expectErrMsg: "validation fail",
		},
		{
			name: "repository_delete_fail",
			id:   1,
			setupMock: func() {
				suite.mockRepo.EXPECT().
					Delete(ctx, uint(1)).
					Return(int64(0), errors.New("database error")).
					Times(1)
			},
			expectErrMsg: "internal fail",
		},
		{
			name: "todo_not_found_no_rows_affected",
			id:   999,
			setupMock: func() {
				suite.mockRepo.EXPECT().
					Delete(ctx, uint(999)).
					Return(int64(0), nil). // 0 rows affected means not found
					Times(1)
			},
			expectErrMsg: "not found",
		},
		{
			name: "success_delete_todo",
			id:   1,
			setupMock: func() {
				suite.mockRepo.EXPECT().
					Delete(ctx, uint(1)).
					Return(int64(1), nil). // 1 row affected means success
					Times(1)
			},
			expectErrMsg: "",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.setupMock()

			// Execute
			err := suite.uc.DeleteTodo(ctx, tt.id)

			// Verify
			if tt.expectErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErrMsg)
			}
		})
	}
}
