# Routing API Contract

**Feature**: React Router Implementation
**Type**: Frontend Component API
**Date**: 2025-09-30

## Overview
This document defines the API contracts for routing components, hooks, and utilities. These are frontend component interfaces, not HTTP endpoints.

---

## 1. ProtectedRoute Component

### Purpose
Wrapper component that guards routes requiring authentication.

### Interface
```typescript
interface ProtectedRouteProps {
  redirectTo?: string; // Where to redirect unauthenticated users (default: "/login")
  children?: React.ReactNode; // Optional: for wrapper pattern
}

const ProtectedRoute: React.FC<ProtectedRouteProps>
```

### Behavior
**Input**: `redirectTo` path (optional)
**Output**:
- If `isLoading = true`: Renders loading spinner
- If `isAuthenticated = true`: Renders `<Outlet />` (child routes)
- If `isAuthenticated = false`: Renders `<Navigate to={redirectTo} replace />`

### Dependencies
- `useAuth()` hook for authentication state
- `react-router-dom` for `Outlet` and `Navigate`

### Test Contract
```typescript
describe('ProtectedRoute', () => {
  it('redirects to login when not authenticated')
  it('renders child routes when authenticated')
  it('shows loading state while checking auth')
  it('uses custom redirect path when provided')
})
```

---

## 2. useTypedNavigate Hook

### Purpose
Type-safe navigation utility with autocomplete for route keys.

### Interface
```typescript
type RouteKey = 'HOME' | 'LOGIN' | 'AUTH_CALLBACK' | 'DASHBOARD' | 'TASKS' | 'PROFILE' | 'NOT_FOUND';

interface TypedNavigateReturn {
  navigateTo: (key: RouteKey, options?: NavigateOptions) => void;
  navigateToPath: (path: string, options?: NavigateOptions) => void;
  goBack: () => void;
  goForward: () => void;
}

const useTypedNavigate: () => TypedNavigateReturn
```

### Behavior
**Methods**:
- `navigateTo(key, options)`: Navigate using route key constant
- `navigateToPath(path, options)`: Navigate using string path (for dynamic routes)
- `goBack()`: Navigate back in history
- `goForward()`: Navigate forward in history

**NavigateOptions**:
```typescript
interface NavigateOptions {
  replace?: boolean; // Replace current history entry
  state?: any; // State to pass to destination
}
```

### Usage Example
```typescript
const { navigateTo, goBack } = useTypedNavigate();

// Type-safe navigation with autocomplete
navigateTo('DASHBOARD'); // ✓ Valid
navigateTo('INVALID_ROUTE'); // ✗ TypeScript error

// Programmatic back navigation
goBack();
```

### Test Contract
```typescript
describe('useTypedNavigate', () => {
  it('navigates to route using key')
  it('navigates to path using string')
  it('navigates back in history')
  it('passes options to navigate function')
  it('provides type safety for route keys')
})
```

---

## 3. Navigation Component

### Purpose
Top navigation menu with links and user info.

### Interface
```typescript
interface NavigationProps {
  className?: string; // Optional CSS class
}

const Navigation: React.FC<NavigationProps>
```

### Behavior
**Renders**:
- App branding/logo
- Navigation links (Dashboard, Tasks, Profile) with active states
- User display name
- Logout button

**Interactions**:
- Clicking nav link navigates to route
- Active route link has `active` CSS class
- Clicking logout triggers `authService.logout()`
- After logout, redirects to login page

**Visibility Rules**:
- Only renders when user is authenticated
- Shows user name from auth context
- All links visible to authenticated users

### Dependencies
- `useAuth()` for user state and logout
- `NavLink` from react-router-dom for links with active states
- `ROUTES` constants for paths

### Test Contract
```typescript
describe('Navigation', () => {
  it('renders navigation links when authenticated')
  it('does not render when not authenticated')
  it('highlights active route')
  it('displays user name')
  it('calls logout when logout button clicked')
  it('redirects to login after logout')
})
```

---

## 4. NotFound Page Component

### Purpose
404 error page for invalid routes.

### Interface
```typescript
const NotFound: React.FC
```

### Behavior
**Renders**:
- "404" error code
- Current attempted path (from useLocation)
- Helpful message
- Links to Dashboard and Home

**Interactions**:
- "Go to Dashboard" link navigates to /dashboard
- "Go to Home" link navigates to /

### Dependencies
- `useLocation()` to display attempted path
- `Link` from react-router-dom for navigation links
- `ROUTES` constants for paths

### Test Contract
```typescript
describe('NotFound', () => {
  it('displays 404 error message')
  it('shows attempted path')
  it('provides link to dashboard')
  it('provides link to home')
})
```

---

## 5. AuthProvider Component

### Purpose
Wraps app with authentication context.

### Interface
```typescript
interface AuthProviderProps {
  children: React.ReactNode;
}

const AuthProvider: React.FC<AuthProviderProps>
```

### Behavior
**Provides**:
- Authentication state to all child components via Context
- `user`, `isAuthenticated`, `isLoading`, `error`
- `login()`, `logout()`, `refreshSession()`, `validateSession()` methods

**Lifecycle**:
- On mount: Validates existing session
- Every 5 minutes: Refreshes session
- On unmount: Cleans up refresh interval

### Dependencies
- `useAuthState()` hook for state management
- `AuthContext` from useAuth hook

### Test Contract
```typescript
describe('AuthProvider', () => {
  it('provides auth context to children')
  it('validates session on mount')
  it('refreshes session periodically')
  it('exposes login and logout functions')
  it('updates auth state after login')
  it('updates auth state after logout')
})
```

---

## 6. ROUTES Configuration

