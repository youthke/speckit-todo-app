
# Implementation Plan: React Router Implementation

**Branch**: `006-react-router` | **Date**: 2025-09-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-react-router/spec.md`

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
Implement client-side routing using React Router to enable navigation between Login, Main Todo List, OAuth Callback, and 404 pages. Protected routes will redirect unauthenticated users to login. The system will provide visible navigation UI and support browser back/forward navigation while maintaining URL-based bookmarking.

## Technical Context
**Language/Version**: TypeScript 5.9.2, React 19.1.1
**Primary Dependencies**: React Router (version TBD - v6 or v7), React DOM 19.1.1, Vite 6.0.11
**Storage**: N/A (frontend routing only, no data persistence)
**Testing**: Vitest 3.1.6, @vitest/ui 3.1.6, jsdom 25.0.1
**Target Platform**: Modern web browsers (Chrome, Firefox, Safari, Edge)
**Project Type**: web (frontend + backend)
**Performance Goals**: <100ms route transition, instant navigation for cached routes
**Constraints**: Must work with existing auth flow, maintain browser history API compatibility
**Scale/Scope**: 4 primary routes (Login, Todo List, OAuth Callback, 404), visible navigation menu, protected route mechanism

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Status**: ✅ PASS (No constitution file violations detected)

Since the constitution file is a template placeholder, this routing feature follows general best practices:
- ✅ Uses established library (React Router) rather than custom routing
- ✅ Frontend-only change, no new backend services or complexity
- ✅ Test-driven approach will be followed (tests before implementation)
- ✅ Integration tests for routing flows
- ✅ Simple, focused scope (page-level routing only)

## Project Structure

### Documentation (this feature)
```
specs/006-react-router/
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
├── [No changes required for routing feature]

frontend/
├── src/
│   ├── components/
│   │   ├── auth/            # Existing auth components
│   │   ├── TaskForm/        # Existing task components
│   │   ├── TaskList/
│   │   ├── TaskItem/
│   │   ├── navigation/      # NEW: Navigation menu component
│   │   └── routes/          # NEW: Route protection components
│   ├── pages/
│   │   ├── Login.tsx        # Existing login page
│   │   ├── TodoList.tsx     # NEW: Main todo list page (refactored from App.tsx)
│   │   ├── NotFound.tsx     # NEW: 404 page
│   │   └── AuthCallback.tsx # Existing (move from components/auth/)
│   ├── hooks/
│   │   └── useAuth.ts       # Existing auth hook
│   ├── services/
│   │   ├── api.ts           # Existing API service
│   │   └── auth.ts          # Existing auth service
│   ├── App.tsx              # MODIFIED: Setup router, define routes
│   ├── index.tsx            # Existing entry point
│   └── types/               # Existing type definitions
└── tests/
    ├── routes/              # NEW: Route tests
    │   ├── navigation.test.tsx
    │   ├── protected-routes.test.tsx
    │   └── route-transitions.test.tsx
    └── components/          # Existing component tests
```

**Structure Decision**: Web application structure (frontend + backend). This feature only modifies the frontend with new routing components, pages directory, and navigation UI. The backend requires no changes as routing is purely client-side.

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

1. **Setup & Configuration Tasks**:
   - Install React Router v6 and testing dependencies
   - Update Vitest configuration for routing tests
   - Create route constants file (ROUTES)
   - Create test utilities (renderWithRouter)

2. **Core Component Tasks** (from routing-api.md contract):
   - Create ProtectedRoute component
   - Create AuthProvider wrapper component
   - Create MainLayout component with Outlet
   - Create NotFound (404) page component
   - Create Navigation menu component

3. **Test Tasks** (TDD approach):
   - Write ProtectedRoute component tests (3 scenarios)
   - Write Navigation component tests (6 scenarios)
   - Write NotFound page tests (4 scenarios)
   - Write route integration tests (protected routes)
   - Write route integration tests (navigation flows)
   - Write route integration tests (browser back/forward)

4. **Routing Configuration Tasks**:
   - Update App.tsx with BrowserRouter
   - Define route structure (public, protected, 404)
   - Wrap app with AuthProvider
   - Configure nested protected routes with MainLayout

5. **Page Refactoring Tasks**:
   - Refactor Login page to pages directory (if needed)
   - Move AuthCallback to pages directory
   - Create Dashboard page (refactor from App.tsx)
   - Create TaskList page (refactor from App.tsx)
   - Create Profile page (new)

6. **Hook Tasks**:
   - Create useTypedNavigate hook for type-safe navigation
   - Update useAuth hook to work with routing (if needed)

7. **Integration & Testing Tasks**:
   - Run all tests and verify they pass
   - Execute quickstart manual test scenarios
   - Fix any failing tests
   - Performance testing (route transition timing)

**Ordering Strategy**:
- **Phase 1**: Setup (tasks 1-4) → Can run in parallel [P]
- **Phase 2**: Write tests for components (tasks 5-9) → Sequential (TDD)
- **Phase 3**: Implement components to make tests pass (tasks 10-14) → Can run in parallel [P]
- **Phase 4**: Configure routing (tasks 15-18) → Sequential (dependency order)
- **Phase 5**: Refactor pages (tasks 19-23) → Can run in parallel [P]
- **Phase 6**: Utility hooks (tasks 24-25) → Can run in parallel [P]
- **Phase 7**: Integration (tasks 26-29) → Sequential (verification)

**Task Dependencies**:
```
Setup (1-4) → Tests (5-9) → Components (10-14)
                              ↓
                         Routing Config (15-18)
                              ↓
                         Page Refactoring (19-23)
                              ↓
                         Hooks (24-25)
                              ↓
                         Integration (26-29)
```

**Estimated Output**: ~28-32 numbered, ordered tasks in tasks.md

**Key Principles**:
- TDD: Write tests before implementation
- Parallel execution: Mark independent tasks with [P]
- Dependencies: Ensure components exist before configuring routes
- Verification: Final tasks validate entire feature

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

No complexity deviations or violations detected. Implementation follows standard React Router patterns with existing authentication system.


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [x] Phase 3: Tasks generated (/tasks command) - 32 tasks created
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved (React Router v6 chosen)
- [x] Complexity deviations documented (none)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
