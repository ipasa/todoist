package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/todoist/backend/api-gateway/internal/config"
	"github.com/todoist/backend/api-gateway/internal/middleware"
	"github.com/todoist/backend/pkg/logger"
)

func main() {
	// Initialize logger
	log, err := logger.New("api-gateway")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer log.Sync()

	// Load configuration
	cfg := config.Load()

	// Create router
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logging(log))

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Setup service proxies
	authServiceURL, _ := url.Parse(cfg.AuthServiceURL)
	taskServiceURL, _ := url.Parse(cfg.TaskServiceURL)
	projectServiceURL, _ := url.Parse(cfg.ProjectServiceURL)
	notificationServiceURL, _ := url.Parse(cfg.NotificationServiceURL)

	// Auth service routes
	authProxy := httputil.NewSingleHostReverseProxy(authServiceURL)
	authProxy.Director = func(req *http.Request) {
		req.URL.Scheme = authServiceURL.Scheme
		req.URL.Host = authServiceURL.Host
		// Strip /v1 prefix and keep /auth
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/v1")
		req.Host = authServiceURL.Host
	}
	r.PathPrefix("/v1/auth").Handler(authProxy)

	// Task service routes (protected)
	taskProxy := httputil.NewSingleHostReverseProxy(taskServiceURL)
	taskProxy.Director = func(req *http.Request) {
		req.URL.Scheme = taskServiceURL.Scheme
		req.URL.Host = taskServiceURL.Host
		// Strip /v1 prefix and keep /tasks
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/v1")
		req.Host = taskServiceURL.Host
	}
	r.PathPrefix("/v1/tasks").Handler(middleware.Auth(cfg.JWTSecret)(taskProxy))

	// Project service routes (protected)
	projectProxy := httputil.NewSingleHostReverseProxy(projectServiceURL)
	projectProxy.Director = func(req *http.Request) {
		req.URL.Scheme = projectServiceURL.Scheme
		req.URL.Host = projectServiceURL.Host
		// Strip /v1 prefix and keep /projects
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/v1")
		req.Host = projectServiceURL.Host
	}
	r.PathPrefix("/v1/projects").Handler(middleware.Auth(cfg.JWTSecret)(projectProxy))

	// Notification service routes (protected)
	notifProxy := httputil.NewSingleHostReverseProxy(notificationServiceURL)
	notifProxy.Director = func(req *http.Request) {
		req.URL.Scheme = notificationServiceURL.Scheme
		req.URL.Host = notificationServiceURL.Host
		// Strip /v1 prefix and keep /notifications
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/v1")
		req.Host = notificationServiceURL.Host
	}
	r.PathPrefix("/v1/notifications").Handler(middleware.Auth(cfg.JWTSecret)(notifProxy))

	// Start server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.WithFields(map[string]interface{}{
			"port": cfg.Port,
		}).Info("starting API Gateway")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("server forced to shutdown")
	}

	log.Info("server stopped")
}

