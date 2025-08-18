package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"itmrchow/go-todolist-service/internal/domain/entity"
	"itmrchow/go-todolist-service/internal/domain/repository"
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
