package utils

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
	Success bool   `json:"success"`
}

func WriteJSON[T any](w http.ResponseWriter, message string, statusCode int, data *T) error {
	w.Header().Set("Content-Type", "application/json")
	res := &Response[T]{
		Message: message,
		Data:    data,
		Success: true,
	}
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(res)
}

func WriteMessage(w http.ResponseWriter, message string, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	res := map[string]any{
		"message": message,
		"success": true,
	}
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(res)
}
