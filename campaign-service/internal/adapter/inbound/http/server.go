package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/likhon22/ad-system/campaign-service/config"
)

type App struct {
	server *http.Server
}

func NewApp(cfg *config.Config) *App {
	// mw := middleware.NewMiddleware()
	// validate := utils.NewValidator()

	return &App{
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      http.NewServeMux(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}
}

func (a *App) Run() error {
	fmt.Printf("auth-service starting on %s\n", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
