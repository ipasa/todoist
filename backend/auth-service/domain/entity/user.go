package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	AvatarURL    string
	Provider     string // email, google, github
	ProviderID   string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Version      int // for optimistic locking
}

// NewUser creates a new user entity
func NewUser(email, password, fullName, provider string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		Provider:     provider,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		Version:      1,
	}, nil
}

// NewOAuthUser creates a new user from OAuth provider
func NewOAuthUser(email, fullName, provider, providerID, avatarURL string) *User {
	now := time.Now()
	return &User{
		ID:         uuid.New(),
		Email:      email,
		FullName:   fullName,
		AvatarURL:  avatarURL,
		Provider:   provider,
		ProviderID: providerID,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
		Version:    1,
	}
}

// Authenticate verifies the password
func (u *User) Authenticate(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

// ChangePassword updates the user's password
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if err := u.Authenticate(oldPassword); err != nil {
		return err
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	u.PasswordHash = hashedPassword
	u.UpdatedAt = time.Now()
	u.Version++
	return nil
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(fullName, avatarURL string) {
	if fullName != "" {
		u.FullName = fullName
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}
	u.UpdatedAt = time.Now()
	u.Version++
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
	u.Version++
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
