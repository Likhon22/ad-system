package outbound

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PasswordResetRepository interface {
	CreateToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (uuid.UUID, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteAllForUser(ctx context.Context, userID uuid.UUID) error
}
