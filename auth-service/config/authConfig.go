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
}

func loadAuthConfig() (*AuthConfig, error) {
	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")

	accessDurStr := getEnv("ACCESS_TOKEN_DURATION", "15m")
	refreshDurStr := getEnv("REFRESH_TOKEN_DURATION", "168h")

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

	return &AuthConfig{
		JWTAccessSecret:      accessSecret,
		JWTRefreshSecret:     refreshSecret,
		AccessTokenDuration:  accessDur,
		RefreshTokenDuration: refreshDur,
	}, nil
}
