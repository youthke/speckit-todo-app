# Tasks: Google Account Signup

**Feature**: 007-google | **Date**: 2025-10-01
**Input**: Design documents from `/specs/007-google/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Summary

This task list implements Google OAuth 2.0 signup functionality following Test-Driven Development. The implementation spans backend (Go/Gin) and frontend (React/TypeScript) with strict ordering: Setup → Database → Tests → Implementation → Validation.

**Total Tasks**: 35
**Estimated Time**: 12-16 hours
**Parallel Opportunities**: 18 tasks marked [P]

---

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- All file paths are absolute from repository root

---

## Phase 3.1: Setup & Dependencies

### T001: Install Go OAuth dependencies ✅
**Files**: `backend/go.mod`, `backend/go.sum`
**Action**:
```bash
cd backend
go get golang.org/x/oauth2
go get golang.org/x/oauth2/google
go get github.com/golang-jwt/jwt/v5
```
**Acceptance**: Dependencies added to go.mod

---

### T002: Configure environment variables ✅
**Files**: `backend/.env.example` (create), `backend/.env` (local)
**Action**:
1. Create `.env.example`:
```
GOOGLE_CLIENT_ID=your_client_id_here
GOOGLE_CLIENT_SECRET=your_client_secret_here
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
JWT_SECRET=your_jwt_secret_here
```
2. Copy to `.env` and fill with actual values (developer must obtain from Google Console)

**Acceptance**: Environment variables documented and loadable

---

### T003 [P]: Add frontend dependencies ✅
**Files**: `frontend/package.json`
**Action**:
```bash
cd frontend
# No new dependencies needed - using existing axios, react-router-dom
```
**Acceptance**: Verify existing dependencies are sufficient (already installed)

---

## Phase 3.2: Database Migration

### T004: Create database migration for Google OAuth ✅
**Files**: `backend/migrations/007_add_google_oauth.sql` (create)
**Action**: Create migration file with:
```sql
-- Up Migration
ALTER TABLE users ADD COLUMN auth_method VARCHAR(50) NOT NULL DEFAULT 'password';

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

CREATE INDEX idx_google_user_id ON google_identities(google_user_id);
CREATE INDEX idx_google_email ON google_identities(email);

-- Down Migration (commented)
-- DROP INDEX IF EXISTS idx_google_email;
-- DROP INDEX IF EXISTS idx_google_user_id;
-- DROP TABLE IF EXISTS google_identities;
-- ALTER TABLE users DROP COLUMN auth_method;
```
**Acceptance**: Migration file created and syntax-valid

**Dependencies**: None (blocks T005-T035)

---

### T005: Run database migration ✅
**Files**: Database file `backend/todo.db`
**Action**: Execute migration (method depends on project's migration tool)
```bash
cd backend
# Apply migration using your migration runner
# Example: go run cmd/migrate/main.go up
```
**Acceptance**:
- `google_identities` table exists
- `users.auth_method` column exists
- Indexes created

**Dependencies**: T004

---

## Phase 3.3: Backend Models

### T006 [P]: Create GoogleIdentity model ✅
**Files**: `backend/models/google_identity.go` (create)
**Action**: Implement GORM model:
```go
package models

import "time"

