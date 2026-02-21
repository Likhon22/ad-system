package utils

import (
	"net/http"
	"time"
)

const refreshTokenCookieName = "refresh_token"

func SetRefreshTokenCookie(w http.ResponseWriter, token string, duration time.Duration, isProduction bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(duration.Seconds()),
		HttpOnly: isProduction,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}
func ClearRefreshTokenCookie(w http.ResponseWriter, isProduction bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: isProduction,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}
func GetRefreshTokenCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
