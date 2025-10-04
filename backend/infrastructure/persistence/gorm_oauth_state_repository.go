package persistence

import (
	"context"
	"time"

	"domain/auth/entities"
	"domain/auth/repositories"
	"gorm.io/gorm"
)

// GormOAuthStateRepository implements OAuthStateRepository using GORM
type GormOAuthStateRepository struct {
	db *gorm.DB
}

// NewGormOAuthStateRepository creates a new GORM-based OAuth state repository
func NewGormOAuthStateRepository(db *gorm.DB) repositories.OAuthStateRepository {
	return &GormOAuthStateRepository{db: db}
}

// Create persists a new OAuth state
func (r *GormOAuthStateRepository) Create(ctx context.Context, state *entities.OAuthState) error {
	return r.db.WithContext(ctx).Create(state).Error
}

// FindByStateToken retrieves an OAuth state by its unique state token
func (r *GormOAuthStateRepository) FindByStateToken(ctx context.Context, stateToken string) (*entities.OAuthState, error) {
	var state entities.OAuthState
	err := r.db.WithContext(ctx).Where("state_token = ?", stateToken).First(&state).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}

// MarkAsUsed marks an OAuth state as used to prevent replay attacks
func (r *GormOAuthStateRepository) MarkAsUsed(ctx context.Context, stateToken string) error {
	return r.db.WithContext(ctx).
		Model(&entities.OAuthState{}).
		Where("state_token = ?", stateToken).
		Update("used", true).Error
}

// DeleteExpired removes all expired OAuth states (cleanup job)
func (r *GormOAuthStateRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.OAuthState{})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
