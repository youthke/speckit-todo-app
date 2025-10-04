package persistence

import (
	"context"
	"time"

	"domain/auth/entities"
	"domain/auth/repositories"
	"gorm.io/gorm"
)

// GormSessionRepository implements SessionRepository using GORM
type GormSessionRepository struct {
	db *gorm.DB
}

// NewGormSessionRepository creates a new GORM-based session repository
func NewGormSessionRepository(db *gorm.DB) repositories.SessionRepository {
	return &GormSessionRepository{db: db}
}

// Create persists a new authentication session
func (r *GormSessionRepository) Create(ctx context.Context, session *entities.AuthenticationSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// FindByID retrieves a session by its unique session ID
func (r *GormSessionRepository) FindByID(ctx context.Context, sessionID string) (*entities.AuthenticationSession, error) {
	var session entities.AuthenticationSession
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindByUserID retrieves all active sessions for a given user
func (r *GormSessionRepository) FindByUserID(ctx context.Context, userID string) ([]*entities.AuthenticationSession, error) {
	var sessions []*entities.AuthenticationSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// Update modifies an existing session
func (r *GormSessionRepository) Update(ctx context.Context, session *entities.AuthenticationSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// Delete removes a session by ID
func (r *GormSessionRepository) Delete(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Delete(&entities.AuthenticationSession{}).Error
}

// DeleteExpired removes all expired sessions (cleanup job)
func (r *GormSessionRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.AuthenticationSession{})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
