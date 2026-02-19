package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Version     string
	ServiceName string
	Addr        string
	DB          *DBConfig
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		// .env is for local dev only — ignore error in production
		godotenv.Load()
		instance = &Config{
			Version:     os.Getenv("VERSION"),
			ServiceName: os.Getenv("SERVICE_NAME"),
			Addr:        os.Getenv("ADDR"),
			DB:          loadDBConfig(),
		}
		if err := instance.validate(); err != nil {
			panic(fmt.Sprintf("invalid config: %v", err))
		}
	})
	return instance
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
