package httpserver

import (
	"net/http"

	"github.com/likhon22/ad-system/auth-service/internal/adapter/inbound/http/handler"
)

func setupRouter(authHandler *handler.AuthHandler) http.Handler {
	mux := http.NewServeMux()
	v1 := http.NewServeMux()
	v1.HandleFunc("POST /auth/register", authHandler.Register)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))
	return mux
}
