package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError represents a custom application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code int, message string, details ...string) *AppError {
	appErr := &AppError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		appErr.Details = details[0]
	}
	return appErr
}

// ErrorHandler returns a middleware that handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 處理在處理過程中產生的錯誤
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			switch e := err.Err.(type) {
			case *AppError:
				// 處理自定義應用程式錯誤
				c.JSON(e.Code, gin.H{
					"error": gin.H{
						"code":    e.Code,
						"message": e.Message,
						"details": e.Details,
					},
				})
			default:
				// 處理其他未預期的錯誤
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code":    http.StatusInternalServerError,
						"message": "Internal server error",
					},
				})
			}

			// 確保回應已經發送
			c.Abort()
		}
	}
}
