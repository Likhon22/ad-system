// @title           Auth Service API
// @version         1.0
// @description     Authentication service for the Ad System
// @host            localhost:5000
// @BasePath        /api/v1

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/likhon22/ad-system/auth-service/config"
	httpserver "github.com/likhon22/ad-system/auth-service/internal/adapter/inbound/http"
	"github.com/likhon22/ad-system/auth-service/internal/adapter/outbound/postgres"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}
	pool, err := postgres.NewPool(cfg.DB)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer pool.Close()

	app := httpserver.NewApp(pool, cfg)
	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server exited cleanly")
}
