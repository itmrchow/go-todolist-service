package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"itmrchow/go-todolist-service/internal/infrastructure/config"
)

var _ ServerService = &ServerImpl{}

// ServerImpl implements the ServerService interface.
type ServerImpl struct {
	httpServer *http.Server
	addr       string
	ctx        context.Context
}

// NewServer creates a new server instance.
func NewServer() *ServerImpl {
	return &ServerImpl{}
}

// Start starts the HTTP server with the given configuration and router.
func (s *ServerImpl) Start(ctx context.Context, config *config.APIServerConfig, router *gin.Engine) error {
	s.ctx = ctx
	s.addr = fmt.Sprintf(":%d", config.ServerPort)

	s.httpServer = &http.Server{
		Addr:           s.addr,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 監聽 context 取消信號，自動關閉服務器
	go s.monitorContext()

	return s.httpServer.ListenAndServe()
}

// monitorContext 監聽 context 取消信號並自動關閉服務器
func (s *ServerImpl) monitorContext() {
	<-s.ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		// 這裡無法使用 log，因為可能會造成循環依賴
		// 在實際專案中可以考慮使用內建的 log 包或其他方式
	}
}

// Stop gracefully stops the HTTP server.
func (s *ServerImpl) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	return s.httpServer.Shutdown(ctx)
}

// GetAddr returns the server address.
func (s *ServerImpl) GetAddr() string {
	return s.addr
}
