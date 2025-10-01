package jobs

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
	"todo-app/internal/models"
)

// SessionCleanupJob handles cleanup of expired authentication sessions
type SessionCleanupJob struct {
	db       *gorm.DB
	interval time.Duration
	done     chan bool
}

// NewSessionCleanupJob creates a new session cleanup job
func NewSessionCleanupJob(db *gorm.DB, interval time.Duration) *SessionCleanupJob {
	if interval == 0 {
		interval = 1 * time.Hour // Default to hourly cleanup
	}

	return &SessionCleanupJob{
		db:       db,
		interval: interval,
		done:     make(chan bool),
	}
}

// Start begins the session cleanup job
func (j *SessionCleanupJob) Start(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	log.Printf("Session cleanup job started (interval: %v)", j.interval)

	// Run cleanup immediately on start
	j.cleanup(ctx)

	for {
		select {
		case <-ticker.C:
			j.cleanup(ctx)
		case <-ctx.Done():
			log.Println("Session cleanup job stopped")
			j.done <- true
			return
		}
	}
}

// Stop stops the session cleanup job
func (j *SessionCleanupJob) Stop() {
	<-j.done
}

// cleanup removes expired authentication sessions
func (j *SessionCleanupJob) cleanup(ctx context.Context) {
	startTime := time.Now()

	// Delete expired sessions
	result := j.db.WithContext(ctx).
		Where("session_expires_at < ?", time.Now()).
		Delete(&models.AuthenticationSession{})

	if result.Error != nil {
		log.Printf("Error cleaning up expired sessions: %v", result.Error)
		return
	}

	duration := time.Since(startTime)
	if result.RowsAffected > 0 {
		log.Printf("Session cleanup completed: removed %d expired sessions in %v",
			result.RowsAffected, duration)
	}

	// Also cleanup inactive sessions (no activity for 7 days)
	j.cleanupInactiveSessions(ctx)
}

// cleanupInactiveSessions removes sessions with no activity for extended period
func (j *SessionCleanupJob) cleanupInactiveSessions(ctx context.Context) {
	inactivityThreshold := time.Now().Add(-7 * 24 * time.Hour) // 7 days

	result := j.db.WithContext(ctx).
		Where("last_activity < ?", inactivityThreshold).
		Delete(&models.AuthenticationSession{})

	if result.Error != nil {
		log.Printf("Error cleaning up inactive sessions: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("Removed %d inactive sessions (no activity for 7+ days)",
			result.RowsAffected)
	}
}

// RunOnce executes cleanup once (useful for testing or manual execution)
func (j *SessionCleanupJob) RunOnce(ctx context.Context) error {
	j.cleanup(ctx)
	return nil
}

// GetStats returns statistics about authentication sessions
func (j *SessionCleanupJob) GetStats(ctx context.Context) (map[string]interface{}, error) {
	var totalCount int64
	var expiredCount int64
	var activeCount int64
	var oauthCount int64

	// Count total sessions
	if err := j.db.WithContext(ctx).
		Model(&models.AuthenticationSession{}).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Count expired sessions
	if err := j.db.WithContext(ctx).
		Model(&models.AuthenticationSession{}).
		Where("session_expires_at < ?", time.Now()).
		Count(&expiredCount).Error; err != nil {
		return nil, err
	}

	// Count active sessions
	if err := j.db.WithContext(ctx).
		Model(&models.AuthenticationSession{}).
		Where("session_expires_at >= ?", time.Now()).
		Count(&activeCount).Error; err != nil {
		return nil, err
	}

	// Count OAuth sessions
	if err := j.db.WithContext(ctx).
		Model(&models.AuthenticationSession{}).
		Where("access_token IS NOT NULL").
		Count(&oauthCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_sessions":   totalCount,
		"active_sessions":  activeCount,
		"expired_sessions": expiredCount,
		"oauth_sessions":   oauthCount,
	}, nil
}

// CleanupSessionsByUserID removes all sessions for a specific user
func CleanupSessionsByUserID(db *gorm.DB, userID uint) error {
	result := db.Where("user_id = ?", userID).
		Delete(&models.AuthenticationSession{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Printf("Removed %d sessions for user %d", result.RowsAffected, userID)
	}

	return nil
}

// CleanupSessionsOlderThan removes sessions older than specified duration
func CleanupSessionsOlderThan(db *gorm.DB, duration time.Duration) error {
	cutoffTime := time.Now().Add(-duration)

	result := db.Where("created_at < ?", cutoffTime).
		Delete(&models.AuthenticationSession{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Printf("Removed %d sessions older than %v", result.RowsAffected, duration)
	}

	return nil
}