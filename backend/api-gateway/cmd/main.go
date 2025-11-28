package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
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
	r.PathPrefix("/v1/auth").Handler(createReverseProxy(authServiceURL, "/auth"))

	// Task service routes (protected)
	taskProxy := http.StripPrefix("/v1/tasks", createReverseProxy(taskServiceURL, "/tasks"))
	r.PathPrefix("/v1/tasks").Handler(middleware.Auth(cfg.JWTSecret)(taskProxy))

	// Project service routes (protected)
	projectProxy := http.StripPrefix("/v1/projects", createReverseProxy(projectServiceURL, "/projects"))
	r.PathPrefix("/v1/projects").Handler(middleware.Auth(cfg.JWTSecret)(projectProxy))

	// Notification service routes (protected)
	notifProxy := http.StripPrefix("/v1/notifications", createReverseProxy(notificationServiceURL, "/notifications"))
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

func createReverseProxy(target *url.URL, pathPrefix string) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = pathPrefix + req.URL.Path[len(pathPrefix):]
		req.Host = target.Host
	}

	return proxy
}
