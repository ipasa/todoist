package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/pkg/validator"
	"github.com/todoist/backend/task-service/infrastructure/config"
	"github.com/todoist/backend/task-service/infrastructure/persistence/postgres"
	"github.com/todoist/backend/task-service/interface/http/handler"
	"github.com/todoist/backend/task-service/interface/http/router"
)

func main() {
	// Initialize logger
	log, err := logger.New("task-service")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer log.Sync()

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.WithError(err).Fatal("failed to ping database")
	}
	log.Info("connected to database")

	// Initialize dependencies
	taskRepo := postgres.NewTaskRepository(db)
	validatorInstance := validator.New()

	// Initialize handlers
	taskHandler := handler.NewTaskHandler(validatorInstance, log)

	// Initialize router
	r := router.NewRouter(taskHandler, log)

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
		}).Info("starting HTTP server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("failed to start server")
		}
	}()

	// Prevent unused variable error
	_ = taskRepo

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
