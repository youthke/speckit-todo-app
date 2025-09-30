package jobs

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
	"todo-app/models"
)

// OAuthCleanupJob handles cleanup of expired OAuth state records
type OAuthCleanupJob struct {
	db       *gorm.DB
	interval time.Duration
	done     chan bool
}

// NewOAuthCleanupJob creates a new OAuth cleanup job
func NewOAuthCleanupJob(db *gorm.DB, interval time.Duration) *OAuthCleanupJob {
	if interval == 0 {
		interval = 5 * time.Minute // Default to 5 minutes
	}

	return &OAuthCleanupJob{
		db:       db,
		interval: interval,
		done:     make(chan bool),
	}
}

// Start begins the OAuth cleanup job
func (j *OAuthCleanupJob) Start(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	log.Printf("OAuth cleanup job started (interval: %v)", j.interval)

	// Run cleanup immediately on start
	j.cleanup(ctx)

	for {
		select {
		case <-ticker.C:
			j.cleanup(ctx)
		case <-ctx.Done():
			log.Println("OAuth cleanup job stopped")
			j.done <- true
			return
		}
	}
}

// Stop stops the OAuth cleanup job
func (j *OAuthCleanupJob) Stop() {
	<-j.done
}

// cleanup removes expired OAuth state records
func (j *OAuthCleanupJob) cleanup(ctx context.Context) {
	startTime := time.Now()

	// Delete expired OAuth states
	result := j.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.OAuthState{})

	if result.Error != nil {
		log.Printf("Error cleaning up OAuth states: %v", result.Error)
		return
	}

	duration := time.Since(startTime)
	if result.RowsAffected > 0 {
		log.Printf("OAuth cleanup completed: removed %d expired states in %v",
			result.RowsAffected, duration)
	}
}

// RunOnce executes cleanup once (useful for testing or manual execution)
func (j *OAuthCleanupJob) RunOnce(ctx context.Context) error {
	j.cleanup(ctx)
	return nil
}

// GetStats returns statistics about OAuth state records
func (j *OAuthCleanupJob) GetStats(ctx context.Context) (map[string]interface{}, error) {
	var totalCount int64
	var expiredCount int64

	// Count total OAuth states
	if err := j.db.WithContext(ctx).
		Model(&models.OAuthState{}).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Count expired OAuth states
	if err := j.db.WithContext(ctx).
		Model(&models.OAuthState{}).
		Where("expires_at < ?", time.Now()).
		Count(&expiredCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_states":   totalCount,
		"expired_states": expiredCount,
		"active_states":  totalCount - expiredCount,
	}, nil
}

// CleanupOAuthStatesOlderThan removes OAuth states older than specified duration
func CleanupOAuthStatesOlderThan(db *gorm.DB, duration time.Duration) error {
	cutoffTime := time.Now().Add(-duration)

	result := db.Where("created_at < ?", cutoffTime).
		Delete(&models.OAuthState{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Printf("Removed %d OAuth states older than %v", result.RowsAffected, duration)
	}

	return nil
}