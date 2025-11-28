package dto

import "github.com/google/uuid"

// RegisterUserDTO represents registration request data
type RegisterUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

// LoginDTO represents login request data
type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponseDTO represents user response data
type UserResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
	Provider  string    `json:"provider"`
}

// AuthResponseDTO represents authentication response
type AuthResponseDTO struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	ExpiresIn    int64           `json:"expires_in"`
	User         UserResponseDTO `json:"user"`
}

// RefreshTokenDTO represents refresh token request
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateProfileDTO represents profile update request
type UpdateProfileDTO struct {
	FullName  string `json:"full_name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}
