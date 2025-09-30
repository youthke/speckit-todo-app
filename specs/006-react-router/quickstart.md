# Quickstart: React Router Implementation

**Feature**: React Router Implementation
**Date**: 2025-09-30
**Estimated Time**: 15 minutes

## Overview
This quickstart guide validates the React Router implementation by walking through core user scenarios. Follow these steps to verify routing, authentication, and navigation work as expected.

---

## Prerequisites

1. **Dependencies Installed**:
   ```bash
   cd frontend
   npm install
   ```

2. **Backend Running**:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```
   Backend should be accessible at `http://localhost:8080`

3. **Frontend Running**:
   ```bash
   cd frontend
   npm run dev
   ```
   Frontend should be accessible at `http://localhost:3000`

---

## Scenario 1: Initial Load and Public Routes

### Test: Root Path Redirect
1. Open browser to `http://localhost:3000/`
2. **Expected**: Automatically redirects to `/login`
3. **Verify**: URL bar shows `http://localhost:3000/login`
4. **Verify**: Login page is displayed with "Login with Google" button

### Test: Direct Login Page Access
1. Open browser to `http://localhost:3000/login`
2. **Expected**: Login page renders directly
3. **Verify**: No redirect occurs
4. **Verify**: Page shows branding and login button

### Test: Invalid Route (404)
1. Navigate to `http://localhost:3000/invalid-route`
2. **Expected**: 404 NotFound page displayed
3. **Verify**: Page shows:
   - "404" error code
   - "Page Not Found" message
   - Attempted path: `/invalid-route`
   - "Go to Dashboard" link
   - "Go to Home" link

---

## Scenario 2: Protected Route Access (Unauthenticated)

### Test: Direct Protected Route Access
1. Ensure you are logged out (clear cookies if needed)
2. Navigate to `http://localhost:3000/tasks`
3. **Expected**: Immediate redirect to `/login`
4. **Verify**: URL changes to `http://localhost:3000/login`
5. **Verify**: No flash of protected content

### Test: Dashboard Access
1. While logged out, navigate to `http://localhost:3000/dashboard`
2. **Expected**: Redirect to `/login`
3. **Verify**: Login page displayed

### Test: Profile Access
1. While logged out, navigate to `http://localhost:3000/profile`
2. **Expected**: Redirect to `/login`
3. **Verify**: Login page displayed

---

## Scenario 3: OAuth Authentication Flow

### Test: Google OAuth Login
1. From login page, click "Login with Google" button
2. **Expected**: Redirect to Google OAuth consent screen
3. **Verify**: URL changes to `accounts.google.com`
4. Complete Google login
5. **Expected**: Redirect to `http://localhost:3000/auth/callback`
6. **Expected**: Briefly see "Processing authentication..." message
7. **Expected**: Redirect to `/dashboard`
8. **Verify**: Dashboard page displays
9. **Verify**: Navigation menu is visible with user name

### Test: OAuth Callback Error Handling
1. Navigate to `http://localhost:3000/auth/callback?error=access_denied`
2. **Expected**: Error message displayed
3. **Expected**: Link to return to login page
4. **Verify**: No crash or blank page

---

## Scenario 4: Navigation Menu

### Test: Navigation Menu Visibility
1. Log in successfully
2. **Expected**: Navigation menu appears at top of page
3. **Verify**: Menu shows:
   - App branding/logo
   - "Dashboard" link
   - "Tasks" link
   - "Profile" link
   - User name
   - "Logout" button

### Test: Active Route Highlighting
1. Navigate to `/dashboard`
2. **Verify**: "Dashboard" link has `active` CSS class
3. Navigate to `/tasks`
4. **Verify**: "Tasks" link has `active` CSS class
5. **Verify**: "Dashboard" link no longer has `active` class

