package httpserver

import (
	"net/http"

	"github.com/likhon22/ad-system/auth-service/internal/adapter/inbound/http/handler"
	"github.com/likhon22/ad-system/auth-service/internal/adapter/inbound/http/middleware"
	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"

	_ "github.com/likhon22/ad-system/auth-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func setupRouter(authHandler *handler.AuthHandler, healthHandler *handler.HealthHandler, tokenMaker outbound.TokenMaker) http.Handler {
	mux := http.NewServeMux()
	v1 := http.NewServeMux()
	v1.HandleFunc("POST /auth/register", authHandler.Register)
	v1.HandleFunc("POST /auth/login", authHandler.Login)
	v1.HandleFunc("POST /auth/refresh-token", authHandler.RefreshToken)
	//protected route
	requireAuth := middleware.RequiredAuth(tokenMaker)
	v1.Handle("GET /auth/me", requireAuth(http.HandlerFunc(authHandler.GetMe)))
	mux.HandleFunc("GET /healthz", healthHandler.Liveness)
	mux.HandleFunc("GET /readyz", healthHandler.Readiness)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))
	return mux
}
