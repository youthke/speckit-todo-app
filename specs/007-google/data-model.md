# Data Model: Google Account Signup

**Feature**: 007-google | **Date**: 2025-10-01

## Entities

### 1. User (Extended)
**Description**: Existing user entity extended to support multiple authentication methods

**Fields**:
- `id` (INTEGER, PRIMARY KEY): Unique identifier
- `email` (VARCHAR(255), UNIQUE, NOT NULL): User's email address
- `password_hash` (VARCHAR(255), NULLABLE): Hashed password (null for OAuth-only users)
- `auth_method` (VARCHAR(50), NOT NULL): Authentication method ("password" | "google" | "hybrid")
- `created_at` (TIMESTAMP, DEFAULT NOW()): Account creation timestamp
- `updated_at` (TIMESTAMP, DEFAULT NOW()): Last update timestamp

**Changes for This Feature**:
- Add `auth_method` field to distinguish authentication type
- Make `password_hash` nullable to support OAuth-only users
- Ensure email remains unique constraint

**Validation Rules**:
- Email must be valid format (RFC 5322)
- `auth_method` must be one of allowed values
- If `auth_method = "password"`, `password_hash` must not be null
- If `auth_method = "google"`, `password_hash` may be null

**Relationships**:
- One-to-One with GoogleIdentity (optional, only if `auth_method` includes "google")
- One-to-Many with AuthenticationSession

---

### 2. GoogleIdentity (New)
**Description**: Links a User account to their Google account credentials

**Fields**:
- `id` (INTEGER, PRIMARY KEY): Unique identifier
- `user_id` (INTEGER, UNIQUE, NOT NULL, FOREIGN KEY): Reference to User
- `google_user_id` (VARCHAR(255), UNIQUE, NOT NULL): Google's unique user identifier (sub claim)
- `email` (VARCHAR(255), NOT NULL): Email address from Google
- `email_verified` (BOOLEAN, NOT NULL): Email verification status from Google
- `created_at` (TIMESTAMP, DEFAULT NOW()): Link creation timestamp
- `updated_at` (TIMESTAMP, DEFAULT NOW()): Last update timestamp

**Validation Rules**:
- `google_user_id` must be unique across all records (prevents duplicate Google accounts)
- `user_id` must be unique (one Google identity per user)
- `email` must match user's email
- `email_verified` must be true at creation (enforced by signup logic)

**Relationships**:
- One-to-One with User (required)

**Indexes**:
- `idx_google_user_id` on `google_user_id` (lookup during signup to detect duplicates)
- `idx_email` on `email` (faster email-based queries)

---

### 3. AuthenticationSession (Extended)
**Description**: Existing session entity extended to support 7-day expiration

**Fields**:
- `id` (INTEGER, PRIMARY KEY): Unique identifier
- `user_id` (INTEGER, NOT NULL, FOREIGN KEY): Reference to User
- `token_hash` (VARCHAR(255), UNIQUE, NOT NULL): Hashed session token (JWT hash)
- `expires_at` (TIMESTAMP, NOT NULL): Session expiration time
- `created_at` (TIMESTAMP, DEFAULT NOW()): Session creation timestamp
- `last_accessed_at` (TIMESTAMP, DEFAULT NOW()): Last access timestamp
- `ip_address` (VARCHAR(45), NULLABLE): IP address of session creation
- `user_agent` (VARCHAR(255), NULLABLE): User agent string

**Changes for This Feature**:
- Ensure `expires_at` calculation supports 7-day duration
- `expires_at = created_at + 7 days` for Google OAuth signup sessions

**Validation Rules**:
- `expires_at` must be in the future
- `token_hash` must be unique
- `user_id` must reference existing user

**Relationships**:
- Many-to-One with User

**State Transitions**:
- Active: `expires_at > NOW()`
- Expired: `expires_at <= NOW()`

---

## Entity Relationships Diagram

```
┌─────────────────┐
│     User        │
├─────────────────┤
│ id (PK)         │
│ email (UNIQUE)  │◄─┐
│ password_hash   │  │
│ auth_method     │  │  One-to-One
│ created_at      │  │
│ updated_at      │  │
└─────────────────┘  │
         │           │
         │ One-to-Many
         │           │
         ▼           │
┌─────────────────┐  │
│AuthSession      │  │
├─────────────────┤  │
│ id (PK)         │  │
│ user_id (FK)    │  │
│ token_hash      │  │
│ expires_at      │  │
│ created_at      │  │
└─────────────────┘  │
                     │
         ┌───────────┘
         │
         ▼
┌─────────────────────┐
│ GoogleIdentity      │
├─────────────────────┤
│ id (PK)             │
│ user_id (FK,UNIQUE) │
│ google_user_id      │
│ email               │
│ email_verified      │
│ created_at          │
│ updated_at          │
└─────────────────────┘
```