### Test: Navigation Link Clicks
1. From dashboard, click "Tasks" link
2. **Expected**: Navigate to `/tasks` page
3. **Verify**: URL changes to `http://localhost:3000/tasks`
4. **Verify**: Tasks page renders
5. Click "Profile" link
6. **Expected**: Navigate to `/profile`
7. **Verify**: Profile page renders

---

## Scenario 5: Browser Navigation Controls

### Test: Back Button
1. Navigate: `/dashboard` → `/tasks` → `/profile`
2. Click browser back button
3. **Expected**: Navigate back to `/tasks`
4. **Verify**: Tasks page renders
5. Click back button again
6. **Expected**: Navigate back to `/dashboard`
7. **Verify**: Dashboard page renders

### Test: Forward Button
1. After using back button twice, click browser forward button
2. **Expected**: Navigate forward to `/tasks`
3. **Verify**: Tasks page renders
4. Click forward again
5. **Expected**: Navigate to `/profile`
6. **Verify**: Profile page renders

### Test: Back from Protected to Login
1. Log out (should redirect to `/login`)
2. Click browser back button
3. **Expected**: Attempt to navigate to protected route
4. **Expected**: Immediately redirect back to `/login`
5. **Verify**: Cannot access protected routes via back button when logged out

---

## Scenario 6: Bookmarking and Direct URLs

### Test: Bookmark Protected Route
1. While logged in, navigate to `/tasks`
2. Bookmark the page (`Cmd+D` or `Ctrl+D`)
3. Close browser completely
4. Reopen browser and click bookmark
5. **Expected**: One of two behaviors:
   - If session valid: Tasks page loads directly
   - If session expired: Redirect to login, then back to tasks after login
6. **Verify**: Eventually reach `/tasks` page

### Test: Share Link While Authenticated
1. While logged in at `/profile`
2. Copy URL: `http://localhost:3000/profile`
3. Open URL in incognito/private window
4. **Expected**: Redirect to `/login` (new session, not authenticated)
5. **Verify**: Cannot bypass authentication with direct URL

---

## Scenario 7: Logout Flow

### Test: Logout Button
1. While logged in, click "Logout" button in navigation
2. **Expected**: Call to logout endpoint
3. **Expected**: Redirect to `/login`
4. **Verify**: URL changes to `http://localhost:3000/login`
5. **Verify**: Navigation menu disappears
6. **Verify**: Auth cookie cleared

### Test: Post-Logout Protected Route Access
1. After logout, navigate to `http://localhost:3000/tasks`
2. **Expected**: Redirect to `/login`
3. **Verify**: Cannot access protected routes
4. **Verify**: No user data in navigation

---

## Scenario 8: Session Validation

### Test: Existing Session on Load
1. Log in successfully
2. Close browser tab (not entire browser)
3. Open new tab to `http://localhost:3000/tasks`
4. **Expected**: Session cookie still valid
5. **Expected**: Tasks page loads directly without redirect to login
6. **Verify**: User remains authenticated

### Test: Expired Session Handling
1. Log in successfully
2. Wait for session to expire (or manually delete session cookie)
3. Navigate to `/tasks`
4. **Expected**: Auth check fails
5. **Expected**: Redirect to `/login`
6. **Verify**: User must log in again

---

## Scenario 9: OAuth Redirect Preservation

### Test: Intended Destination Redirect
1. Log out completely
2. Navigate directly to `http://localhost:3000/tasks`
3. **Expected**: Redirect to `/login`
4. Click "Login with Google"
5. Complete OAuth flow
6. **Expected**: After successful auth, redirect to `/tasks` (original intended destination)
7. **Verify**: End up on tasks page, not dashboard

### Test: Multiple Redirect Attempts
1. Log out
2. Try to access `/profile` → redirected to login
3. Before logging in, try to access `/tasks` → redirected to login (overwrites first redirect)
4. Log in
5. **Expected**: Redirect to `/tasks` (most recent attempt)
6. **Verify**: Does not redirect to `/profile`

---

## Scenario 10: Error Recovery

