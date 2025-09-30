# Research: Google Account Login

**Feature**: Google OAuth 2.0 Authentication Integration
**Date**: 2025-09-29

## Technology Decisions

### OAuth 2.0 Library Selection (Go Backend)

**Decision**: Use `golang.org/x/oauth2` and `golang.org/x/oauth2/google`
**Rationale**:
- Official Go OAuth 2.0 implementation by Google
- Direct Google provider support
- Well-maintained with security updates
- Integrates naturally with Gin framework
**Alternatives considered**:
- Third-party OAuth libraries (rejected: unnecessary complexity)
- Custom OAuth implementation (rejected: security risks)

### Session Management Strategy

**Decision**: JWT tokens with 24-hour expiration and refresh token mechanism
**Rationale**:
- Stateless session validation
- Automatic refresh during active use
- Compatible with existing GORM user model
- Supports immediate revocation detection
**Alternatives considered**:
- Server-side sessions (rejected: scales poorly)
- Long-lived tokens (rejected: security concern)

### Google OAuth Scopes

**Decision**: Request `openid email profile` scopes only
**Rationale**:
- Minimal permission approach (privacy-first)
- Sufficient for account linking via email
- Includes Google account ID for reliable identification
- Matches clarified requirements from spec
**Alternatives considered**:
- Additional Google service scopes (rejected: not needed)
- Email-only scope (rejected: insufficient for reliable linking)

### Account Linking Strategy

**Decision**: Automatic merge based on email address matching
**Rationale**:
- Seamless user experience
- Prevents duplicate accounts
- Matches existing user identification pattern
- Clarified in spec requirements
**Alternatives considered**:
- User confirmation prompt (rejected: adds friction)
- Separate account creation (rejected: data fragmentation)

### Frontend OAuth Flow

**Decision**: Authorization Code flow with PKCE (Proof Key for Code Exchange)
**Rationale**:
- Most secure OAuth flow for web applications
- Prevents code interception attacks
- Industry standard for SPAs
- Supported by Google OAuth
**Alternatives considered**:
- Implicit flow (rejected: deprecated security practice)
- Device code flow (rejected: not applicable to web)

### Error Handling Approach

**Decision**: Graceful degradation with clear user messaging
**Rationale**:
- Google service outages handled with informative errors
- Network failures show retry options
- Authorization denials redirect to login with context
- Access revocation triggers immediate logout
**Alternatives considered**:
- Silent failures (rejected: poor UX)
- Generic error messages (rejected: confusing to users)

## Integration Patterns

### Backend Integration Points
- Extend existing user model with oauth_provider and oauth_id fields
- Add middleware for OAuth token validation
- Implement revocation webhook endpoint for Google
- Create service layer for OAuth operations

### Frontend Integration Points
- Add Google login button to existing login page
- Create OAuth callback component for authorization code handling
- Implement token storage in secure HTTP-only cookies
- Add automatic token refresh on API calls

### Database Schema Extensions
- Users table: add `google_id`, `oauth_provider`, `oauth_created_at` columns
- Sessions table: add `refresh_token`, `token_expires_at` columns
- Maintain backward compatibility with existing authentication

## Security Considerations

### Token Storage
- Use HTTP-only, Secure, SameSite cookies for token storage
- Separate access and refresh token storage
- Implement token rotation on refresh

### CSRF Protection
- State parameter validation in OAuth flow
- PKCE implementation for additional security
- Origin validation on callback endpoints

### Data Privacy
- Minimal data collection (email, name, Google ID only)
- Clear user consent during OAuth authorization
- Implement data retention policies per existing system