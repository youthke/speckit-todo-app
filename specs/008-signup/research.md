# Research: Signup Page Implementation

**Feature**: 008-signup
**Date**: 2025-10-01
**Purpose**: Document technical decisions and research findings for Google OAuth signup page

## Research Questions & Findings

### 1. Existing Google OAuth Implementation Review

**Question**: How does the current Google OAuth flow work and where should we modify it for auto-login?

**Current Implementation** (`backend/internal/handlers/google_oauth_handler.go`):
- Line 28-55: `GoogleLogin()` - Initiates OAuth flow, generates state, redirects to Google
- Line 57-140: `GoogleCallback()` - Handles Google callback
  - Lines 60-79: Validates state (CSRF protection) and handles OAuth errors
  - Lines 81-94: Exchanges code for user info and validates email verification
  - Lines 96-109: **KEY ISSUE** - Checks for existing user and redirects to `/login` (line 107)
  - Lines 111-126: Creates new user from Google info
  - Lines 119-140: Creates session and sets cookie, redirects to home

**Decision**: Modify `GoogleCallback` handler behavior at lines 104-109

**Implementation Approach**:
- **Current behavior** (lines 104-109):
  ```go
  if existingUser != nil {
      log.Printf("User already exists with Google ID: %s, redirecting to login", userInfo.GoogleUserID)
      c.Redirect(http.StatusFound, "http://localhost:3000/login")
      return
  }
  ```

- **New behavior** (replace return with continue to session creation):
  ```go
  if existingUser != nil {
      log.Printf("User already exists with Google ID: %s, auto-logging in", userInfo.GoogleUserID)
      user = existingUser  // Use existing user for session creation
      // Continue to session creation (lines 119-140)
  } else {
      // Create new user from Google info (lines 111-117)
      user, err := h.oauthService.CreateUserFromGoogle(userInfo)
      // ...
  }
  ```

**Rationale**:
- Minimal code change (remove early return, add else block)
- Reuses existing session creation logic (lines 119-140)
- Maintains CSRF protection and email verification
- No breaking changes to other authentication flows

**Alternatives Considered**:
1. **Create separate `/api/v1/auth/google/signup` endpoint**
   - Pros: Clean separation of signup vs login
   - Cons: Code duplication, need to duplicate OAuth flow logic, state management complexity
   - Rejected: Unnecessary complexity for minimal behavior difference

2. **Add query parameter to distinguish signup vs login**
   - Pros: Single endpoint, query param controls behavior
   - Cons: More complex state management, Google OAuth callback URL constraints
   - Rejected: OAuth callback URL must match exactly, query params complicate state validation

### 2. Rate Limiting Strategy

**Question**: What's the best approach for IP-based rate limiting in Go/Gin?

**Options Evaluated**:

1. **golang.org/x/time/rate (Token Bucket)**
   - Pros: Standard library, efficient, well-tested, in-memory (fast)
   - Cons: State lost on server restart, no distributed support
   - Use case: Small-medium scale, single server

2. **Redis-based rate limiter (e.g., go-redis/redis_rate)**
   - Pros: Distributed, persistent, shared across instances
   - Cons: External dependency, network latency, overkill for current scale
   - Use case: Large scale, multi-server deployments

3. **Gin middleware packages (e.g., gin-contrib/limiter)**
   - Pros: Easy integration, pre-built
   - Cons: Less flexible, may not support IP-based limiting cleanly
   - Use case: Quick prototyping

**Decision**: Use `golang.org/x/time/rate` with in-memory store per IP

**Rationale**:
- Current application is single-server (no load balancing mentioned)
- SQLite database indicates development/small-scale deployment
- Existing codebase has no Redis dependency
- In-memory is sufficient for rate limiting (state loss acceptable on restart)
- Token bucket algorithm handles bursts well

**Implementation**:
```go
package middleware

import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type IPRateLimiter struct {
    ips map[string]*rate.Limiter
    mu  sync.RWMutex
    r   rate.Limit  // requests per second
    b   int         // bucket size (burst)
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
    return &IPRateLimiter{
        ips: make(map[string]*rate.Limiter),
        r:   r,
        b:   b,
    }
}

// 10 attempts per 15 minutes = 10/(15*60) = 0.0111 req/sec
// Burst of 10 allows initial burst, then throttles
```

**Rate Limit Configuration**:
- **Limit**: 10 requests per 15 minutes per IP
- **Calculation**: 10 / (15 × 60) = 0.0111 requests/second
- **Burst**: 10 (allows 10 immediate requests, then throttle)
- **Response**: HTTP 429 with `Retry-After` header

**Cleanup Strategy**:
- Periodic cleanup goroutine to remove inactive IPs (prevent memory leak)
- Clean IPs with no requests in last 30 minutes

**Alternatives Considered**:
- **Redis**: Rejected due to added complexity and no existing Redis infrastructure
- **Database-backed counters**: Rejected due to performance overhead and SQLite limitations
- **Fixed window counters**: Rejected in favor of token bucket for better burst handling

### 3. Frontend Signup Page Design

**Question**: How should the signup page be structured and integrated?

**Existing Components Review**:
- `GoogleSignupButton.tsx`: Already exists, redirects to `/api/v1/auth/google/login`
- `frontend/src/pages/`: Contains other page components
- `frontend/src/routes/`: React Router configuration

