package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"todo-app/models"
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
func (s *SessionService) CreateSession(req CreateSessionRequest) (*models.AuthenticationSession, string, error) {
	// Calculate session expiration (24 hours)
	sessionExpiresAt := time.Now().Add(24 * time.Hour)

	var session *models.AuthenticationSession

	if req.IsOAuth && req.AccessToken != "" {
		// Create OAuth session
		tokenExpiry := time.Now().Add(1 * time.Hour) // Default 1 hour
		if req.TokenExpiry != nil {
			tokenExpiry = *req.TokenExpiry
		}

		session = models.NewOAuthSession(
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
		session = models.NewSession(
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
func (s *SessionService) ValidateSession(tokenString string) (*models.SessionValidationResult, error) {
	// Validate JWT token
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return &models.SessionValidationResult{
			Valid: false,
			Error: "invalid token: " + err.Error(),
		}, nil
	}

	// Find session in database
	var session models.AuthenticationSession
	result := s.db.Preload("User").Where("id = ?", claims.SessionID).First(&session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &models.SessionValidationResult{
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
		return &models.SessionValidationResult{
			Valid: false,
			Error: "session expired",
		}, nil
	}

	// Update last activity
	session.UpdateActivity()
	s.db.Save(&session)

	// Check if OAuth tokens need refresh
	needsRefresh := session.NeedsRefresh()

	return &models.SessionValidationResult{
		Valid:        true,
		Session:      &session,
		User:         &session.User,
		NeedsRefresh: needsRefresh,
	}, nil
}

// RefreshSession refreshes a session and extends its expiration
func (s *SessionService) RefreshSession(sessionID string) (*models.AuthenticationSession, string, error) {
	var session models.AuthenticationSession

	// Find session
	result := s.db.Preload("User").Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		return nil, "", result.Error
	}

	// Check if session is expired
	if session.IsExpired() {
		return nil, "", errors.New("session has expired")
	}

	// Extend session
	if err := session.ExtendSession(); err != nil {
		return nil, "", err
	}

	// Generate new JWT token
	jwtToken, err := s.jwtService.GenerateToken(
		session.UserID,
		session.User.Email,
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
	var session models.AuthenticationSession

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
	result := s.db.Where("user_id = ?", userID).Delete(&models.AuthenticationSession{})
	return result.Error
}

// GetUserSessions retrieves all active sessions for a user
func (s *SessionService) GetUserSessions(userID uint) ([]models.AuthenticationSession, error) {
	var sessions []models.AuthenticationSession

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
	result := s.db.Where("session_expires_at <= ?", time.Now()).Delete(&models.AuthenticationSession{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// GetSession retrieves a session by ID
func (s *SessionService) GetSession(sessionID string) (*models.AuthenticationSession, error) {
	var session models.AuthenticationSession

	result := s.db.Preload("User").Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

// UpdateSessionActivity updates the last activity timestamp of a session
func (s *SessionService) UpdateSessionActivity(sessionID string) error {
	return s.db.Model(&models.AuthenticationSession{}).
		Where("id = ?", sessionID).
		Update("last_activity", time.Now()).Error
}

// GetSessionByToken retrieves a session by its JWT token
func (s *SessionService) GetSessionByToken(tokenString string) (*models.AuthenticationSession, error) {
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
	err := s.db.Model(&models.AuthenticationSession{}).
		Where("id = ? AND session_expires_at > ?", sessionID, time.Now()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}