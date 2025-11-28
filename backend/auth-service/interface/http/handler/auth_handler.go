package handler

import (
	"encoding/json"
	"net/http"

	"github.com/todoist/backend/auth-service/application/dto"
	"github.com/todoist/backend/auth-service/application/usecase"
	"github.com/todoist/backend/pkg/errors"
	"github.com/todoist/backend/pkg/logger"
	"github.com/todoist/backend/pkg/validator"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	registerUseCase *usecase.RegisterUserUseCase
	loginUseCase    *usecase.LoginUserUseCase
	validator       *validator.Validator
	logger          *logger.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	registerUseCase *usecase.RegisterUserUseCase,
	loginUseCase *usecase.LoginUserUseCase,
	validator *validator.Validator,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
		validator:       validator,
		logger:          logger,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, errors.NewBadRequestError("invalid request body"))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		validationErrors := validator.GetValidationErrors(err)
		h.sendValidationError(w, validationErrors)
		return
	}

	response, err := h.registerUseCase.Execute(r.Context(), req)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusCreated, response)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, errors.NewBadRequestError("invalid request body"))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		validationErrors := validator.GetValidationErrors(err)
		h.sendValidationError(w, validationErrors)
		return
	}

	ipAddress := getIPAddress(r)
	userAgent := r.UserAgent()

	response, err := h.loginUseCase.Execute(r.Context(), req, ipAddress, userAgent)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, response)
}

// Helper methods

func (h *AuthHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) sendError(w http.ResponseWriter, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.NewInternalError("an unexpected error occurred", err)
	}

	h.logger.WithError(err).Error("request error")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	})
}

func (h *AuthHandler) sendValidationError(w http.ResponseWriter, validationErrors []validator.ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request data",
			"details": validationErrors,
		},
	})
}

func getIPAddress(r *http.Request) string {
	// Try to get real IP from X-Forwarded-For or X-Real-IP headers
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
