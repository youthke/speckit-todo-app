package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"domain/auth/entities"
	"todo-app/internal/models"
)

// SessionService handles session management operations
type SessionService struct {
	db         *gorm.DB
	jwtService *JWTService
}

// NewSessionService creates a new session service
func NewSessionService(db *gorm.DB, jwtService *JWTService) *SessionService {
	return &SessionService{
		db:         db,
		jwtService: jwtService,
	}
}

// CreateSessionRequest represents the data needed to create a session
type CreateSessionRequest struct {
	UserID       uint
	Email        string
	UserAgent    string
	IPAddress    string
	IsOAuth      bool
	AccessToken  string
	RefreshToken string
	TokenExpiry  *time.Time
}

// CreateSession creates a new authentication session
func (s *SessionService) CreateSession(req CreateSessionRequest) (*entities.AuthenticationSession, string, error) {
	// Calculate session expiration (24 hours)
	sessionExpiresAt := time.Now().Add(24 * time.Hour)

	var session *entities.AuthenticationSession

	if req.IsOAuth && req.AccessToken != "" {
		// Create OAuth session
		tokenExpiry := time.Now().Add(1 * time.Hour) // Default 1 hour
		if req.TokenExpiry != nil {
			tokenExpiry = *req.TokenExpiry
		}

		session = entities.NewOAuthSession(
			req.UserID,
			"", // JWT token will be set below
			req.AccessToken,
			req.RefreshToken,
			tokenExpiry,
			sessionExpiresAt,
			req.UserAgent,
			req.IPAddress,
		)
	} else {
		// Create regular session
		session = entities.NewSession(
			req.UserID,
			"", // JWT token will be set below
			sessionExpiresAt,
			req.UserAgent,
			req.IPAddress,
		)
	}

	// Generate JWT token
	jwtToken, err := s.jwtService.GenerateToken(req.UserID, req.Email, session.ID, req.IsOAuth)
	if err != nil {
		return nil, "", err
	}

	session.SessionToken = jwtToken

	// Save session to database
	if err := s.db.Create(session).Error; err != nil {
		return nil, "", err
	}

	return session, jwtToken, nil
}

// ValidateSession validates a session token and returns the session
func (s *SessionService) ValidateSession(tokenString string) (*entities.SessionValidationResult, error) {
	// Validate JWT token
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return &entities.SessionValidationResult{
			Valid: false,
			Error: "invalid token: " + err.Error(),
		}, nil
	}

	// Find session in database
	var session entities.AuthenticationSession
	result := s.db.Where("id = ?", claims.SessionID).First(&session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &entities.SessionValidationResult{
				Valid: false,
				Error: "session not found",
			}, nil
		}
		return nil, result.Error
	}

	// Check if session is expired
	if session.IsExpired() {
		// Delete expired session
		s.db.Delete(&session)
		return &entities.SessionValidationResult{
			Valid: false,
			Error: "session expired",
		}, nil
	}

	// Load user separately as simple model
	var user models.User
	if err := s.db.Where("id = ?", session.UserID).First(&user).Error; err != nil {
		return &entities.SessionValidationResult{
			Valid: false,
			Error: "user not found",
		}, nil
	}

	// Update last activity
	session.UpdateActivity()
	s.db.Save(&session)

	// Check if OAuth tokens need refresh
	needsRefresh := session.NeedsRefresh()

	return &entities.SessionValidationResult{
		Valid:        true,
		Session:      &session,
		User:         &user,
		NeedsRefresh: needsRefresh,
	}, nil
}

// RefreshSession refreshes a session and extends its expiration
func (s *SessionService) RefreshSession(sessionID string) (*entities.AuthenticationSession, string, error) {
	var session entities.AuthenticationSession

	// Find session
	result := s.db.Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		return nil, "", result.Error
	}

	// Check if session is expired
	if session.IsExpired() {
		return nil, "", errors.New("session has expired")
	}

	// Get user for JWT generation
	var user models.User
	if err := s.db.Where("id = ?", session.UserID).First(&user).Error; err != nil {
		return nil, "", err
	}

	// Extend session
	if err := session.ExtendSession(); err != nil {
		return nil, "", err
	}

	// Generate new JWT token
	jwtToken, err := s.jwtService.GenerateToken(
		session.UserID,
		user.Email,
		session.ID,
		session.IsOAuthSession(),
	)
	if err != nil {
		return nil, "", err
	}

	session.SessionToken = jwtToken

	// Save updated session
	if err := s.db.Save(&session).Error; err != nil {
		return nil, "", err
	}

	return &session, jwtToken, nil
}

// TerminateSession terminates a session
func (s *SessionService) TerminateSession(sessionID string) error {
	var session entities.AuthenticationSession

	// Find session
	result := s.db.Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // Already terminated
		}
		return result.Error
	}

	// Delete session
	if err := s.db.Delete(&session).Error; err != nil {
		return err
	}

	return nil
}

// TerminateAllUserSessions terminates all sessions for a user
func (s *SessionService) TerminateAllUserSessions(userID uint) error {
	result := s.db.Where("user_id = ?", userID).Delete(&entities.AuthenticationSession{})
	return result.Error
}

// GetUserSessions retrieves all active sessions for a user
func (s *SessionService) GetUserSessions(userID uint) ([]entities.AuthenticationSession, error) {
	var sessions []entities.AuthenticationSession

	result := s.db.Where("user_id = ? AND session_expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&sessions)

	if result.Error != nil {
		return nil, result.Error
	}

	return sessions, nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *SessionService) CleanupExpiredSessions() (int64, error) {
	result := s.db.Where("session_expires_at <= ?", time.Now()).Delete(&entities.AuthenticationSession{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// GetSession retrieves a session by ID
func (s *SessionService) GetSession(sessionID string) (*entities.AuthenticationSession, error) {
	var session entities.AuthenticationSession

	result := s.db.Preload("User").Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

// UpdateSessionActivity updates the last activity timestamp of a session
func (s *SessionService) UpdateSessionActivity(sessionID string) error {
	return s.db.Model(&entities.AuthenticationSession{}).
		Where("id = ?", sessionID).
		Update("last_activity", time.Now()).Error
}

// GetSessionByToken retrieves a session by its JWT token
func (s *SessionService) GetSessionByToken(tokenString string) (*entities.AuthenticationSession, error) {
	// Extract session ID from token
	sessionID, err := s.jwtService.ExtractSessionID(tokenString)
	if err != nil {
		return nil, err
	}

	return s.GetSession(sessionID)
}

// IsSessionValid checks if a session is valid without full validation
func (s *SessionService) IsSessionValid(sessionID string) (bool, error) {
	var count int64
	err := s.db.Model(&entities.AuthenticationSession{}).
		Where("id = ? AND session_expires_at > ?", sessionID, time.Now()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}