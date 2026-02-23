package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/likhon22/ad-system/auth-service/internal/domain"
	"github.com/likhon22/ad-system/auth-service/internal/port/inbound"
	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"
	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

type authService struct {
	userRepo      outbound.UserRepository
	tokenMaker    outbound.TokenMaker
	oauthProvider outbound.OAuthProvider
	resetRepo     outbound.PasswordResetRepository
	emailSender   outbound.EmailSender
	frontendURL   string
}

func NewAuthService(
	userRepo outbound.UserRepository,
	tokenMaker outbound.TokenMaker,
	oauthProvider outbound.OAuthProvider,
	resetRepo outbound.PasswordResetRepository,
	emailSender outbound.EmailSender,
	frontendURL string) inbound.AuthService {
	return &authService{
		userRepo:      userRepo,
		tokenMaker:    tokenMaker,
		oauthProvider: oauthProvider,
		resetRepo:     resetRepo,
		emailSender:   emailSender,
		frontendURL:   frontendURL,
	}
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

func (s *authService) GoogleInitiate(ctx context.Context) (string, string, error) {
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate state: %w", err)
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)
	authURL := s.oauthProvider.BuildAuthURL(state)
	return authURL, state, nil
}

func (s *authService) GoogleCallback(ctx context.Context, code, state, storedState, flow string, role domain.Role) (inbound.LoginResponse, error) {
	// CSRF check — state must match what we stored in cookie
	if code == "" || state == "" || storedState == "" || state != storedState {
		return inbound.LoginResponse{}, domain.ErrInvalidCredentials
	}

	googleUser, err := s.oauthProvider.ExchangeCode(ctx, code)
	if err != nil {
		return inbound.LoginResponse{}, fmt.Errorf("OAuth exchange failed: %w", err)
	}

	user, err := s.userRepo.FindByProviderID(ctx, "google", googleUser.Sub)
	userExists := err == nil
	userNotFound := errors.Is(err, domain.ErrUserNotFound)

	if err != nil && !userNotFound {
		return inbound.LoginResponse{}, err
	}

	if flow == "login" {
		if userNotFound {
			existing, emailErr := s.userRepo.FindByEmail(ctx, googleUser.Email)
			if emailErr == nil && existing.Provider == "local" {
				return inbound.LoginResponse{}, domain.ErrEmailConflict
			}
			return inbound.LoginResponse{}, domain.ErrUserNotFound
		}
	} else {

		if userExists {
			return inbound.LoginResponse{}, domain.ErrEmailAlreadyExists
		}

		// Validate role
		if role != domain.RoleAdvertiser && role != domain.RolePublisher {
			return inbound.LoginResponse{}, domain.ErrInvalidRole
		}

		existing, emailErr := s.userRepo.FindByEmail(ctx, googleUser.Email)
		if emailErr == nil && existing.Provider == "local" {
			return inbound.LoginResponse{}, domain.ErrEmailConflict
		}

		// Create new user
		now := time.Now()
		newUser := &domain.User{
			ID:         uuid.New(),
			Email:      googleUser.Email,
			Name:       googleUser.Name,
			AvatarURL:  googleUser.Picture,
			Provider:   "google",
			ProviderID: googleUser.Sub,
			Role:       role,
			Status:     domain.StatusActive,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		user, err = s.userRepo.Create(ctx, newUser)
		if err != nil {
			return inbound.LoginResponse{}, err
		}
	}

	if user.Status != domain.StatusActive {
		return inbound.LoginResponse{}, domain.ErrUserSuspended
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

func (s *authService) ChangePassword(ctx context.Context, email, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user.Provider != "local" {
		return domain.ErrOAuthUser
	}
	matched, err := utils.ComparePassword(currentPassword, user.PasswordHash)
	if err != nil {
		return err
	}
	if !matched {
		return domain.ErrInvalidCredentials
	}
	samePassword, err := utils.ComparePassword(newPassword, user.PasswordHash)
	if err != nil {
		return err
	}
	if samePassword {
		return domain.ErrSamePassword
	}
	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePassword(ctx, user.ID, newHash)
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)

	if err != nil {
		return nil
	}
	if user.Provider != "local" {
		return nil
	}
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	plainToken := base64.URLEncoding.EncodeToString(tokenBytes)
	tokenHash := utils.HashToken(plainToken)
	_ = s.resetRepo.DeleteAllForUser(ctx, user.ID)

	expiresAt := time.Now().Add(15 * time.Minute)

	if err := s.resetRepo.CreateToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return err
	}

	resetLink := s.frontendURL + "/reset-password?token=" + plainToken

	return s.emailSender.SendPasswordResetEmail(ctx, user.Email, resetLink)

}
func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	tokenHash := utils.HashToken(token)

	userID, err := s.resetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return domain.ErrInvalidOrExpiredToken
	}

	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, newHash); err != nil {
		return err
	}

	return s.resetRepo.DeleteAllForUser(ctx, userID)
}
