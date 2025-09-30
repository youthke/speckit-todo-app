# Quickstart: Google Account Login

**Feature**: Google OAuth 2.0 Authentication Integration
**Date**: 2025-09-29

## Quick Validation Steps

### Prerequisites Verification
1. ✅ Google Cloud Console project configured with OAuth 2.0 credentials
2. ✅ Authorized redirect URIs configured (`http://localhost:8080/api/v1/auth/google/callback`)
3. ✅ Backend server running on `http://localhost:8080`
4. ✅ Frontend server running on `http://localhost:3000`

### End-to-End User Flow Test

#### Test Case 1: New User Google Sign-In
**Steps:**
1. Navigate to `http://localhost:3000/login`
2. Click "Sign in with Google" button
3. Verify redirect to `accounts.google.com` with consent screen
4. Grant permissions (email, profile, openid)
5. Verify redirect back to application with authentication success
6. Check user is logged in and can access protected routes
7. Verify new user account created in database with `google_id` populated

**Expected Results:**
- User redirected to Google OAuth consent screen
- After consent, user returned to app logged in
- User can access todo dashboard
- Database has new user record with Google OAuth data

#### Test Case 2: Existing User Account Linking
**Pre-condition:** User already has account with email `test@gmail.com` using password auth

**Steps:**
1. Navigate to `http://localhost:3000/login`
2. Click "Sign in with Google" button
3. Sign in with Google account using same email `test@gmail.com`
4. Grant OAuth permissions
5. Verify redirect back to application
6. Check user is logged in with existing account data
7. Verify database shows `google_id` and `oauth_provider` added to existing user

**Expected Results:**
- Google OAuth links to existing account (no duplicate created)
- User retains all existing todo data
- User can now log in with either method

#### Test Case 3: Session Management
**Steps:**
1. Complete Test Case 1 (sign in with Google)
2. Wait for 30 minutes of inactivity
3. Make API request to protected endpoint
4. Verify automatic token refresh occurs
5. Check session remains valid with updated expiration
6. Close browser and reopen
7. Verify session persistence via secure cookies

**Expected Results:**
- Session automatically refreshed during activity
- 24-hour session maintained with activity
- Secure cookie persistence across browser sessions

#### Test Case 4: Error Handling
**Steps:**
1. Navigate to `http://localhost:3000/login`
2. Click "Sign in with Google" button
3. On Google consent screen, click "Deny" or "Cancel"
4. Verify redirect back to login page with error message
5. Disconnect internet during OAuth flow
6. Verify appropriate error handling and messaging

**Expected Results:**
- Graceful handling of user denial
- Clear error messages for network issues
- User remains on login page with retry option

### API Contract Validation

#### Test OAuth Initiation Endpoint
```bash
curl -v http://localhost:8080/api/v1/auth/google/login \
  -H "Accept: application/json"
```
**Expected:** 302 redirect to Google OAuth with state parameter

#### Test Session Validation Endpoint
```bash
curl -v http://localhost:8080/api/v1/auth/session/validate \
  -H "Cookie: session_token=<valid_token>" \
  -H "Accept: application/json"
```
**Expected:** 200 response with user and session information

#### Test OAuth Callback Endpoint
```bash
curl -v "http://localhost:8080/api/v1/auth/google/callback?code=test&state=invalid" \
  -H "Accept: application/json"
```
**Expected:** 400 error for invalid state parameter

### Database State Validation

#### Verify User Account Structure
```sql
SELECT id, email, name, google_id, oauth_provider, oauth_created_at
FROM users
WHERE oauth_provider = 'google';
```
**Expected:** Google OAuth users have populated google_id and oauth_provider fields

#### Verify Session Management
```sql
SELECT session_token, user_id, token_expires_at, session_expires_at
FROM authentication_sessions
WHERE user_id = <test_user_id>;
```
**Expected:** Active sessions with proper expiration times

#### Verify OAuth State Cleanup
```sql
SELECT COUNT(*) FROM oauth_states
WHERE expires_at < NOW();
```
**Expected:** 0 (expired states should be automatically cleaned)

### Security Validation

#### Test CSRF Protection
1. Attempt OAuth callback with invalid state parameter
2. Verify request is rejected with 400 error
3. Check logs for security event recording

#### Test Token Security
1. Inspect browser cookies after login
2. Verify session_token is HttpOnly and Secure
3. Confirm tokens are not exposed in browser storage
4. Test that tokens are encrypted at rest in database

#### Test Session Termination
1. Log in via Google OAuth
2. From Google account settings, revoke app access
3. Make API request with session token
4. Verify immediate session termination and logout

### Performance Validation

#### Test OAuth Flow Timing
1. Measure time from "Sign in with Google" click to successful login
2. Target: Complete flow in <3 seconds (excluding user interaction time)
3. Verify database queries are optimized (indexed lookups)

#### Test Concurrent Login Handling
1. Simulate 10 concurrent Google OAuth flows
2. Verify all succeed without race conditions
3. Check database consistency under load

### Integration Test Scenarios

#### Frontend Integration
1. Verify Google login button displays correctly
2. Test OAuth flow with JavaScript enabled/disabled
3. Verify proper error handling in React components
4. Test responsive design on mobile devices

#### Backend Integration
1. Verify middleware properly validates OAuth sessions
2. Test API endpoints require proper authentication
3. Verify CORS headers for cross-origin requests
4. Test rate limiting on OAuth endpoints

### Rollback Testing

#### Test Graceful Degradation
1. Temporarily disable Google OAuth (remove credentials)
2. Verify existing password authentication still works
3. Verify appropriate error messages for OAuth unavailability
4. Test system recovery when OAuth credentials restored

## Success Criteria Checklist

- [ ] New users can sign in with Google account
- [ ] Existing users can link Google accounts to existing profiles
- [ ] Sessions persist for 24 hours with activity-based refresh
- [ ] User denial and errors handled gracefully
- [ ] All API contracts respond according to specification
- [ ] Database properly stores OAuth user data
- [ ] Security measures (CSRF, token encryption) function correctly
- [ ] Performance meets targets (<3 second OAuth flow)
- [ ] Frontend and backend integration works seamlessly
- [ ] Access revocation immediately terminates sessions
- [ ] System maintains functionality with OAuth disabled

## Troubleshooting Common Issues

### "Invalid redirect URI" Error
- Verify Google Cloud Console redirect URIs match exactly
- Check for trailing slashes or protocol mismatches
- Confirm localhost ports match running services

### "Invalid state parameter" Error
- Check OAuth state generation and validation logic
- Verify state cookies are properly set and readable
- Confirm state cleanup isn't removing active states too quickly

### Session Not Persisting
- Verify cookie settings (HttpOnly, Secure, SameSite)
- Check domain and path configuration for cookies
- Confirm database sessions table is properly configured

### Account Linking Not Working
- Verify email matching logic in user lookup
- Check unique constraints on google_id field
- Confirm OAuth user creation vs. linking logic