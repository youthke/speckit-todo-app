# Data Model: Import Path Cleanup

**Feature**: 009-resolve-it-1
**Date**: 2025-10-02

## Overview

This document defines the **target** data model structure after the import path cleanup. This is a refactoring feature, so the models themselves don't change functionally - only their location and organization within the DDD structure.

## Domain Contexts

The refactored codebase will have four bounded contexts:

1. **User Domain** (existing, well-established)
2. **Task Domain** (existing, well-established)
3. **Auth Domain** (new, consolidating scattered auth models)
4. **Health Domain** (new, extracted from internal/models)

## 1. Auth Domain (NEW)

### 1.1 Entities

#### AuthenticationSession

**Purpose**: Represents an active user session with authentication state

**Location**: `domain/auth/entities/authentication_session.go`

**Fields**:
```go
type AuthenticationSession struct {
    SessionID        string        // Unique session identifier (PK)
    UserID           string        // Foreign key to User
    AccessToken      string        // JWT or opaque token
    RefreshToken     string        // Token for session renewal
    ExpiresAt        time.Time     // Session expiration
    RefreshExpiresAt time.Time     // Refresh token expiration
    CreatedAt        time.Time     // Creation timestamp
    LastAccessedAt   time.Time     // Last activity timestamp
    IPAddress        string        // Client IP (security audit)
    UserAgent        string        // Client user agent
    IsActive         bool          // Soft delete flag
}
```

**Relationships**:
- Many-to-One with User (one user can have multiple sessions)

**Validation Rules**:
- SessionID must be unique, non-empty
- UserID must reference valid User
- ExpiresAt must be in the future for active sessions
- AccessToken must not be empty
- IsActive=false invalidates session regardless of expiration

**State Transitions**:
- `Created` → `Active` (user logs in)
- `Active` → `Expired` (time-based)
- `Active` → `Revoked` (user logs out or admin action)
- `Expired` → `Refreshed` (new session created via refresh token)

**Database Mapping** (GORM):
- Table: `authentication_sessions` (existing table, name preserved)
- Indexes: `user_id`, `session_id` (unique), `expires_at`

---

#### OAuthState

**Purpose**: Temporary state token for OAuth 2.0 PKCE flow security

**Location**: `domain/auth/entities/oauth_state.go`

**Fields**:
```go
type OAuthState struct {
    StateToken      string    // Unique state token (PK)
    RedirectURI     string    // Post-auth redirect target
    PKCEVerifier    string    // PKCE code verifier (OAuth 2.0)
    PKCEChallenge   string    // PKCE code challenge (derived)
    CreatedAt       time.Time // Creation timestamp
    ExpiresAt       time.Time // Short TTL (e.g., 10 minutes)
    Used            bool      // One-time use flag
}
```

**Relationships**:
- None (ephemeral, deleted after OAuth callback)

**Validation Rules**:
- StateToken must be cryptographically random, unique
- PKCEVerifier must be 43-128 characters (RFC 7636)
- PKCEChallenge must be base64url(SHA256(PKCEVerifier))
- ExpiresAt must be CreatedAt + 10 minutes (configurable)
- Used flag prevents replay attacks

**State Transitions**:
- `Created` → `Pending` (user redirected to OAuth provider)
- `Pending` → `Consumed` (callback received, Used=true)
- `Pending` → `Expired` (10 minutes elapsed)

**Database Mapping** (GORM):
- Table: `oauth_states` (existing table, name preserved)
- Indexes: `state_token` (unique), `expires_at`

---

### 1.2 Value Objects

#### SessionToken
**Location**: `domain/auth/valueobjects/session_token.go`
**Purpose**: Immutable wrapper for session token strings with validation
**Fields**: `value string`
**Validation**: Non-empty, min 32 chars, matches token format

#### PKCEVerifier
**Location**: `domain/auth/valueobjects/pkce_verifier.go`
**Purpose**: PKCE code verifier (RFC 7636 compliant)
**Fields**: `value string`
**Validation**: 43-128 characters, base64url charset

#### StateToken
**Location**: `domain/auth/valueobjects/state_token.go`
**Purpose**: OAuth state token value object
**Fields**: `value string`
**Validation**: Non-empty, cryptographically random, min 32 chars

---

### 1.3 Repositories (Interfaces)

#### SessionRepository
**Location**: `domain/auth/repositories/session_repository.go`

```go
type SessionRepository interface {
    Create(ctx context.Context, session *entities.AuthenticationSession) error
    FindByID(ctx context.Context, sessionID string) (*entities.AuthenticationSession, error)
    FindByUserID(ctx context.Context, userID string) ([]*entities.AuthenticationSession, error)
    Update(ctx context.Context, session *entities.AuthenticationSession) error
    Delete(ctx context.Context, sessionID string) error
    DeleteExpired(ctx context.Context) (int64, error) // Cleanup job
}
```

#### OAuthStateRepository
**Location**: `domain/auth/repositories/oauth_state_repository.go`

```go
type OAuthStateRepository interface {
    Create(ctx context.Context, state *entities.OAuthState) error
    FindByStateToken(ctx context.Context, stateToken string) (*entities.OAuthState, error)
    MarkAsUsed(ctx context.Context, stateToken string) error
    DeleteExpired(ctx context.Context) (int64, error) // Cleanup job
}
```

---

## 2. User Domain (EXISTING - No Changes)

**Location**: `domain/user/`

**Entities**:
- `entities/user.go` - User aggregate root

