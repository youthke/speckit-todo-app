# Data Model: React Router Implementation

**Feature**: React Router Implementation
**Date**: 2025-09-30

## Overview
This document defines the data structures and state management for the routing feature. Since routing is a frontend-only feature with no backend persistence, this focuses on client-side state and navigation structures.

---

## 1. Route Configuration

### Entity: Route
Represents a URL path mapping to a page component.

**Fields**:
- `path` (string): URL path pattern (e.g., "/login", "/tasks", "*")
- `component` (React.Component): Page component to render
- `isProtected` (boolean): Whether route requires authentication
- `layout` (React.Component | null): Optional layout wrapper
- `meta` (RouteMeta): Metadata for the route

**Relationships**:
- Routes can be nested (parent-child relationship)
- Protected routes depend on authentication state

**Validation Rules**:
- Path must start with "/" or be "*" for catch-all
- Protected routes must be wrapped in ProtectedRoute component
- Each path must be unique

---

## 2. Navigation State

### Entity: NavigationState
Managed by React Router's internal state.

**Fields**:
- `location` (Location): Current location object
  - `pathname` (string): Current path
  - `search` (string): Query string
  - `hash` (string): URL hash
  - `state` (any): Location state
- `history` (History): Navigation history stack
  - `length` (number): Number of entries
  - `action` (string): Last navigation action (PUSH, REPLACE, POP)

**State Transitions**:
1. `IDLE` → `NAVIGATING`: User clicks link or calls navigate()
2. `NAVIGATING` → `IDLE`: Route transition completes
3. `NAVIGATING` → `REDIRECTING`: Protected route check triggers redirect
4. `REDIRECTING` → `IDLE`: Redirect completes to login page

---

## 3. Authentication Context

### Entity: AuthState
Existing structure from useAuth hook.

**Fields**:
- `user` (User | null): Current authenticated user
  - `id` (number): User ID
  - `email` (string): User email
  - `name` (string): User display name
- `isAuthenticated` (boolean): Whether user is logged in
- `isLoading` (boolean): Whether auth check is in progress
- `error` (Error | null): Authentication error if any

**Usage in Routing**:
- `isAuthenticated` determines access to protected routes
- `isLoading` prevents premature redirects during auth check
- `user` data displayed in navigation menu

---

## 4. Route Metadata

### Entity: RouteMeta
Optional metadata for routes.

**Fields**:
- `title` (string): Page title for document.title
- `requiresAuth` (boolean): Whether authentication is required
- `redirectTo` (string): Where to redirect if access denied
- `icon` (string): Icon name for navigation menu
- `showInNav` (boolean): Whether to display in navigation menu

**Example**:
```typescript
{
  path: "/tasks",
  meta: {
    title: "My Tasks",
    requiresAuth: true,
    redirectTo: "/login",
    icon: "task-list",
    showInNav: true
  }
}
```

---

## 5. Navigation Menu Item

### Entity: NavItem
Represents a link in the navigation menu.

**Fields**:
- `label` (string): Display text
- `path` (string): Route path to navigate to
- `icon` (string): Icon identifier
- `isActive` (boolean): Whether this is the current route
- `isVisible` (boolean): Whether item should be displayed

**Computed Properties**:
- `isActive`: Computed from current location.pathname
- `isVisible`: Computed from auth state and route meta

**Example**:
```typescript
[
  { label: "Dashboard", path: "/dashboard", icon: "home", isVisible: isAuthenticated },
  { label: "Tasks", path: "/tasks", icon: "check-square", isVisible: isAuthenticated },
  { label: "Profile", path: "/profile", icon: "user", isVisible: isAuthenticated },
  { label: "Login", path: "/login", icon: "log-in", isVisible: !isAuthenticated }
]
```

---

## 6. OAuth Redirect State

### Entity: OAuthRedirectState
Temporary state during OAuth flow.

**Fields**:
- `intendedDestination` (string): Path user wanted before redirect to login
- `timestamp` (number): When redirect was initiated
- `state` (string): OAuth state parameter for security

**Storage**:
- Stored in sessionStorage as temporary data
- Key: `oauth_redirect`
- Cleared after successful authentication

**Lifecycle**:
1. User attempts to access protected route while unauthenticated
2. System stores current path in sessionStorage
3. User redirects to OAuth login
4. After successful auth, user redirects to intended destination
5. sessionStorage cleared

---

## 7. Loading State

### Entity: RouteLoadingState
Loading state during route transitions.

