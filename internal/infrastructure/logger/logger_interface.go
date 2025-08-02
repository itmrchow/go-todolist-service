package logger

import (
	"github.com/rs/zerolog"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

// Logger defines the interface for logger management.
type Logger interface {
	NewLogger(config *config.LogConfig) (zerolog.Logger, error)
}
