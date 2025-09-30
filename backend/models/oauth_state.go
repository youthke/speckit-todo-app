package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"
)

// OAuthState represents temporary state for OAuth flow CSRF protection
type OAuthState struct {
	StateToken    string    `json:"state_token" gorm:"primaryKey;type:varchar(255)"`
	PKCEVerifier  string    `json:"-" gorm:"type:varchar(255);not null"` // Keep secret
	RedirectURI   string    `json:"redirect_uri" gorm:"type:varchar(1000);not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime;index"`
	ExpiresAt     time.Time `json:"expires_at" gorm:"not null;index"`
}

// TableName specifies the table name for the OAuthState model
func (OAuthState) TableName() string {
	return "oauth_states"
}

// BeforeCreate hook to validate OAuth state before creation
func (s *OAuthState) BeforeCreate(tx *gorm.DB) error {
	return s.Validate()
}

// BeforeUpdate hook to validate OAuth state before update
func (s *OAuthState) BeforeUpdate(tx *gorm.DB) error {
	return s.Validate()
}

// Validate performs validation on the OAuthState model
func (s *OAuthState) Validate() error {
	if len(s.StateToken) < 32 {
		return errors.New("state_token must be at least 32 characters")
	}

	if s.PKCEVerifier == "" {
		return errors.New("pkce_verifier cannot be empty")
	}

	// Validate redirect URI format
	if !isValidURL(s.RedirectURI) {
		return errors.New("redirect_uri must be a valid URL")
	}

	// Validate redirect URI is whitelisted
	if !ValidateRedirectURI(s.RedirectURI) {
		return errors.New("redirect_uri not in whitelist")
	}

	// State cannot be expired
	if s.ExpiresAt.Before(time.Now()) || s.ExpiresAt.Equal(time.Now()) {
		return errors.New("state cannot be expired")
	}

	// State cannot be valid for more than 5 minutes
	maxExpiry := time.Now().Add(6 * time.Minute) // Small buffer for clock skew
	if s.ExpiresAt.After(maxExpiry) {
		return errors.New("expires_at cannot exceed 5 minutes")
	}

	return nil
}

// IsExpired returns true if the OAuth state has expired
func (s *OAuthState) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now()) || s.ExpiresAt.Equal(time.Now())
}

// GenerateOAuthState creates a new OAuth state with PKCE
func GenerateOAuthState(redirectURI string) (*OAuthState, error) {
	if !ValidateRedirectURI(redirectURI) {
		return nil, errors.New("invalid redirect URI")
	}

	stateToken := generateSecureRandomString(40) // 40 characters for extra security
	pkceVerifier := GeneratePKCEVerifier()

	state := &OAuthState{
		StateToken:   stateToken,
		PKCEVerifier: pkceVerifier,
		RedirectURI:  redirectURI,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}

	return state, state.Validate()
}

// ValidateRedirectURI validates that a redirect URI is allowed
func ValidateRedirectURI(uri string) bool {
	// Allowed redirect URIs for the application
	allowedURIs := []string{
		"http://localhost:3000/",
		"http://localhost:3000/dashboard",
		"http://localhost:3000/auth/callback",
	}

	// In production, this should be configurable and include production domains
	for _, allowed := range allowedURIs {
		if strings.HasPrefix(uri, allowed) {
			return true
		}
	}

	return false
}

// GeneratePKCEVerifier generates a cryptographically random PKCE verifier
func GeneratePKCEVerifier() string {
	// PKCE verifier: 43-128 characters, URL-safe
	bytes := make([]byte, 32) // Will result in 43 characters when base64url encoded
	rand.Read(bytes)
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
}

// GeneratePKCEChallenge generates the PKCE challenge from the verifier
func (s *OAuthState) GeneratePKCEChallenge() string {
	// For simplicity, using plain challenge method (S256 is more secure)
	// In production, should use SHA256 hash: base64url(sha256(verifier))
	return s.PKCEVerifier
}

// generateSecureRandomString generates a cryptographically secure random string
func generateSecureRandomString(length int) string {
	bytes := make([]byte, length/2+1) // Hex encoding doubles the length
	rand.Read(bytes)
	return strings.ToUpper(base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes))[:length]
}

// isValidURL checks if a string is a valid URL
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// CleanupExpiredOAuthStates removes expired OAuth states from database
func CleanupExpiredOAuthStates(db *gorm.DB) int64 {
	result := db.Where("expires_at <= ?", time.Now()).Delete(&OAuthState{})
	return result.RowsAffected
}

// OAuthStateResponse represents OAuth state data returned in API responses
type OAuthStateResponse struct {
	StateToken  string    `json:"state_token"`
	RedirectURI string    `json:"redirect_uri"`
	ExpiresAt   time.Time `json:"expires_at"`
	Challenge   string    `json:"pkce_challenge"`
}

// ToResponse converts OAuthState model to OAuthStateResponse
func (s *OAuthState) ToResponse() OAuthStateResponse {
	return OAuthStateResponse{
		StateToken:  s.StateToken,
		RedirectURI: s.RedirectURI,
		ExpiresAt:   s.ExpiresAt,
		Challenge:   s.GeneratePKCEChallenge(),
	}
}

// CreateOAuthStateRequest represents the request for creating OAuth state
type CreateOAuthStateRequest struct {
	RedirectURI string `json:"redirect_uri" binding:"required"`
}

// ValidateOAuthStateRequest represents the request for validating OAuth state
type ValidateOAuthStateRequest struct {
	StateToken   string `json:"state_token" binding:"required"`
	Code         string `json:"code" binding:"required"`
}

// OAuthStateValidationResult represents the result of OAuth state validation
type OAuthStateValidationResult struct {
	Valid        bool        `json:"valid"`
	State        *OAuthState `json:"state,omitempty"`
	PKCEVerifier string      `json:"pkce_verifier,omitempty"`
	RedirectURI  string      `json:"redirect_uri,omitempty"`
	Error        string      `json:"error,omitempty"`
}

// ValidateAndConsume validates an OAuth state and removes it from database
func ValidateAndConsume(db *gorm.DB, stateToken string) (*OAuthStateValidationResult, error) {
	var state OAuthState

	// Find the state
	result := db.Where("state_token = ?", stateToken).First(&state)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &OAuthStateValidationResult{
				Valid: false,
				Error: "invalid state token",
			}, nil
		}
		return nil, result.Error
	}

	// Check if expired
	if state.IsExpired() {
		// Clean up expired state
		db.Delete(&state)
		return &OAuthStateValidationResult{
			Valid: false,
			Error: "state token expired",
		}, nil
	}

	// State is valid, delete it to prevent reuse
	db.Delete(&state)

	return &OAuthStateValidationResult{
		Valid:        true,
		State:        &state,
		PKCEVerifier: state.PKCEVerifier,
		RedirectURI:  state.RedirectURI,
	}, nil
}

// CreateAndSave creates a new OAuth state and saves it to database
func CreateAndSave(db *gorm.DB, redirectURI string) (*OAuthState, error) {
	state, err := GenerateOAuthState(redirectURI)
	if err != nil {
		return nil, err
	}

	result := db.Create(state)
	if result.Error != nil {
		return nil, result.Error
	}

	return state, nil
}