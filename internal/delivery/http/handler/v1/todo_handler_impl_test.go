package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	v1 "itmrchow/go-todolist-service/internal/delivery/http/dto/v1"
	"itmrchow/go-todolist-service/internal/domain/usecase"
	"itmrchow/go-todolist-service/internal/utils/dto"
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

	suite.handler = NewTodoHandlerImpl(zerolog.New(os.Stdout), suite.mockTodoUc)
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
					Return(&usecase.CreateTodoResponse{ID: 1}, nil)
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

			// create gin context
			c, w := CreateGinContext("/create-todo", tt.body)

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

func (suite *TodoHandlerImplTestSuite) TestTodoHandlerImpl_FindTodo() {
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
			name: "UseCase FindTodo Fail",
			body: map[string]interface{}{
				"keyword": "test",
				"status":  "pending",
				"pagination": map[string]interface{}{
					"page":       1,
					"page_size":  10,
					"sort_by":    "created_at",
					"sort_order": "desc",
				},
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					FindTodo(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal database error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: map[string]interface{}{
				"error": "internal server error",
			},
		},
		{
			name: "Success",
			body: map[string]interface{}{
				"keyword": "test",
				"status":  "pending",
				"pagination": map[string]interface{}{
					"page":       1,
					"page_size":  10,
					"sort_by":    "created_at",
					"sort_order": "desc",
				},
			},
			mockSetup: func() {
				now := time.Now()
				suite.mockTodoUc.EXPECT().
					FindTodo(gomock.Any(), gomock.Any()).
					Return(&usecase.FindTodoResponse{
						Todos: []usecase.TodoResponse{
							{
								ID:          1,
								Title:       "test todo",
								Description: stringPtr("test description"),
								Status:      "pending",
								DueDate:     &now,
								CreatedAt:   now,
								UpdatedAt:   now,
							},
						},
						Pagination: dto.PaginationResp{
							Page:       1,
							PageSize:   10,
							TotalCount: 1,
							TotalPages: 1,
						},
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedResp: func() v1.FindTodoResponse {
				now := time.Now()
				return v1.FindTodoResponse{
					Todos: []v1.TodoItem{
						{
							ID:          1,
							Title:       "test todo",
							Description: stringPtr("test description"),
							Status:      "pending",
							DueDate:     &now,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
					},
					Pagination: dto.PaginationResp{
						Page:       1,
						PageSize:   10,
						TotalCount: 1,
						TotalPages: 1,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Setup
			tt.mockSetup()

			// Create request
			c, w := CreateGinContext("/find-todo", tt.body)

			// Execute
			suite.handler.FindTodo(c)

			// Assert
			assert.Equal(suite.T(), tt.expectedCode, w.Code)

			// Handle different response types for assertion
			if httpResp, ok := tt.expectedResp.(v1.FindTodoResponse); ok {
				var resp v1.FindTodoResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(suite.T(), err)
				
				// Compare non-time fields first
				assert.Equal(suite.T(), len(httpResp.Todos), len(resp.Todos))
				if len(httpResp.Todos) > 0 && len(resp.Todos) > 0 {
					expectedTodo := httpResp.Todos[0]
					actualTodo := resp.Todos[0]
					
					assert.Equal(suite.T(), expectedTodo.ID, actualTodo.ID)
					assert.Equal(suite.T(), expectedTodo.Title, actualTodo.Title)
					assert.Equal(suite.T(), expectedTodo.Description, actualTodo.Description)
					assert.Equal(suite.T(), expectedTodo.Status, actualTodo.Status)
					
					// For time fields, just check they are not zero values
					if expectedTodo.DueDate != nil {
						assert.NotNil(suite.T(), actualTodo.DueDate)
					}
					assert.False(suite.T(), actualTodo.CreatedAt.IsZero())
					assert.False(suite.T(), actualTodo.UpdatedAt.IsZero())
				}
				assert.Equal(suite.T(), httpResp.Pagination, resp.Pagination)
			} else {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResp, resp)
			}
		})
	}
}

func (suite *TodoHandlerImplTestSuite) TestTodoHandlerImpl_UpdateTodo() {
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
			name: "Missing ID",
			body: map[string]interface{}{
				"title":       "updated title",
				"description": "updated description",
				"status":      "doing",
			},
			mockSetup: func() {
				// no mock setup needed for validation error
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": "invalid request format",
			},
		},
		{
			name: "Invalid Status",
			body: map[string]interface{}{
				"id":          1,
				"title":       "updated title",
				"description": "updated description", 
				"status":      "invalid_status",
			},
			mockSetup: func() {
				// no mock setup needed for validation error
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": "invalid request format",
			},
		},
		{
			name: "UseCase Not Found Error",
			body: map[string]interface{}{
				"id":          999,
				"title":       "updated title",
				"description": "updated description",
				"status":      "doing",
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					UpdateTodo(gomock.Any(), gomock.Any()).
					Return(errors.New("not found: todo not found")).
					Times(1)
			},
			expectedCode: http.StatusNotFound,
			expectedResp: map[string]interface{}{
				"error": "not found: todo not found",
			},
		},
		{
			name: "UseCase Validation Error",
			body: map[string]interface{}{
				"id":          1,
				"title":       "valid title",
				"description": "updated description",
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					UpdateTodo(gomock.Any(), gomock.Any()).
					Return(errors.New("validation fail: due date must be in the future")).
					Times(1)
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": "validation fail: due date must be in the future",
			},
		},
		{
			name: "UseCase Internal Error",
			body: map[string]interface{}{
				"id":          1,
				"title":       "valid title",
				"description": "updated description",
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					UpdateTodo(gomock.Any(), gomock.Any()).
					Return(errors.New("internal fail: database connection error")).
					Times(1)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: map[string]interface{}{
				"error": "internal server error",
			},
		},
		{
			name: "Success - Full Update",
			body: map[string]interface{}{
				"id":          1,
				"title":       "updated title",
				"description": "updated description",
				"status":      "doing",
				"due_date":    "2024-12-31T23:59:59Z",
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					UpdateTodo(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedCode: http.StatusNoContent,
			expectedResp: nil, // No response body for 204
		},
		{
			name: "Success - Partial Update (no status)",
			body: map[string]interface{}{
				"id":          2,
				"title":       "updated title only",
				"description": "updated description",
			},
			mockSetup: func() {
				suite.mockTodoUc.EXPECT().
					UpdateTodo(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedCode: http.StatusNoContent,
			expectedResp: nil, // No response body for 204
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.mockSetup()

			// Create request
			var body []byte
			var err error
			if str, ok := tt.body.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.body)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/update-todo", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			// Call handler
			suite.handler.UpdateTodo(ctx)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedResp != nil {
				var actualResp interface{}
				err = json.Unmarshal(w.Body.Bytes(), &actualResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, actualResp)
			} else {
				// For 204 responses, body should be empty
				assert.Empty(t, w.Body.String())
			}
		})
	}
}

func CreateGinContext(target string, body interface{}) (*gin.Context,
	*httptest.ResponseRecorder) {
	// Create request
	var reqBody []byte
	if str, ok := body.(string); ok {
		reqBody = []byte(str)
	} else {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c, w
}

// stringPtr is a helper function to create a pointer to string
func stringPtr(s string) *string {
	return &s
}
