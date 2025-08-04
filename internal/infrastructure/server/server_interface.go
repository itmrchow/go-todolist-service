package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

// ServerService defines the interface for HTTP server management.
type ServerService interface {
	Start(ctx context.Context, config *config.APIServerConfig, router *gin.Engine) error
	Stop(ctx context.Context) error
	GetAddr() string
}