**Fields**:
- `isNavigating` (boolean): Whether route transition is in progress
- `isPending` (boolean): Whether data is being loaded for route
- `progress` (number): Optional progress indicator (0-100)

**Usage**:
- Display loading spinner during transitions
- Prevent duplicate navigation attempts
- Show progress for slow-loading pages

---

## 8. Route Guards

### Entity: RouteGuard
Logic for determining route access.

**Type**: Function
**Signature**: `(authState: AuthState) => boolean | string`

**Returns**:
- `true`: Access granted, render route
- `false`: Access denied, redirect to default login
- `string`: Access denied, redirect to specified path

**Example Guards**:
```typescript
// Requires authentication
const requireAuth = (authState: AuthState) => {
  if (authState.isLoading) return false; // Wait for auth check
  return authState.isAuthenticated ? true : "/login";
};

// Requires no authentication (login page)
const requireGuest = (authState: AuthState) => {
  return !authState.isAuthenticated ? true : "/dashboard";
};
```

---

## 9. Data Flow Diagrams

### Protected Route Access Flow
```
User navigates to /tasks
  ↓
ProtectedRoute component checks authState
  ↓
Is isLoading = true?
  ↓ Yes → Show loading spinner
  ↓ No
Is isAuthenticated = true?
  ↓ Yes → Render <Outlet /> (tasks page)
  ↓ No → <Navigate to="/login" replace />
```

### OAuth Redirect Flow
```
Unauthenticated user tries /tasks
  ↓
Store "/tasks" in sessionStorage
  ↓
Redirect to /login
  ↓
User clicks "Login with Google"
  ↓
Redirect to backend OAuth endpoint
  ↓
Backend redirects to Google
  ↓
User approves
  ↓
Google redirects to /auth/callback
  ↓
AuthCallback component processes
  ↓
Retrieve "/tasks" from sessionStorage
  ↓
Navigate to "/tasks"
  ↓
Clear sessionStorage
```

---

## 10. State Management Strategy

### Authentication State
- **Source**: useAuth hook (React Context)
- **Scope**: Global (entire app)
- **Persistence**: None (session managed by HTTP-only cookies)
- **Updates**: On login, logout, session refresh

### Navigation State
- **Source**: React Router (useLocation, useNavigate)
- **Scope**: Global (routing layer)
- **Persistence**: Browser history API
- **Updates**: On route transitions

### Component State
- **Source**: useState, useRef (local component state)
- **Scope**: Component-specific
- **Persistence**: None (ephemeral)
- **Updates**: User interactions, prop changes

---

## 11. Type Definitions

### TypeScript Interfaces

```typescript
// Route configuration
interface RouteConfig {
  path: string;
  component: React.ComponentType;
  isProtected: boolean;
  layout?: React.ComponentType;
  meta?: RouteMeta;
  children?: RouteConfig[];
}

// Route metadata
interface RouteMeta {
  title?: string;
  requiresAuth?: boolean;
  redirectTo?: string;
  icon?: string;
  showInNav?: boolean;
}

// Navigation item
interface NavItem {
  label: string;
  path: string;
  icon?: string;
  isActive?: boolean;
  isVisible: boolean;
}

// OAuth redirect state
interface OAuthRedirectState {
  intendedDestination: string;
  timestamp: number;
  state: string;
}

// Route guard function type
type RouteGuard = (authState: AuthState) => boolean | string;

// Auth state (from useAuth hook)
interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: Error | null;
}

interface User {
  id: number;
  email: string;
  name: string;
}
```

---

## 12. Validation Rules

### Route Configuration
- All paths must be unique
- Catch-all route ("*") must be last in configuration
- Protected routes must use ProtectedRoute wrapper
- Path must start with "/" or be "*"

### Navigation
- Cannot navigate to same route consecutively
- Protected route navigation while unauthenticated triggers redirect
- Invalid paths trigger 404 NotFound page

### OAuth Redirect
- Intended destination must be a valid route path
- Redirect state expires after 10 minutes
- Only one redirect state stored at a time

---

## 13. Performance Considerations

### Route Lazy Loading
- Not implemented initially (app is small)
- Consider when individual pages exceed 50KB
- Use React.lazy() + Suspense for code splitting

### Navigation Caching
- React Router v6 caches component instances
- No additional caching needed for route transitions
- Browser history API handles back/forward efficiently

### State Updates
- Auth state updates trigger protected route re-evaluation
- Use React Context to minimize re-renders
- Navigation state changes don't trigger unnecessary component updates

---

**Status**: ✅ Data model complete, ready for contract definition