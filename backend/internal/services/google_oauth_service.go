package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"todo-app/internal/config"
	"todo-app/internal/models"
)

// GoogleUserInfo contains user information from Google OAuth
type GoogleUserInfo struct {
	GoogleUserID  string
	Email         string
	EmailVerified bool
	Name          string
}

// GoogleOAuthService handles Google OAuth authentication
type GoogleOAuthService struct {
	config *oauth2.Config
	db     *gorm.DB
}

// NewGoogleOAuthService creates a new Google OAuth service
func NewGoogleOAuthService(db *gorm.DB) *GoogleOAuthService {
	return &GoogleOAuthService{
		config: config.GetGoogleOAuthConfig(),
		db:     db,
	}
}

// GenerateAuthURL creates OAuth URL with state token for CSRF protection
func (s *GoogleOAuthService) GenerateAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges authorization code for user info
func (s *GoogleOAuthService) ExchangeCode(ctx context.Context, code string) (*GoogleUserInfo, error) {
	// Exchange code for token
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Google
	client := s.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
	}

	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &GoogleUserInfo{
		GoogleUserID:  googleUser.ID,
		Email:         googleUser.Email,
		EmailVerified: googleUser.VerifiedEmail,
		Name:          googleUser.Name,
	}, nil
}

// CreateUserFromGoogle creates a new user and GoogleIdentity from Google OAuth info
func (s *GoogleOAuthService) CreateUserFromGoogle(info *GoogleUserInfo) (*models.User, error) {
	// Validate email is verified
	if !info.EmailVerified {
		return nil, errors.New("email not verified")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	user := models.User{
		Email:      info.Email,
		Name:       info.Name,
		AuthMethod: "google",
		IsActive:   true,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create Google identity link
	googleIdentity := models.GoogleIdentity{
		UserID:        user.ID,
		GoogleUserID:  info.GoogleUserID,
		Email:         info.Email,
		EmailVerified: info.EmailVerified,
	}

	if err := tx.Create(&googleIdentity).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create Google identity: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

// FindUserByGoogleID checks if a user with the given Google ID already exists
func (s *GoogleOAuthService) FindUserByGoogleID(googleUserID string) (*models.User, error) {
	var googleIdentity models.GoogleIdentity
	if err := s.db.Where("google_user_id = ?", googleUserID).First(&googleIdentity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No user found, not an error
		}
		return nil, fmt.Errorf("failed to query Google identity: %w", err)
	}

	// Load the associated user
	var user models.User
	if err := s.db.First(&user, googleIdentity.UserID).Error; err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	return &user, nil
}
