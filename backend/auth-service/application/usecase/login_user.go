package usecase

import (
	"context"
	"time"

	"github.com/todoist/backend/auth-service/application/dto"
	"github.com/todoist/backend/auth-service/application/mapper"
	"github.com/todoist/backend/auth-service/domain/repository"
	"github.com/todoist/backend/pkg/errors"
	"github.com/todoist/backend/pkg/events"
	"github.com/todoist/backend/pkg/jwt"
)

// LoginUserUseCase handles user login
type LoginUserUseCase struct {
	userRepo       repository.UserRepository
	jwtService     *jwt.Service
	eventPublisher EventPublisher
}

// NewLoginUserUseCase creates a new LoginUserUseCase
func NewLoginUserUseCase(
	userRepo repository.UserRepository,
	jwtService *jwt.Service,
	eventPublisher EventPublisher,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:       userRepo,
		jwtService:     jwtService,
		eventPublisher: eventPublisher,
	}
}

// Execute logs in a user
func (uc *LoginUserUseCase) Execute(ctx context.Context, req dto.LoginDTO, ipAddress, userAgent string) (*dto.AuthResponseDTO, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.NewUnauthorizedError("account is deactivated")
	}

	// Verify password
	if err := user.Authenticate(req.Password); err != nil {
		return nil, errors.NewUnauthorizedError("invalid email or password")
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

	// Publish UserLoggedIn event
	event := events.NewUserLoggedIn(user.ID, ipAddress, userAgent)
	if err := uc.eventPublisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the login
	}

	// Return response
	return &dto.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(15 * time.Minute.Seconds()),
		User:         mapper.ToUserResponseDTO(user),
	}, nil
}
