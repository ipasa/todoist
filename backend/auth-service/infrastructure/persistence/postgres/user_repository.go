package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/todoist/backend/auth-service/domain/entity"
	pkgErrors "github.com/todoist/backend/pkg/errors"
)

// UserRepository implements the repository interface using PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save creates a new user
func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, avatar_url, provider, provider_id, is_active, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.AvatarURL,
		user.Provider,
		user.ProviderID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
		user.Version,
	)
	if err != nil {
		return pkgErrors.NewInternalError("failed to save user", err)
	}
	return nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, provider, provider_id, is_active, created_at, updated_at, version
		FROM users
		WHERE id = $1
	`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.AvatarURL,
		&user.Provider,
		&user.ProviderID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NewNotFoundError("user not found")
		}
		return nil, pkgErrors.NewInternalError("failed to find user", err)
	}
	return user, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, provider, provider_id, is_active, created_at, updated_at, version
		FROM users
		WHERE email = $1
	`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.AvatarURL,
		&user.Provider,
		&user.ProviderID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NewNotFoundError("user not found")
		}
		return nil, pkgErrors.NewInternalError("failed to find user", err)
	}
	return user, nil
}

// FindByProvider retrieves a user by OAuth provider and provider ID
func (r *UserRepository) FindByProvider(ctx context.Context, provider, providerID string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, provider, provider_id, is_active, created_at, updated_at, version
		FROM users
		WHERE provider = $1 AND provider_id = $2
	`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.AvatarURL,
		&user.Provider,
		&user.ProviderID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NewNotFoundError("user not found")
		}
		return nil, pkgErrors.NewInternalError("failed to find user", err)
	}
	return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, full_name = $4, avatar_url = $5,
		    provider = $6, provider_id = $7, is_active = $8, updated_at = $9, version = $10
		WHERE id = $1 AND version = $11
	`
	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.AvatarURL,
		user.Provider,
		user.ProviderID,
		user.IsActive,
		user.UpdatedAt,
		user.Version,
		user.Version-1, // optimistic locking
	)
	if err != nil {
		return pkgErrors.NewInternalError("failed to update user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkgErrors.NewInternalError("failed to check rows affected", err)
	}
	if rowsAffected == 0 {
		return pkgErrors.NewConflictError("user was modified by another process")
	}

	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return pkgErrors.NewInternalError("failed to delete user", err)
	}
	return nil
}

// Exists checks if a user with the given email exists
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, pkgErrors.NewInternalError("failed to check user existence", err)
	}
	return exists, nil
}
