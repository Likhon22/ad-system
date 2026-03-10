package config

import (
	"fmt"
	"os"
	"strconv"
)

type DBConfig struct {
	URL          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func loadDBConfig() *DBConfig {
	url := os.Getenv("DB_URL")
	if url == "" {
		panic("DB_URL is required")
	}

	return &DBConfig{
		URL:          url,
		MaxOpenConns: envInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: envInt("DB_MAX_IDLE_CONNS", 25),
		MaxIdleTime:  getEnv("DB_MAX_IDLE_TIME", "15m"),
	}
}

func envInt(key string, defaultVal int) int {
	s := os.Getenv(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("%s must be a valid integer, got: %s", key, s))
	}
	return v
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
