package utils

import (
	"net/http"
	"time"
)

const refreshTokenCookieName = "refresh_token"

func SetRefreshTokenCookie(w http.ResponseWriter, token string, duration time.Duration) {

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(duration.Seconds()),
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

}
func ClearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
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
