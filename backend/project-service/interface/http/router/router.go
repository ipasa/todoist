package router

import (
	"github.com/gorilla/mux"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/project-service/interface/http/handler"
	"github.com/todoist/backend/project-service/interface/http/middleware"
)

func NewRouter(projectHandler *handler.ProjectHandler, log *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logging(log))

	// Health check
	r.HandleFunc("/health", projectHandler.HealthCheck).Methods("GET")

	// Project routes
	r.HandleFunc("/projects", projectHandler.CreateProject).Methods("POST")
	r.HandleFunc("/projects", projectHandler.GetUserProjects).Methods("GET")
	r.HandleFunc("/projects/{id}", projectHandler.GetProject).Methods("GET")
	r.HandleFunc("/projects/{id}", projectHandler.UpdateProject).Methods("PUT")
	r.HandleFunc("/projects/{id}", projectHandler.DeleteProject).Methods("DELETE")

	return r
}
