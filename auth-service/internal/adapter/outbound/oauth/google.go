package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleOAuthProvider struct {
	config *oauth2.Config
}

func NewGoogleOAuthProvider(clientID, clientSecret, redirectURL string) outbound.OAuthProvider {

	return &googleOAuthProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"openid",
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}

}
func (g *googleOAuthProvider) BuildAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}
func (g *googleOAuthProvider) ExchangeCode(ctx context.Context, code string) (*outbound.GoogleUserInfo, error) {

	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info from Google: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google userinfo returned status %d", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var userInfo outbound.GoogleUserInfo

	if err := decoder.Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse Google userinfo: %w", err)
	}

	if userInfo.Sub == "" {
		return nil, fmt.Errorf("Google userinfo missing 'sub' — invalid response")
	}

	return &userInfo, nil
}
