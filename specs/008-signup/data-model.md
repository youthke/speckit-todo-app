# Data Model: Signup Page

**Feature**: 008-signup
**Date**: 2025-10-01
**Purpose**: Document data entities and their relationships for Google OAuth signup

## Overview

The signup feature **does not introduce new entities**. It reuses existing entities from feature 007-google (Google OAuth implementation) with no schema changes required.

## Existing Entities

### User
**Table**: `users`
**Purpose**: Represents a registered user account in the system

**Fields**:
- `id` (INTEGER, PRIMARY KEY): Unique user identifier
- `email` (TEXT, UNIQUE, NOT NULL): User email address from Google
- `google_id` (TEXT, UNIQUE): Google account identifier (sub claim from JWT)
- `name` (TEXT, NULLABLE): User display name from Google
- `profile_picture_url` (TEXT, NULLABLE): URL to user's Google profile picture
- `created_at` (DATETIME, NOT NULL): Account creation timestamp
- `updated_at` (DATETIME, NOT NULL): Last modification timestamp
- `email_verified` (BOOLEAN, NOT NULL, DEFAULT FALSE): Email verification status

**Constraints**:
- Email must be unique across all users
- Google ID must be unique (one Google account = one user)
- Email is required (cannot be NULL or empty)

**Indexes**:
- `idx_users_email`: Index on email for fast lookup
- `idx_users_google_id`: Index on google_id for OAuth callback

**Source**: Existing table from features 001-todo, 005-google, 007-google

---

### Session
**Table**: `sessions` (or in-memory, depending on SessionService implementation)
**Purpose**: Manages user authentication sessions

**Fields**:
- `id` (INTEGER, PRIMARY KEY): Unique session identifier
- `user_id` (INTEGER, FOREIGN KEY → users.id): Reference to user
- `token` (TEXT, UNIQUE, NOT NULL): Session token (stored in HTTP-only cookie)
- `expires_at` (DATETIME, NOT NULL): Session expiration timestamp
- `created_at` (DATETIME, NOT NULL): Session creation timestamp

**Constraints**:
- Token must be unique and cryptographically random
- Session expires after 7 days (604800 seconds)
- One user can have multiple sessions (multi-device support)

**Source**: Existing from SessionService in feature 007-google

---

### GoogleIdentity (Transient)
**Type**: Struct (not persisted)
**Purpose**: Represents user information received from Google OAuth

**Fields**:
- `google_user_id` (string): Google account identifier (sub claim)
- `email` (string): User email address
- `name` (string): User display name
- `picture` (string): Profile picture URL
- `email_verified` (bool): Email verification status from Google

**Usage**:
- Received from Google OAuth token exchange
- Used to create or match User entity
- Not stored directly (fields mapped to User table)

**Source**: `backend/internal/models/google_identity.go`

## Data Flow

### New User Signup Flow
```
1. User clicks "Sign up with Google" on /signup page
2. Backend generates OAuth state token
3. User redirects to Google OAuth consent screen
4. User authorizes → Google redirects to /callback with code
5. Backend exchanges code for GoogleIdentity
6. Backend validates:
   - EmailVerified = true
   - Email is not empty
7. Backend checks if google_id exists in users table
   - NOT FOUND → Create new User record
   - Email from GoogleIdentity → users.email
   - GoogleUserID → users.google_id
   - Name → users.name (optional)
   - Picture → users.profile_picture_url (optional)
   - EmailVerified → users.email_verified
   - created_at, updated_at = now()
8. Backend creates Session:
   - user_id = new user ID
   - token = random 32-byte base64 string
   - expires_at = now() + 7 days
9. Backend sets session_token cookie (HttpOnly, 7 day expiration)
10. Backend redirects to http://localhost:3000/
```

