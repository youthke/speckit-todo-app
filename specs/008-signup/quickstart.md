# Quickstart: Signup Page Testing

**Feature**: 008-signup
**Date**: 2025-10-01
**Purpose**: Manual testing guide for Google OAuth signup feature

## Prerequisites

Before testing, ensure:
1. ✓ Backend server is running on `http://localhost:8080`
2. ✓ Frontend dev server is running on `http://localhost:3000`
3. ✓ Google OAuth credentials configured in `.env`:
   - `GOOGLE_CLIENT_ID=your_client_id`
   - `GOOGLE_CLIENT_SECRET=your_client_secret`
   - `GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/auth/google/callback`
4. ✓ Database migrations applied (users table exists with google_id column)

## Quick Start Commands

```bash
# Terminal 1: Start backend server
cd backend
go run cmd/server/main.go

# Terminal 2: Start frontend dev server
cd frontend
npm run dev

# Verify servers are running
curl http://localhost:8080/api/v1/health  # Should return 200 OK
curl http://localhost:3000                 # Should return React app
```

## Test Scenario 1: New User Signup

**Objective**: Verify new user can sign up with Google and is automatically logged in.

**Steps**:
1. Open browser in **incognito mode** (to ensure fresh session)
2. Navigate to `http://localhost:3000/signup`
3. Verify page displays:
   - "Sign Up" heading
   - "Sign up with Google" button
   - "Already have an account? Log in" link
4. Click **"Sign up with Google"** button
5. Verify redirect to Google OAuth consent screen:
   - URL starts with `https://accounts.google.com/o/oauth2/v2/auth`
   - Shows app name and requested permissions (email, profile)
6. Select Google account (use test account not previously signed up)
7. Click **"Continue"** or **"Allow"** to authorize
8. Verify redirect to `http://localhost:3000/` (home page)
9. Verify user is logged in:
   - Check for user profile/avatar in UI
   - Verify `session_token` cookie exists in browser DevTools → Application → Cookies
   - Cookie should have:
     - `HttpOnly` flag = true
     - `Max-Age` = 604800 (7 days)
     - `Path` = /

**Expected Outcome**:
- ✓ New user created in database
- ✓ Session created and cookie set
- ✓ User redirected to home page
- ✓ User has access to protected features

**Database Verification**:
```bash
# Check user was created
cd backend
sqlite3 todo.db "SELECT id, email, google_id, created_at FROM users ORDER BY id DESC LIMIT 1;"

# Check session was created
sqlite3 todo.db "SELECT user_id, expires_at FROM sessions WHERE user_id = (SELECT id FROM users ORDER BY id DESC LIMIT 1);"
```

---

## Test Scenario 2: Existing User Auto-Login

**Objective**: Verify existing user is automatically logged in when attempting to sign up again.

**Prerequisites**: Complete Test Scenario 1 first.

**Steps**:
1. **Logout** from the application:
   - Click logout button in UI, OR
   - Clear `session_token` cookie manually in DevTools
2. Navigate to `http://localhost:3000/signup` (not login page)
3. Click **"Sign up with Google"** button
4. Verify redirect to Google OAuth consent screen
5. Select the **same Google account** used in Test Scenario 1
6. Authorize (may skip consent screen if already authorized)
7. Verify redirect to `http://localhost:3000/` (home page)
8. Verify user is logged in (check session cookie and UI)

**Expected Outcome**:
- ✓ No new user created (existing user reused)
- ✓ New session created for existing user
- ✓ User redirected to home page (NOT to login page)
- ✓ User has access to protected features

**Database Verification**:
```bash
# Verify no duplicate user created (count should be same as before)
cd backend
sqlite3 todo.db "SELECT COUNT(*) FROM users;"

# Verify new session created for same user
sqlite3 todo.db "SELECT id, user_id, created_at FROM sessions WHERE user_id = (SELECT id FROM users WHERE email = 'your_test_email@example.com') ORDER BY id DESC LIMIT 2;"
# Should show 2 sessions (old + new)
```

