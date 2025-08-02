package main

import (
	"github.com/rs/zerolog/log"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

func main() {
	// Config
	config := &config.ConfigImpl{}
	configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal().Err(configErr).Str("module", "config").Msg("config init error")
	}

	// Server

}
