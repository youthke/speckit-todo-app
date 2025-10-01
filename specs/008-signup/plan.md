# Implementation Plan: Signup Page

**Branch**: `008-signup` | **Date**: 2025-10-01 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/008-signup/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path ✓
2. Fill Technical Context (scan for NEEDS CLARIFICATION) ✓
   → Project Type: web (frontend + backend)
   → Structure Decision: backend/ + frontend/ layout
3. Fill Constitution Check section (template placeholder)
4. Evaluate Constitution Check section
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, CLAUDE.md
7. Re-evaluate Constitution Check section
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 9. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Implement a Google OAuth signup page for the todo-app that allows new users to create accounts and existing users to automatically log in. The feature leverages the existing Google OAuth infrastructure (feature 007-google) but modifies the behavior to auto-login existing users instead of redirecting them to a separate login page. The signup page includes rate limiting per IP address and requires email from Google to complete registration.

## Technical Context
**Language/Version**: Backend: Go 1.24.0, Frontend: TypeScript 5.9.2 + React 19.1.1
**Primary Dependencies**: Backend: Gin web framework, GORM ORM, Google OAuth 2.0 libraries (`golang.org/x/oauth2`); Frontend: React Router 6.30.1, Axios 1.12.2, Vite 6.0.11
**Storage**: SQLite (development), existing database schema with users table and Google OAuth support
**Testing**: Backend: Go testing framework with testify; Frontend: React Testing Library
**Target Platform**: Web application (backend server + React SPA)
**Project Type**: web - frontend + backend architecture
**Performance Goals**: OAuth flow completion < 3 seconds, rate limiter should handle 100 req/s
**Constraints**: Must reuse existing Google OAuth service (007-google), maintain backward compatibility with login flow
**Scale/Scope**: Signup page component + backend handler modification + rate limiter middleware

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Constitution Status**: Template placeholder - no specific constitutional requirements defined for this project.

**Initial Assessment**: PASS
- No custom constitutional principles to validate
- Following standard web application patterns
- Reusing existing authentication infrastructure

## Project Structure

### Documentation (this feature)
```
specs/008-signup/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
backend/
├── internal/
│   ├── handlers/        # google_oauth_handler.go (MODIFY)
│   ├── services/        # google_oauth_service.go (existing)
│   └── models/          # user.go, google_identity.go (existing)
├── middleware/          # rate_limiter.go (NEW)
└── tests/
    ├── contract/        # signup API contract tests
    └── integration/     # signup flow integration tests

frontend/
├── src/
│   ├── pages/           # SignupPage.tsx (NEW)
│   ├── components/      # GoogleSignupButton.tsx (exists, may need updates)
│   └── services/        # auth.ts (existing)
└── tests/               # signup page component tests
```

**Structure Decision**: Using existing web application structure (backend/ + frontend/). The signup feature integrates into the current Google OAuth implementation (feature 007-google) with minimal changes to existing handlers and services.

## Phase 0: Outline & Research

### Research Tasks
1. **Review existing Google OAuth implementation (007-google)**:
   - Understand current flow in `google_oauth_handler.go`
   - Identify modification points for auto-login behavior
   - Review session management in `SessionService`

2. **Rate limiting patterns for Gin framework**:
   - Research IP-based rate limiting middleware
   - Best practices for rate limit configuration (e.g., 10 attempts per 15 minutes)
   - Storage mechanism for rate limit counters (in-memory vs database)

3. **Frontend signup page patterns**:
   - Review existing components (GoogleSignupButton.tsx)
   - React Router integration for `/signup` route
   - Error handling and user feedback patterns

4. **Email requirement validation**:
   - Verify Google OAuth email field is always present
   - Review existing email verification logic (line 90 in handler)
   - Error messaging for missing email

### Research Output Format
**Output**: `research.md` with decisions:
- **Decision**: Modify existing `GoogleCallback` handler to auto-login existing users
- **Rationale**: Minimal code changes, reuses existing session creation logic
- **Alternatives considered**: Create separate signup endpoint vs. modify existing handler

- **Decision**: Use golang.org/x/time/rate for IP-based rate limiting
- **Rationale**: Standard library, efficient token bucket algorithm
- **Alternatives considered**: Redis-based rate limiter (overkill for current scale)

