package handler

import (
	"encoding/json"
	"net/http"

	"github.com/todoist/backend/notification-service/application/dto"
	"github.com/todoist/backend/notification-service/application/usecase"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/pkg/validator"
)

type NotificationHandler struct {
	sendNotificationUseCase *usecase.SendNotificationUseCase
	validator               *validator.Validator
	logger                  *logger.Logger
}

func NewNotificationHandler(
	sendNotificationUseCase *usecase.SendNotificationUseCase,
	v *validator.Validator,
	log *logger.Logger,
) *NotificationHandler {
	return &NotificationHandler{
		sendNotificationUseCase: sendNotificationUseCase,
		validator:               v,
		logger:                  log,
	}
}

func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req dto.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.sendNotificationUseCase.Execute(req); err != nil {
		h.logger.WithError(err).Error("failed to send notification")
		h.respondWithError(w, http.StatusInternalServerError, "failed to send notification")
		return
	}

	h.respondWithJSON(w, http.StatusOK, dto.SendNotificationResponse{
		Success: true,
		Message: "notification sent successfully",
	})
}

func (h *NotificationHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *NotificationHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *NotificationHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
