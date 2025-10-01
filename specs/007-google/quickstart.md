# Quickstart: Google Account Signup

**Feature**: 007-google | **Date**: 2025-10-01

## Purpose
This quickstart validates the Google Account Signup feature by walking through the complete OAuth flow and verifying all functional requirements.

---

## Prerequisites

### 1. Google OAuth Credentials
```bash
# Ensure these environment variables are set in backend/.env
GOOGLE_CLIENT_ID=your_client_id_here
GOOGLE_CLIENT_SECRET=your_client_secret_here
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

### 2. Services Running
```bash
# Terminal 1: Start backend
cd backend
go run cmd/main.go

# Terminal 2: Start frontend
cd frontend
npm run dev
```

### 3. Clean Database State
```bash
# Run migrations to ensure schema is up to date
cd backend
# Run migration command (adjust based on your migration tool)
# Example: go run cmd/migrate.go up
```

---

## Test Scenarios

### Scenario 1: Successful New User Signup
**Validates**: FR-001, FR-002, FR-003, FR-004, FR-005, FR-009, FR-010

**Steps**:
1. Open browser to `http://localhost:3000/signup`
2. Verify "Sign up with Google" button is visible
3. Click "Sign up with Google" button
4. **Expected**: Redirected to Google OAuth consent page
5. Sign in with a Google account that has a verified email
6. Grant permissions when prompted
7. **Expected**: Redirected back to `http://localhost:3000/` (home page)
8. **Expected**: Session cookie `session_token` is set (check browser DevTools → Application → Cookies)
9. Open DevTools → Network, make a request to `/api/auth/me`
10. **Expected Response**:
    ```json
    {
      "id": 1,
      "email": "your-google-email@gmail.com",
      "auth_method": "google",
      "created_at": "2025-10-01T12:00:00Z"
    }
    ```

**Database Validation**:
```sql
-- Check user was created
SELECT * FROM users WHERE email = 'your-google-email@gmail.com';
-- Should return 1 row with auth_method = 'google'

-- Check Google identity link was created
SELECT * FROM google_identities WHERE email = 'your-google-email@gmail.com';
-- Should return 1 row with email_verified = true

-- Check session was created
SELECT * FROM authentication_sessions WHERE user_id = 1;
-- Should return 1 row with expires_at = created_at + 7 days
```

**Pass Criteria**:
- ✅ User account created with `auth_method = 'google'`
- ✅ Google identity record created with correct `google_user_id` and `email`
- ✅ Session created with 7-day expiration
- ✅ User redirected to home page and logged in
- ✅ `/api/auth/me` returns correct user data

---

### Scenario 2: Duplicate Signup Prevention
**Validates**: FR-008

**Precondition**: Complete Scenario 1 first (user account exists)

**Steps**:
1. Log out from the application
2. Navigate to `http://localhost:3000/signup`
3. Click "Sign up with Google" button
4. Sign in with the **same Google account** used in Scenario 1
5. Grant permissions
6. **Expected**: Redirected to `http://localhost:3000/login` (not home page)
7. **Expected**: No new user or google_identity records created

**Database Validation**:
```sql
-- Count users with this email
SELECT COUNT(*) FROM users WHERE email = 'your-google-email@gmail.com';
-- Should return 1 (no duplicate created)

-- Count Google identities
SELECT COUNT(*) FROM google_identities WHERE email = 'your-google-email@gmail.com';
-- Should return 1 (no duplicate created)
```

**Pass Criteria**:
- ✅ User redirected to login page (not home page)
- ✅ No duplicate user account created
- ✅ No duplicate google_identity record created
- ✅ User can log in normally after redirect

---

### Scenario 3: Unverified Email Rejection
**Validates**: FR-007

**Note**: This scenario requires a Google account with an unverified email, which is difficult to obtain. Consider mocking this in automated tests instead.

**Mocked Test Approach**:
1. Modify OAuth service to accept test mode flag
2. In test mode, return mocked ID token with `email_verified = false`
3. Attempt signup
4. **Expected**: Redirected to signup page with error query parameter
5. **Expected**: Error message displayed: "Authentication failed"
6. **Expected**: No user or google_identity records created

**Pass Criteria**:
- ✅ Signup rejected when `email_verified = false`
- ✅ Generic error message shown
- ✅ No database records created

---

### Scenario 4: OAuth Error Handling
**Validates**: FR-006

**Steps**:
1. Navigate to `http://localhost:3000/signup`
2. Click "Sign up with Google"
3. On Google consent page, click "Cancel" or "Deny"
4. **Expected**: Redirected to `http://localhost:3000/signup?error=authentication_failed`
5. **Expected**: Error message displayed: "Authentication failed"
6. **Expected**: No user or google_identity records created

