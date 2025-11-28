package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/pkg/validator"
	"github.com/todoist/backend/project-service/application/dto"
)

type ProjectHandler struct {
	validator *validator.Validator
	logger    *logger.Logger
}

func NewProjectHandler(v *validator.Validator, log *logger.Logger) *ProjectHandler {
	return &ProjectHandler{
		validator: v,
		logger:    log,
	}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement project creation logic
	h.respondWithJSON(w, http.StatusCreated, map[string]string{"message": "project created"})
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["id"]

	// TODO: Implement get project logic
	h.respondWithJSON(w, http.StatusOK, map[string]string{"id": projectID})
}

func (h *ProjectHandler) GetUserProjects(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement get user projects logic
	h.respondWithJSON(w, http.StatusOK, []dto.ProjectResponse{})
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["id"]

	var req dto.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// TODO: Implement update project logic
	h.respondWithJSON(w, http.StatusOK, map[string]string{"id": projectID, "message": "project updated"})
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["id"]

	// TODO: Implement delete project logic
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "project deleted", "id": projectID})
}

func (h *ProjectHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *ProjectHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *ProjectHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
