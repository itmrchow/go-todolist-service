package database

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

func TestDatabaseImpl_generateDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.DatabaseConfig
		expected string
	}{
		{
			name: "Complete configuration",
			config: &config.DatabaseConfig{
				Host:      "localhost",
				Port:      3306,
				Name:      "todolist_db",
				Account:   "root",
				Password:  "password123",
				URLSuffix: "?charset=utf8mb4&parseTime=True&loc=Local",
			},
			expected: "root:password123@tcp(localhost:3306)/todolist_db?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name: "Empty password",
			config: &config.DatabaseConfig{
				Host:      "localhost",
				Port:      3306,
				Name:      "test_db",
				Account:   "user",
				Password:  "",
				URLSuffix: "?charset=utf8mb4",
			},
			expected: "user:@tcp(localhost:3306)/test_db?charset=utf8mb4",
		},
		{
			name: "Different port",
			config: &config.DatabaseConfig{
				Host:      "db.example.com",
				Port:      3307,
				Name:      "app_db",
				Account:   "app_user",
				Password:  "secret",
				URLSuffix: "",
			},
			expected: "app_user:secret@tcp(db.example.com:3307)/app_db",
		},
		{
			name: "Production-like config",
			config: &config.DatabaseConfig{
				Host:      "prod-mysql.internal",
				Port:      3306,
				Name:      "production_db",
				Account:   "prod_user",
				Password:  "super_secure_password",
				URLSuffix: "?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s",
			},
			expected: "prod_user:super_secure_password@tcp(prod-mysql.internal:3306)/production_db?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &MySQLDBImpl{}
			result := db.generateDSN(tt.config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDatabaseImpl_Migrate_NoConnection(t *testing.T) {
	db := &MySQLDBImpl{}

	// 測試在沒有建立連接的情況下呼叫 Migrate
	err := db.Migrate(struct{}{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection not established")
}

func TestDatabaseImpl_Close_NoConnection(t *testing.T) {
	db := &MySQLDBImpl{}

	// 測試在沒有建立連接的情況下呼叫 Close
	err := db.Close()
	assert.NoError(t, err, "Close should not return error when no connection exists")
}
