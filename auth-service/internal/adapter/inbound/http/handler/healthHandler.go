package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	pool *pgxpool.Pool
}

func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})

}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	w.Header().Set("Content-Type", "application/json")

	if err := h.pool.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable) // 503
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "unavailable",
			"reason": "database unreachable",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
