package events

import (
	"time"

	"github.com/google/uuid"
)

// BaseEvent represents the common fields for all domain events
type BaseEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	UserID    uuid.UUID `json:"user_id"`
}

// NewBaseEvent creates a new base event
func NewBaseEvent(eventType string, userID uuid.UUID) BaseEvent {
	return BaseEvent{
		EventID:   uuid.New(),
		EventType: eventType,
		Timestamp: time.Now(),
		UserID:    userID,
	}
}
