package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleOAuthConfig holds the Google OAuth 2.0 configuration
type GoogleOAuthConfig struct {
	config       *oauth2.Config
	clientID     string
	clientSecret string
	redirectURI  string
}

// GoogleUserInfo represents the user information returned by Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// NewGoogleOAuthConfig creates a new Google OAuth configuration from environment variables
func NewGoogleOAuthConfig() (*GoogleOAuthConfig, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	if clientID == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID environment variable is not set")
	}

	if clientSecret == "" {
		return nil, errors.New("GOOGLE_CLIENT_SECRET environment variable is not set")
	}

	if redirectURI == "" {
		return nil, errors.New("GOOGLE_REDIRECT_URI environment variable is not set")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthConfig{
		config:       config,
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}, nil
}

// GetAuthURL generates the Google OAuth authorization URL with state and PKCE challenge
func (g *GoogleOAuthConfig) GetAuthURL(state, codeChallenge string) string {
	options := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("prompt", "consent"),
	}

	// Add PKCE challenge if provided
	if codeChallenge != "" {
		options = append(options,
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", "plain"),
		)
	}

	return g.config.AuthCodeURL(state, options...)
}

// ExchangeCode exchanges the authorization code for an access token
func (g *GoogleOAuthConfig) ExchangeCode(ctx context.Context, code string, codeVerifier string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("authorization code cannot be empty")
	}

	options := []oauth2.AuthCodeOption{}

	// Add PKCE verifier if provided
	if codeVerifier != "" {
		options = append(options, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	}

	token, err := g.config.Exchange(ctx, code, options...)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// RefreshToken refreshes an OAuth access token using the refresh token
func (g *GoogleOAuthConfig) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token cannot be empty")
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := g.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// GetUserInfo retrieves user information from Google using the access token
func (g *GoogleOAuthConfig) GetUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	if accessToken == "" {
		return nil, errors.New("access token cannot be empty")
	}

	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	client := g.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get user info from Google")
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// RevokeToken revokes an OAuth token (access or refresh)
func (g *GoogleOAuthConfig) RevokeToken(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("token cannot be empty")
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://oauth2.googleapis.com/revoke",
		strings.NewReader("token="+token),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("failed to revoke token")
	}

	return nil
}

// ValidateToken validates an access token by making a test API call
func (g *GoogleOAuthConfig) ValidateToken(ctx context.Context, accessToken string) (bool, error) {
	if accessToken == "" {
		return false, errors.New("access token cannot be empty")
	}

	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	client := g.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=" + accessToken)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}

// GetConfig returns the underlying oauth2.Config
func (g *GoogleOAuthConfig) GetConfig() *oauth2.Config {
	return g.config
}

// GetClientID returns the Google OAuth client ID
func (g *GoogleOAuthConfig) GetClientID() string {
	return g.clientID
}

// GetRedirectURI returns the configured redirect URI
func (g *GoogleOAuthConfig) GetRedirectURI() string {
	return g.redirectURI
}