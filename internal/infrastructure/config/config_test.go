package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfigImpl tests the ConfigImpl implementation.
func TestConfigImpl(t *testing.T) {
	// Create Config
	config := &ConfigImpl{}
	configErr := config.LoadConfig()

	// assert error
	assert.NoError(t, configErr, "Config should load without error%s", configErr)

	// assert API server config info
	apiServerConfig := config.GetAPIServerConfig()
	assert.Equal(t, apiServerConfig.ServerPort, 8080, "API server port should be 8080")

	// assert Database config info
	dbConfig := config.GetDatabaseConfig()
	assert.Equal(t, dbConfig.URLSuffix, "?charset=utf8mb4&parseTime=True&loc=Local", "Database URL suffix should match expected value")
	assert.Equal(t, dbConfig.Host, "localhost", "Database host should be localhost")
	assert.Equal(t, dbConfig.Port, 3306, "Database port should be 3306")
	assert.Equal(t, dbConfig.Name, "todolist_db", "Database name should be todolist")
	assert.Equal(t, dbConfig.Account, "", "Database account should be empty")
	assert.Equal(t, dbConfig.Password, "", "Database password should be empty")

	// assert Log config info
	logConfig := config.GetLogConfig()
	assert.Equal(t, logConfig.Level, "debug", "Log level should be debug")
}
