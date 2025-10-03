# Implementation Plan: Import Path Cleanup

**Branch**: `009-resolve-it-1` | **Date**: 2025-10-02 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/009-resolve-it-1/spec.md`

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
This feature resolves import path inconsistencies preventing backend compilation. The primary requirement is to fix undefined model references in `services/auth/oauth.go` and `internal/config/database.go`, deprecate legacy models, and reorganize the codebase following Domain-Driven Design (DDD) principles. The technical approach involves scanning the entire codebase for import issues, refactoring import paths to point to the new DDD structure, removing deprecated legacy models, and ensuring all tests pass after the refactoring.

## Technical Context
**Language/Version**: Go 1.24.7
**Primary Dependencies**: GORM ORM, Gin web framework, golang.org/x/oauth2
**Storage**: SQLite (development)
**Testing**: Go testing framework (go test), testify
**Target Platform**: Backend server (darwin/arm64 development, Linux deployment)
**Project Type**: Web application (frontend + backend)
**Performance Goals**: Build time < 10 seconds, no regression in test execution time
**Constraints**: Zero downtime refactoring (backward compatibility required), no functional changes
**Scale/Scope**: ~100 Go files, existing DDD structure in domain/ directory

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Note**: The constitution file is currently a template with placeholders. Proceeding with standard software engineering best practices:

### Standard Quality Gates
- [x] **No Breaking Changes**: Refactoring maintains existing functionality and API contracts
- [x] **Test Coverage**: All existing tests must continue to pass
- [x] **Documentation**: Import path changes documented in quickstart guide
- [x] **Code Quality**: Follow existing Go conventions and project structure
- [x] **Backward Compatibility**: Public APIs remain unchanged

### Complexity Assessment
- **Low Complexity**: This is a refactoring task with clear boundaries
- **No New Features**: Only reorganizing existing code
- **Well-Scoped**: Limited to import path corrections and model organization

## Project Structure

### Documentation (this feature)
```
specs/009-resolve-it-1/
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
├── cmd/
│   └── server/          # Main application entry point
├── domain/              # DDD domain layer (target structure)
│   ├── user/
│   │   ├── entities/
│   │   ├── valueobjects/
│   │   ├── repositories/
│   │   └── services/
│   └── task/
│       ├── entities/
│       ├── valueobjects/
│       ├── repositories/
│       └── services/
├── application/         # Application services
│   ├── user/
│   └── task/
├── infrastructure/      # Infrastructure layer
│   └── persistence/
├── presentation/        # Presentation layer
│   └── http/
├── internal/            # Internal packages (to be refactored)
│   ├── models/          # Legacy flat models (to be deprecated)
│   ├── services/
│   ├── handlers/
│   └── config/
├── services/            # Legacy services directory
│   └── auth/            # Has broken imports to models
├── models/              # Root-level legacy models (to be deprecated)
├── handlers/
├── middleware/
├── utils/
└── tests/
    ├── contract/
    ├── integration/
    └── unit/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/
