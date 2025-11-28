package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/todoist/backend/notification-service/application/usecase"
	"github.com/todoist/backend/notification-service/infrastructure/config"
	"github.com/todoist/backend/notification-service/infrastructure/messaging"
	"github.com/todoist/backend/notification-service/interface/http/handler"
	"github.com/todoist/backend/notification-service/interface/http/router"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/pkg/validator"
)

func main() {
	// Initialize logger
	log, err := logger.New("notification-service")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer log.Sync()

	// Load configuration
	cfg := config.Load()

	// Initialize RabbitMQ consumer
	consumer, err := messaging.NewRabbitMQConsumer(cfg.RabbitMQURL, "notifications")
	if err != nil {
		log.WithError(err).Fatal("failed to initialize event consumer")
	}
	defer consumer.Close()
	log.Info("connected to RabbitMQ")

	// Initialize dependencies
	validatorInstance := validator.New()
	sendNotificationUseCase := usecase.NewSendNotificationUseCase(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)

	// Initialize handlers
	notificationHandler := handler.NewNotificationHandler(sendNotificationUseCase, validatorInstance, log)

	// Initialize router
	r := router.NewRouter(notificationHandler, log)

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

	// Start consuming messages in a goroutine
	go func() {
		log.Info("starting message consumer")
		if err := consumer.Consume(func(body []byte) error {
			log.WithFields(map[string]interface{}{
				"message": string(body),
			}).Info("received notification event")
			// TODO: Process notification events
			return nil
		}); err != nil {
			log.WithError(err).Error("consumer error")
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
