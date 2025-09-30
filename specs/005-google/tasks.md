# Tasks: Google Account Login

**Input**: Design documents from `/specs/005-google/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Web app**: `backend/src/`, `frontend/src/`
- Paths shown below match implementation plan structure

## Phase 3.1: Setup
- [x] T001 Install Google OAuth 2.0 dependencies: `golang.org/x/oauth2` and `golang.org/x/oauth2/google` in backend/go.mod
- [x] T002 Add JWT library dependency `github.com/golang-jwt/jwt/v5` for session management in backend/go.mod
- [x] T003 [P] Configure OAuth environment variables in backend/.env (GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, JWT_SECRET)

## Phase 3.2: Database Schema & Models (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T004 [P] Database migration for users table OAuth fields in backend/migrations/005_add_oauth_to_users.sql
- [x] T005 [P] Database migration for authentication_sessions table in backend/migrations/006_create_authentication_sessions.sql
- [x] T006 [P] Database migration for oauth_states table in backend/migrations/007_create_oauth_states.sql
- [x] T007 [P] User model unit test for OAuth fields validation in backend/tests/unit/models/user_test.go
- [x] T008 [P] AuthenticationSession model unit test in backend/tests/unit/models/session_test.go
- [x] T009 [P] OAuthState model unit test in backend/tests/unit/models/oauth_state_test.go

## Phase 3.3: API Contract Tests (TDD) ⚠️ MUST COMPLETE BEFORE 3.4
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T010 [P] Contract test GET /auth/google/login in backend/tests/contract/auth_google_login_test.go
- [x] T011 [P] Contract test GET /auth/google/callback in backend/tests/contract/auth_google_callback_test.go
- [x] T012 [P] Contract test GET /auth/session/validate in backend/tests/contract/auth_session_validate_test.go
- [x] T013 [P] Contract test POST /auth/session/refresh in backend/tests/contract/auth_session_refresh_test.go
- [x] T014 [P] Contract test POST /auth/logout in backend/tests/contract/auth_logout_test.go
- [x] T015 [P] Contract test POST /auth/revoke-webhook in backend/tests/contract/auth_revoke_webhook_test.go

## Phase 3.4: Integration Tests (TDD) ⚠️ MUST COMPLETE BEFORE 3.5
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T016 [P] Integration test: Google OAuth flow for new user in backend/tests/integration/oauth_new_user_test.go
- [x] T017 [P] Integration test: Google OAuth account linking for existing user in backend/tests/integration/oauth_account_linking_test.go
- [x] T018 [P] Integration test: Session management and automatic refresh in backend/tests/integration/session_management_test.go
- [x] T019 [P] Integration test: OAuth access revocation handling in backend/tests/integration/oauth_revocation_test.go
- [x] T020 [P] Integration test: Error handling for Google service unavailable in backend/tests/integration/oauth_error_handling_test.go

## Phase 3.5: Core Models Implementation (ONLY after tests are failing)
- [x] T021 [P] User model with OAuth fields in backend/models/user.go
- [x] T022 [P] AuthenticationSession model in backend/models/session.go
- [x] T023 [P] OAuthState model in backend/models/oauth_state.go

## Phase 3.6: OAuth Services Implementation
- [x] T024 [P] Google OAuth configuration service in backend/services/auth/google.go
- [x] T025 [P] OAuth flow service (initiate, callback, token exchange) in backend/services/auth/oauth.go
- [x] T026 [P] Session management service (create, validate, refresh) in backend/services/auth/session.go
- [x] T027 [P] JWT token service (generate, validate, refresh) in backend/services/auth/jwt.go
- [x] T028 User service extension for OAuth account linking in backend/services/user/user.go

## Phase 3.7: API Handlers Implementation
- [x] T029 GET /auth/google/login handler (initiate OAuth flow) in backend/handlers/auth.go
- [x] T030 GET /auth/google/callback handler (process OAuth callback) in backend/handlers/auth.go
- [x] T031 GET /auth/session/validate handler (validate current session) in backend/handlers/auth.go
- [x] T032 POST /auth/session/refresh handler (refresh OAuth tokens) in backend/handlers/auth.go
- [x] T033 POST /auth/logout handler (terminate session) in backend/handlers/auth.go
- [x] T034 POST /auth/revoke-webhook handler (handle Google revocation) in backend/handlers/auth.go

## Phase 3.8: Middleware & Security
- [x] T035 OAuth session validation middleware in backend/middleware/auth.go
- [x] T036 CSRF protection for OAuth state validation in backend/middleware/security.go
- [x] T037 Rate limiting for OAuth endpoints in backend/middleware/rate_limit.go
- [x] T038 Secure cookie configuration for session tokens in backend/utils/cookies.go

## Phase 3.9: Frontend Implementation
- [x] T039 [P] Google OAuth login button component in frontend/src/components/auth/GoogleLoginButton.tsx
- [x] T040 [P] OAuth callback handler component in frontend/src/components/auth/AuthCallback.tsx
- [x] T041 [P] Authentication service for OAuth flow in frontend/src/services/auth.ts
- [x] T042 [P] Session management hooks in frontend/src/hooks/useAuth.ts
- [x] T043 Login page integration with Google OAuth in frontend/src/pages/Login.tsx

## Phase 3.10: Frontend Tests
- [ ] T044 [P] GoogleLoginButton component test in frontend/src/components/auth/__tests__/GoogleLoginButton.test.tsx
- [ ] T045 [P] AuthCallback component test in frontend/src/components/auth/__tests__/AuthCallback.test.tsx
- [ ] T046 [P] Auth service unit tests in frontend/src/services/__tests__/auth.test.ts
- [ ] T047 [P] Integration test for complete OAuth flow in frontend/src/tests/integration/oauth-flow.test.ts

## Phase 3.11: Database Integration & Utilities
- [x] T048 Database connection configuration for OAuth tables in backend/config/database.go
- [x] T049 OAuth state cleanup job (remove expired states) in backend/jobs/oauth_cleanup.go
- [x] T050 Session cleanup job (remove expired sessions) in backend/jobs/session_cleanup.go
- [x] T051 Token encryption/decryption utilities in backend/utils/crypto.go

## Phase 3.12: Polish & Performance
- [ ] T052 [P] Performance optimization for OAuth flow (<500ms target) in backend/services/auth/
- [ ] T053 [P] Session validation performance optimization (<200ms target) in backend/middleware/auth.go
- [ ] T054 [P] Add comprehensive logging for OAuth events in backend/utils/logger.go
- [ ] T055 [P] Add metrics collection for OAuth success/failure rates in backend/utils/metrics.go
- [ ] T056 [P] Security audit for token storage and transmission in backend/security/
- [ ] T057 [P] Update API documentation with OAuth endpoints in docs/api.md
- [ ] T058 Manual testing using quickstart.md validation scenarios

## Dependencies
- Database & Models (T004-T009) before Contract Tests (T010-T015)
- Contract Tests (T010-T015) before Integration Tests (T016-T020)
- Integration Tests (T016-T020) before Implementation (T021-T051)
- Core Models (T021-T023) before Services (T024-T028)
- Services (T024-T028) before Handlers (T029-T034)
- T035 (auth middleware) blocks T031-T034 (protected endpoints)
- Backend implementation (T021-T038) before Frontend (T039-T043)
- Implementation (T021-T051) before Polish (T052-T058)

## Parallel Example
```bash
# Launch database migration tasks together (Phase 3.2):
Task: "Database migration for users table OAuth fields in backend/migrations/005_add_oauth_to_users.sql"
Task: "Database migration for authentication_sessions table in backend/migrations/006_create_authentication_sessions.sql"
Task: "Database migration for oauth_states table in backend/migrations/007_create_oauth_states.sql"

