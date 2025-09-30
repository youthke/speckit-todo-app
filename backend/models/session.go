package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"
)

// AuthenticationSession represents an active user session with OAuth token management
type AuthenticationSession struct {
	ID        string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID    uint   `json:"user_id" gorm:"not null;index"`
	User      User   `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	// Session tokens
	SessionToken string `json:"-" gorm:"type:text;uniqueIndex;not null"`

	// OAuth tokens (encrypted at rest)
	RefreshToken   string     `json:"-" gorm:"type:text"`
	AccessToken    string     `json:"-" gorm:"type:text"`
	TokenExpiresAt *time.Time `json:"token_expires_at"`

	// Session management
	SessionExpiresAt time.Time `json:"session_expires_at" gorm:"not null;index"`
	LastActivity     time.Time `json:"last_activity" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Audit fields
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UserAgent string    `json:"user_agent" gorm:"type:text"`
	IPAddress string    `json:"ip_address" gorm:"type:varchar(45)"`
}

// TableName specifies the table name for the AuthenticationSession model
func (AuthenticationSession) TableName() string {
	return "authentication_sessions"
}

// BeforeCreate hook to validate session before creation
func (s *AuthenticationSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = generateSessionID()
	}
	if s.LastActivity.IsZero() {
		s.LastActivity = time.Now()
	}
	return s.Validate()
}

// BeforeUpdate hook to validate session before update
func (s *AuthenticationSession) BeforeUpdate(tx *gorm.DB) error {
	return s.Validate()
}

// Validate performs validation on the AuthenticationSession model
func (s *AuthenticationSession) Validate() error {
	if s.UserID == 0 {
		return errors.New("user_id is required")
	}

	if s.SessionToken == "" {
		return errors.New("session_token cannot be empty")
	}

	if s.SessionExpiresAt.Before(time.Now()) {
		return errors.New("session cannot be expired")
	}

	// Session cannot be longer than 24 hours
	maxSessionTime := time.Now().Add(24 * time.Hour)
	if s.SessionExpiresAt.After(maxSessionTime) {
		return errors.New("session_expires_at cannot exceed 24 hours")
	}

	// If access_token is present, token_expires_at is required
	if s.AccessToken != "" && s.TokenExpiresAt == nil {
		return errors.New("token_expires_at required when access_token present")
	}

	return nil
}

// IsExpired returns true if the session has expired
func (s *AuthenticationSession) IsExpired() bool {
	return s.SessionExpiresAt.Before(time.Now()) || s.SessionExpiresAt.Equal(time.Now())
}

// IsTokenExpired returns true if the OAuth tokens have expired
func (s *AuthenticationSession) IsTokenExpired() bool {
	if s.TokenExpiresAt == nil {
		return false
	}
	return s.TokenExpiresAt.Before(time.Now()) || s.TokenExpiresAt.Equal(time.Now())
}

// NeedsRefresh returns true if OAuth tokens need to be refreshed soon
func (s *AuthenticationSession) NeedsRefresh() bool {
	if s.AccessToken == "" || s.TokenExpiresAt == nil {
		return false
	}

	// Refresh if tokens expire within 5 minutes
	refreshThreshold := time.Now().Add(5 * time.Minute)
	return s.TokenExpiresAt.Before(refreshThreshold)
}

// UpdateActivity updates the last activity timestamp
func (s *AuthenticationSession) UpdateActivity() {
	s.LastActivity = time.Now()
}

// ExtendSession extends the session expiry if within allowed time
func (s *AuthenticationSession) ExtendSession() error {
	// Only extend if session is still valid and user has been active
	if s.IsExpired() {
		return errors.New("cannot extend expired session")
	}

	// Extend session by 24 hours from now
	s.SessionExpiresAt = time.Now().Add(24 * time.Hour)
	s.UpdateActivity()

	return s.Validate()
}

// UpdateOAuthTokens updates the OAuth access and refresh tokens
func (s *AuthenticationSession) UpdateOAuthTokens(accessToken, refreshToken string, expiresAt time.Time) error {
	s.AccessToken = accessToken
	s.RefreshToken = refreshToken
	s.TokenExpiresAt = &expiresAt
	s.UpdateActivity()

	return s.Validate()
}

// ClearOAuthTokens removes OAuth tokens from the session
func (s *AuthenticationSession) ClearOAuthTokens() {
	s.AccessToken = ""
	s.RefreshToken = ""
	s.TokenExpiresAt = nil
	s.UpdateActivity()
}

// IsOAuthSession returns true if this session has OAuth tokens
func (s *AuthenticationSession) IsOAuthSession() bool {
	return s.AccessToken != "" || s.RefreshToken != ""
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return "sess_" + hex.EncodeToString(bytes)
}

// CreateSessionRequest represents the request for creating a new session
type CreateSessionRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	UserAgent   string `json:"user_agent"`
	IPAddress   string `json:"ip_address"`
	AccessToken string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
}

// SessionResponse represents the session data returned in API responses
type SessionResponse struct {
	SessionID     string     `json:"session_id"`
	ExpiresAt     time.Time  `json:"expires_at"`
	LastActivity  time.Time  `json:"last_activity"`
	IsOAuth       bool       `json:"is_oauth"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
}

// ToResponse converts AuthenticationSession model to SessionResponse
func (s *AuthenticationSession) ToResponse() SessionResponse {
	return SessionResponse{
		SessionID:      s.ID,
		ExpiresAt:      s.SessionExpiresAt,
		LastActivity:   s.LastActivity,
		IsOAuth:        s.IsOAuthSession(),
		TokenExpiresAt: s.TokenExpiresAt,
	}
}

// SessionValidationResult represents the result of session validation
type SessionValidationResult struct {
	Valid         bool                   `json:"valid"`
	Session       *AuthenticationSession `json:"session,omitempty"`
	User          *User                  `json:"user,omitempty"`
	NeedsRefresh  bool                   `json:"needs_refresh"`
	Error         string                 `json:"error,omitempty"`
}

// NewSession creates a new authentication session
func NewSession(userID uint, sessionToken string, expiresAt time.Time, userAgent, ipAddress string) *AuthenticationSession {
	return &AuthenticationSession{
		ID:               generateSessionID(),
		UserID:           userID,
		SessionToken:     sessionToken,
		SessionExpiresAt: expiresAt,
		LastActivity:     time.Now(),
		UserAgent:        userAgent,
		IPAddress:        ipAddress,
	}
}

// NewOAuthSession creates a new OAuth authentication session
func NewOAuthSession(userID uint, sessionToken string, accessToken, refreshToken string, tokenExpiresAt, sessionExpiresAt time.Time, userAgent, ipAddress string) *AuthenticationSession {
	session := NewSession(userID, sessionToken, sessionExpiresAt, userAgent, ipAddress)
	session.AccessToken = accessToken
	session.RefreshToken = refreshToken
	session.TokenExpiresAt = &tokenExpiresAt
	return session
}