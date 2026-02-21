package handler

import (
	"net/http"
	"time"

	"github.com/likhon22/ad-system/auth-service/internal/domain"
	"github.com/likhon22/ad-system/auth-service/internal/port/inbound"
	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

type AuthHandler struct {
	authService     inbound.AuthService
	validate        *utils.Validator
	refreshDuration time.Duration
	isProduction    bool
}

func NewHandler(service inbound.AuthService, validator *utils.Validator, refreshDuration time.Duration, isProduction bool) *AuthHandler {
	return &AuthHandler{
		authService:     service,
		validate:        validator,
		refreshDuration: refreshDuration,
		isProduction:    isProduction,
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
// @Summary      Login a user
// @Description  Login with email and password. Returns access token in body, refresh token in httpOnly cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login payload"
// @Success      200 {object} utils.Response[inbound.LoginResponse]
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
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
	result, err := h.authService.Login(r.Context(), inbound.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}
	utils.SetRefreshTokenCookie(w, result.RefreshToken, h.refreshDuration, h.isProduction)
	if err := utils.WriteJSON(w, "user logged successfully", http.StatusOK, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

}