```

**Structure Decision**: This is a web application with separate backend/ and frontend/ directories. The backend follows a Domain-Driven Design structure with domain/, application/, infrastructure/, and presentation/ layers. However, there are legacy directories (internal/models/, services/auth/, models/) that have inconsistent import paths. This feature will consolidate all models into the DDD domain structure and update all import references.

## Phase 0: Outline & Research

### Research Tasks
1. **Scan for all import path issues**:
   - Run `go build` to identify all compilation errors
   - Use `go list -json ./...` to analyze package dependencies
   - Grep for import statements referencing undefined packages
   - Document all files with broken imports

2. **Analyze current model structure**:
   - Catalog all model definitions across directories
   - Identify models in: domain/, internal/models/, models/
   - Map relationships between models
   - Determine which are "legacy" vs current DDD models

3. **Research DDD package organization best practices**:
   - Review Go DDD project structures
   - Identify entity vs value object patterns
   - Research repository pattern implementations
   - Find guidance on migration strategies

4. **Identify backward compatibility requirements**:
   - List all public API endpoints
   - Document model types exposed in APIs
   - Identify any external consumers
   - Plan alias/wrapper strategy if needed

### Output
✅ **Complete**: `research.md` created with comprehensive analysis:
- 14 compilation errors documented across 5 files
- 34 files identified with `todo-app/internal/models` imports
- Model duplication analysis (internal/models vs domain/)
- DDD best practices research
- Backward compatibility strategy
- Full impact analysis of ~50-60 files
- Decision: Create `domain/auth/` and `domain/health/` bounded contexts

## Phase 1: Design & Contracts
*Prerequisites: research.md complete ✅*

### 1. Data Model Design

✅ **Complete**: `data-model.md` created with:

**New Domain Contexts**:
- **Auth Domain** (NEW):
  - Entities: `AuthenticationSession`, `OAuthState`
  - Value Objects: `SessionToken`, `PKCEVerifier`, `StateToken`
  - Repositories: `SessionRepository`, `OAuthStateRepository`
  - Location: `domain/auth/`

- **Health Domain** (NEW):
  - Entities: `HealthStatus`
  - Location: `domain/health/entities/`

**Existing Domains** (No Changes):
- User Domain (`domain/user/`) - Well-established
- Task Domain (`domain/task/`) - Well-established

**Deprecated Models** (To Be Removed):
- `backend/models/session.go` → Move to `domain/auth/entities/authentication_session.go`
- `backend/models/oauth_state.go` → Move to `domain/auth/entities/oauth_state.go`
- `internal/models/user.go` → Delete (use `domain/user/entities/user.go`)
- `internal/models/task.go` → Delete (use `domain/task/entities/task.go`)

**Key Design Decisions**:
- GORM table names preserved (no database migration needed)
- JSON struct tags preserved (API compatibility maintained)
- Repository pattern follows existing conventions

---

### 2. API Contracts

✅ **Complete**: `contracts/import-paths.md` created with:

**Import Path Standards**:
```go
// Correct DDD imports
import "todo-app/domain/[context]/[layer]"

