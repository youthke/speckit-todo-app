# Data Model: Google Account Login

**Feature**: Google OAuth 2.0 Authentication Integration
**Date**: 2025-09-29

## Core Entities

### User Account (Extended)

**Purpose**: Represents a user in the system with Google authentication support

**Fields**:
- `id` (Primary Key): Unique identifier
- `email` (Unique): User's email address (used for account linking)
- `name`: Display name from Google profile
- `password_hash` (Optional): For users with traditional auth (nullable for Google-only users)
- `google_id` (Optional, Unique): Google account identifier for OAuth users
- `oauth_provider` (Optional): Authentication provider ("google" or null)
- `oauth_created_at` (Optional): Timestamp when OAuth was first linked
- `created_at`: Account creation timestamp
- `updated_at`: Last modification timestamp
- `is_active`: Account status flag

**Validation Rules**:
- Email must be valid email format and unique across system
- Google ID must be unique when present
- Either password_hash OR google_id must be present (not both null)
- oauth_provider must be "google" when google_id is present
- Name must be present for active accounts

**State Transitions**:
- Traditional user → OAuth linked: google_id and oauth_provider added
- New OAuth user → Created with google_id, no password_hash
- Account deactivation → is_active set to false, sessions terminated

**Relationships**:
- One-to-Many with AuthenticationSession
- One-to-Many with TodoItems (existing)

### Authentication Session

**Purpose**: Tracks active user sessions with OAuth token management

**Fields**:
- `id` (Primary Key): Session identifier
- `user_id` (Foreign Key): Reference to User Account
- `session_token`: Unique session identifier (JWT)
- `refresh_token` (Optional): OAuth refresh token (encrypted)
- `access_token` (Optional): OAuth access token (encrypted)
- `token_expires_at` (Optional): OAuth token expiration time
- `session_expires_at`: Session expiration (24 hours from creation)
- `last_activity`: Last user activity timestamp
- `created_at`: Session creation timestamp
- `user_agent`: Browser/client information
- `ip_address`: Client IP address

**Validation Rules**:
- session_token must be unique and JWT format
- user_id must reference existing active user
- session_expires_at must be within 24 hours of creation
- OAuth tokens (refresh/access) must be present for OAuth sessions
- token_expires_at required when access_token present

**State Transitions**:
- Session creation → expires_at set to 24h, tokens stored
- Activity refresh → last_activity updated, extend if needed
- Token refresh → new access_token and expires_at
- Session termination → record deleted, tokens invalidated
- Forced logout (revocation) → immediate termination

**Relationships**:
- Many-to-One with User Account
- One-to-Many with SessionActivity (optional audit trail)

### OAuth State (Temporary)

**Purpose**: Temporary storage for OAuth flow state validation

**Fields**:
- `state_token` (Primary Key): Random state parameter
- `pkce_verifier`: PKCE code verifier for security
- `redirect_uri`: Post-auth redirect destination
- `created_at`: Creation timestamp
- `expires_at`: Expiration (5 minutes from creation)

**Validation Rules**:
- state_token must be cryptographically random (32+ chars)
- pkce_verifier must meet PKCE specification
- expires_at must be 5 minutes from creation
- redirect_uri must be whitelisted application URL

**State Transitions**:
- OAuth initiation → Created with 5-minute TTL
- OAuth callback → Validated and deleted
- Expiration → Automatically purged

**Relationships**:
- Standalone entity (no foreign keys)

## Database Schema Extensions

### Users Table Modifications
```sql
ALTER TABLE users ADD COLUMN google_id VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN oauth_provider VARCHAR(50);
ALTER TABLE users ADD COLUMN oauth_created_at TIMESTAMP;
ALTER TABLE users MODIFY COLUMN password_hash VARCHAR(255) NULL;

CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_oauth_provider ON users(oauth_provider);
```

### New Authentication Sessions Table
```sql
CREATE TABLE authentication_sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id INT NOT NULL,
    session_token TEXT NOT NULL UNIQUE,
    refresh_token TEXT,
    access_token TEXT,
    token_expires_at TIMESTAMP,
    session_expires_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_agent TEXT,
    ip_address VARCHAR(45),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_session_token (session_token),
    INDEX idx_user_sessions (user_id),
    INDEX idx_session_expires (session_expires_at)
);
```

### OAuth State Table
```sql
CREATE TABLE oauth_states (
    state_token VARCHAR(255) PRIMARY KEY,
    pkce_verifier VARCHAR(255) NOT NULL,
    redirect_uri VARCHAR(1000) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,

    INDEX idx_oauth_expires (expires_at)
);
```

## Data Flow Patterns

### Account Creation/Linking Flow
1. Google OAuth callback provides email + google_id
2. Query users by email address
3. If exists: Add google_id and oauth_provider to existing record
4. If not exists: Create new user with google_id, no password_hash
5. Create authentication_session record with OAuth tokens

### Session Management Flow
1. Store session_token as HTTP-only cookie
2. API requests validate session_token → user lookup
3. Check session_expires_at and token_expires_at
4. Auto-refresh if tokens expire but session valid
5. Update last_activity on each validated request

### Revocation Detection Flow
1. Google webhook or API call failure indicates revocation
2. Find authentication_sessions by access_token/refresh_token
3. Delete session records immediately
4. User must re-authenticate on next request

## Data Privacy & Security

### Encryption at Rest
- `refresh_token` and `access_token` encrypted using application key
- `session_token` is JWT, signed but not encrypted
- Sensitive OAuth data not logged or cached

### Data Retention
- OAuth state records auto-expire (5 minutes)
- Sessions auto-expire (24 hours)
- Inactive sessions purged daily
- User account data retained per existing policies

### Audit Considerations
- Track OAuth linking events in application logs
- Log session creation/termination with IP/user_agent
- Monitor for suspicious OAuth callback patterns
- Track token refresh frequency for anomaly detection