### Existing User Signup Flow (Auto-Login)
```
1-6. Same as above
7. Backend checks if google_id exists in users table
   - FOUND → Use existing User record (no insert)
8. Backend creates NEW Session for existing user:
   - user_id = existing user ID
   - token = random 32-byte base64 string
   - expires_at = now() + 7 days
9. Backend sets session_token cookie (HttpOnly, 7 day expiration)
10. Backend redirects to http://localhost:3000/ (user is logged in)
```

**Key Difference**: Step 7 branches on existing user check
- **Old behavior**: Redirect to /login page (line 107 in handler)
- **New behavior**: Continue to session creation (auto-login)

## State Transitions

### User Entity State
```
[No Account]
    ↓ (Google OAuth + email verified)
[Account Created] → email, google_id, created_at set
    ↓ (Session created)
[Authenticated] → active session exists
    ↓ (Logout / session expires)
[Account Created] → can re-authenticate anytime
```

### Session State
```
[No Session]
    ↓ (User signs up/logs in)
[Active Session] → token valid, expires_at > now()
    ↓ (User activity)
[Active Session] → session refreshed (optional)
    ↓ (expires_at reached OR user logs out)
[Expired/Deleted] → user must re-authenticate
```

## Validation Rules

### Email Validation
1. **Source**: Google OAuth (email scope)
2. **Required**: YES - signup fails if email is empty
3. **Format**: Validated by Google (assumed valid if provided)
4. **Verification**: Must be verified by Google (EmailVerified = true)
5. **Uniqueness**: Must not exist in users table (duplicate check)

**Error Cases**:
- Email is empty → `authentication_failed` error
- Email not verified → `authentication_failed` error
- Email already exists with different google_id → potential conflict (not handled by spec)

### Google ID Validation
1. **Source**: Google OAuth (sub claim from ID token)
2. **Required**: YES - unique identifier for user
3. **Format**: String, typically numeric (e.g., "1234567890123456789")
4. **Uniqueness**: One Google account = one User record

**Duplicate Handling**:
- If google_id exists → auto-login (new behavior)
- If google_id doesn't exist → create new user

## Rate Limiting Data

**Type**: In-memory map (not persisted)
**Purpose**: Track signup attempts per IP address

**Structure**:
```go
type IPRateLimiter struct {
    ips map[string]*rate.Limiter  // IP → rate limiter
    mu  sync.RWMutex               // Thread-safe access
}
```

**Per-IP State**:
- Token bucket: 10 tokens, refill rate 0.0111/sec
- Bucket refills to 10 over 15 minutes
- No persistence (resets on server restart)

**Cleanup**:
- Periodic goroutine removes IPs inactive for 30+ minutes
- Prevents memory leak from IP accumulation

## Database Migrations

**No new migrations required** - feature uses existing schema.

**Relevant Existing Migrations**:
- `005_add_oauth_to_users.sql`: Added google_id, profile_picture_url to users
- `006_create_authentication_sessions.sql`: Created sessions table
- `008_add_google_oauth.sql`: Google OAuth entities (if applicable)

## Data Integrity

**Constraints Enforced**:
1. User email uniqueness (UNIQUE constraint)
2. User google_id uniqueness (UNIQUE constraint)
3. Session token uniqueness (UNIQUE constraint)
4. Foreign key: sessions.user_id → users.id (CASCADE on delete)

**Edge Cases**:
1. **Same email, different Google ID**:
   - Google accounts can have duplicate emails (e.g., gsuite vs gmail)
   - Current design: email UNIQUE constraint will fail
   - Recommendation: Consider composite key (email + auth_provider) in future

2. **Session token collision**:
   - Extremely unlikely (32-byte random = 2^256 space)
   - Database UNIQUE constraint will catch collision

3. **Multiple sessions per user**:
   - Allowed - user can log in from multiple devices
   - Old sessions expire after 7 days

## Summary

**New Entities**: None
**Modified Entities**: None (schema unchanged)
**Modified Behavior**: User lookup logic in signup flow
**Data Consistency**: Maintained by existing constraints
**Performance**: No new indexes or queries required

All data requirements satisfied by existing schema from features 005-google and 007-google.
