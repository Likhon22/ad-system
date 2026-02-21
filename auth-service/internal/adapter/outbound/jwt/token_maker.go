package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/likhon22/ad-system/auth-service/internal/domain"
	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"
)

type jwtClaims struct {
	UserID uuid.UUID   `json:"user_id"`
	Email  string      `json:"email"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type JWTMaker struct {
	accessSecret    string
	refreshSecret   string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewJWTMaker(
	accessSecret string,
	refreshSecret string,
	accessDuration time.Duration,
	refreshDuration time.Duration,
) outbound.TokenMaker {
	return &JWTMaker{
		accessSecret:    accessSecret,
		refreshSecret:   refreshSecret,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

func (m *JWTMaker) createToken(userID uuid.UUID, role domain.Role, email string, secret string, duration time.Duration) (string, error) {
	claims := &jwtClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err

	}
	return signedToken, nil
}

func (m *JWTMaker) verifyToken(tokenStr string, secret string) (*outbound.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return &outbound.Claims{
		UserId: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

func (m *JWTMaker) CreateAccessToken(userID uuid.UUID, role domain.Role, email string) (string, error) {
	return m.createToken(userID, role, email, m.accessSecret, m.accessDuration)
}

func (m *JWTMaker) CreateRefreshToken(userID uuid.UUID, role domain.Role, email string) (string, error) {
	return m.createToken(userID, role, email, m.refreshSecret, m.refreshDuration)
}

func (m *JWTMaker) VerifyAccessToken(tokenStr string) (*outbound.Claims, error) {
	return m.verifyToken(tokenStr, m.accessSecret)
}

func (m *JWTMaker) VerifyRefreshToken(tokenStr string) (*outbound.Claims, error) {
	return m.verifyToken(tokenStr, m.refreshSecret)
}
