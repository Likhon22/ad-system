package outbound

import (
	"github.com/google/uuid"
	"github.com/likhon22/ad-system/auth-service/internal/domain"
)

type Claims struct {
	UserID uuid.UUID
	Email  string
	Role   domain.Role
}

type TokenMaker interface {
	CreateAccessToken(userID uuid.UUID, role domain.Role, email string) (string, error)
	CreateRefreshToken(userID uuid.UUID, role domain.Role, email string) (string, error)
	VerifyAccessToken(token string) (*Claims, error)
	VerifyRefreshToken(token string) (*Claims, error)
}