**Alternative - Network Error Simulation**:
1. Start signup flow
2. Disable network before callback completes
3. **Expected**: Timeout and error message "Authentication failed"

**Pass Criteria**:
- ✅ User returned to signup page on OAuth denial
- ✅ Generic error message displayed
- ✅ No database records created

---

### Scenario 5: Session Expiration
**Validates**: FR-010

**Steps**:
1. Complete Scenario 1 (create session)
2. Check session cookie expiration:
   ```javascript
   // In browser console
   document.cookie.split('; ').find(row => row.startsWith('session_token='))
   ```
3. **Expected**: Cookie has `Max-Age=604800` (7 days in seconds)
4. Query database for session:
   ```sql
   SELECT expires_at, created_at,
          (expires_at - created_at) AS duration_seconds
   FROM authentication_sessions
   WHERE user_id = 1
   ORDER BY created_at DESC
   LIMIT 1;
   ```
5. **Expected**: `duration_seconds = 604800` (7 days)

**Pass Criteria**:
- ✅ Session cookie set with 7-day max age
- ✅ Database record has `expires_at = created_at + 7 days`

---

### Scenario 6: Cancelled OAuth Flow
**Validates**: FR-006

**Steps**:
1. Navigate to `http://localhost:3000/signup`
2. Click "Sign up with Google"
3. On Google's OAuth page, close the tab/window before completing
4. Return to original tab
5. Click "Sign up with Google" again
6. This time, complete the OAuth flow normally
7. **Expected**: Signup succeeds (previous cancelled flow doesn't interfere)

**Pass Criteria**:
- ✅ Cancelled flow doesn't create partial database records
- ✅ New flow completes successfully
- ✅ No orphaned data in database

---

## Integration Test Checklist

Manual validation of all acceptance scenarios from spec.md:

- [ ] **AS-1**: New user signup creates account and logs in ✓ (Scenario 1)
- [ ] **AS-2**: Existing user redirected to login ✓ (Scenario 2)
- [ ] **AS-3**: Cancelled flow returns to signup page ✓ (Scenario 6)
- [ ] **AS-4**: OAuth errors show "Authentication failed" ✓ (Scenario 4)
- [ ] **AS-5**: Unverified email rejected ✓ (Scenario 3)

---

## Performance Validation

### OAuth Flow Timing
```bash
# Use browser DevTools Network tab to measure:
# 1. Time from button click to Google redirect: <500ms
# 2. Time from Google callback to home page redirect: <3s (includes user creation)
```

**Pass Criteria**:
- Initial redirect: <500ms
- Complete signup flow: <3 seconds
- Session validation (`/api/auth/me`): <50ms

---

## Security Validation

### Session Cookie Security
1. Inspect session cookie in DevTools → Application → Cookies
2. Verify flags:
   - ✅ `HttpOnly` (prevents JavaScript access)
   - ✅ `Secure` (HTTPS only, if in production)
   - ✅ `SameSite=Strict` (CSRF protection)

### State Parameter Validation
1. Inspect network request to `/api/auth/google/login`
2. Verify `oauth_state` cookie is set
3. Inspect callback request to `/api/auth/google/callback`
4. Verify `state` query parameter matches cookie value

**Pass Criteria**:
- ✅ State token generated and validated (CSRF protection)
- ✅ Session cookie has security flags
- ✅ No sensitive data in URL parameters (except code/state)

---

## Cleanup

After testing:
```sql
-- Remove test data
DELETE FROM authentication_sessions WHERE user_id IN (
  SELECT id FROM users WHERE email = 'your-google-email@gmail.com'
);

DELETE FROM google_identities WHERE email = 'your-google-email@gmail.com';

DELETE FROM users WHERE email = 'your-google-email@gmail.com';
```

---

## Automated Test Execution

Once implemented, run automated tests:
```bash
# Backend contract tests
cd backend
go test ./tests/contract/google_oauth_test.go -v

# Backend integration tests
go test ./tests/integration/signup_flow_test.go -v

# Frontend component tests
cd ../frontend
npm test -- GoogleSignupButton.test.tsx
```

---

## Success Criteria Summary

All scenarios must pass for feature acceptance:
1. ✅ New user signup flow (Scenario 1)
2. ✅ Duplicate prevention (Scenario 2)
3. ✅ Email verification check (Scenario 3)
4. ✅ Error handling (Scenario 4)
5. ✅ 7-day session duration (Scenario 5)
6. ✅ Cancelled flow handling (Scenario 6)
7. ✅ Performance targets met
8. ✅ Security validations passed
