package repositories

import (
	"context"
	"domain/auth/entities"
)

// SessionRepository defines the interface for session persistence operations
type SessionRepository interface {
	// Create persists a new authentication session
	Create(ctx context.Context, session *entities.AuthenticationSession) error

	// FindByID retrieves a session by its unique session ID
	FindByID(ctx context.Context, sessionID string) (*entities.AuthenticationSession, error)

	// FindByUserID retrieves all active sessions for a given user
	FindByUserID(ctx context.Context, userID string) ([]*entities.AuthenticationSession, error)

	// Update modifies an existing session
	Update(ctx context.Context, session *entities.AuthenticationSession) error

	// Delete removes a session by ID
	Delete(ctx context.Context, sessionID string) error

	// DeleteExpired removes all expired sessions (cleanup job)
	DeleteExpired(ctx context.Context) (int64, error)
}