---

## Test Scenario 3: Error Handling - Denied Permission

**Objective**: Verify proper error handling when user denies Google authorization.

**Steps**:
1. Open new incognito browser window
2. Navigate to `http://localhost:3000/signup`
3. Click **"Sign up with Google"** button
4. On Google consent screen, click **"Cancel"** or **"Deny"**
5. Verify redirect to `http://localhost:3000/signup?error=authentication_failed`
6. Verify error message displayed on signup page:
   - "Authentication failed. Please try again."
7. Verify no session cookie created
8. Click **"Sign up with Google"** again to retry
9. This time, authorize → verify successful signup

**Expected Outcome**:
- ✓ Error displayed when permission denied
- ✓ No user/session created
- ✓ User can retry signup
- ✓ Retry succeeds after authorization

---

## Test Scenario 4: Rate Limiting

**Objective**: Verify rate limiting prevents abuse (10 attempts per 15 minutes per IP).

**Steps**:
1. Open browser (can use same session)
2. Navigate to `http://localhost:3000/signup`
3. Click **"Sign up with Google"** button repeatedly
4. On Google consent screen, click **Cancel** each time
5. Repeat steps 2-4 for **10 times** within a few minutes
6. On the **11th attempt**, verify:
   - Backend returns HTTP 429 (check Network tab in DevTools)
   - Redirect to `http://localhost:3000/signup?error=rate_limit_exceeded`
   - Error message: "Too many attempts. Please try again later."
7. Wait **15 minutes** (or restart backend server to clear in-memory limiter)
8. Retry → verify signup works again

**Expected Outcome**:
- ✓ First 10 attempts succeed (reach Google OAuth screen)
- ✓ 11th attempt blocked with HTTP 429
- ✓ Rate limit error displayed in UI
- ✓ Rate limit resets after 15 minutes

**cURL Testing** (alternative to manual clicks):
```bash
# Rapidly hit the endpoint 11 times
for i in {1..11}; do
  curl -i http://localhost:8080/api/v1/auth/google/login
  echo "Attempt $i"
done

# Expected: First 10 return 302, 11th returns 429
```

---

## Test Scenario 5: Email Verification Required

**Objective**: Verify signup fails if Google account has unverified email.

**Note**: This is difficult to test manually as most Google accounts have verified emails.

**Steps** (if you have unverified Google account):
1. Navigate to `http://localhost:3000/signup`
2. Click **"Sign up with Google"**
3. Use Google account with **unverified email**
4. Authorize
5. Verify redirect to `http://localhost:3000/signup?error=authentication_failed`
6. Verify error message displayed

**Alternative Test** (modify backend code temporarily):
```go
// In google_oauth_handler.go, temporarily set:
if !userInfo.EmailVerified {
    // Force this condition by commenting out the check
    // This will allow you to test error handling
}
```

**Expected Outcome**:
- ✓ Unverified email rejected
- ✓ Error message displayed
- ✓ No user/session created

---

## Test Scenario 6: Navigation Between Signup and Login

**Objective**: Verify users can navigate between signup and login pages.

**Steps**:
1. Navigate to `http://localhost:3000/signup`
2. Verify "Already have an account? Log in" link exists
3. Click **"Log in"** link
4. Verify redirect to `http://localhost:3000/login`
5. Verify login page displays (with Google login button)
6. Navigate back to `http://localhost:3000/signup`
7. Verify signup page displays correctly

**Expected Outcome**:
- ✓ Navigation links work
- ✓ Pages render correctly
- ✓ No console errors in DevTools

---

## Test Scenario 7: Session Persistence

**Objective**: Verify session persists across page reloads and browser tabs.

**Prerequisites**: Complete Test Scenario 1 (user logged in).

