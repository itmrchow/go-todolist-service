package main

import (
	"github.com/rs/zerolog/log"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
	"itmrchow/go-todolist-service/internal/infrastructure/logger"
)

func main() {
	// Config
	config := &config.ConfigImpl{}
	configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal().Err(configErr).Str("module", "config").Msg("config init error")
	}

	// Logger
	logger := logger.LoggerImpl{}
	_, loggerErr := logger.NewLogger(config.GetLogConfig())
	if loggerErr != nil {
		log.Fatal().Err(loggerErr).Str("module", "logger").Msg("logger init error")
	}

	// DB

	// Router

}
