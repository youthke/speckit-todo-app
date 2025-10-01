package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"todo-app/internal/config"
)

// SessionService handles JWT session management
type SessionService struct {
	jwtSecret string
}

// NewSessionService creates a new session service
func NewSessionService() *SessionService {
	return &SessionService{
		jwtSecret: config.GetJWTSecret(),
	}
}

// CreateSession generates a JWT token with 7-day expiration
func (s *SessionService) CreateSession(userID uint) (string, error) {
	// Set expiration to 7 days from now
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateSession verifies a JWT token and returns the user ID
func (s *SessionService) ValidateSession(tokenString string) (uint, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Extract user ID
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id in token")
	}

	return uint(userIDFloat), nil
}

// GetSessionMaxAge returns the max age in seconds for session cookies (7 days)
func (s *SessionService) GetSessionMaxAge() int {
	return 7 * 24 * 60 * 60 // 604800 seconds
}
