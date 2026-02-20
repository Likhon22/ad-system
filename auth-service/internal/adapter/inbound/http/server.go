package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/likhon22/ad-system/auth-service/config"
	"github.com/likhon22/ad-system/auth-service/internal/adapter/inbound/http/handler"
	postgres "github.com/likhon22/ad-system/auth-service/internal/adapter/outbound/postgres"
	"github.com/likhon22/ad-system/auth-service/internal/application"
	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

type App struct {
	server *http.Server
}

func NewApp(pool *pgxpool.Pool, cfg *config.Config) *App {
	validate := utils.NewValidator()
	// outbound adapters (real implementations)
	userRepo := postgres.NewUserRepository(pool)

	// application (business logic — gets interfaces, not real implementations)
	authSvc := application.NewAuthService(userRepo)

	// inbound adapters (HTTP handlers — gets interfaces)
	authHandler := handler.NewHandler(authSvc, validate)

	// wire routes
	mux := setupRouter(authHandler)

	return &App{
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}
}

func (a *App) Run() error {
	fmt.Printf("auth-service starting on %s\n", a.server.Addr)
	return a.server.ListenAndServe()
}
