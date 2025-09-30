# Research: React Router Implementation

**Feature**: React Router Implementation
**Date**: 2025-09-30
**Status**: Complete

## Overview
This document contains research findings for implementing React Router in a TypeScript React 19.1.1 + Vite application with existing OAuth authentication.

---

## 1. React Router Version

### Decision
**React Router v6 (latest v6.x, specifically v6.29.1)**

### Rationale
- **React 19.1.1 Compatibility**: Both v6 and v7 are compatible, but v6 is more stable and battle-tested
- **SPA-First Architecture**: Application is client-side SPA with no SSR requirements; v6 is purpose-built for SPAs while v7 focuses on framework mode with SSR
- **Minimal Configuration**: v6 requires no additional config files or build tool changes, works seamlessly with existing Vite setup
- **Existing Patterns**: Codebase already uses v6 patterns (`useNavigate`, `useSearchParams`)
- **Production Stability**: v6 has been production-ready since 2021 with extensive real-world usage
- **No Migration Work**: Staying with v6 avoids breaking changes while maintaining all needed functionality

### Alternatives Considered
- **React Router v7**:
  - Pros: Latest features, enhanced TypeScript support, React 19 optimizations, 15% smaller bundle
  - Cons: Requires replacing `@vitejs/plugin-react`, additional configuration, framework-mode overhead for SPA-only apps, newer with fewer production deployments
  - Rejected: Unnecessary complexity for pure SPA; framework mode adds overhead

- **TanStack Router**:
  - Pros: Excellent TypeScript support, type-safe routing
  - Cons: Different ecosystem, complete rewrite required, smaller community
  - Rejected: Unnecessary disruption for marginal benefits

---

## 2. Protected Routes Pattern

### Decision
**Outlet-based Protected Route Component**

```typescript
// Pseudocode structure
ProtectedRoute component:
  - Check authentication status from useAuth hook
  - If loading: Show loading spinner
  - If authenticated: Render Outlet (child routes)
  - If not authenticated: Navigate to login with replace
```

### Rationale
- Uses `Outlet` component, maintaining React Router v6 best practices
- Integrates with existing `useAuth` hook seamlessly
- Handles loading states to prevent flash of unauthenticated content
- Uses `replace` prop to avoid polluting browser history
- Separation of concerns: auth logic in hook, routing logic in component

### Alternatives Considered
- **HOC Pattern**: More verbose, less idiomatic in v6
- **Wrapper Component with Children**: Less flexible for nested routes

---

## 3. Type-Safe Route Definitions

### Decision
**Centralized Route Configuration with TypeScript Constants**

```typescript
// Pseudocode
ROUTES constant:
  HOME: '/'
  LOGIN: '/login'
  AUTH_CALLBACK: '/auth/callback'
  TASKS: '/tasks'
  NOT_FOUND: '*'

Type-safe navigation helper:
  getPath(key) -> returns route path
  navigate(key, navigateFunction) -> navigates to route
```

### Rationale
- Provides autocomplete and compile-time validation
- Single source of truth prevents typos
- Easy to refactor routes across application
- TypeScript's `as const` ensures literal types

### Alternatives Considered
- **React Router v7 typegen**: Would require v7 migration
- **Manual string literals**: Prone to typos, no autocomplete

---

## 4. Testing Approach

### Decision
**MemoryRouter with Custom Test Utilities using Vitest**

### Key Components
1. **Test Setup**: Configure jsdom environment, mock window.matchMedia
2. **Render Utility**: Wrap components with MemoryRouter + AuthContext for testing
3. **Mock Auth Context**: Provide configurable auth state for tests

### Rationale
- `MemoryRouter` provides full routing context without browser dependencies
- Custom render utility eliminates test boilerplate
- Mocks essential browser APIs to prevent errors
- Uses Vitest's native mocking
- `@testing-library/jest-dom` provides familiar matchers

### Testing Patterns
- Route navigation: Use `initialEntries` to set starting route
- Protected routes: Mock different auth states
- Link testing: Use `userEvent.click()` to test navigation
- useNavigate hook: Mock and assert correct calls

### Alternatives Considered
- **Browser Mode Testing**: More complex, slower tests
- **Manual Router Context**: More verbose and error-prone

---

## 5. Auth Flow Integration

### Decision
**HTTP-Only Cookie-Based Authentication with Credentials**

### Current Implementation
- Uses `credentials: 'include'` to send HTTP-only cookies
- Session-based authentication (no localStorage/sessionStorage)
- OAuth flow with Google via backend proxy

### Integration Pattern
```typescript
// Pseudocode
App structure:
  BrowserRouter
    AuthProvider (wraps all routes with auth context)
      Routes
        Public routes (login, callback)
        Protected routes (wrapped with ProtectedRoute)
          Protected pages (tasks, profile)
        404 catch-all
```

