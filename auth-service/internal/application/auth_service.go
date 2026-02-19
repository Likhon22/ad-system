package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/likhon22/ad-system/auth-service/internal/domain"
	"github.com/likhon22/ad-system/auth-service/internal/port/inbound"
	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"
	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

type authService struct {
	userRepo outbound.UserRepository // interface — knows nothing about postgres
}

// Returns inbound.AuthService interface — handler never sees the concrete struct
func NewAuthService(userRepo outbound.UserRepository) inbound.AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, input inbound.RegisterInput) (*domain.User, error) {

	if input.Role != domain.RoleAdvertiser && input.Role != domain.RolePublisher {
		return nil, domain.ErrInvalidRole
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: hash,
		Role:         input.Role,
		Status:       domain.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.PasswordHash = "" // scrub before returning — never expose hash
	return user, nil
}
