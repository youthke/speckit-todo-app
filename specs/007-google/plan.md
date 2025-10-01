
# Implementation Plan: Google Account Signup

**Branch**: `007-google` | **Date**: 2025-10-01 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/007-google/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Enable users to sign up for the todo-app using their Google account credentials via OAuth 2.0 flow. Upon successful authentication, the system creates a user account linked to their Google identity, extracts their email address, and establishes a 7-day session. The feature handles duplicate signups by redirecting to login, validates email verification, and provides standardized error handling.

## Technical Context
**Language/Version**: Backend: Go 1.24.0, Frontend: TypeScript 5.9.2 + React 19.1.1
**Primary Dependencies**: Backend: Gin web framework, GORM ORM, Google OAuth 2.0 libraries; Frontend: React Router 6.30.1, Axios 1.12.2, Vite 6.0.11
**Storage**: SQLite database (development), existing schema to be extended with Google OAuth entities
**Testing**: Backend: testify framework, Frontend: Vitest 3.1.6 + Testing Library
**Target Platform**: Web application (browser-based frontend + HTTP API backend)
**Project Type**: web (frontend + backend)
**Performance Goals**: OAuth flow completion <3 seconds, session validation <50ms
**Constraints**: Email verification required, 7-day session duration, generic error messages only
**Scale/Scope**: Small-to-medium scale todo app, adding OAuth as alternative signup method alongside existing authentication

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Note**: The constitution file is currently a template placeholder with no specific principles defined. Proceeding with standard software engineering best practices:

✅ **Standard Best Practices Applied**:
- Clean separation of concerns (OAuth logic in dedicated service layer)
- Test-driven development approach (tests before implementation)
- Security-first design (email verification, secure session management)
- RESTful API design patterns
- Error handling and validation at boundaries
- Database migrations for schema changes

**Status**: PASS (no constitutional violations - constitution not yet defined for this project)

## Project Structure

### Documentation (this feature)
```
specs/007-google/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
│   └── google-oauth-api.yaml
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
backend/
├── models/
│   ├── user.go                    # Extended with OAuth fields
│   └── google_identity.go         # New: Google account link
├── services/
│   └── google_oauth_service.go    # New: OAuth flow logic
├── handlers/
│   └── auth_handler.go            # Extended with OAuth endpoints
├── middleware/
│   └── session_middleware.go      # Extended for 7-day sessions
├── migrations/
│   └── 00X_add_google_oauth.sql   # New: Schema migration
└── tests/
    ├── contract/
    │   └── google_oauth_test.go   # New: Contract tests
    ├── integration/
    │   └── signup_flow_test.go    # New: End-to-end tests
    └── unit/
        └── google_oauth_service_test.go

frontend/
├── src/
│   ├── components/
│   │   └── GoogleSignupButton.tsx # New: OAuth trigger UI
│   ├── pages/
│   │   ├── SignupPage.tsx         # Extended with Google option
│   │   └── OAuthCallbackPage.tsx  # New: Handle OAuth redirect
│   └── services/
│       └── authService.ts         # Extended with Google OAuth
└── tests/
    └── components/
        └── GoogleSignupButton.test.tsx
```

**Structure Decision**: Web application structure selected based on existing `backend/` and `frontend/` directories. Backend follows domain-driven design with clear separation between models, services, handlers, and middleware. Frontend uses component-based React architecture with service layer for API communication.

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh claude`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
The `/tasks` command will generate a dependency-ordered task list from the design artifacts:

1. **From data-model.md**:
   - Database migration task (create `google_identities` table, extend `users` table)
   - Go model definition tasks (GoogleIdentity, User extension)

2. **From contracts/google-oauth-api.yaml**:
   - Contract test tasks for each endpoint:
     - `GET /api/auth/google/login` [P]
     - `GET /api/auth/google/callback` [P]
     - `GET /api/auth/me` (verify Google OAuth compatibility) [P]

3. **From research.md**:
   - OAuth service implementation tasks:
     - Config setup (environment variables)
     - Token exchange logic
     - User creation logic
     - Session management

4. **From quickstart.md**:
   - Integration test tasks for each scenario:
     - Scenario 1: Successful signup [P]
     - Scenario 2: Duplicate prevention [P]
     - Scenario 3: Unverified email rejection [P]
     - Scenario 4: OAuth error handling [P]
     - Scenario 5: Session expiration validation [P]

5. **Frontend tasks**:
   - Component: GoogleSignupButton.tsx [P]
   - Page: OAuthCallbackPage.tsx [P]
   - Service: authService.ts extension
   - Component tests [P]

**Ordering Strategy**:
1. **Phase 0 - Setup**: Environment variables, dependencies
2. **Phase 1 - Database**: Migration scripts (blocking for all)
3. **Phase 2 - Models**: Go structs for entities [P - parallel safe]
4. **Phase 3 - Tests First (TDD)**:
   - Contract tests [P]
   - Integration tests [P]
   - All tests should FAIL initially
5. **Phase 4 - Backend Implementation**:
   - OAuth service
   - Auth handlers
   - Session middleware
6. **Phase 5 - Frontend Implementation**:
   - Components [P]
   - Pages [P]
   - Service layer
7. **Phase 6 - Validation**:
   - Run all tests (should PASS)
   - Execute quickstart scenarios

**Dependency Graph**:
```
Environment Setup
    ↓
Database Migration
    ↓
┌───────────────┬───────────────┐
Models [P]      Tests [P]       │
    ↓               ↓           │
Backend Impl    Frontend [P] ←─┘
    ↓               ↓
    └───→ Validation ←───┘
```

**Estimated Output**: 30-35 numbered tasks in tasks.md

**Task Template Structure**:
- Task number and description
- Dependencies (which tasks must complete first)
- Acceptance criteria (how to verify completion)
- Files to create/modify
- [P] marker for parallel-safe tasks

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [x] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (none)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
