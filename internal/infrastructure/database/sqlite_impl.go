package database

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

var _ Database = &SQLiteDBImpl{}

// SQLiteDBImpl implements the Database interface.
type SQLiteDBImpl struct {
	db  *gorm.DB
	ctx context.Context
}

// Connect establishes a connection to the SQLite database using GORM.
func (d *SQLiteDBImpl) Connect(ctx context.Context, config *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := d.generateDSN(config)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	d.db = db
	d.ctx = ctx

	// 監聽 context 取消信號，自動關閉連線
	go d.monitorContext()

	return db, nil
}

// monitorContext 監聽 context 取消信號並自動關閉資料庫連線
func (d *SQLiteDBImpl) monitorContext() {
	<-d.ctx.Done()
	if err := d.Close(); err != nil {
		// 這裡無法使用 log，因為可能會造成循環依賴
		// 在實際專案中可以考慮使用內建的 log 包或其他方式
	}
}

// Migrate runs GORM auto migration for the given models.
func (d *SQLiteDBImpl) Migrate(models ...interface{}) error {
	if d.db == nil {
		return fmt.Errorf("database connection not established")
	}

	for _, model := range models {
		if err := d.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}

	return nil
}

// Close closes the database connection.
func (d *SQLiteDBImpl) Close() error {
	if d.db == nil {
		return nil
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	return sqlDB.Close()
}

// generateDSN generates the Data Source Name for SQLite connection.
func (d *SQLiteDBImpl) generateDSN(config *config.DatabaseConfig) string {
	// 對於 SQLite，我們可以使用 config.Name 作為檔案路徑
	// 或者使用 ":memory:" 來建立記憶體資料庫
	if config.Name == ":memory:" {
		return ":memory:"
	}
	
	// 如果需要檔案路徑，可以使用 config.Name
	return config.Name
}