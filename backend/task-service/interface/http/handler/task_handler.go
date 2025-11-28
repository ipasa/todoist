package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/pkg/validator"
	"github.com/todoist/backend/task-service/application/dto"
)

type TaskHandler struct {
	validator *validator.Validator
	logger    *logger.Logger
}

func NewTaskHandler(v *validator.Validator, log *logger.Logger) *TaskHandler {
	return &TaskHandler{
		validator: v,
		logger:    log,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement task creation logic
	h.respondWithJSON(w, http.StatusCreated, map[string]string{"message": "task created"})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	// TODO: Implement get task logic
	h.respondWithJSON(w, http.StatusOK, map[string]string{"id": taskID})
}

func (h *TaskHandler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement get user tasks logic
	h.respondWithJSON(w, http.StatusOK, []dto.TaskResponse{})
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// TODO: Implement update task logic
	h.respondWithJSON(w, http.StatusOK, map[string]string{"id": taskID, "message": "task updated"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	// TODO: Implement delete task logic
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "task deleted", "id": taskID})
}

func (h *TaskHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *TaskHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *TaskHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
