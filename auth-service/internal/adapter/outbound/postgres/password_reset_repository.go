package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/likhon22/ad-system/auth-service/internal/domain"
	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"
)

type passwordResetRepository struct {
	db *pgxpool.Pool
}

func NewPasswordResetRepository(db *pgxpool.Pool) outbound.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) CreateToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, queryCreateResetToken, userID, tokenHash, expiresAt, time.Now())
	return err
}
func (r *passwordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.db.QueryRow(ctx, queryFindByTokenHash, tokenHash).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, domain.ErrUserNotFound
	}
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
func (r *passwordResetRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx, queryDeleteByTokenHash, tokenHash)
	return err
}
func (r *passwordResetRepository) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, queryDeleteAllForUser, userID)
	return err
}
