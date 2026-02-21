package postgres

import (
	"context"
	"errors"

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

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {

	_, err := r.pool.Exec(ctx, queryCreateUser,
		user.ID, user.Email, user.Name, user.PasswordHash, user.Provider,
		user.Role, user.Status, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, queryFindByEmail, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Provider,
		&user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	return user, err
}

// func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
// 	query := `
//         SELECT id, email, password_hash, role, status, created_at, updated_at
//         FROM users WHERE id = $1
//     `
// 	user := &domain.User{}
// 	err := r.pool.QueryRow(ctx, query, id).Scan(
// 		&user.ID, &user.Email, &user.PasswordHash,
// 		&user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt,
// 	)
// 	if errors.Is(err, pgx.ErrNoRows) {
// 		return nil, domain.ErrUserNotFound
// 	}
// 	return user, err
// }