# Launch contract test tasks together (Phase 3.3):
Task: "Contract test GET /auth/google/login in backend/tests/contract/auth_google_login_test.go"
Task: "Contract test GET /auth/google/callback in backend/tests/contract/auth_google_callback_test.go"
Task: "Contract test GET /auth/session/validate in backend/tests/contract/auth_session_validate_test.go"

# Launch model implementation tasks together (Phase 3.5):
Task: "User model with OAuth fields in backend/models/user.go"
Task: "AuthenticationSession model in backend/models/session.go"
Task: "OAuthState model in backend/models/oauth_state.go"

# Launch service implementation tasks together (Phase 3.6):
Task: "Google OAuth configuration service in backend/services/auth/google.go"
Task: "OAuth flow service (initiate, callback, token exchange) in backend/services/auth/oauth.go"
Task: "Session management service (create, validate, refresh) in backend/services/auth/session.go"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify all tests fail before implementing corresponding features
- Follow TDD strictly: Red → Green → Refactor
- Commit after each completed task
- Run full test suite before proceeding to next phase
- Environment variables must be configured before running OAuth tests
- Database migrations must be applied before running any database-related tests

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts** (auth-api.yaml):
   - 6 endpoints → 6 contract test tasks [P] (T010-T015)
   - 6 endpoints → 6 implementation tasks (T029-T034)

2. **From Data Model**:
   - 3 entities → 3 model tasks [P] (T021-T023)
   - 3 entities → 3 migration tasks [P] (T004-T006)

3. **From User Stories** (quickstart.md):
   - 5 test scenarios → 5 integration tests [P] (T016-T020)
   - Manual validation → quickstart execution task (T058)

4. **Ordering**:
   - Setup → Database → Tests → Models → Services → Handlers → Integration → Polish
   - Dependencies strictly enforced for TDD compliance

## Validation Checklist
*GATE: Checked before task execution*

- [x] All contracts have corresponding tests (T010-T015 cover auth-api.yaml)
- [x] All entities have model tasks (T021-T023 cover User, Session, OAuthState)
- [x] All tests come before implementation (T004-T020 before T021+)
- [x] Parallel tasks truly independent (different files, no shared dependencies)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] TDD workflow enforced (tests must fail before implementation)
- [x] All quickstart scenarios covered by integration tests