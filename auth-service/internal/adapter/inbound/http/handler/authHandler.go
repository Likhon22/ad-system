package handler

import (
	"net/http"

	"github.com/likhon22/ad-system/auth-service/internal/domain"
	"github.com/likhon22/ad-system/auth-service/internal/port/inbound"
	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

type AuthHandler struct {
	authService inbound.AuthService
	validate    *utils.Validator
}

func NewHandler(service inbound.AuthService, validator *utils.Validator) *AuthHandler {
	return &AuthHandler{
		authService: service,
		validate:    validator,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role"     validate:"required,oneof=advertiser publisher"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account with email, password and role
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Register payload"
// @Success      201 {object} utils.Response[domain.User]
// @Failure      400 {object} utils.ErrorResponse
// @Failure      409 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validate.ValidateStruct(req); err != nil {
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
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Login godoc
// @Summary      Login a new user
// @Description  Login a user account with email, password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login payload"
// @Success      201 {object} utils.Response[domain.User]
// @Failure      400 {object} utils.ErrorResponse
// @Failure      409 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validate.ValidateStruct(req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.authService.Login(r.Context(), inbound.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}
	if err := utils.WriteJSON(w, "user logged successfully", http.StatusOK, user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

}
