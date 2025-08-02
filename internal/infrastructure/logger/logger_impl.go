package logger

import (
	"os"

	"github.com/rs/zerolog"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

var _ Logger = &LoggerImpl{}

// LoggerImpl implements the LoggerService interface.
type LoggerImpl struct {
	logger *zerolog.Logger
}

// InitLogger initializes the logger with the given configuration.
func (l *LoggerImpl) NewLogger(config *config.LogConfig) (zerolog.Logger, error) {
	// 設定日誌等級
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	// 設定輸出格式為 JSON
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// 建立 logger 實例
	logger := zerolog.New(os.Stdout).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	// 儲存 logger 實例
	l.logger = &logger

	return logger, nil
}
