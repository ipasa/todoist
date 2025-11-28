package events

import "github.com/google/uuid"

// UserRegistered event published when a new user registers
type UserRegistered struct {
	BaseEvent
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Provider string `json:"provider"` // email, google, github
}

// NewUserRegistered creates a new UserRegistered event
func NewUserRegistered(userID uuid.UUID, email, fullName, provider string) UserRegistered {
	return UserRegistered{
		BaseEvent: NewBaseEvent("auth.user.registered", userID),
		Email:     email,
		FullName:  fullName,
		Provider:  provider,
	}
}

// UserLoggedIn event published when a user logs in
type UserLoggedIn struct {
	BaseEvent
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

// NewUserLoggedIn creates a new UserLoggedIn event
func NewUserLoggedIn(userID uuid.UUID, ipAddress, userAgent string) UserLoggedIn {
	return UserLoggedIn{
		BaseEvent: NewBaseEvent("auth.user.logged_in", userID),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
}
