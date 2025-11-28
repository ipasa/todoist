package usecase

import (
	"fmt"
	"net/smtp"

	"github.com/todoist/backend/notification-service/application/dto"
)

type SendNotificationUseCase struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	smtpFrom     string
}

func NewSendNotificationUseCase(host, port, username, password, from string) *SendNotificationUseCase {
	return &SendNotificationUseCase{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		smtpFrom:     from,
	}
}

func (uc *SendNotificationUseCase) Execute(req dto.SendNotificationRequest) error {
	// Construct email message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		uc.smtpFrom, req.Email, req.Subject, req.Body)

	// SMTP authentication
	auth := smtp.PlainAuth("", uc.smtpUsername, uc.smtpPassword, uc.smtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", uc.smtpHost, uc.smtpPort)
	err := smtp.SendMail(addr, auth, uc.smtpFrom, []string{req.Email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
