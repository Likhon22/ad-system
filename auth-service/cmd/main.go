// @title           Auth Service API
// @version         1.0
// @description     Authentication service for the Ad System
// @host            localhost:5000
// @BasePath        /api/v1

package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/likhon22/ad-system/auth-service/config"
	httpserver "github.com/likhon22/ad-system/auth-service/internal/adapter/inbound/http"
	"github.com/likhon22/ad-system/auth-service/internal/adapter/outbound/postgres"
	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(utils.NewContextHandler(jsonHandler)))

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error("invalid config", "err", err)
		os.Exit(1)
	}
	pool, err := postgres.NewPool(cfg.DB)
	if err != nil {
		slog.Error("cannot connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	app := httpserver.NewApp(pool, cfg)

	go func() {
		if err := app.Run(); err != nil {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "err", err)
		os.Exit(1)
	}
	slog.Info("server exited cleanly")
}
