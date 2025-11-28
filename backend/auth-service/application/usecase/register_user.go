package usecase

import (
	"context"
	"time"

	"github.com/todoist/backend/auth-service/application/dto"
	"github.com/todoist/backend/auth-service/application/mapper"
	"github.com/todoist/backend/auth-service/domain/entity"
	"github.com/todoist/backend/auth-service/domain/repository"
	"github.com/todoist/backend/pkg/errors"
	"github.com/todoist/backend/pkg/events"
	"github.com/todoist/backend/pkg/jwt"
)

// RegisterUserUseCase handles user registration
type RegisterUserUseCase struct {
	userRepo   repository.UserRepository
	jwtService *jwt.Service
	eventPublisher EventPublisher
}

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}

// NewRegisterUserUseCase creates a new RegisterUserUseCase
func NewRegisterUserUseCase(
	userRepo repository.UserRepository,
	jwtService *jwt.Service,
	eventPublisher EventPublisher,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:       userRepo,
		jwtService:     jwtService,
		eventPublisher: eventPublisher,
	}
}

// Execute registers a new user
func (uc *RegisterUserUseCase) Execute(ctx context.Context, req dto.RegisterUserDTO) (*dto.AuthResponseDTO, error) {
	// Check if user already exists
	exists, err := uc.userRepo.Exists(ctx, req.Email)
	if err != nil {
		return nil, errors.NewInternalError("failed to check user existence", err)
	}
	if exists {
		return nil, errors.NewConflictError("user with this email already exists")
	}

	// Create new user entity
	user, err := entity.NewUser(req.Email, req.Password, req.FullName, "email")
	if err != nil {
		return nil, errors.NewInternalError("failed to create user", err)
	}

	// Save user to database
	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, errors.NewInternalError("failed to save user", err)
	}

	// Generate tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.NewInternalError("failed to generate access token", err)
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.NewInternalError("failed to generate refresh token", err)
	}

	// Publish UserRegistered event
	event := events.NewUserRegistered(user.ID, user.Email, user.FullName, user.Provider)
	if err := uc.eventPublisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the registration
		// In production, consider using a retry mechanism or dead-letter queue
	}

	// Return response
	return &dto.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(15 * time.Minute.Seconds()),
		User:         mapper.ToUserResponseDTO(user),
	}, nil
}
