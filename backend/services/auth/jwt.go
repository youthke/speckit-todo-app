package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token operations
type JWTService struct {
	secretKey     []byte
	expiresHours  int
	issuer        string
}

// JWTClaims represents the claims stored in the JWT token
type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
	IsOAuth   bool   `json:"is_oauth"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service from environment variables
func NewJWTService() (*JWTService, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET environment variable is not set")
	}

	expiresHoursStr := os.Getenv("JWT_EXPIRES_HOURS")
	if expiresHoursStr == "" {
		expiresHoursStr = "24" // Default to 24 hours
	}

	expiresHours, err := strconv.Atoi(expiresHoursStr)
	if err != nil {
		return nil, errors.New("JWT_EXPIRES_HOURS must be a valid integer")
	}

	return &JWTService{
		secretKey:    []byte(secretKey),
		expiresHours: expiresHours,
		issuer:       "todo-app",
	}, nil
}

// GenerateToken generates a new JWT token for a user session
func (s *JWTService) GenerateToken(userID uint, email, sessionID string, isOAuth bool) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.expiresHours) * time.Hour)

	claims := JWTClaims{
		UserID:    userID,
		Email:     email,
		SessionID: sessionID,
		IsOAuth:   isOAuth,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   strconv.FormatUint(uint64(userID), 10),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        sessionID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token with extended expiration
func (s *JWTService) RefreshToken(oldTokenString string) (string, error) {
	// Validate the old token
	claims, err := s.ValidateToken(oldTokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same user info but new expiration
	return s.GenerateToken(claims.UserID, claims.Email, claims.SessionID, claims.IsOAuth)
}

// ExtractUserID extracts the user ID from a JWT token without full validation
func (s *JWTService) ExtractUserID(tokenString string) (uint, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token claims")
}

// ExtractSessionID extracts the session ID from a JWT token without full validation
func (s *JWTService) ExtractSessionID(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims.SessionID, nil
	}

	return "", errors.New("invalid token claims")
}

// IsExpired checks if a token is expired
func (s *JWTService) IsExpired(tokenString string) (bool, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		// If validation fails, consider it expired
		return true, err
	}

	return claims.ExpiresAt.Before(time.Now()), nil
}

// GetExpirationTime returns the expiration time of a token
func (s *JWTService) GetExpirationTime(tokenString string) (time.Time, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}

// RevokeToken marks a token as revoked (implementation depends on token blacklist strategy)
// For now, this is a placeholder - actual revocation is handled by deleting the session
func (s *JWTService) RevokeToken(tokenString string) error {
	// In a production system, you would add the token to a blacklist/revocation list
	// For this implementation, tokens are revoked by deleting the session from the database
	return nil
}

// TokenValidationResult represents the result of token validation
type TokenValidationResult struct {
	Valid     bool        `json:"valid"`
	Claims    *JWTClaims  `json:"claims,omitempty"`
	ExpiresAt time.Time   `json:"expires_at,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// ValidateAndParse validates a token and returns a detailed result
func (s *JWTService) ValidateAndParse(tokenString string) *TokenValidationResult {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return &TokenValidationResult{
			Valid: false,
			Error: err.Error(),
		}
	}

	return &TokenValidationResult{
		Valid:     true,
		Claims:    claims,
		ExpiresAt: claims.ExpiresAt.Time,
	}
}