// Examples:
import "todo-app/domain/auth/entities"
import "todo-app/domain/user/valueobjects"
import "todo-app/domain/task/repositories"
```

**Forbidden Patterns** (Zero tolerance):
```go
// MUST NOT exist after refactoring
import "todo-app/internal/models"
import "todo-app/models"
import legacymodels "todo-app/internal/models"
```

**Verification Commands**:
- `grep -r "todo-app/internal/models" backend/` → Must return 0 matches
- `go build ./...` → Must succeed with 0 errors
- All 51+ tests → Must pass

**Contract Test**: Import path compliance test defined to enforce standards

---

### 3. Quickstart Guide

✅ **Complete**: `quickstart.md` created with:

**Quick Verification** (5 minutes):
1. Check branch: `009-resolve-it-1`
2. Verify build: `go build ./...` succeeds
3. Check forbidden imports: Zero matches
4. Run tests: All pass
5. Verify DDD structure exists
6. Verify legacy cleanup complete

**Detailed Verification** (15 minutes):
- Session management tests
- OAuth flow tests
- Integration tests
- Contract tests (API compatibility)
- Import path compliance script

**Acceptance Scenarios** (from spec.md):
1. ✅ Codebase scan identifies all import issues
2. ✅ Compilation succeeds without errors
3. ✅ All tests execute without failures
4. ✅ CI/CD build produces deployable artifacts

**Troubleshooting Guide**: Common issues and fixes documented

**Success Checklist**: 25-item verification checklist covering:
- Compilation (3 items)
- Import paths (4 items)
- Directory structure (4 items)
- Testing (4 items)
- Functionality (4 items)
- Documentation (6 items)

---

### 4. Agent Context Update

✅ **Complete**: `CLAUDE.md` updated incrementally via `update-agent-context.sh`:
- Added Go 1.24.7 to technologies
- Added GORM ORM, Gin framework, golang.org/x/oauth2
- Added SQLite (development) database
- Added to recent changes (keeps last 3 features)
- Preserved manual additions between markers

**File Size**: Kept under 150 lines for token efficiency

---

### 5. Test Strategy (No new tests for refactoring)

**Approach**: Rely on existing comprehensive test suite

**Existing Test Coverage**:
- 51+ tests across unit, integration, contract layers
- Contract tests validate API JSON compatibility
- Integration tests validate auth flows
- Unit tests validate entity behavior

**Test Execution Strategy**:
- Run tests after each import path change
- Ensure compilation succeeds at each step
- No new test files needed (refactoring only)

**Validation**:
- All existing tests must pass unchanged
- Zero new test failures introduced
- Test execution time must not regress

---

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

### Task Generation Strategy

The `/tasks` command will generate a sequenced task list following TDD and dependency order principles:

**Input Sources**:
1. `research.md` - List of 50-60 files requiring updates
2. `data-model.md` - New domain structure to create
3. `contracts/import-paths.md` - Validation criteria

**Task Categories**:

#### Category 1: Domain Structure Creation (Parallel)
- Create `domain/auth/entities/` directory
- Create `domain/auth/valueobjects/` directory
- Create `domain/auth/repositories/` directory
- Create `domain/auth/services/` directory
- Create `domain/health/entities/` directory

**Estimated**: 5 tasks, all `[P]` (parallel - independent directories)

#### Category 2: Model Migration (Sequential)
- Move `backend/models/session.go` → `domain/auth/entities/authentication_session.go`
- Move `backend/models/oauth_state.go` → `domain/auth/entities/oauth_state.go`
- Move `internal/models/health.go` → `domain/health/entities/health_status.go`
- Move `internal/models/google_identity.go` → `domain/auth/valueobjects/google_identity.go`
- Ensure GORM `TableName()` methods preserve table names

**Estimated**: 5 tasks, sequential (models depend on directories)

#### Category 3: High-Priority Import Fixes (Compilation Blockers)
Priority order from research.md:
1. Fix `services/auth/oauth.go` (11 errors)
2. Fix `services/auth/session.go` (10 errors)
3. Fix `jobs/session_cleanup.go` (10 errors)
4. Fix `jobs/oauth_cleanup.go` (4 errors)
5. Fix `internal/config/database.go` (2 errors - remove legacymodels)

**Estimated**: 5 tasks, sequential (verify compilation after each)

#### Category 4: Repository Interface Creation
- Create `domain/auth/repositories/session_repository.go` interface
- Create `domain/auth/repositories/oauth_state_repository.go` interface
- Create infrastructure implementations in `infrastructure/persistence/`

**Estimated**: 4 tasks, `[P]` (interfaces are independent)

#### Category 5: Remaining Import Updates (34 files)
Batch by directory for efficiency:
- Update `internal/handlers/` imports (3 files)
- Update `internal/services/` imports (3 files)
- Update `tests/contract/` imports (23 files)
- Update `tests/integration/` imports (13 files)
- Update `tests/unit/` imports (8 files)
- Update other files (4 files)

**Estimated**: 6-12 tasks, `[P]` where possible (tests can be updated in parallel)

#### Category 6: Legacy Cleanup
- Delete `backend/models/` directory (verify empty first)
- Delete `backend/internal/models/user.go`
- Delete `backend/internal/models/task.go`
- Delete `backend/services/auth/` directory (moved to domain)
- Clean up `backend/internal/models/` (delete if empty, keep only non-domain files)

**Estimated**: 5 tasks, sequential (only after all imports updated)

#### Category 7: Verification & Validation
- Run `go build ./...` (must succeed)
- Run `go vet ./...` (must pass)
- Run full test suite `go test ./...` (all pass)
- Run import compliance check script
- Verify directory structure matches data-model.md
- Manual smoke test (start server, test auth flow)

**Estimated**: 6 tasks, sequential (gates at end)

---

### Task Ordering Strategy

**Phase Order**: 1 → 2 → 3 → 4 → (5 || 6 partial) → 6 complete → 7

**Rationale**:
1. Create structure first (can't move files to non-existent directories)
2. Move models second (enables import updates)
3. Fix compilation blockers third (enables testing ASAP)
4. Create repositories fourth (needed by some imports)
5. Update remaining imports in parallel batches
6. Clean up legacy code only after all imports migrated
7. Validate everything works at the end

**Parallelization Markers**:
- `[P]` on directory creation tasks (independent)
- `[P]` on repository interface creation (independent)
- `[P]` on test file imports within same directory (share no state)
- Sequential for model moves and validation tasks

**TDD Adherence**:
- No new test files (refactoring only)
- Existing tests act as regression suite
- Each task must not break currently passing tests
- "Red-Green-Refactor": Refactor while keeping green (all tests pass)

---

### Estimated Task Count

| Category | Task Count | Parallelizable |
|----------|-----------|----------------|
| 1. Domain Structure | 5 | Yes (all) |
| 2. Model Migration | 5 | No (sequential) |
| 3. High-Priority Fixes | 5 | No (sequential) |
| 4. Repository Creation | 4 | Yes (all) |
| 5. Remaining Imports | 6-12 | Partial (by batch) |
| 6. Legacy Cleanup | 5 | No (sequential) |
| 7. Verification | 6 | No (sequential) |
| **TOTAL** | **36-42** | ~15 parallel, ~25 sequential |

**Execution Time Estimate** (with parallelization):
- Sequential tasks: ~2-3 minutes each → 50-75 minutes
- Parallel tasks (run in batches): ~5-10 minutes total
- Verification/testing: ~10-15 minutes
- **Total**: ~65-100 minutes (1-1.5 hours)

---

### Task Dependencies

```
Phase 1: Create Directories [P] [P] [P] [P] [P]
                ↓
