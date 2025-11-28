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
	"github.com/todoist/backend/pkg/jwt"
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

	// Initialize database schema
	if err := initializeDatabase(db); err != nil {
		log.WithError(err).Fatal("failed to initialize database")
	}
	log.Info("database initialized")

	// Initialize dependencies
	taskRepo := postgres.NewTaskRepository(db)
	// Parse JWT expiry strings to time.Duration
	accessTokenExpiry, _ := time.ParseDuration(cfg.JWTExpiry)
	refreshTokenExpiry, _ := time.ParseDuration(cfg.RefreshTokenExpiry)
	jwtService := jwt.NewService(cfg.JWTSecret, accessTokenExpiry, refreshTokenExpiry)
	validatorInstance := validator.New()

	// Initialize handlers
	taskHandler := handler.NewTaskHandler(validatorInstance, log, taskRepo, jwtService)

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

// initializeDatabase creates the tables needed for the task service
func initializeDatabase(db *sql.DB) error {
	// Create tasks table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(255) NOT NULL,
		description TEXT,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		priority INTEGER NOT NULL DEFAULT 1,
		user_id UUID NOT NULL,
		project_id UUID,
		due_date TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
	`

	// Execute the SQL
	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}

	return nil
}
