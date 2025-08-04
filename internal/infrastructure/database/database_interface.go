package database

import (
	"context"

	"gorm.io/gorm"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

// Database defines the interface for database management.
type Database interface {
	Connect(ctx context.Context, config *config.DatabaseConfig) (*gorm.DB, error)
	Migrate(models ...interface{}) error
	Close() error
}
