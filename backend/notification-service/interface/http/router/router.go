package router

import (
	"github.com/gorilla/mux"
	"github.com/todoist/backend/notification-service/interface/http/handler"
	"github.com/todoist/backend/notification-service/interface/http/middleware"
	"github.com/todoist/backend/pkg/logger"
)

func NewRouter(notificationHandler *handler.NotificationHandler, log *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logging(log))

	// Health check
	r.HandleFunc("/health", notificationHandler.HealthCheck).Methods("GET")

	// Notification routes
	r.HandleFunc("/api/v1/notifications/send", notificationHandler.SendNotification).Methods("POST")

	return r
}