### Rationale
- **Security**: HTTP-only cookies prevent XSS attacks
- **Seamless Integration**: Existing `useAuth` hook provides all needed state
- **Automatic Session Management**: Hook validates sessions on mount and refreshes periodically
- **Redirect Handling**: OAuth redirect destination stored in sessionStorage

### Alternatives Considered
- **localStorage for JWT**: Vulnerable to XSS
- **Token in Redux/Zustand**: Unnecessary overhead
- **Authorization Header**: Requires manual token management

---

## 6. Navigation Menu Component

### Decision
**NavLink with Active State Styling**

### Key Features
- Uses `NavLink` component with callback className prop
- Receives `{ isActive }` for conditional styling
- Type-safe navigation with custom hook
- User info display and logout button

### Rationale
- `NavLink` automatically adds active class when route matches
- Callback prop provides flexible conditional styling
- Type-safe hook prevents navigation to non-existent routes
- Clean separation of navigation logic

### Alternatives Considered
- **useLocation + Link**: More manual, requires custom active state logic
- **CSS-only**: Less flexible, doesn't work with dynamic routes

---

## 7. 404 Handling

### Decision
**Wildcard Route with Custom NotFound Component**

### Implementation
```typescript
// Route configuration
<Route path="*" element={<NotFound />} />

// NotFound component shows:
// - Error code (404)
// - Attempted path
// - Links to common destinations (home, dashboard)
```

### Rationale
- Wildcard path (`*`) catches all unmatched routes
- Must be placed last in route configuration
- `useLocation` displays attempted path for clarity
- Clear actions for quick user recovery
- Consistent styling with app design

### Alternatives Considered
- **Redirect to Home**: Loses error information
- **Error Boundary**: For runtime errors, not 404s
- **Nested 404s**: For section-specific errors (not needed yet)

---

## 8. Dependencies

### Required Packages
```
react-router-dom@^6.29.1
@types/react-router-dom (dev)
```

### Testing Dependencies
```
@testing-library/react@^16.1.0 (dev)
@testing-library/user-event@^14.5.2 (dev)
@testing-library/jest-dom@^6.6.3 (dev)
```

### Existing Compatible Dependencies
- react: ^19.1.1 ✓
- react-dom: ^19.1.1 ✓
- typescript: ^5.9.2 ✓
- vite: ^6.0.11 ✓
- vitest: ^3.1.6 ✓

---

## 9. Performance Considerations

### Route Transition Performance
- **Target**: <100ms route transition
- **Strategy**: React Router v6 handles transitions efficiently for SPA
- **Optimization**: Lazy loading for large pages (optional, not needed initially)

### Code Splitting
- **When**: Consider after app grows (pages >50KB)
- **Pattern**: Use React.lazy() + Suspense for route-based splitting
- **Current Assessment**: App is small enough that code splitting isn't necessary yet

---

## 10. Security Considerations

1. **HTTP-Only Cookies**: Current implementation uses `credentials: 'include'` ✓
2. **No Token in localStorage**: Avoids XSS vulnerability ✓
3. **CSRF Protection**: Ensure backend implements CSRF tokens for state-changing requests
4. **Redirect URL Validation**: OAuth redirect URLs should be validated to prevent open redirects

---

## 11. Migration Timeline

**Phase 1: Setup** (1-2 hours)
- Install dependencies
- Create route configuration
- Create ProtectedRoute component
- Create AuthProvider component
- Update Vitest configuration

**Phase 2: Core Routing** (2-3 hours)
- Wrap app with BrowserRouter
- Define main routes
- Implement protected routes
- Add 404 catch-all

**Phase 3: Navigation UI** (1-2 hours)
- Create Navigation component
- Add active state styling
- Implement logout

**Phase 4: Testing** (2-3 hours)
- Setup test utilities
- Write component tests
- Write integration tests

**Phase 5: Polish** (1 hour)
- Add loading states
- Performance testing

**Total Estimated Time**: 7-11 hours

---

## 12. Key Decisions Summary

| Area | Choice | Primary Reason |
|------|--------|----------------|
| Router Version | React Router v6 | SPA-first, stable, already used |
| Protected Routes | Outlet-based | v6 best practice, flexible |
| Route Definitions | TypeScript constants | Type safety, autocomplete |
| Testing | Vitest + Testing Library | Already configured, modern |
| Auth Integration | HTTP-only cookies + Context | Security, existing pattern |
| Navigation | NavLink with callbacks | Built-in active state |
| 404 Handling | Wildcard route | Standard, clear UX |

---

## Resources

- React Router v6 Docs: https://reactrouter.com/
- TypeScript Guide: https://reactrouter.com/explanation/type-safety
- Testing Guide: https://testing-library.com/docs/example-react-router/
- Vitest Docs: https://vitest.dev/

---

**Status**: ✅ All unknowns resolved, ready for Phase 1 (Design & Contracts)