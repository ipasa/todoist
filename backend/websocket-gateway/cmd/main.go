package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/websocket-gateway/domain"
	"github.com/todoist/backend/websocket-gateway/infrastructure/config"
	wsHandler "github.com/todoist/backend/websocket-gateway/interface/websocket"
)

func main() {
	// Initialize logger
	log, err := logger.New("websocket-gateway")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer log.Sync()

	// Load configuration
	cfg := config.Load()

	// Initialize WebSocket hub
	hub := domain.NewHub()
	go hub.Run()

	// Initialize handlers
	handler := wsHandler.NewHandler(hub, log)

	// Initialize router
	r := mux.NewRouter()
	r.HandleFunc("/ws", handler.HandleWebSocket)
	r.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Start HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithFields(map[string]interface{}{
			"port": cfg.Port,
		}).Info("starting WebSocket server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("server forced to shutdown")
	}

	log.Info("server stopped")
}
