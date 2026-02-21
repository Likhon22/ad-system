package middleware

import (
	"net/http"
	"strings"

	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"
	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

func RequiredAuth(tokenMaker outbound.TokenMaker) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return

			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				utils.WriteError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}
			claims, err := tokenMaker.VerifyAccessToken(parts[1])
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}
			ctx := utils.SetClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