---

## Migration Strategy

### Migration File: `00X_add_google_oauth.sql`

**Up Migration**:
```sql
-- Step 1: Add auth_method column to users table
ALTER TABLE users ADD COLUMN auth_method VARCHAR(50) NOT NULL DEFAULT 'password';

-- Step 2: Make password_hash nullable (if not already)
-- SQLite doesn't support ALTER COLUMN, so recreate table if needed
-- For SQLite: This may require table recreation with data migration

-- Step 3: Create google_identities table
CREATE TABLE IF NOT EXISTS google_identities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    google_user_id VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Step 4: Create indexes
CREATE INDEX idx_google_user_id ON google_identities(google_user_id);
CREATE INDEX idx_google_email ON google_identities(email);

-- Step 5: Verify sessions table supports expires_at (should already exist)
-- No changes needed if sessions table already exists with expires_at field
```

**Down Migration**:
```sql
-- Rollback: Remove Google OAuth support
DROP INDEX IF EXISTS idx_google_email;
DROP INDEX IF EXISTS idx_google_user_id;
DROP TABLE IF EXISTS google_identities;
ALTER TABLE users DROP COLUMN auth_method;
-- Restore password_hash NOT NULL constraint if needed
```

---

## Data Flow for Google Signup

1. **User Initiates Signup**:
   - Frontend: User clicks "Sign up with Google"
   - No data persisted yet

2. **OAuth Flow**:
   - Backend generates state token (temporary, not persisted)
   - Redirects to Google

3. **Google Callback**:
   - Backend receives authorization code
   - Exchanges code for ID token
   - Extracts claims: `sub` (google_user_id), `email`, `email_verified`

4. **Duplicate Check**:
   - Query: `SELECT * FROM google_identities WHERE google_user_id = ?`
   - If found: Redirect to login (no data changes)
   - If not found: Continue to user creation

5. **User Creation** (Transaction):
   ```sql
   BEGIN TRANSACTION;

   -- Create user
   INSERT INTO users (email, auth_method, password_hash)
   VALUES (?, 'google', NULL);

   -- Get user_id
   SET @user_id = LAST_INSERT_ID();

   -- Create Google identity link
   INSERT INTO google_identities (user_id, google_user_id, email, email_verified)
   VALUES (@user_id, ?, ?, true);

   COMMIT;
   ```

6. **Session Creation**:
   - Generate JWT token with claims: `user_id`, `google_user_id`, `exp` (7 days)
   - Hash token for storage
   - Insert into `authentication_sessions`:
     ```sql
     INSERT INTO authentication_sessions
       (user_id, token_hash, expires_at, ip_address, user_agent)
     VALUES (?, ?, NOW() + INTERVAL 7 DAY, ?, ?);
     ```

7. **Response**:
   - Set HTTP-only cookie with JWT
   - Return user data to frontend

---

## Validation Summary

| Requirement | Validation | Enforcement Point |
|-------------|------------|-------------------|
| FR-004: Extract email | `email` field in GoogleIdentity | OAuth callback handler |
| FR-005: Link Google ID | `google_user_id` field | OAuth callback handler |
| FR-007: Email verified | `email_verified = true` | Pre-creation validation |
| FR-008: No duplicates | UNIQUE constraint on `google_user_id` | Database + pre-check |
| FR-010: 7-day session | `expires_at = created_at + 7 days` | Session creation |

---

## Performance Considerations

- **Index on `google_user_id`**: Enables O(log n) duplicate lookup (FR-008)
- **Index on `email`**: Faster user lookups by email
- **Unique constraint**: Database-level duplicate prevention
- **Transaction for signup**: Ensures atomicity (user + google_identity created together)
- **Session expiration index**: Consider adding index on `expires_at` for cleanup queries

---

## Security Considerations

- **Password hash nullable**: Prevents account takeover via password reset for OAuth-only users
- **Email verification enforced**: `email_verified` must be true, validated before creation
- **Foreign key with CASCADE**: Deleting user automatically removes Google identity
- **Unique constraints**: Prevents duplicate Google accounts and duplicate links