Phase 2: Move Models [sequential, 5 tasks]
                ↓
Phase 3: Fix Compilation Blockers [sequential, 5 tasks]
                ↓
Phase 4: Create Repositories [P] [P] [P] [P]
                ↓
Phase 5: Update Remaining Imports [batched parallel]
                ↓
Phase 6: Delete Legacy Code [sequential, 5 tasks]
                ↓
Phase 7: Verify & Validate [sequential, 6 tasks]
```

---

### Success Criteria (from quickstart.md)

All tasks complete when:
- ✅ `go build ./...` succeeds (0 errors)
- ✅ Zero matches for forbidden import patterns
- ✅ At least 50 DDD-style imports exist
- ✅ `domain/auth/` and `domain/health/` exist with proper structure
- ✅ `backend/models/` deleted
- ✅ All 51+ tests pass
- ✅ Manual smoke test passes (auth flow works)

---

**IMPORTANT**: This phase is executed by the `/tasks` command, NOT by `/plan`.

The `/tasks` command will:
1. Load this plan.md
2. Load research.md for file list
3. Load data-model.md for structure
4. Generate tasks.md with 36-42 numbered tasks
5. Mark parallel tasks with `[P]`
6. Include verification commands for each task

---

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution - `/tasks` command creates `tasks.md` with detailed steps

**Phase 4**: Implementation - Execute each task in tasks.md:
- Follow constitutional principles (test-first where applicable)
- Maintain backward compatibility
- Verify compilation after each change
- Run tests frequently

**Phase 5**: Validation - Final verification:
- Run full quickstart.md verification script
- Perform manual smoke testing
- Update documentation if needed
- Create pull request with summary

---

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

**No violations identified**. This is a straightforward refactoring with no architectural complexity.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | N/A | N/A |

---

## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (N/A - none)

**Artifacts Generated**:
- [x] research.md (Phase 0)
- [x] data-model.md (Phase 1)
- [x] contracts/import-paths.md (Phase 1)
- [x] quickstart.md (Phase 1)
- [x] CLAUDE.md updated (Phase 1)
- [ ] tasks.md (Phase 2 - awaits /tasks command)

---

## Next Steps

**Ready for**: `/tasks` command

**What to do**:
1. Review this plan.md to understand the strategy
2. Review research.md for detailed findings
3. Review data-model.md for target architecture
4. Run `/tasks` to generate tasks.md
5. Execute tasks in order, marking [P] tasks for parallelization

**Expected Duration**: 1-1.5 hours for full implementation

---

*Implementation plan complete. Ready for task generation.*

*Based on Constitution template - See `/.specify/memory/constitution.md`*
