package dto

type SendNotificationRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

type SendNotificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
