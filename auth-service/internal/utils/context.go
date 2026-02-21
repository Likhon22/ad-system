package utils

import (
	"context"

	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	claimsKey    contextKey = "claims"
)

func ContextWithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

func SetClaims(ctx context.Context, claims *outbound.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)

}
func GetClaims(ctx context.Context) (*outbound.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*outbound.Claims)
	return claims, ok
}