**Steps**:
1. Verify user is logged in at `http://localhost:3000/`
2. **Reload page** (F5 or Ctrl+R)
3. Verify user still logged in (session cookie persists)
4. Open **new tab** → navigate to `http://localhost:3000/`
5. Verify user logged in in new tab (same session)
6. Close all tabs → reopen browser → navigate to `http://localhost:3000/`
7. Verify user still logged in (session cookie persists in browser)
8. Wait **7 days** (or manually expire session in database)
9. Reload page → verify session expired and user logged out

**Expected Outcome**:
- ✓ Session persists across reloads
- ✓ Session shared across tabs
- ✓ Session persists after browser restart (within 7 days)
- ✓ Session expires after 7 days

---

## Test Scenario 8: End-to-End Flow

**Objective**: Complete realistic user journey from signup to using the app.

**Steps**:
1. Start fresh (clear database or use new test account)
2. Navigate to `http://localhost:3000/` (home page)
3. If not logged in, click **"Sign Up"** button (should navigate to `/signup`)
4. Click **"Sign up with Google"**
5. Complete Google OAuth flow
6. Verify redirect to home page with user logged in
7. Create a new TODO item (if feature exists)
8. Refresh page → verify TODO persists and user still logged in
9. Logout → verify redirect to login page
10. Navigate to `/signup` again → sign up with same account
11. Verify auto-login → redirect to home page → TODO items still visible

**Expected Outcome**:
- ✓ Complete signup flow works
- ✓ User can access protected features
- ✓ Data persists across sessions
- ✓ Auto-login for existing users works

---

## Troubleshooting

### Issue: "Authentication failed" error immediately after clicking signup button

**Possible Causes**:
1. Backend server not running
2. Google OAuth credentials missing/invalid in `.env`
3. CSRF state cookie not being set (check cookie settings)

**Debug Steps**:
```bash
# Check backend logs
cd backend
tail -f logs/app.log  # or check console output

# Verify OAuth endpoint responds
curl -i http://localhost:8080/api/v1/auth/google/login
# Should return 302 with Location header
```

### Issue: Redirect to Google fails (blank page or error)

**Possible Causes**:
1. Invalid `GOOGLE_REDIRECT_URI` in `.env`
2. Redirect URI not whitelisted in Google Cloud Console
3. OAuth consent screen not configured

**Fix**:
1. Go to Google Cloud Console → Credentials
2. Add `http://localhost:8080/api/v1/auth/google/callback` to authorized redirect URIs
3. Save and wait 5 minutes for changes to propagate

### Issue: "Email not verified" error

**Possible Causes**:
1. Test Google account has unverified email
2. OAuth scope doesn't include email

**Fix**:
1. Verify email in Google account settings
2. Check backend code: OAuth scope should include `openid email profile`

### Issue: Session cookie not set

**Possible Causes**:
1. Browser blocking third-party cookies
2. `HttpOnly` flag preventing JavaScript access (this is correct behavior)
3. Session creation failed in backend

**Debug Steps**:
1. Check DevTools → Application → Cookies → `http://localhost:3000`
2. Should see `session_token` cookie with HttpOnly flag
3. Check backend logs for session creation errors

---

## Success Criteria

All test scenarios should pass with these outcomes:

- [x] New user can sign up with Google and is automatically logged in
- [x] Existing user attempting signup is automatically logged in (not redirected to login page)
- [x] User can deny permission and retry signup
- [x] Rate limiting blocks excessive signup attempts (11th attempt within 15 min)
- [x] Unverified email is rejected with error message
- [x] Navigation between signup and login pages works
- [x] Session persists across reloads and tabs for 7 days
- [x] End-to-end flow completes successfully

## Next Steps

After completing manual testing:
1. Run automated contract tests: `cd backend && go test ./tests/contract/signup_*_test.go`
2. Run integration tests: `cd backend && go test ./tests/integration/signup_*_test.go`
3. Run frontend tests: `cd frontend && npm test -- SignupPage.test.tsx`
4. Review test coverage report
5. Document any bugs found in GitHub issues

---

**Last Updated**: 2025-10-01
**Feature**: 008-signup
**Test Environment**: localhost (backend:8080, frontend:3000)
