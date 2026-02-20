package middleware

import (
	"log/slog"
)

type Middleware struct {
	logger *slog.Logger
}

func NewMiddleware() *Middleware {
	return &Middleware{
		logger: slog.Default(),
	}

}
