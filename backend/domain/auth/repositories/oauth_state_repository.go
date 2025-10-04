package repositories

import (
	"context"
	"domain/auth/entities"
)

// OAuthStateRepository defines the interface for OAuth state persistence operations
type OAuthStateRepository interface {
	// Create persists a new OAuth state
	Create(ctx context.Context, state *entities.OAuthState) error

	// FindByStateToken retrieves an OAuth state by its unique state token
	FindByStateToken(ctx context.Context, stateToken string) (*entities.OAuthState, error)

	// MarkAsUsed marks an OAuth state as used to prevent replay attacks
	MarkAsUsed(ctx context.Context, stateToken string) error

	// DeleteExpired removes all expired OAuth states (cleanup job)
	DeleteExpired(ctx context.Context) (int64, error)
}
