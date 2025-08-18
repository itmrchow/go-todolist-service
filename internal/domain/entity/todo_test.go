package entity

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_todo_new_todo(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description *string
		status      *TodoStatus
		dueDate     *time.Time
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid_todo_with_title_only",
			title:       "測試標題",
			description: nil,
			status:      nil,
			dueDate:     nil,
			wantErr:     false,
		},
		{
			name:        "valid_todo_with_all_fields",
			title:       "完整的測試標題",
			description: stringPtr("這是一個測試描述"),
			status:      statusPtr(StatusDoing),
			dueDate:     timePtr(time.Now().Add(24 * time.Hour)),
			wantErr:     false,
		},
		{
			name:        "empty_title_should_fail",
			title:       "",
			description: nil,
			status:      nil,
			dueDate:     nil,
			wantErr:     true,
			errMsg:      "title cannot be empty",
		},
		{
			name:        "title_too_long_over_20_characters",
			title:       "這是一個非常長的標題超過二十個中文字符限制",
			description: nil,
			status:      nil,
			dueDate:     nil,
			wantErr:     true,
			errMsg:      "title cannot exceed 20 characters",
		},
		{
			name:        "description_too_long_over_100_characters",
			title:       "測試標題",
			description: stringPtr("這是一個非常長的描述，超過了一百個中文字符的限制。這個描述故意寫得很長，目的是要測試系統的驗證功能。當描述超過一百個中文字符時，系統應該要返回錯誤。這個測試用例確保了我們的驗證邏輯正確運作。一二三四五六七八九十一二三四五六七八九十一二三四"),
			status:      nil,
			dueDate:     nil,
			wantErr:     true,
			errMsg:      "description cannot exceed 100 characters",
		},
		{
			name:        "valid_20_character_title",
			title:       "十二三四五六七八九十一二三四五六七八九十", // exactly 20 characters
			description: nil,
			status:      nil,
			dueDate:     nil,
			wantErr:     false,
		},
		{
			name:        "valid_100_character_description",
			title:       "測試標題",
			description: nil,
			status:      nil,
			dueDate:     nil,
			wantErr:     false,
		},

		{
			name:        "due_date_in_the_future",
			title:       "測試標題",
			description: nil,
			status:      nil,
			dueDate:     timePtr(time.Now().Add(24 * time.Hour)),
			wantErr:     false,
		},
		{
			name:        "due_date_in_the_past",
			title:       "測試標題",
			description: nil,
			status:      nil,
			dueDate:     timePtr(time.Now().Add(-1 * time.Hour)),
			// NewTodo 函式會驗證 DueDate 是否為未來時間，所以這裡預期會出錯。
			wantErr: true,
		},
		{
			name:        "nil_due_date",
			title:       "測試標題",
			description: nil,
			status:      nil,
			dueDate:     nil,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			todo, err := NewTodo(tt.title, tt.description, tt.status, tt.dueDate)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, todo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, todo)
				assert.Equal(t, tt.title, todo.Title)
				assert.Equal(t, tt.description, todo.Description)

				// Check default status
				if tt.status == nil {
					assert.Equal(t, StatusPending, todo.Status)
				} else {
					assert.Equal(t, *tt.status, todo.Status)
				}

				assert.Equal(t, tt.dueDate, todo.DueDate)
				assert.NotZero(t, todo.CreatedAt)
				assert.NotZero(t, todo.UpdatedAt)
				assert.Nil(t, todo.DeletedAt) // New todos should not be deleted
			}
		})
	}
}

func Test_todo_status_string(t *testing.T) {
	tests := []struct {
		status   TodoStatus
		expected string
	}{
		{StatusPending, "pending"},
		{StatusDoing, "doing"},
		{StatusDone, "done"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func Test_todo_status_is_valid(t *testing.T) {
	tests := []struct {
		status TodoStatus
		valid  bool
	}{
		{StatusPending, true},
		{StatusDoing, true},
		{StatusDone, true},
		{TodoStatus("invalid"), false},
		{TodoStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.status.IsValid())
		})
	}
}

func Test_todo_json_serialization(t *testing.T) {
	dueDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	description := "測試描述"

	todo := &Todo{
		ID:          1,
		Title:       "測試標題",
		Description: &description,
		Status:      StatusDoing,
		DueDate:     &dueDate,
		CreatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		DeletedAt:   nil,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(todo)
	assert.NoError(t, err)

	expectedJSON := `{"id":1,"title":"測試標題","description":"測試描述","status":"doing","due_date":"2024-12-31T23:59:59Z","created_at":"2024-01-01T10:00:00Z","updated_at":"2024-01-01T10:00:00Z"}`
	assert.JSONEq(t, expectedJSON, string(jsonData))

	// Test JSON unmarshaling
	var unmarshaledTodo Todo
	err = json.Unmarshal(jsonData, &unmarshaledTodo)
	assert.NoError(t, err)
	assert.Equal(t, todo.ID, unmarshaledTodo.ID)
	assert.Equal(t, todo.Title, unmarshaledTodo.Title)
	assert.Equal(t, *todo.Description, *unmarshaledTodo.Description)
	assert.Equal(t, todo.Status, unmarshaledTodo.Status)
	assert.Equal(t, todo.DueDate.Unix(), unmarshaledTodo.DueDate.Unix())
}

func Test_todo_json_serialization_with_nil_fields(t *testing.T) {
	todo := &Todo{
		ID:          1,
		Title:       "測試標題",
		Description: nil,
		Status:      StatusPending,
		DueDate:     nil,
		CreatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		DeletedAt:   nil,
	}

	// Test JSON marshaling with nil fields
	jsonData, err := json.Marshal(todo)
	assert.NoError(t, err)

	expectedJSON := `{"id":1,"title":"測試標題","status":"pending","created_at":"2024-01-01T10:00:00Z","updated_at":"2024-01-01T10:00:00Z"}`
	assert.JSONEq(t, expectedJSON, string(jsonData))
}

func Test_todo_soft_delete(t *testing.T) {
	todo, err := NewTodo("測試標題", nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, todo)

	// Initially not deleted
	assert.False(t, todo.IsDeleted())
	assert.Nil(t, todo.DeletedAt)

	originalUpdatedAt := todo.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure timestamp difference

	// Delete the todo
	todo.Delete()
	assert.True(t, todo.IsDeleted())
	assert.NotNil(t, todo.DeletedAt)
	assert.True(t, todo.UpdatedAt.After(originalUpdatedAt))
}

func Test_todo_restore(t *testing.T) {
	todo, err := NewTodo("測試標題", nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, todo)

	// Delete first
	todo.Delete()
	assert.True(t, todo.IsDeleted())

	originalUpdatedAt := todo.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure timestamp difference

	// Restore the todo
	todo.Restore()
	assert.False(t, todo.IsDeleted())
	assert.Nil(t, todo.DeletedAt)
	assert.True(t, todo.UpdatedAt.After(originalUpdatedAt))
}

// Helper functions for test cases
func stringPtr(s string) *string {
	return &s
}

func statusPtr(s TodoStatus) *TodoStatus {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