### Test: Navigation During Auth Check
1. Log in
2. Refresh page while on `/dashboard`
3. **Expected**: Brief loading state while validating session
4. **Expected**: Dashboard page renders after validation
5. **Verify**: No redirect to login if session valid

### Test: 404 Page Actions
1. Navigate to `http://localhost:3000/nonexistent`
2. On 404 page, click "Go to Dashboard" button
3. **Expected**: Navigate to `/dashboard`
4. **Verify**: Dashboard page renders
5. Navigate back to `/nonexistent`
6. Click "Go to Home" button
7. **Expected**: Navigate to `/` → redirects to `/login` or `/dashboard` (depending on auth state)

---

## Performance Validation

### Test: Route Transition Speed
1. Log in and navigate to `/dashboard`
2. Click "Tasks" link
3. **Measure**: Time from click to page render
4. **Expected**: <100ms for transition
5. **Verify**: Transition feels instant, no noticeable delay

### Test: Initial Load Time
1. Clear browser cache
2. Open `http://localhost:3000/`
3. **Measure**: Time from navigation to fully rendered login page
4. **Expected**: <2 seconds on normal network
5. **Verify**: Acceptable load time

---

## Checklist Summary

- [ ] Root path redirects to login
- [ ] Invalid routes show 404 page
- [ ] Protected routes redirect unauthenticated users to login
- [ ] OAuth login flow works end-to-end
- [ ] Navigation menu displays when authenticated
- [ ] Active route is highlighted in navigation
- [ ] Browser back/forward buttons work correctly
- [ ] Bookmarked protected routes work (after login)
- [ ] Logout button works and redirects to login
- [ ] Session validation prevents unauthorized access
- [ ] Intended destination is preserved through OAuth flow
- [ ] 404 page provides working links to recover
- [ ] Route transitions are fast (<100ms)

---

## Troubleshooting

### Issue: Redirects Don't Work
**Symptom**: Clicking links doesn't change URL or page
**Solution**:
1. Check browser console for errors
2. Verify React Router installed: `npm list react-router-dom`
3. Ensure `<BrowserRouter>` wraps app in App.tsx
4. Check no conflicting routing libraries

### Issue: Protected Routes Not Redirecting
**Symptom**: Can access protected routes while logged out
**Solution**:
1. Verify `ProtectedRoute` component wraps protected routes
2. Check `useAuth` hook returns correct `isAuthenticated` value
3. Ensure auth context is provided via `<AuthProvider>`
4. Check browser console for JavaScript errors

### Issue: Navigation Menu Not Showing
**Symptom**: Menu doesn't appear after login
**Solution**:
1. Verify user logged in: check auth state in React DevTools
2. Ensure `<Navigation />` component included in layout
3. Check CSS display properties aren't hiding menu
4. Verify `MainLayout` component rendering

### Issue: 404 Page Not Showing
**Symptom**: Blank page or error for invalid routes
**Solution**:
1. Verify wildcard route (`path="*"`) exists in router config
2. Ensure wildcard route is LAST in route configuration
3. Check `NotFound` component has no errors
4. Look for TypeScript or import errors in console

### Issue: OAuth Callback Fails
**Symptom**: Redirect back to login after OAuth
**Solution**:
1. Check backend OAuth configuration
2. Verify callback URL matches backend expectation
3. Check browser console and network tab for error details
4. Ensure `AuthCallback` component processes response correctly

---

## Success Criteria

✅ **Feature is ready for production when**:
1. All scenarios pass without errors
2. No console errors during normal usage
3. Route transitions feel instant (<100ms)
4. Protected routes cannot be accessed while logged out
5. Browser navigation (back/forward) works as expected
6. OAuth flow completes and redirects correctly
7. 404 page displays for invalid routes
8. Navigation menu shows correct active states
9. Logout flow completes and clears session
10. Tests pass (run `npm test`)

---

**Next Steps**: After validating with this quickstart, run automated tests with `npm test` to ensure all functionality works programmatically.