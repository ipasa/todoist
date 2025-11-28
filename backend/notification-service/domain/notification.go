package domain

import "time"

type NotificationType string

const (
	NotificationTypeEmail    NotificationType = "email"
	NotificationTypePush     NotificationType = "push"
	NotificationTypeInApp    NotificationType = "in_app"
	NotificationTypeWebhook  NotificationType = "webhook"
)

type Notification struct {
	ID        string
	UserID    string
	Type      NotificationType
	Subject   string
	Body      string
	CreatedAt time.Time
	SentAt    *time.Time
}
