package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/likhon22/ad-system/auth-service/internal/domain"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, statusCode int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := ErrorResponse{
		Success: false,
		Message: err,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}

}

func ErrorHandler(w http.ResponseWriter, err error) {
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			WriteError(w, http.StatusConflict, "email already registered")
		case errors.Is(err, domain.ErrInvalidRole):
			WriteError(w, http.StatusBadRequest, "role must be advertiser or publisher")
		case errors.Is(err, domain.ErrInvalidCredentials):
			WriteError(w, http.StatusUnauthorized, "invalid email or password")
		case errors.Is(err, domain.ErrUserNotFound):
			WriteError(w, http.StatusUnauthorized, "invalid email or password")
		case errors.Is(err, domain.ErrUserSuspended):
			WriteError(w, http.StatusForbidden, "account is suspended")
		default:
			WriteError(w, http.StatusInternalServerError, "internal server error")
		}
	}
}