type GoogleIdentity struct {
    ID             uint      `gorm:"primaryKey"`
    UserID         uint      `gorm:"uniqueIndex;not null"`
    GoogleUserID   string    `gorm:"uniqueIndex;size:255;not null"`
    Email          string    `gorm:"size:255;not null"`
    EmailVerified  bool      `gorm:"not null;default:false"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
    User           User      `gorm:"foreignKey:UserID"`
}
```
**Acceptance**: Model compiles, GORM tags correct

**Dependencies**: T005

---

### T007 [P]: Extend User model for OAuth ✅
**Files**: `backend/models/user.go` (modify)
**Action**: Add `AuthMethod` field and relationship:
```go
type User struct {
    // ... existing fields ...
    AuthMethod     string          `gorm:"size:50;not null;default:'password'"` // "password" | "google" | "hybrid"
    GoogleIdentity *GoogleIdentity `gorm:"foreignKey:UserID"`
}
```
**Acceptance**: User model compiles with new field

**Dependencies**: T005

---

## Phase 3.4: Backend Configuration

### T008: Create OAuth config helper ✅
**Files**: `backend/config/oauth.go` (create)
**Action**: Implement OAuth2 config loader:
```go
package config

import (
    "os"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

func GetGoogleOAuthConfig() *oauth2.Config {
    return &oauth2.Config{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
        Endpoint: google.Endpoint,
    }
}
```
**Acceptance**: Config loads from environment variables

**Dependencies**: T002

---

## Phase 3.5: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE IMPLEMENTATION

**CRITICAL**: All tests in this phase MUST be written and MUST FAIL before proceeding to Phase 3.6

---

### T009 [P]: Contract test for /api/auth/google/login ✅
**Files**: `backend/tests/contract/google_oauth_login_test.go` (create)
**Action**: Test GET /api/auth/google/login:
```go
func TestGoogleOAuthLogin(t *testing.T) {
    // Setup test server
    // Send GET to /api/auth/google/login
    // Assert: 302 redirect
    // Assert: Location header contains accounts.google.com
    // Assert: oauth_state cookie is set
}
```
**Acceptance**: Test written, currently FAILS (endpoint not implemented)

**Dependencies**: T005

---

### T010 [P]: Contract test for /api/auth/google/callback ✅
**Files**: `backend/tests/contract/google_oauth_callback_test.go` (create)
**Action**: Test GET /api/auth/google/callback with mock:
```go
func TestGoogleOAuthCallback_Success(t *testing.T) {
    // Mock Google OAuth token exchange
    // Send GET with code & state params
    // Assert: 302 redirect to frontend
    // Assert: session_token cookie set with 7-day expiration
}

func TestGoogleOAuthCallback_InvalidState(t *testing.T) {
    // Send GET with invalid state
    // Assert: 400 error with "Authentication failed"
}

func TestGoogleOAuthCallback_UnverifiedEmail(t *testing.T) {
    // Mock Google response with email_verified=false
    // Assert: Redirect to signup with error
}
```
**Acceptance**: Tests written, currently FAIL (endpoint not implemented)

**Dependencies**: T005

---

### T011 [P]: Contract test for /api/auth/me ✅ with Google OAuth
**Files**: `backend/tests/contract/google_oauth_me_test.go` (create)
**Action**: Test GET /api/auth/me returns Google user:
```go
func TestAuthMe_GoogleUser(t *testing.T) {
    // Create test user with auth_method="google"
    // Create valid session token
    // Send GET /api/auth/me with cookie
    // Assert: 200 OK
    // Assert: Response contains user with auth_method="google"
}
```
**Acceptance**: Test written, currently FAILS (Google user creation not implemented)

**Dependencies**: T006, T007

---

### T012 [P]: Integration test - Successful new user signup ✅
**Files**: `backend/tests/integration/google_signup_success_test.go` (create)
**Action**: Test complete OAuth flow (Scenario 1 from quickstart.md):
```go
func TestGoogleSignup_NewUser_Success(t *testing.T) {
    // Mock Google OAuth server
    // Initiate login flow
    // Simulate Google callback with valid verified email
    // Assert: User created in database
    // Assert: GoogleIdentity created
    // Assert: Session created with 7-day expiration
    // Assert: User redirected to frontend home
}
```
**Acceptance**: Test written, currently FAILS (full flow not implemented)

**Dependencies**: T006, T007

---

### T013 [P]: Integration test - Duplicate signup prevention ✅
**Files**: `backend/tests/integration/google_signup_duplicate_test.go` (create)
**Action**: Test duplicate detection (Scenario 2 from quickstart.md):
```go
func TestGoogleSignup_DuplicateUser_RedirectsToLogin(t *testing.T) {
    // Create existing user with Google identity
    // Attempt signup with same google_user_id
    // Assert: No new user created
    // Assert: Redirected to /login (not /home)
}
```
**Acceptance**: Test written, currently FAILS (duplicate detection not implemented)

**Dependencies**: T006, T007

---

### T014 [P]: Integration test - Unverified email rejection ✅
**Files**: `backend/tests/integration/google_signup_unverified_test.go` (create)
**Action**: Test email verification check (Scenario 3 from quickstart.md):
```go
func TestGoogleSignup_UnverifiedEmail_Rejected(t *testing.T) {
    // Mock Google response with email_verified=false
    // Attempt signup
    // Assert: No user created
    // Assert: Redirected to signup with error
    // Assert: Error message is "Authentication failed"
}
```
**Acceptance**: Test written, currently FAILS (email verification check not implemented)

**Dependencies**: T006, T007

---

### T015 [P]: Integration test - OAuth error handling ✅
**Files**: `backend/tests/integration/google_signup_error_test.go` (create)
**Action**: Test error scenarios (Scenario 4 from quickstart.md):
```go
func TestGoogleSignup_OAuthDenied_ShowsError(t *testing.T) {
    // Simulate user denying permission
    // Assert: Redirected to signup with error
    // Assert: Generic error message shown
}

func TestGoogleSignup_NetworkError_ShowsError(t *testing.T) {
    // Simulate network timeout
    // Assert: Generic "Authentication failed" message
}
```
**Acceptance**: Tests written, currently FAIL (error handling not implemented)

**Dependencies**: T006, T007

---

### T016 [P]: Integration test - Session expiration ✅
**Files**: `backend/tests/integration/google_session_test.go` (create)
**Action**: Test 7-day session duration (Scenario 5 from quickstart.md):
```go
func TestGoogleSignup_SessionDuration_SevenDays(t *testing.T) {
    // Complete signup flow
    // Extract session cookie
    // Assert: MaxAge = 604800 (7 days in seconds)
    // Query database for session
    // Assert: expires_at = created_at + 7 days
}
```
**Acceptance**: Test written, currently FAILS (session creation not implemented)

**Dependencies**: T006, T007

---

## Phase 3.6: Backend Service Layer (ONLY after tests are failing)

### T017: Create Google OAuth service ✅
**Files**: `backend/services/google_oauth_service.go` (create)
**Action**: Implement OAuth service with methods:
```go
package services

type GoogleOAuthService struct {
    config *oauth2.Config
    db     *gorm.DB
}

// GenerateAuthURL creates OAuth URL with state token
func (s *GoogleOAuthService) GenerateAuthURL(state string) string { ... }

// ExchangeCode exchanges authorization code for user info
func (s *GoogleOAuthService) ExchangeCode(code string) (*GoogleUserInfo, error) { ... }

// CreateUserFromGoogle creates user and Google identity
func (s *GoogleOAuthService) CreateUserFromGoogle(info *GoogleUserInfo) (*models.User, error) { ... }

// FindUserByGoogleID checks for existing Google account
func (s *GoogleOAuthService) FindUserByGoogleID(googleUserID string) (*models.User, error) { ... }
```
**Acceptance**:
- Service compiles
- All methods implemented
- Email verification check enforced in CreateUserFromGoogle

**Dependencies**: T006, T007, T008, T009-T016 (tests must exist and fail)

---

### T018: Create session service for JWT ✅
**Files**: `backend/services/session_service.go` (create or extend)
**Action**: Implement JWT session creation:
```go
type SessionService struct {
    jwtSecret string
}

// CreateSession generates JWT with 7-day expiration
func (s *SessionService) CreateSession(userID uint) (string, error) {
    expiresAt := time.Now().Add(7 * 24 * time.Hour)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     expiresAt.Unix(),
    })
    return token.SignedString([]byte(s.jwtSecret))
}

// ValidateSession verifies JWT token
func (s *SessionService) ValidateSession(tokenString string) (uint, error) { ... }
```
**Acceptance**:
- Service creates JWT with 7-day expiration
- Tokens are valid and verifiable

**Dependencies**: T002, T009-T016 (tests must exist)

---

## Phase 3.7: Backend HTTP Handlers

### T019: Implement /api/auth/google/login handler ✅
**Files**: `backend/handlers/auth_handler.go` (extend)
**Action**: Add GoogleLogin handler:
```go
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
    // Generate random state token
    state := generateRandomState()
    // Store state in session cookie (10 min expiration)
    c.SetCookie("oauth_state", state, 600, "/", "", false, true)
    // Generate OAuth URL
    url := h.oauthService.GenerateAuthURL(state)
    // Redirect to Google
    c.Redirect(http.StatusFound, url)
}
```
**Acceptance**:
- Handler redirects to Google OAuth
- State cookie set correctly
- T009 test passes

**Dependencies**: T017

---

### T020: Implement /api/auth/google/callback handler ✅
**Files**: `backend/handlers/auth_handler.go` (extend)
**Action**: Add GoogleCallback handler:
```go
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
    // Validate state parameter
    code := c.Query("code")
    state := c.Query("state")
    savedState, _ := c.Cookie("oauth_state")
    if state != savedState {
        c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
        return
    }

    // Exchange code for user info
    userInfo, err := h.oauthService.ExchangeCode(code)
    if err != nil || !userInfo.EmailVerified {
        c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
        return
    }

    // Check for duplicate
    existingUser, _ := h.oauthService.FindUserByGoogleID(userInfo.GoogleUserID)
    if existingUser != nil {
        c.Redirect(http.StatusFound, "http://localhost:3000/login")
        return
    }

    // Create user
    user, err := h.oauthService.CreateUserFromGoogle(userInfo)
    if err != nil {
        c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
        return
    }

    // Create session
    token, _ := h.sessionService.CreateSession(user.ID)
    c.SetCookie("session_token", token, 604800, "/", "", false, true) // 7 days

    // Redirect to home
    c.Redirect(http.StatusFound, "http://localhost:3000/")
}
```
**Acceptance**:
- Handler implements all logic from contract
- All error cases return generic "Authentication failed"
- Duplicate users redirected to login
- T010, T012-T016 tests pass

**Dependencies**: T017, T018

---

### T021: Register Google OAuth routes ✅
**Files**: `backend/cmd/main.go` or `backend/routes/routes.go` (extend)
**Action**: Add routes:
```go
authHandler := handlers.NewAuthHandler(oauthService, sessionService)
router.GET("/api/auth/google/login", authHandler.GoogleLogin)
router.GET("/api/auth/google/callback", authHandler.GoogleCallback)
```
**Acceptance**: Routes accessible, return expected responses

**Dependencies**: T019, T020

---

## Phase 3.8: Frontend Components

### T022 [P]: Create GoogleSignupButton component ✅
**Files**: `frontend/src/components/GoogleSignupButton.tsx` (create)
**Action**: Implement button component:
```tsx
import React from 'react';

const GoogleSignupButton: React.FC = () => {
  const handleClick = () => {
    window.location.href = 'http://localhost:8080/api/auth/google/login';
  };

  return (
    <button onClick={handleClick} className="google-signup-btn">
      <img src="/google-icon.svg" alt="Google" />
      Sign up with Google
    </button>
  );
};

export default GoogleSignupButton;
```
**Acceptance**: Component renders button, onClick navigates to backend

**Dependencies**: T021 (backend must be ready for manual testing)

---

### T023 [P]: Create OAuthCallbackPage ✅
**Files**: `frontend/src/pages/OAuthCallbackPage.tsx` (create)
**Action**: Handle OAuth redirect:
```tsx
import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

const OAuthCallbackPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    const error = searchParams.get('error');
    if (error) {
      navigate('/signup?error=authentication_failed');
    } else {
      // Cookie already set by backend, just redirect
      navigate('/');
    }
  }, [searchParams, navigate]);

  return <div>Processing authentication...</div>;
};

export default OAuthCallbackPage;
```
**Acceptance**: Page handles both success and error cases

**Dependencies**: None (independent frontend work)

---

### T024: Integrate GoogleSignupButton into SignupPage ⚠️
**Files**: `frontend/src/pages/SignupPage.tsx` (modify)
**Action**: Add Google signup option:
```tsx
import GoogleSignupButton from '../components/GoogleSignupButton';

// Inside SignupPage component
<div className="signup-options">
  <GoogleSignupButton />
  <div className="separator">OR</div>
  {/* Existing email/password signup form */}
</div>
```
**Acceptance**: Button visible on signup page

**Dependencies**: T022

---

### T025: Add OAuth callback route ⚠️
**Files**: `frontend/src/App.tsx` or routing config (modify)
**Action**: Add route for callback:
```tsx
<Route path="/auth/google/callback" element={<OAuthCallbackPage />} />
```
**Acceptance**: Route registered and accessible

**Dependencies**: T023

---

### T026 [P]: Create GoogleSignupButton test ⚠️
**Files**: `frontend/tests/components/GoogleSignupButton.test.tsx` (create)
**Action**: Test button behavior:
```tsx
import { render, screen, fireEvent } from '@testing-library/react';
import GoogleSignupButton from '../../src/components/GoogleSignupButton';

test('renders button with correct text', () => {
  render(<GoogleSignupButton />);
  expect(screen.getByText(/Sign up with Google/i)).toBeInTheDocument();
});

test('navigates to OAuth endpoint on click', () => {
  delete window.location;
  window.location = { href: '' } as any;

  render(<GoogleSignupButton />);
  fireEvent.click(screen.getByText(/Sign up with Google/i));

  expect(window.location.href).toBe('http://localhost:8080/api/auth/google/login');
});
```
**Acceptance**: Tests pass

**Dependencies**: T022

---

## Phase 3.9: Error Handling & Validation

### T027: Add error display on signup page ⚠️
**Files**: `frontend/src/pages/SignupPage.tsx` (modify)
**Action**: Show error message from query params:
```tsx
const [searchParams] = useSearchParams();
const error = searchParams.get('error');

return (
  <div>
    {error === 'authentication_failed' && (
      <div className="error-message">Authentication failed</div>
    )}
    {/* Rest of signup page */}
  </div>
);
```
**Acceptance**: Error message displays when error param present

**Dependencies**: T024

---

### T028: Add backend error logging ✅
**Files**: `backend/services/google_oauth_service.go`, `backend/handlers/auth_handler.go` (modify)
**Action**: Add detailed logging for all error cases:
```go
import "log"

// In error handlers
if err != nil {
    log.Printf("OAuth error: %v", err) // Internal logging
    // Return generic message to user
}
```
**Acceptance**: Errors logged to console with details

**Dependencies**: T017, T020

---

## Phase 3.10: Polish & Documentation

### T029 [P]: Add unit tests for OAuth config
**Files**: `backend/tests/unit/oauth_config_test.go` (create)
**Action**: Test config loading:
```go
func TestOAuthConfig_LoadsFromEnv(t *testing.T) {
    os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
    config := GetGoogleOAuthConfig()
    assert.Equal(t, "test-client-id", config.ClientID)
}
```
**Acceptance**: Tests pass

**Dependencies**: T008

---

### T030 [P]: Add unit tests for session service
**Files**: `backend/tests/unit/session_service_test.go` (create)
**Action**: Test JWT creation and validation:
```go
func TestSessionService_CreateSession_SevenDayExpiration(t *testing.T) {
    service := NewSessionService("test-secret")
    token, err := service.CreateSession(1)
    assert.NoError(t, err)

    // Parse token and check expiration
    claims := parseToken(token)
    assert.InDelta(t, time.Now().Add(7*24*time.Hour).Unix(), claims.Exp, 5)
}
```
**Acceptance**: Tests pass

**Dependencies**: T018

---

### T031: Run all backend tests
**Files**: All backend test files
**Action**: Execute test suite:
```bash
cd backend
go test ./tests/contract/... -v
go test ./tests/integration/... -v
go test ./tests/unit/... -v
```
**Acceptance**: All tests pass (green)

**Dependencies**: T009-T016, T029-T030

---

### T032: Run all frontend tests
**Files**: All frontend test files
**Action**: Execute test suite:
```bash
cd frontend
npm test
```
**Acceptance**: All tests pass

**Dependencies**: T026

---

### T033: Execute quickstart scenarios
**Files**: `specs/007-google/quickstart.md`
**Action**: Manually execute all 6 scenarios from quickstart:
1. Successful new user signup
2. Duplicate signup prevention
3. Unverified email rejection (mocked)
4. OAuth error handling
5. Session expiration validation
6. Cancelled OAuth flow

**Acceptance**: All scenarios pass with expected behavior

**Dependencies**: T031, T032 (all automated tests pass first)

---

### T034: Performance validation
**Files**: Browser DevTools, backend logs
**Action**: Measure performance:
1. OAuth flow completion time: Must be <3 seconds
2. Session validation (/api/auth/me): Must be <50ms
3. Initial redirect: Must be <500ms

**Acceptance**: All performance targets met

**Dependencies**: T033

---

### T035: Update feature documentation
**Files**: `specs/007-google/plan.md` (update Progress Tracking)
**Action**: Mark implementation phases complete:
```markdown
**Phase Status**:
- [x] Phase 0: Research complete
- [x] Phase 1: Design complete
- [x] Phase 2: Task planning complete
- [x] Phase 3: Tasks generated
- [x] Phase 4: Implementation complete
- [x] Phase 5: Validation passed
```

**Acceptance**: Documentation reflects completed state

**Dependencies**: T034

---

## Dependencies Summary

```
Setup (T001-T003)
    ↓
Database (T004-T005)
    ↓
Models (T006-T007) [P]
    ↓
Config (T008) + Tests (T009-T016) [P]
    ↓
Services (T017-T018)
    ↓
Handlers (T019-T021)
    ↓
Frontend (T022-T027) [P] + Error Handling (T028)
    ↓
Unit Tests (T029-T030) [P]
    ↓
Test Execution (T031-T032)
    ↓
Manual Validation (T033-T034)
    ↓
Documentation (T035)
```

---

## Parallel Execution Examples

### Batch 1: Install Dependencies (After T001)
```bash
# Can run simultaneously
Task T002: Configure environment variables
Task T003: Add frontend dependencies
```

### Batch 2: Create Models (After T005)
```bash
# Can run simultaneously - different files
Task T006: Create GoogleIdentity model
Task T007: Extend User model
```

### Batch 3: Write Tests (After T007)
```bash
# Can run simultaneously - all independent test files
Task T009: Contract test for /api/auth/google/login
Task T010: Contract test for /api/auth/google/callback
Task T011: Contract test for /api/auth/me
Task T012: Integration test - New user signup
Task T013: Integration test - Duplicate prevention
Task T014: Integration test - Unverified email
Task T015: Integration test - Error handling
Task T016: Integration test - Session expiration
```

### Batch 4: Frontend Components (After T021)
```bash
# Can run simultaneously - different files
Task T022: Create GoogleSignupButton
Task T023: Create OAuthCallbackPage
Task T026: GoogleSignupButton test
```

### Batch 5: Unit Tests (After T028)
```bash
# Can run simultaneously - different files
Task T029: OAuth config unit tests
Task T030: Session service unit tests
```

---

## Validation Checklist

Before marking feature complete, verify:

- [x] All contracts have corresponding tests (T009-T011)
- [x] All entities have model tasks (T006-T007)
- [x] All tests written before implementation (T009-T016 before T017-T021)
- [x] Parallel tasks are truly independent (verified - all [P] tasks use different files)
- [x] Each task specifies exact file path (verified)
- [x] No task modifies same file as another [P] task (verified)
- [x] All functional requirements covered:
  - FR-001: Initiate Google signup ✓ (T022-T024)
  - FR-002: Redirect to Google ✓ (T019)
  - FR-003: Create user account ✓ (T017, T020)
  - FR-004: Extract email ✓ (T017)
  - FR-005: Link Google ID ✓ (T006, T017)
  - FR-006: Handle failures ✓ (T020, T028)
  - FR-007: Verify email ✓ (T017, T020)
  - FR-008: Prevent duplicates ✓ (T017, T020)
  - FR-009: Auto login ✓ (T020)
  - FR-010: 7-day session ✓ (T018, T020)

---

## Notes

- **TDD Enforcement**: Tasks T009-T016 MUST be completed and failing before starting T017
- **Environment Setup**: Developer must obtain Google OAuth credentials from Google Cloud Console
- **Local Testing**: Requires both backend (port 8080) and frontend (port 3000) running
- **Database**: SQLite file at `backend/todo.db` will be modified
- **Commit Strategy**: Commit after each task or logical group
- **Estimated Time**:
  - Setup & Database: 1-2 hours (T001-T005)
  - Models & Config: 1 hour (T006-T008)
  - Tests: 3-4 hours (T009-T016)
  - Backend Implementation: 3-4 hours (T017-T021)
  - Frontend: 2-3 hours (T022-T027)
  - Polish & Validation: 2-3 hours (T028-T035)

---

**Ready for Implementation**: Yes
**Next Step**: Execute T001 to begin implementation
