package httpserver

import (
	"net/http"

	"github.com/likhon22/ad-system/auth-service/internal/adapter/inbound/http/handler"
)

func setupRouter(authHandler *handler.AuthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	return mux
}
