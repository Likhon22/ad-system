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

const (
	oauthStateCookie = "oauth_state"
	oauthFlowCookie  = "oauth_flow"
	oauthRoleCookie  = "oauth_role"
)

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role"     validate:"required,oneof=advertiser publisher"`
	Name     string `json:"name" validate:"required"`
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
		Name:     req.Name,
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

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := utils.GetRefreshTokenCookie(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "missing or invalid token")
		return
	}
	result, err := h.authService.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}
	if err := utils.WriteJSON(w, "token refreshed successfully", http.StatusOK, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {

	claims, ok := utils.GetClaims(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := h.authService.GetMe(r.Context(), claims.Email)
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}
	if err := utils.WriteJSON(w, "user fetched successfully", http.StatusOK, user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

func (h *AuthHandler) GoogleInitiate(w http.ResponseWriter, r *http.Request) {
	flow := r.URL.Query().Get("flow")
	if flow != "login" && flow != "register" {
		utils.WriteError(w, http.StatusBadRequest, "flow must be 'login' or 'register'")
		return
	}

	role := r.URL.Query().Get("role")
	if flow == "register" && role != "advertiser" && role != "publisher" {
		utils.WriteError(w, http.StatusBadRequest, "role must be 'advertiser' or 'publisher'")
		return
	}
	authURL, state, err := h.authService.GoogleInitiate(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to initiate Google login")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    state,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     oauthFlowCookie,
		Value:    flow,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	if flow == "register" {
		http.SetCookie(w, &http.Cookie{
			Name:     oauthRoleCookie,
			Value:    role,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   h.isProduction,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})
	}
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Read cookies
	stateCookie, err := r.Cookie(oauthStateCookie)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "missing state cookie")
		return
	}
	flowCookie, err := r.Cookie(oauthFlowCookie)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "missing flow cookie")
		return
	}

	roleCookie := ""
	if flowCookie.Value == "register" {
		rc, err := r.Cookie(oauthRoleCookie)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "missing role cookie")
			return
		}
		roleCookie = rc.Value
	}

	// Delete all cookies immediately — single use
	for _, name := range []string{oauthStateCookie, oauthFlowCookie, oauthRoleCookie} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   h.isProduction,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})
	}

	result, err := h.authService.GoogleCallback(
		r.Context(),
		code,
		state,
		stateCookie.Value,
		flowCookie.Value,
		domain.Role(roleCookie),
	)
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	utils.SetRefreshTokenCookie(w, result.RefreshToken, h.refreshDuration, h.isProduction)
	if err := utils.WriteJSON(w, "logged in successfully", http.StatusOK, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetClaims(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req PasswordChangeRequest
	utils.ReadJson(w, r, &req)
	err := h.authService.ChangePassword(r.Context(), claims.Email, req.CurrentPassword, req.NewPassword)
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}
	if err := utils.WriteMessage(w, "password changed successfully", http.StatusOK); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}
