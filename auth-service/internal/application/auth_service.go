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
	userRepo   outbound.UserRepository
	tokenMaker outbound.TokenMaker
}

// Returns inbound.AuthService interface — handler never sees the concrete struct
func NewAuthService(userRepo outbound.UserRepository, tokenMaker outbound.TokenMaker) inbound.AuthService {
	return &authService{userRepo: userRepo, tokenMaker: tokenMaker}
}

func (s *authService) Register(ctx context.Context, input inbound.RegisterInput) (*inbound.UserSummary, error) {

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
		Name:         input.Name,
		PasswordHash: hash,
		Provider:     "local",
		Role:         input.Role,
		Status:       domain.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	newUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &inbound.UserSummary{
		ID:     newUser.ID,
		Email:  newUser.Email,
		Role:   newUser.Role,
		Name:   newUser.Name,
		Status: newUser.Status,
	}, nil
}

func (s *authService) Login(ctx context.Context, input inbound.LoginInput) (inbound.LoginResponse, error) {

	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return inbound.LoginResponse{}, err

	}
	if user.Provider != "local" {
		return inbound.LoginResponse{}, domain.ErrOAuthUser
	}
	if user.Status != domain.StatusActive {
		return inbound.LoginResponse{}, domain.ErrUserSuspended
	}

	matched, err := utils.ComparePassword(input.Password, user.PasswordHash)
	if err != nil {
		return inbound.LoginResponse{}, err

	}
	if !matched {
		return inbound.LoginResponse{}, domain.ErrInvalidCredentials

	}

	accessToken, err := s.tokenMaker.CreateAccessToken(user.ID, user.Role, user.Email)
	if err != nil {
		return inbound.LoginResponse{}, err
	}

	refreshToken, err := s.tokenMaker.CreateRefreshToken(user.ID, user.Role, user.Email)
	if err != nil {
		return inbound.LoginResponse{}, err
	}

	return inbound.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &inbound.UserSummary{
			ID:     user.ID,
			Email:  user.Email,
			Role:   user.Role,
			Name:   user.Name,
			Status: user.Status,
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, token string) (inbound.LoginResponse, error) {

	verified, err := s.tokenMaker.VerifyRefreshToken(token)
	if err != nil {
		return inbound.LoginResponse{}, err
	}
	newToken, err := s.tokenMaker.CreateAccessToken(verified.UserID, verified.Role, verified.Email)
	if err != nil {
		return inbound.LoginResponse{}, err
	}
	return inbound.LoginResponse{
		User: &inbound.UserSummary{
			ID:    verified.UserID,
			Email: verified.Email,
			Role:  verified.Role,
		},
		AccessToken: newToken,
	}, nil
}
func (s *authService) GetMe(ctx context.Context, email string) (*inbound.UserSummary, error) {

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &inbound.UserSummary{

		ID:     user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Name:   user.Name,
		Status: user.Status,
	}, nil
}
