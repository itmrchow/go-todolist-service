package v1

type DeleteTodoRequest struct {
	ID uint `json:"id" binding:"required"`
}
