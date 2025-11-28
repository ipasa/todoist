package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/todoist/backend/auth-service/domain/entity"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Save creates a new user
	Save(ctx context.Context, user *entity.User) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByProvider retrieves a user by OAuth provider and provider ID
	FindByProvider(ctx context.Context, provider, providerID string) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// Exists checks if a user with the given email exists
	Exists(ctx context.Context, email string) (bool, error)
}
