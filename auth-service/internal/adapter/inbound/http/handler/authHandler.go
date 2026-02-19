package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/likhon22/ad-system/auth-service/internal/domain"
	"github.com/likhon22/ad-system/auth-service/internal/port/inbound"
	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

type AuthHandler struct {
	authService inbound.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService inbound.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

type registerRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role"     validate:"required,oneof=advertiser publisher"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req registerRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.authService.Register(r.Context(), inbound.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Role:     domain.Role(req.Role),
	})
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}
	if err := utils.WriteJSON(w, "user created successfully", http.StatusCreated, user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error creating student")
		return
	}

}
