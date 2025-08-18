package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	v1 "itmrchow/go-todolist-service/internal/delivery/http/dto/v1"
	"itmrchow/go-todolist-service/internal/domain/usecase"
)

type TodoHandlerImplTestSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	mockTodoUc *usecase.MockTodoUseCase
	handler    *TodoHandlerImpl
}

func TestTodoHandlerImplTestSuite(t *testing.T) {

	suite.Run(t, new(TodoHandlerImplTestSuite))
}

func (suite *TodoHandlerImplTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockTodoUc = usecase.NewMockTodoUseCase(suite.ctrl)

	suite.handler = NewTodoHandlerImpl(suite.mockTodoUc)
}

func (suite *TodoHandlerImplTestSuite) TearDownTest() {
	if suite.ctrl != nil {
		suite.ctrl.Finish()
	}
}

func (suite *TodoHandlerImplTestSuite) TestTodoHandlerImpl_CreateTodo() {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         interface{}
		mockSetup    func()
		expectedCode int
		expectedResp interface{}
	}{
		{
			name: "Invalid JSON",
			body: "invalid json",
			mockSetup: func() {
				// no mock setup needed for JSON parsing error
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": "invalid request format",
			},
		},
		{
			name: "UseCase Validation Fail",
			body: map[string]interface{}{
				"title":       "test",
				"description": "test description",
				"status":      "pending",
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("validation fail: title cannot be empty"))
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": "validation fail: title cannot be empty",
			},
		},
		{
			name: "UseCase Internal Fail",
			body: map[string]interface{}{
				"title":       "test",
				"description": "test description",
				"status":      "pending",
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal fail: database connection error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: map[string]interface{}{
				"error": "internal server error",
			},
		},
		{
			name: "Success: Request not status",
			body: map[string]interface{}{
				"title":       "test",
				"description": "test description",
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("validation fail: title cannot be empty"))
			},
			expectedCode: http.StatusOK,
			expectedResp: v1.CreateTodoResponse{
				ID: 1,
			},
		},
		{
			name: "Success",
			body: map[string]interface{}{
				"title":       "test todo",
				"description": "test description",
				"status":      "pending",
				"due_date":    "2024-12-31T23:59:59Z",
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Return(&usecase.CreateTodoResponse{ID: 1}, nil)
			},
			expectedCode: http.StatusOK,
			expectedResp: v1.CreateTodoResponse{
				ID: 1,
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Setup
			tt.mockSetup()

			// Create request
			var reqBody []byte
			if str, ok := tt.body.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/create-todo", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Create gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			suite.handler.CreateTodo(c)

			// Assert
			assert.Equal(suite.T(), tt.expectedCode, w.Code)

			// Handle different response types for assertion
			if httpResp, ok := tt.expectedResp.(v1.CreateTodoResponse); ok {
				var resp v1.CreateTodoResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), httpResp, resp)
			} else {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResp, resp)
			}
		})
	}
}
