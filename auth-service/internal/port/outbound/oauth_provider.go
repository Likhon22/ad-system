package outbound

import "context"

type GoogleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type OAuthProvider interface {
	BuildAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*GoogleUserInfo, error)
}