**Decision**: Create `/signup` page mirroring existing page structure

**Page Structure**:
```tsx
// frontend/src/pages/SignupPage.tsx
import React from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import GoogleSignupButton from '../components/GoogleSignupButton';

export default function SignupPage() {
  const [searchParams] = useSearchParams();
  const error = searchParams.get('error');

  return (
    <div className="signup-page">
      <h1>Sign Up</h1>
      <GoogleSignupButton />
      {error === 'rate_limit_exceeded' && (
        <p className="error">Too many attempts. Please try again later.</p>
      )}
      {error === 'authentication_failed' && (
        <p className="error">Authentication failed. Please try again.</p>
      )}
      <p>
        Already have an account? <Link to="/login">Log in</Link>
      </p>
    </div>
  );
}
```

**Rationale**:
- Consistent with existing page patterns
- Uses existing `GoogleSignupButton` component (no changes needed)
- Error handling via URL query parameters (matches backend redirect pattern)
- Navigation to login page for existing users who need to switch

**React Router Integration**:
```tsx
// Add to frontend/src/routes/index.tsx (or wherever routes are defined)
<Route path="/signup" element={<SignupPage />} />
```

**Alternatives Considered**:
1. **Modal-based signup**
   - Pros: No navigation required
   - Cons: Doesn't fit existing navigation model, harder to deep-link
   - Rejected: Inconsistent with login page pattern

2. **Modify login page to be login/signup toggle**
   - Pros: Single page for both flows
   - Cons: Confusing UX, harder to maintain separate behaviors
   - Rejected: Less clear user intent, complicates URL structure

### 4. Email Requirement Validation

**Question**: Is email always provided by Google OAuth? How to handle missing email?

**Google OAuth Scopes Review**:
- Current scope: `openid email profile` (confirmed in existing code)
- Email is **required** scope, but users can have unverified emails
- Google may not provide email if user account has no email address (rare edge case)

**Existing Validation** (`google_oauth_handler.go:90-94`):
```go
// Validate email is verified
if !userInfo.EmailVerified {
    log.Printf("Email not verified for user: %s", userInfo.Email)
    c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
    return
}
```

**Decision**: Keep existing email verification check, no additional changes needed

**Rationale**:
- Existing code already validates `EmailVerified` flag
- Specification requires email (FR-002)
- Error handling already redirects with `authentication_failed`
- Edge case (no email) is already handled by existing validation

**Additional Safeguard** (optional):
Add explicit nil/empty email check before verification check:
```go
if userInfo.Email == "" {
    log.Printf("No email provided by Google for user: %s", userInfo.GoogleUserID)
    c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
    return
}
```

**Error Message Enhancement**:
Frontend could distinguish error types:
- `authentication_failed`: Generic error
- `email_required`: Specific error for missing email (if we add query param)

**Decision**: Use existing generic error handling, add specific email check in backend

## Technology Stack Confirmation

Based on codebase analysis:

**Backend**:
- Language: Go 1.24.0
- Framework: Gin web framework
- ORM: GORM
- Database: SQLite (development)
- OAuth: `golang.org/x/oauth2` library
- Testing: Go testing + testify assertions

**Frontend**:
- Language: TypeScript 5.9.2
- Framework: React 19.1.1
- Routing: React Router 6.30.1
- Build Tool: Vite 6.0.11
- HTTP Client: Axios 1.12.2
- Testing: React Testing Library

**Infrastructure**:
- Development: localhost:8080 (backend), localhost:3000 (frontend)
- Session Management: HTTP-only cookies, 7-day expiration
- CSRF Protection: State parameter in OAuth flow

## Performance Considerations

**OAuth Flow Latency**:
- Google OAuth redirect: ~500ms - 1s (network dependent)
- Code exchange: ~200-500ms (Google API call)
- Database operations: <50ms (SQLite, single user lookup + insert)
- Session creation: <10ms (in-memory + cookie)
- **Total estimated**: 1-2 seconds (within 3-second goal)

**Rate Limiter Performance**:
- Token bucket operations: O(1) time complexity
- Memory per IP: ~100 bytes (rate.Limiter struct)
- Concurrent access: sync.RWMutex for thread safety
- **Capacity**: 100+ req/s (far exceeds requirement)

**Bottlenecks**:
- Network latency to Google OAuth (external, unavoidable)
- SQLite write lock contention (unlikely at small scale)

## Security Considerations

**Existing Security Measures** (maintained):
- CSRF protection via state parameter
- HttpOnly cookies (prevents XSS)
- Email verification requirement
- Session expiration (7 days)

**New Security Measures**:
- Rate limiting (10 attempts per 15 min per IP)
- Prevents brute force / DoS on OAuth endpoint

**Potential Enhancements** (future consideration):
- HTTPS enforcement (currently localhost, need for production)
- Secure cookie flag (requires HTTPS)
- Session rotation on sensitive operations
- IP allowlist/blocklist for rate limiter

## Summary

All research questions resolved. Key decisions:
1. **Modify existing callback handler** - minimal code change, reuse logic
2. **golang.org/x/time/rate for rate limiting** - standard library, sufficient for scale
3. **Dedicated `/signup` page** - consistent UX, clear separation
4. **Keep existing email validation** - already handles requirements

No NEEDS CLARIFICATION remaining. Ready to proceed to Phase 1 (Design & Contracts).
