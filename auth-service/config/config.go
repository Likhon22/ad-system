package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Version     string
	ServiceName string
	Addr        string
	DBCnf       *DBConfig
}

var (
	config *Config
	once   sync.Once
)

func loadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env file found")
	}
	version := os.Getenv("VERSION")
	serviceName := os.Getenv("SERVICE_NAME")
	addr := os.Getenv("ADDR")
	config = &Config{
		Version:     version,
		ServiceName: serviceName,
		Addr:        addr,
		DBCnf:       LoadDBConfig(),
	}
	validateMainConfig(config)
}

func GetConfig() *Config {
	once.Do(loadConfig)
	return config

}

func validateMainConfig(cnf *Config) {
	if cnf.Version == "" {
		log.Fatal("Missing version")
	}
	if cnf.ServiceName == "" {
		log.Fatal("Missing service name")
	}
	if cnf.Addr == "" {
		log.Fatal("MIssing address/port")
	}
}
