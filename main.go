package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"itmrchow/go-todolist-service/internal/delivery/http/handler"
	v1 "itmrchow/go-todolist-service/internal/delivery/http/handler/v1"
	"itmrchow/go-todolist-service/internal/domain/usecase"
	"itmrchow/go-todolist-service/internal/infrastructure/config"
	"itmrchow/go-todolist-service/internal/infrastructure/database"
	"itmrchow/go-todolist-service/internal/infrastructure/database/model"
	"itmrchow/go-todolist-service/internal/infrastructure/logger"
	"itmrchow/go-todolist-service/internal/infrastructure/repository"
	"itmrchow/go-todolist-service/internal/infrastructure/router"
	"itmrchow/go-todolist-service/internal/infrastructure/server"
)

func main() {
	// 建立根 context 和 cancel 函數
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 設定信號監聽 channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

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

	// DB - 傳入 context，讓資料庫可以監聽取消信號
	db := &database.MySQLDBImpl{}
	gormDb, dbErr := db.Connect(ctx, config.GetDatabaseConfig())
	if dbErr != nil {
		log.Fatal().Err(dbErr).Str("module", "database").Msg("database connection error")
	}

	// Run database migrations
	migrationErr := db.Migrate(&model.Todo{})
	if migrationErr != nil {
		log.Fatal().Err(migrationErr).Str("module", "database").Msg("database migration error")
	}
	log.Info().Str("module", "database").Msg("Database migration completed successfully")

	// Repository
	todoRepo := repository.NewTodoRepository(gormDb, &logger)

	// Usecase
	todoUc := usecase.NewTodoUseCaseImpl(todoRepo)

	// Router handlers
	healthHandler := handler.NewHealthHandler()
	todoV1Handler := v1.NewTodoHandlerImpl(todoUc) // 假設有一個 TodoUseCase

	// Router
	appRouter := router.NewRouter(
		healthHandler,
		todoV1Handler,
	)
	engine := appRouter.SetupRoutes()

	// Server - 傳入 context，讓服務器可以監聽取消信號
	httpServer := server.NewServer()

	// 在 goroutine 中啟動服務器，避免阻塞
	go func() {
		log.Info().Str("module", "server").Msg("Starting HTTP server...")
		if err := httpServer.Start(ctx, config.GetAPIServerConfig(), engine); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal().Err(err).Str("module", "server").Msg("server start error")
			}
		}
	}()

	// 等待關閉信號
	<-quit
	log.Info().Str("module", "server").Msg("Shutting down server...")

	// 呼叫 cancel()，通知所有模組開始關閉
	var closeErr error
	if closeErr = db.Close(); closeErr != nil {
		log.Error().Err(closeErr).Str("module", "close").Msg("database close error")
	}
	if closeErr = httpServer.Stop(ctx); closeErr != nil {
		log.Error().Err(closeErr).Str("module", "close").Msg("http server close error")
	}

	// 給予一些時間讓各模組完成關閉
	time.Sleep(2 * time.Second)

	log.Info().Str("module", "server").Msg("Server exited")
}
