package main

import (
	"log"

	"github.com/likhon22/ad-system/auth-service/config"
	httpserver "github.com/likhon22/ad-system/auth-service/internal/adapter/inbound/http"
	"github.com/likhon22/ad-system/auth-service/internal/adapter/outbound/postgres"
)

func main() {
	cfg := config.GetConfig()
  
	pool, err := postgres.NewPool(cfg.DB)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer pool.Close()

	app := httpserver.NewApp(pool, cfg)
	if err := app.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
