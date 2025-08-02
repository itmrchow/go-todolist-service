package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health checks the health status of the service
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "todolist-service",
	})
}

// Version returns the version information of the service
func (h *HealthHandler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":   "1.0.0",
		"service":   "todolist-service",
		"build_env": "development",
	})
}
