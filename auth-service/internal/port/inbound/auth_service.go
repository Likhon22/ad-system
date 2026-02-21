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
	Name     string
}

type LoginInput struct {
	Email    string
	Password string
}
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"-"`
	User         *UserSummary `json:"user"`
}

type UserSummary struct {
	ID     uuid.UUID     `json:"id"`
	Email  string        `json:"email"`
	Role   domain.Role   `json:"role"`
	Name   string        `json:"name"`
	Status domain.Status `json:"status"`
}
type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*UserSummary, error)
	Login(ctx context.Context, input LoginInput) (LoginResponse, error)
	RefreshToken(ctx context.Context, token string) (LoginResponse, error)
	GetMe(ctx context.Context, email string) (*UserSummary, error)
}
