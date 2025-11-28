package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todoist/backend/auth-service/interface/http/handler"
	"github.com/todoist/backend/auth-service/interface/http/middleware"
	"github.com/todoist/backend/pkg/logger"
)

// NewRouter creates a new HTTP router
func NewRouter(authHandler *handler.AuthHandler, log *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logging(log))

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Auth routes
	r.HandleFunc("/auth/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)

	return r
}
