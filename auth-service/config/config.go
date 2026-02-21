package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Version      string
	ServiceName  string
	Addr         string
	IsProduction bool
	DB           *DBConfig
	Auth         *AuthConfig
}

var (
	instance *Config
	initErr  error
	once     sync.Once
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		// .env is for local dev only — ignore error in production
		_ = godotenv.Load()
		instance = &Config{
			Version:      os.Getenv("VERSION"),
			ServiceName:  os.Getenv("SERVICE_NAME"),
			Addr:         os.Getenv("ADDR"),
			IsProduction: os.Getenv("ENV") == "production",
			DB:           loadDBConfig(),
		}
		authCfg, err := loadAuthConfig()
		if err != nil {
			initErr = err
			return
		}
		instance.Auth = authCfg
		initErr = instance.validate()

	})
	return instance, initErr
}

func (c *Config) validate() error {
	if c.Version == "" {
		return fmt.Errorf("VERSION is required")
	}
	if c.ServiceName == "" {
		return fmt.Errorf("SERVICE_NAME is required")
	}
	if c.Addr == "" {
		return fmt.Errorf("ADDR is required")
	}
	return nil
}
