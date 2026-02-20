package httpserver

import (
	"net/http"

	"github.com/likhon22/ad-system/auth-service/internal/adapter/inbound/http/handler"
	httpSwagger "github.com/swaggo/http-swagger"
)

func setupRouter(authHandler *handler.AuthHandler, healthHandler *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()
	v1 := http.NewServeMux()
	v1.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("GET /healthz", healthHandler.Liveness)
	mux.HandleFunc("GET /readyz", healthHandler.Readiness)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))
	return mux
}
