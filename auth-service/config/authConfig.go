package config

import (
	"fmt"
	"os"
	"time"
)

type AuthConfig struct {
	JWTAccessSecret      string
	JWTRefreshSecret     string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	GoogleClientID       string
	GoogleClientSecret   string
	GoogleRedirectURL    string
	FrontendURL          string
}

func loadAuthConfig() (*AuthConfig, error) {
	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")

	accessDurStr := getEnv("ACCESS_TOKEN_DURATION", "15m")
	refreshDurStr := getEnv("REFRESH_TOKEN_DURATION", "168h")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	if accessSecret == "" {
		return nil, fmt.Errorf("ACCESS_TOKEN_SECRET is required")
	}
	if refreshSecret == "" {
		return nil, fmt.Errorf("REFRESH_TOKEN_SECRET is required")
	}

	accessDur, err := time.ParseDuration(accessDurStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_DURATION: %w", err)
	}

	refreshDur, err := time.ParseDuration(refreshDurStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_DURATION: %w", err)
	}
	if googleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}
	if googleClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is required")
	}
	if googleRedirectURL == "" {
		return nil, fmt.Errorf("GOOGLE_REDIRECT_URL is required")
	}
	if frontendURL == "" {
		return nil, fmt.Errorf("fronted_url is required")
	}

	return &AuthConfig{
		JWTAccessSecret:      accessSecret,
		JWTRefreshSecret:     refreshSecret,
		AccessTokenDuration:  accessDur,
		RefreshTokenDuration: refreshDur,
		GoogleClientID:       googleClientID,
		GoogleClientSecret:   googleClientSecret,
		GoogleRedirectURL:    googleRedirectURL,
		FrontendURL:          frontendURL,
	}, nil
}
