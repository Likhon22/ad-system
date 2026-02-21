package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/likhon22/ad-system/auth-service/internal/domain"
)

type RegisterInput struct {
	Email    string
	Password string
	Role     domain.Role
}

type LoginInput struct {
	Email    string
	Password string
}
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *UserSummary `json:"user"`
}

type UserSummary struct {
	ID    uuid.UUID   `json:"id"`
	Email string      `json:"email"`
	Role  domain.Role `json:"role"`
}
type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*domain.User, error)
	Login(ctx context.Context, input LoginInput) (LoginResponse, error)
}
