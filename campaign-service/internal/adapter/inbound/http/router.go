package httpserver

import (
	"net/http"

	"github.com/likhon22/ad-system/campaign-service/internal/adapter/inbound/http/handler"
)

func setupRouter(healthHandler *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()
	v1 := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthHandler.Liveness)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))
	return mux
}
