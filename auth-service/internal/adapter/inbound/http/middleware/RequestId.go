package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

const RequestIDHeader = "X-Request-ID"

func (mw *Middleware) RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		w.Header().Set(RequestIDHeader, requestID)
		ctx := utils.ContextWithRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
