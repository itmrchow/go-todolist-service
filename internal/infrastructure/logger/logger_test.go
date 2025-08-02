package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

func TestLoggerService_InitLogger(t *testing.T) {
	tests := []struct {
		name        string
		logConfig   *config.LogConfig
		expectError bool
	}{
		{
			name: "Valid debug level",
			logConfig: &config.LogConfig{
				Level: "debug",
			},
			expectError: false,
		},
		{
			name: "Valid info level",
			logConfig: &config.LogConfig{
				Level: "info",
			},
			expectError: false,
		},
		{
			name: "Valid warn level",
			logConfig: &config.LogConfig{
				Level: "warn",
			},
			expectError: false,
		},
		{
			name: "Valid error level",
			logConfig: &config.LogConfig{
				Level: "error",
			},
			expectError: false,
		},
		{
			name: "Invalid log level",
			logConfig: &config.LogConfig{
				Level: "invalid",
			},
			expectError: true,
		},
		{
			name: "Empty log level should use default",
			logConfig: &config.LogConfig{
				Level: "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &LoggerImpl{}
			mLog, err := logger.NewLogger(tt.logConfig)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, mLog)
			}
		})
	}
}

func TestLoggerService_LoggerConfiguration(t *testing.T) {
	loggerService := &LoggerImpl{}

	// 測試初始化後可以正常使用日誌
	logConfig := &config.LogConfig{
		Level: "info",
	}

	logger, err := loggerService.NewLogger(logConfig)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// 測試日誌記錄不會產生錯誤
	logger.Info().Msg("test info message")
	logger.Debug().Msg("test debug message")
	logger.Error().Msg("test error message")
}

func TestLoggerService_JSONOutput(t *testing.T) {
	// 測試 JSON 格式輸出
	loggerService := &LoggerImpl{}

	logConfig := &config.LogConfig{
		Level: "info",
	}

	logger, err := loggerService.NewLogger(logConfig)
	assert.NoError(t, err)

	// 確認日誌輸出為 JSON 格式（通過檢查 zerolog 實例）
	assert.IsType(t, zerolog.Logger{}, logger)
}
