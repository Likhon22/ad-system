package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/likhon22/ad-system/auth-service/internal/domain"
	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) outbound.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {

	var providerID *string
	if user.ProviderID != "" {
		providerID = &user.ProviderID
	}

	var avatarURL *string
	if user.AvatarURL != "" {
		avatarURL = &user.AvatarURL
	}

	var passwordHash *string
	if user.PasswordHash != "" {
		passwordHash = &user.PasswordHash
	}

	_, err := r.pool.Exec(ctx, queryCreateUser,
		user.ID,
		user.Email,
		user.Name,
		passwordHash,
		user.Provider,
		providerID,
		avatarURL,
		user.Role,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("userRepository.Create: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	var passwordHash *string
	var providerID *string
	var avatarURL *string

	err := r.pool.QueryRow(ctx, queryFindByEmail, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&passwordHash,
		&user.Provider,
		&providerID,
		&avatarURL,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByEmail: %w", err)
	}

	// Safely dereference only if the pointer is non-nil
	if passwordHash != nil {
		user.PasswordHash = *passwordHash
	}
	if providerID != nil {
		user.ProviderID = *providerID
	}
	if avatarURL != nil {
		user.AvatarURL = *avatarURL
	}

	return &user, nil
}

func (r *userRepository) FindByProviderID(ctx context.Context, provider, providerID string) (*domain.User, error) {

	var user domain.User
	var dbProviderID *string
	var avatarURL *string
	var passwordHash *string

	err := r.pool.QueryRow(ctx, queryFindByProviderID, provider, providerID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&passwordHash,
		&user.Provider,
		&dbProviderID,
		&avatarURL,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByProviderID: %w", err)
	}

	if passwordHash != nil {
		user.PasswordHash = *passwordHash
	}
	if dbProviderID != nil {
		user.ProviderID = *dbProviderID
	}
	if avatarURL != nil {
		user.AvatarURL = *avatarURL
	}

	return &user, nil
}
