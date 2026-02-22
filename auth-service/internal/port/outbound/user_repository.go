package outbound

import (
	"context"

	"github.com/likhon22/ad-system/auth-service/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByProviderID(ctx context.Context, provider, providerID string) (*domain.User, error)
}
