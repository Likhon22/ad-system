package utils

import (
	"encoding/json"
	"net/http"
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

	json.NewEncoder(w).Encode(resp)

}
