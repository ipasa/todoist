package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todoist/backend/pkg/jwt"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/pkg/validator"
	"github.com/todoist/backend/task-service/application/dto"
	"github.com/todoist/backend/task-service/application/usecase"
	"github.com/todoist/backend/task-service/domain"
)

type TaskHandler struct {
	validator       *validator.Validator
	logger         *logger.Logger
	createTaskUC   *usecase.CreateTaskUseCase
	getTaskUC      *usecase.GetTaskUseCase
	getUserTasksUC *usecase.GetUserTasksUseCase
	updateTaskUC   *usecase.UpdateTaskUseCase
	deleteTaskUC   *usecase.DeleteTaskUseCase
	jwtService     *jwt.Service
}

func NewTaskHandler(
	v *validator.Validator,
	log *logger.Logger,
	taskRepo domain.TaskRepository,
	jwtService *jwt.Service,
) *TaskHandler {
	return &TaskHandler{
		validator:       v,
		logger:         log,
		createTaskUC:   usecase.NewCreateTaskUseCase(taskRepo),
		getTaskUC:      usecase.NewGetTaskUseCase(taskRepo),
		getUserTasksUC: usecase.NewGetUserTasksUseCase(taskRepo),
		updateTaskUC:   usecase.NewUpdateTaskUseCase(taskRepo),
		deleteTaskUC:   usecase.NewDeleteTaskUseCase(taskRepo),
		jwtService:     jwtService,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("failed to decode request body")
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(req); err != nil {
		h.logger.WithError(err).Error("validation failed")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user ID from JWT token
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		h.logger.WithError(err).Error("failed to get user ID from token")
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"title":   req.Title,
	}).Info("creating task")

	// Create task
	task, err := h.createTaskUC.Execute(r.Context(), req, userID)
	if err != nil {
		h.logger.WithError(err).Error("failed to create task in use case")
		h.respondWithError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	// Get user ID from JWT token
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get task
	task, err := h.getTaskUC.Execute(r.Context(), taskID, userID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "failed to get task")
		return
	}

	h.respondWithJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse query parameters
	status := r.URL.Query().Get("status")
	priorityStr := r.URL.Query().Get("priority")
	projectID := r.URL.Query().Get("project_id")

	var priority *int
	if priorityStr != "" {
		p := parseInt(priorityStr)
		priority = &p
	}

	var projectIDPtr *string
	if projectID != "" {
		projectIDPtr = &projectID
	}

	// Get tasks
	tasks, err := h.getUserTasksUC.Execute(r.Context(), userID, status, priority, projectIDPtr)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "failed to get tasks")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":  tasks,
		"total": len(tasks),
	})
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get user ID from JWT token
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Update task
	task, err := h.updateTaskUC.Execute(r.Context(), taskID, userID, req)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	h.respondWithJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	// Get user ID from JWT token
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Delete task
	if err := h.deleteTaskUC.Execute(r.Context(), taskID, userID); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

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

func (h *TaskHandler) getUserIDFromToken(r *http.Request) (string, error) {
	// Get user ID from header set by API gateway
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return "", fmt.Errorf("X-User-ID header is required")
	}

	return userID, nil
}

func parseInt(s string) int {
	var i int
	_, _ = fmt.Sscanf(s, "%d", &i)
	return i
}
