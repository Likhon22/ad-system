package inbound

import (
	"context"

	"github.com/likhon22/ad-system/auth-service/internal/domain"
)

type RegisterInput struct {
	Email    string
	Password string
	Role     domain.Role
}

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*domain.User, error)
}
