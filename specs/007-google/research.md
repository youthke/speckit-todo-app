# Research: Google Account Signup

**Feature**: 007-google | **Date**: 2025-10-01

## Research Questions

### 1. Google OAuth 2.0 Library for Go
**Decision**: Use `golang.org/x/oauth2/google` package

**Rationale**:
- Official Google-maintained OAuth 2.0 library for Go
- Well-documented and actively maintained
- Handles token management, refresh tokens, and OAuth flow automatically
- Native integration with Google APIs
- Production-ready with extensive testing

**Alternatives Considered**:
- `github.com/coreos/go-oidc`: OpenID Connect library, more generic but requires more boilerplate for Google-specific flows
- Manual OAuth implementation: Not recommended due to security complexity and maintenance burden

**Implementation Pattern**:
```go
import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

// OAuth2 Config
config := &oauth2.Config{
    ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
    RedirectURL:  "http://localhost:8080/auth/google/callback",
    Scopes: []string{
        "https://www.googleapis.com/auth/userinfo.email",
        "https://www.googleapis.com/auth/userinfo.profile",
    },
    Endpoint: google.Endpoint,
}
```

---

### 2. Session Management with 7-Day Expiration
**Decision**: Use JWT tokens with 7-day expiration stored in HTTP-only cookies

**Rationale**:
- JWT provides stateless authentication, reducing database load
- HTTP-only cookies prevent XSS attacks
- 7-day expiration matches requirement (FR-010)
- Can include user ID and Google account link in token claims
- Supports secure flag for HTTPS-only transmission

**Alternatives Considered**:
- Server-side sessions in database: Higher database load, requires cleanup job
- Local storage: Vulnerable to XSS attacks
- In-memory sessions: Not scalable, lost on server restart

**Implementation Pattern**:
- Use `github.com/golang-jwt/jwt/v5` for token generation
- Store in cookie with `HttpOnly`, `Secure`, `SameSite=Strict` flags
- Set `MaxAge` to 7 days (604800 seconds)

---

### 3. Database Schema Extension
**Decision**: Add `google_identities` table with foreign key to `users` table

**Rationale**:
- Maintains data normalization
- Allows future support for multiple OAuth providers
- Separates authentication method from user core data
- Supports one-to-one relationship (one Google account per user)
- Enables easy lookup for duplicate prevention (FR-008)

**Schema Design**:
```sql
CREATE TABLE google_identities (
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
CREATE INDEX idx_email ON google_identities(email);
```

**Alternatives Considered**:
- Single users table with nullable OAuth columns: Poor normalization, hard to extend for multiple providers
- Polymorphic association: Overly complex for current needs

---

### 4. Frontend OAuth Flow Pattern
**Decision**: Server-initiated redirect flow with callback page

**Rationale**:
- Standard OAuth 2.0 authorization code flow
- Backend controls redirect URL and state parameter for CSRF protection
- Frontend handles initial trigger and callback processing
- Maintains security by keeping client secret on server
- User-friendly: browser handles redirect automatically

**Flow**:
1. User clicks "Sign up with Google" → Frontend navigates to `/api/auth/google/login`
2. Backend generates state token, stores in session, redirects to Google
3. Google authenticates user, redirects to `/auth/google/callback`
4. Frontend callback page extracts code, sends to backend
5. Backend exchanges code for tokens, creates user, returns session cookie
6. Frontend redirects to app home page

**Alternatives Considered**:
- Client-side flow with implicit grant: Deprecated by Google, less secure
- Popup-based flow: Poor UX, blocked by popup blockers

---

### 5. Email Verification Check
**Decision**: Validate `email_verified` claim from Google ID token

**Rationale**:
- Google provides `email_verified` field in ID token claims
- No additional API calls required
- Reliable source of truth from Google
- Meets requirement FR-007

**Implementation**:
```go
// After exchanging code for token
idToken, err := verifier.Verify(ctx, rawIDToken)
if err != nil {
    return err
}

var claims struct {
    Email         string `json:"email"`
    EmailVerified bool   `json:"email_verified"`
}

if err := idToken.Claims(&claims); err != nil {
    return err
}

if !claims.EmailVerified {
    return errors.New("email not verified")
}
```

**Alternatives Considered**:
- Skip verification check: Violates FR-007 requirement
- Send verification email ourselves: Redundant, adds complexity

---

### 6. Error Handling Strategy
**Decision**: Generic error messages for user-facing errors, detailed logging for debugging

**Rationale**:
- Meets requirement FR-006: generic "Authentication failed" message
- Prevents information leakage to potential attackers
- Internal logging provides debugging capability
- Consistent error UX across all failure scenarios

**Error Categories**:
- OAuth errors (denied permission, invalid code): "Authentication failed"
- Network errors: "Authentication failed"
- Unverified email: "Authentication failed"
- Duplicate account: Redirect to login (not an error, per FR-008)

**Alternatives Considered**:
- Specific error messages: Violates requirement, potential security risk
- No logging: Makes debugging impossible

---

### 7. Testing Strategy
**Decision**: Three-layer testing approach

**Contract Tests**:
- Test API endpoint contracts (request/response schemas)
- Mock Google OAuth responses
- Verify error handling paths

**Integration Tests**:
- End-to-end signup flow with mock Google OAuth server
- Duplicate signup detection
- Session creation and validation
- Email verification rejection

**Unit Tests**:
- OAuth service logic (token exchange, user creation)
- Session management utilities
- Database operations

**Tools**:
- Backend: `testify` for assertions, `httptest` for HTTP testing
- Frontend: Vitest + Testing Library for component tests
- Mock OAuth server: Custom test server or `httpmock` library

---

## Summary

All technical unknowns have been resolved with concrete decisions:

1. **OAuth Library**: `golang.org/x/oauth2/google`
2. **Sessions**: JWT with HTTP-only cookies, 7-day expiration
3. **Database**: Separate `google_identities` table with foreign key
4. **Frontend Flow**: Server-initiated redirect with callback page
5. **Email Verification**: Validate from Google ID token claims
6. **Error Handling**: Generic messages, detailed logging
7. **Testing**: Contract + Integration + Unit tests

All decisions align with functional requirements (FR-001 through FR-010) and clarified constraints.