### Purpose
Centralized route path constants.

### Interface
```typescript
const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  AUTH_CALLBACK: '/auth/callback',
  DASHBOARD: '/dashboard',
  TASKS: '/tasks',
  PROFILE: '/profile',
  NOT_FOUND: '*',
} as const;

type RouteKey = keyof typeof ROUTES;
type RoutePath = typeof ROUTES[RouteKey];
```

### Behavior
**Usage**:
- Import ROUTES constant
- Reference route paths using dot notation: `ROUTES.DASHBOARD`
- TypeScript ensures valid route keys
- Single source of truth for all paths

### Test Contract
```typescript
describe('ROUTES', () => {
  it('contains all defined routes')
  it('has unique path values')
  it('paths start with / or are *')
  it('type checking prevents invalid keys')
})
```

---

## 7. App Router Configuration

### Purpose
Main router setup in App.tsx.

### Structure
```typescript
<BrowserRouter>
  <AuthProvider>
    <Routes>
      {/* Root redirect */}
      <Route path="/" element={<Navigate to="/login" replace />} />

      {/* Public routes */}
      <Route path="/login" element={<Login />} />
      <Route path="/auth/callback" element={<AuthCallback />} />

      {/* Protected routes */}
      <Route element={<ProtectedRoute />}>
        <Route element={<MainLayout />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/tasks" element={<TaskList />} />
          <Route path="/profile" element={<Profile />} />
        </Route>
      </Route>

      {/* 404 catch-all */}
      <Route path="*" element={<NotFound />} />
    </Routes>
  </AuthProvider>
</BrowserRouter>
```

### Behavior
**Route Hierarchy**:
1. BrowserRouter (top-level)
2. AuthProvider (provides auth to all routes)
3. Routes container
4. Individual Route definitions

**Route Order**:
1. Exact path matches first
2. Protected routes (nested under ProtectedRoute)
3. Catch-all wildcard last

**Protected Route Pattern**:
- ProtectedRoute as element (not wrapper)
- Nested routes inside
- Optional MainLayout for consistent layout

### Test Contract
```typescript
describe('App Router', () => {
  it('redirects / to /login')
  it('renders login page at /login')
  it('renders auth callback at /auth/callback')
  it('protects dashboard route')
  it('protects tasks route')
  it('protects profile route')
  it('renders 404 for invalid paths')
})
```

---

## 8. MainLayout Component

### Purpose
Shared layout for protected pages with navigation.

### Interface
```typescript
const MainLayout: React.FC
```

### Behavior
**Renders**:
- Navigation component (top)
- Outlet for child routes (main content area)

**Structure**:
```typescript
<div className="main-layout">
  <Navigation />
  <main className="main-content">
    <Outlet />
  </main>
</div>
```

### Dependencies
- `Navigation` component
- `Outlet` from react-router-dom

### Test Contract
```typescript
describe('MainLayout', () => {
  it('renders navigation component')
  it('renders child routes in outlet')
  it('applies correct CSS classes')
})
```

---

## 9. Route Transition Behavior

### Navigation Events

**Browser Back Button**:
- Trigger: User clicks browser back
- Behavior: Navigate to previous route in history
- Protected route check: Re-run on transition

**Browser Forward Button**:
- Trigger: User clicks browser forward
- Behavior: Navigate to next route in history
- Protected route check: Re-run on transition

**Link Click**:
- Trigger: User clicks NavLink or Link component
- Behavior: Navigate to specified route
- History: Push new entry (unless `replace={true}`)

**Programmatic Navigation**:
- Trigger: Code calls `navigate()` or `navigateTo()`
- Behavior: Navigate to specified route
- History: Push new entry (unless `replace: true` in options)

### Redirect Behavior

**Unauthenticated Access to Protected Route**:
- Behavior: `<Navigate to="/login" replace />`
- History: Replace current entry (no back button to protected route)

**OAuth Redirect**:
- Before login: Store intended destination in sessionStorage
- After login: Navigate to intended destination
- History: Replace login page entry

---

## 10. Error Handling

### Invalid Route
- Behavior: Render NotFound component
- HTTP Status: N/A (client-side routing)
- User Action: Links to valid pages

### Auth Check Failure
- Behavior: Render loading spinner, then redirect to login
- Error State: Set in auth context
- Retry: Session refresh attempted automatically

### Navigation Errors
- Invalid path in navigate(): React Router throws error
- Protected route access: Redirect to login
- Concurrent navigations: Last navigation wins

---

## 11. Performance Contracts

### Route Transition Time
- **Target**: <100ms for cached routes
- **Measurement**: Time from navigation trigger to render complete
- **Optimization**: React Router handles efficiently

### Component Mount Time
- **Target**: <200ms for page components
- **Measurement**: Time from route match to first paint
- **Optimization**: Avoid heavy computation in render

### Bundle Size Impact
- **react-router-dom v6**: ~45KB gzipped
- **Additional code**: ~10KB (components, hooks, config)
- **Total impact**: ~55KB

---

## 12. Security Contracts

### Protected Route Enforcement
- **Guarantee**: Unauthenticated users cannot access protected pages
- **Mechanism**: ProtectedRoute component checks before render
- **Redirect**: Always use `replace` to prevent history manipulation

### OAuth State Validation
- **Guarantee**: OAuth callback validates state parameter
- **Mechanism**: Backend validates state, frontend trusts backend
- **Storage**: State stored in sessionStorage, cleared after use

### XSS Prevention
- **Guarantee**: No user-provided content in route config
- **Mechanism**: All routes defined statically in code
- **Dynamic Data**: Displayed via React props, automatically escaped

---

**Status**: ✅ Routing API contract complete