**Value Objects**:
- `valueobjects/email.go` - Email address with validation
- `valueobjects/user_id.go` - Strongly-typed user ID
- `valueobjects/user_profile.go` - Profile information
- `valueobjects/user_preferences.go` - User settings

**Repositories**:
- `repositories/user_repository.go` - User persistence interface

**Services**:
- `services/user_authentication_service.go` - Auth logic
- `services/user_profile_service.go` - Profile management

**Status**: ✅ Well-established, no refactoring needed

---

## 3. Task Domain (EXISTING - No Changes)

**Location**: `domain/task/`

**Entities**:
- `entities/task.go` - Task aggregate root

**Value Objects**:
- `valueobjects/task_id.go` - Strongly-typed task ID
- `valueobjects/task_title.go` - Title with length validation
- `valueobjects/task_description.go` - Description field
- `valueobjects/task_status.go` - Status enum (pending/in_progress/completed)
- `valueobjects/task_priority.go` - Priority enum (low/medium/high)

**Repositories**:
- `repositories/task_repository.go` - Task persistence interface

**Services**:
- `services/task_validation_service.go` - Business rule validation
- `services/task_search_service.go` - Query operations

**Status**: ✅ Well-established, no refactoring needed

---

## 4. Health Domain (NEW)

### 4.1 Entities

#### HealthStatus

**Purpose**: System health check status

**Location**: `domain/health/entities/health_status.go`

**Fields**:
```go
type HealthStatus struct {
    Status         string            // "healthy" | "degraded" | "unhealthy"
    Timestamp      time.Time         // Check timestamp
    Version        string            // Application version
    Dependencies   map[string]string // e.g., {"database": "healthy", "cache": "degraded"}
    Uptime         time.Duration     // Time since startup
}
```

**Relationships**: None (stateless health check)

**Validation Rules**:
- Status must be one of: "healthy", "degraded", "unhealthy"
- Timestamp must not be in the future
- Dependencies map keys must match known service names

**Database Mapping**: None (not persisted, computed on demand)

---

## 5. Deprecated Models (TO BE REMOVED)

### 5.1 internal/models/user.go
**Status**: DEPRECATED
**Replacement**: `domain/user/entities/user.go`
**Reason**: Flat DTO-style model, use rich domain entity

### 5.2 internal/models/task.go
**Status**: DEPRECATED
**Replacement**: `domain/task/entities/task.go`
**Reason**: Flat DTO-style model, use rich domain entity

### 5.3 backend/models/session.go
**Status**: DEPRECATED (ORPHANED)
**Replacement**: `domain/auth/entities/authentication_session.go`
**Reason**: Incorrectly located outside internal/models/, causing import errors

### 5.4 backend/models/oauth_state.go
**Status**: DEPRECATED (ORPHANED)
**Replacement**: `domain/auth/entities/oauth_state.go`
**Reason**: Incorrectly located outside internal/models/, causing import errors

---

## 6. Migration Strategy

### 6.1 Data Migration
**NOT REQUIRED** - This is a code refactoring only. No database schema changes.

**Database Compatibility**:
- GORM `TableName()` methods will preserve existing table names
- Example:
  ```go
  func (AuthenticationSession) TableName() string {
      return "authentication_sessions" // Existing table name
  }
  ```

### 6.2 Code Migration (Import Path Changes)
See research.md section 6 for comprehensive list of ~50-60 files requiring import updates.

**Pattern**:
```go
// OLD
import "todo-app/internal/models"
session := models.AuthenticationSession{}

// NEW
import "todo-app/domain/auth/entities"
session := entities.AuthenticationSession{}
```

---

## 7. Domain Relationships Diagram

```
┌─────────────────┐
│   User Domain   │
│  (Established)  │
└────────┬────────┘
         │
         │ 1:N sessions
         │
┌────────▼────────┐
│   Auth Domain   │
│     (NEW)       │
│                 │
│ • Session       │◄──── OAuth flow creates session
│ • OAuthState    │
└─────────────────┘

┌─────────────────┐
│   Task Domain   │
│  (Established)  │
│                 │
│ (No relations   │
│  with auth)     │
└─────────────────┘

┌─────────────────┐
│  Health Domain  │
│     (NEW)       │
│                 │
│ (Stateless)     │
└─────────────────┘
```

---

## 8. Validation Summary

### Field-Level Validation
- All entities have validation methods
- Value objects enforce invariants in constructor
- Repositories validate entity state before persistence

### Business Rule Validation
- Session expiration checked before use
- OAuth state one-time use enforced
- PKCE challenge/verifier match verified

### Cross-Entity Validation
- Session.UserID must reference valid User (FK constraint)
- No other cross-domain constraints

---

## 9. API Compatibility

### JSON Serialization
All entities maintain existing JSON field names via struct tags:

```go
type AuthenticationSession struct {
    SessionID string `json:"session_id" gorm:"primaryKey"`
    UserID    string `json:"user_id" gorm:"index"`
    // ... (existing tags preserved)
}
```

**Result**: Frontend sees no API changes - JSON contracts are stable.

---

## 10. Testing Strategy

### Unit Tests
- Each entity has constructor validation tests
- Value objects have invariant tests
- Repository interfaces have mock implementations

### Integration Tests
- Database GORM mappings verified
- Table name preservation validated
- Foreign key constraints tested

### Contract Tests
- API JSON shape validated
- Existing 23 contract tests must pass unchanged

---

**Data model complete**. Ready for contract generation and quickstart guide.
