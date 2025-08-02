package database

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

var _ Database = &MySQLDBImpl{}

// MySQLDBImpl implements the Database interface.
type MySQLDBImpl struct {
	db  *gorm.DB
	ctx context.Context
}

// Connect establishes a connection to the MySQL database using GORM.
func (d *MySQLDBImpl) Connect(ctx context.Context, config *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := d.generateDSN(config)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
func (d *MySQLDBImpl) monitorContext() {
	<-d.ctx.Done()
	if err := d.Close(); err != nil {
		// 這裡無法使用 log，因為可能會造成循環依賴
		// 在實際專案中可以考慮使用內建的 log 包或其他方式
	}
}

// Migrate runs GORM auto migration for the given models.
func (d *MySQLDBImpl) Migrate(models ...interface{}) error {
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
func (d *MySQLDBImpl) Close() error {
	if d.db == nil {
		return nil
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	return sqlDB.Close()
}

// generateDSN generates the Data Source Name for MySQL connection.
func (d *MySQLDBImpl) generateDSN(config *config.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s",
		config.Account,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.URLSuffix,
	)
}