- **Decision**: Create dedicated `/signup` page mirroring login page structure
- **Rationale**: Consistent UX, leverages existing components
- **Alternatives considered**: Modal-based signup (doesn't fit navigation model)

## Phase 1: Design & Contracts

### 1. Data Model (`data-model.md`)

**No new entities required** - feature uses existing entities:
- `User` (existing): email, google_id, name, profile_picture_url, created_at, updated_at
- `GoogleIdentity` (existing): google_user_id, email, name, picture, email_verified
- `Session` (existing): user_id, token, expires_at

**New fields**: None - existing schema supports all requirements

**State transitions**:
1. New user: Google OAuth → Create User → Create Session → Redirect to app
2. Existing user: Google OAuth → Find User → Create Session → Redirect to app
3. Rate limited: Check IP → Reject with 429 → User waits

### 2. API Contracts (`contracts/signup-api.yaml`)

**Modified Endpoint** (existing endpoint behavior change):
```yaml
GET /api/v1/auth/google/login
  Query: none
  Response: 302 Redirect to Google OAuth

GET /api/v1/auth/google/callback
  Query:
    - code: string (OAuth code)
    - state: string (CSRF token)
  Response:
    - 302 Redirect to http://localhost:3000/ (success, auto-login)
    - 302 Redirect to http://localhost:3000/signup?error=... (failure)
  Behavior Change: Auto-login existing users instead of redirecting to /login

  Errors:
    - authentication_failed: Generic OAuth error
    - rate_limit_exceeded: Too many signup attempts from IP
```

**Rate Limiting**:
- Apply to `/api/v1/auth/google/login` endpoint
- Limit: 10 attempts per 15 minutes per IP address
- Response: HTTP 429 with Retry-After header

### 3. Contract Tests

Generate contract tests in `backend/tests/contract/`:
- `signup_google_oauth_test.go`: Test OAuth initiation
- `signup_callback_new_user_test.go`: Test new user signup flow
- `signup_callback_existing_user_test.go`: Test existing user auto-login (NEW behavior)
- `signup_rate_limit_test.go`: Test rate limiting enforcement

### 4. Integration Tests

Generate integration tests in `backend/tests/integration/`:
- `signup_flow_success_test.go`: End-to-end new user signup
- `signup_flow_existing_user_test.go`: End-to-end existing user auto-login
- `signup_flow_rate_limited_test.go`: Rate limit enforcement
- `signup_flow_missing_email_test.go`: Google returns no email

Frontend tests in `frontend/src/test/`:
- `SignupPage.test.tsx`: Signup page rendering and interactions
- `signup-navigation.test.tsx`: Navigation to/from signup page

### 5. Quickstart (`quickstart.md`)

**Quickstart test scenario**:
1. Start backend server and frontend dev server
2. Navigate to `http://localhost:3000/signup`
3. Click "Sign up with Google"
4. Authorize with Google account (new user)
5. Verify redirect to `http://localhost:3000/` with active session
6. Logout
7. Navigate to `http://localhost:3000/signup` again
8. Click "Sign up with Google" with same account
9. Verify automatic login and redirect to `http://localhost:3000/`

### 6. Update CLAUDE.md

Run update script:
```bash
.specify/scripts/bash/update-agent-context.sh claude
```

This will add to CLAUDE.md:
- Technology: Google OAuth signup (feature 008-signup)
- Commands: Testing and running signup flow
- Recent changes: Added signup page with auto-login for existing users

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
1. **Backend modifications** (priority: high):
   - Task: Modify `GoogleCallback` handler to auto-login existing users (remove line 107 redirect, continue to session creation)
   - Task: Create rate limiter middleware for IP-based throttling
   - Task: Register rate limiter on OAuth endpoints
   - Task: Update error responses for rate limiting

2. **Frontend additions** (priority: high):
   - Task: Create `SignupPage.tsx` component
   - Task: Add `/signup` route to React Router
   - Task: Update `GoogleSignupButton` component if needed
   - Task: Add error message handling for rate limits

3. **Testing** (priority: high, TDD order):
   - Task: Write contract tests for modified callback behavior [P]
   - Task: Write contract test for rate limiting [P]
   - Task: Write integration test for existing user auto-login
   - Task: Write integration test for rate limit enforcement
   - Task: Write frontend tests for signup page [P]

4. **Documentation** (priority: medium):
   - Task: Update API documentation with rate limiting details
   - Task: Create quickstart.md with manual test steps

**Ordering Strategy**:
- TDD order: Write tests first for new/modified behavior
- Backend before frontend: API changes must be stable before UI consumes them
- Independent tasks marked [P] for parallel execution

**Estimated Output**: 15-20 numbered, ordered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)
**Phase 4**: Implementation (execute tasks.md following constitutional principles)
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

No constitutional violations - using standard patterns and existing infrastructure.

## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command) - research.md created
- [x] Phase 1: Design complete (/plan command) - data-model.md, contracts/signup-api.yaml, quickstart.md, CLAUDE.md updated
- [x] Phase 2: Task planning complete (/plan command - describe approach only) - documented in plan.md
- [x] Phase 3: Tasks generated (/tasks command) - tasks.md created with 24 tasks
- [ ] Phase 4: Implementation complete - NEXT STEP: Execute tasks T001-T024
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS (template placeholder, no violations)
- [x] Post-Design Constitution Check: PASS (no new violations)
- [x] All NEEDS CLARIFICATION resolved (all technical decisions made in research.md)
- [x] Complexity deviations documented (none - using standard patterns)

---
*Based on Constitution v2.1.1 (template) - See `/memory/constitution.md`*
