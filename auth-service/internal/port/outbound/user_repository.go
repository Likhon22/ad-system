package outbound

import (
	"context"

	"github.com/likhon22/ad-system/auth-service/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	// FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}
