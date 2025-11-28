package router

import (
	"github.com/gorilla/mux"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/task-service/interface/http/handler"
	"github.com/todoist/backend/task-service/interface/http/middleware"
)

func NewRouter(taskHandler *handler.TaskHandler, log *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logging(log))

	// Health check
	r.HandleFunc("/health", taskHandler.HealthCheck).Methods("GET")

	// Task routes
	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", taskHandler.GetUserTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	return r
}
