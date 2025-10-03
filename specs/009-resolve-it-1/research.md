# Research: Import Path Cleanup

**Date**: 2025-10-02
**Feature**: Import Path Cleanup (009-resolve-it-1)

## Executive Summary

Analysis of the todo-app backend codebase reveals systematic import path issues preventing compilation. The root cause is a partial migration to Domain-Driven Design (DDD) structure, leaving orphaned models in `backend/models/` that reference non-existent types from `todo-app/internal/models`. This research documents the current state, identifies all affected files, and proposes a consolidation strategy.

## 1. Import Path Issues Scan

### Build Error Analysis
Running `go build ./...` revealed **14 compilation errors** across 3 packages:

#### Package: `todo-app/jobs`
**Files affected**: 2
- `jobs/oauth_cleanup.go` (4 errors): References to `models.OAuthState`
- `jobs/session_cleanup.go` (10 errors): References to `models.AuthenticationSession`

**Root cause**: These files import `"todo-app/internal/models"` but reference types that don't exist in that package. The actual definitions are in `backend/models/` (session.go, oauth_state.go).

#### Package: `todo-app/services/auth`
**Files affected**: 2
- `services/auth/oauth.go` (11 errors):
  - `models.CreateAndSave` (undefined)
  - `models.AuthenticationSession` (undefined)
- `services/auth/session.go` (10 errors):
  - `models.AuthenticationSession` (undefined)
  - `models.SessionValidationResult` (undefined)

**Root cause**: Same as above - imports `todo-app/internal/models` but references types defined in `backend/models/`.

#### Package: `todo-app/internal/config`
**Files affected**: 1
- `internal/config/database.go` (2 errors):
  - `legacymodels.AuthenticationSession` (line 73)
  - `legacymodels.OAuthState` (line 74)

**Root cause**: Uses alias `legacymodels "todo-app/internal/models"` but the types don't exist in that package. This is the **explicit "legacy models" issue** mentioned in the feature spec.

### Summary of Broken Imports
- **Total files with compilation errors**: 5
- **Total error count**: 14+ (build stopped at "too many errors")
- **Undefined types**:
  - `models.OAuthState` (referenced but not in `internal/models/`)
  - `models.AuthenticationSession` (referenced but not in `internal/models/`)
  - `models.CreateAndSave` (function, referenced but not in `internal/models/`)
  - `models.SessionValidationResult` (referenced but not in `internal/models/`)

## 2. Current Model Structure Analysis

### Model Locations Identified

#### A. `backend/internal/models/` (Current "official" location - 4 files)
- `health.go` - Health check models
- `task.go` - Legacy flat Task model
- `user.go` - Legacy flat User model
- `google_identity.go` - GoogleIdentity for OAuth

**Usage**: 34 files import from this location (most common)

#### B. `backend/models/` (Root-level legacy - 2 files)
**THESE ARE THE "ORPHANED" MODELS**:
- `session.go` - Defines `AuthenticationSession` type
- `oauth_state.go` - Defines `OAuthState` type

**Problem**: These files themselves import `internalmodels "todo-app/internal/models"` but export types under the package name `models`, creating a confusing namespace collision.

**Cross-reference**: The types defined here are what `services/auth/*.go` and `jobs/*.go` are trying to reference via `"todo-app/internal/models"` imports, but they're in the wrong location.

#### C. `backend/domain/` (DDD Structure - Target)
This is the **proper DDD structure** that exists and should be the target:

```
domain/
├── user/
│   ├── entities/user.go (DDD User entity)
│   ├── valueobjects/ (email.go, user_id.go, user_profile.go, user_preferences.go)
│   ├── repositories/user_repository.go
│   └── services/ (user_authentication_service.go, user_profile_service.go)
└── task/
    ├── entities/task.go (DDD Task entity)
    ├── valueobjects/ (task_id.go, task_title.go, task_description.go, task_status.go, task_priority.go)
    ├── repositories/task_repository.go
    └── services/ (task_validation_service.go, task_search_service.go)
```

**Status**: This structure is well-established with comprehensive tests. It represents the "new" architecture.

### Model Duplication Analysis

| Model Type | internal/models/ | models/ | domain/ | Notes |
|------------|-----------------|---------|---------|-------|
| User | ✓ (legacy flat) | ✗ | ✓ (DDD entity) | Duplication - flat vs rich domain |
| Task | ✓ (legacy flat) | ✗ | ✓ (DDD entity) | Duplication - flat vs rich domain |
| Health | ✓ | ✗ | ✗ | No duplication |
| GoogleIdentity | ✓ | ✗ | ✗ | OAuth specific |
| AuthenticationSession | ✗ | ✓ (orphaned) | ✗ | **Missing from internal/models/** |
| OAuthState | ✗ | ✓ (orphaned) | ✗ | **Missing from internal/models/** |

**Key Finding**: The "legacy models" to deprecate per clarification are:
1. `backend/models/session.go` (AuthenticationSession)
2. `backend/models/oauth_state.go` (OAuthState)

These should be **moved to domain/ structure** under a new `domain/auth/` or `domain/session/` context, NOT to `internal/models/`.

## 3. DDD Package Organization Research

### Go DDD Best Practices (from community research)

#### Recommended Structure
```
domain/
  ├── [bounded-context]/      # e.g., user, task, auth, order
  │   ├── entities/            # Core domain objects with identity
  │   ├── valueobjects/        # Immutable values without identity
  │   ├── repositories/        # Persistence interfaces (not implementations)
  │   ├── services/            # Domain logic coordinating entities
  │   └── events/              # Domain events (optional)
application/                   # Application services orchestrating domain
infrastructure/                # Technical implementations (persistence, external APIs)
presentation/                  # HTTP handlers, CLI, gRPC, etc.
```

#### Entity vs Value Object Patterns (Applied to our case)

**Entities** (have identity, mutable):
- User (has UserID)
- Task (has TaskID)
- AuthenticationSession (has SessionID) ← **Should be entity**
- OAuthState (has StateToken) ← **Should be entity**

**Value Objects** (no identity, immutable):
- Email, UserProfile, UserPreferences
- TaskTitle, TaskDescription, TaskStatus, TaskPriority
- SessionToken, PKCEVerifier ← **Could extract from AuthenticationSession**

### Applying DDD to Auth/Session Models

**Decision**: Create `domain/auth/` bounded context:
```
domain/auth/
├── entities/
│   ├── authentication_session.go  # From backend/models/session.go
│   └── oauth_state.go             # From backend/models/oauth_state.go
├── valueobjects/
│   ├── session_token.go           # Extract from session
│   ├── pkce_verifier.go           # Extract from oauth_state
│   └── state_token.go             # Extract from oauth_state
├── repositories/
│   ├── session_repository.go      # Interface
│   └── oauth_state_repository.go  # Interface
└── services/
    ├── session_management_service.go
    └── oauth_flow_service.go
```

**Rationale**:
1. Auth/Session is a distinct bounded context separate from User domain
2. Sessions have identity (SessionID) → entities
3. OAuth states are transient but have identity (StateToken) → entities
4. Follows existing task/ and user/ structure conventions
5. Enables proper dependency injection via repository interfaces

### Migration Strategy Best Practices

From Go DDD projects and refactoring guides:

1. **Strangler Fig Pattern**:
   - Create new domain structure alongside old
   - Gradually migrate consumers
   - Remove old structure when usage reaches zero

2. **Alias Strategy** (for backward compatibility):
   ```go
   // Temporary alias during migration
   type LegacyUser = domain.user.entities.User
   ```

3. **Repository Pattern** (already partially implemented):
   - `infrastructure/persistence/` contains GORM implementations
   - Domain repositories are interfaces
   - ✓ This pattern is already followed for user/task

## 4. Backward Compatibility Requirements

### Public API Analysis

Reviewed API endpoints in `backend/cmd/server/main.go` and handler files:

#### Endpoints potentially affected:
- `/api/auth/google/login` - Returns auth URL
- `/api/auth/google/callback` - Processes OAuth callback
- `/api/auth/session/validate` - Validates session
- `/api/auth/session/refresh` - Refreshes session
- `/api/auth/logout` - Destroys session

#### JSON Response Structures:
All responses use Go struct JSON marshaling. Internal model changes won't affect JSON as long as struct tags remain consistent.

**Example** (from research of existing code):
```go
type AuthenticationSession struct {
    SessionID   string    `json:"session_id"`
    UserID      string    `json:"user_id"`
    ExpiresAt   time.Time `json:"expires_at"`
    // ... other fields
}
```

As long as the new `domain/auth/entities/AuthenticationSession` maintains these field names and tags, JSON API compatibility is preserved.

### External Consumers:
- **Frontend** (`frontend/src/`): Consumes auth APIs
- **No external API consumers** identified (internal application)

### Compatibility Strategy:
**Decision**: No aliases needed - direct refactoring is safe because:
1. JSON contracts are stable (based on struct tags, not package names)
2. No external package consumers (all imports are internal)
3. Type changes only affect internal implementation
4. Tests will validate behavior preservation

## 5. Import Path Resolution Strategy

### Root Cause
The compilation errors stem from incorrect assumptions about where models live:

**Current (Broken) State**:
```
services/auth/*.go  →  imports "todo-app/internal/models"
                   →  tries to use models.AuthenticationSession
                   →  ERROR: type not in internal/models/

backend/models/session.go  →  package models
                            →  defines AuthenticationSession
                            →  (orphaned - no one can find it!)
```

**Target (DDD) State**:
```
services/auth/*.go  →  imports "todo-app/domain/auth/entities"
                   →  uses entities.AuthenticationSession
                   →  SUCCESS: type exists in domain structure

domain/auth/entities/authentication_session.go  →  package entities
                                                 →  defines AuthenticationSession
```

### Consolidation Decisions

| Component | Current Location | Target Location | Action |
|-----------|-----------------|-----------------|--------|
| AuthenticationSession | backend/models/ | domain/auth/entities/ | **Move + refactor** |
| OAuthState | backend/models/ | domain/auth/entities/ | **Move + refactor** |
| SessionService | services/auth/session.go | domain/auth/services/ | **Move** |
| OAuthService | services/auth/oauth.go | domain/auth/services/ | **Move** |
| User (flat) | internal/models/ | (delete) | **Deprecate** - use domain/user/entities/User |
| Task (flat) | internal/models/ | (delete) | **Deprecate** - use domain/task/entities/Task |
| Health | internal/models/ | domain/health/entities/ | **Move** (new bounded context) |
| GoogleIdentity | internal/models/ | domain/auth/valueobjects/ | **Move** (auth-related) |

### Import Path Mapping

**Pattern**: All imports should follow DDD structure:
```
OLD: "todo-app/internal/models"
NEW (context-specific):
  - "todo-app/domain/user/entities"
  - "todo-app/domain/task/entities"
  - "todo-app/domain/auth/entities"
  - "todo-app/domain/auth/valueobjects"
  - "todo-app/domain/health/entities"
```

**Deprecated directories to remove**:
- `backend/models/` ← Delete after migration
- `backend/internal/models/` ← Delete after migration (maybe keep config-only models)
- `backend/services/auth/` ← Delete after moving to domain/auth/services/

## 6. Comprehensive File Impact Analysis

### Files Requiring Import Updates (34 files)

Based on grep results showing `todo-app/internal/models` imports:

1. `backend/services/user/user.go`
2. `backend/services/auth/oauth.go` ← **High priority** (compilation error)
3. `backend/services/auth/session.go` ← **High priority** (compilation error)
4. `backend/jobs/session_cleanup.go` ← **High priority** (compilation error)
5. `backend/jobs/oauth_cleanup.go` ← **High priority** (compilation error)
6. `backend/internal/config/database.go` ← **High priority** (legacymodels alias)
7. `backend/cmd/server/main.go`
8. `backend/internal/handlers/google_oauth_handler.go`
9. `backend/internal/handlers/task_handlers.go`
10. `backend/internal/services/health_service.go`
11. `backend/internal/services/task_service.go`
12. `backend/internal/services/google_oauth_service.go`
13. `backend/internal/storage/database.go`
14. All test files (21 files in tests/unit/, tests/integration/, tests/contract/)

### Additional Files Requiring Updates (consumers of moved types)

Files that import the types being moved (deeper analysis needed):
- Any file using `AuthenticationSession` type (estimated 15-20 files)
- Any file using `OAuthState` type (estimated 10 files)
- Handler files in `backend/handlers/`
- Middleware files in `backend/middleware/`

**Estimated total impact**: ~50-60 files

## 7. Risk Assessment

### High Risk Areas
1. **Session management logic**: Any bugs break authentication entirely
2. **OAuth flow**: Broken OAuth = users can't log in
3. **Database migrations**: AuthenticationSession and OAuthState are persisted via GORM

### Mitigation Strategies
1. **Tests first**: Verify all 51 existing tests pass before starting
2. **Incremental migration**: Move one bounded context at a time
3. **Compilation gate**: Each step must compile before proceeding
4. **Test-driven**: Run tests after each file change
5. **Database compatibility**: Use GORM table names to maintain DB schema compatibility

## 8. Decisions Summary

### Key Decisions

| Decision Area | Choice | Rationale |
|--------------|--------|-----------|
| **Package structure** | Domain-Driven Design (DDD) | User clarification + existing codebase momentum |
| **Auth models location** | `domain/auth/` | Separate bounded context from user/task |
| **Legacy models strategy** | Deprecate and remove | User clarification (option C) |
| **Backward compatibility** | No aliases needed | Internal-only codebase, JSON contracts stable |
| **Migration approach** | Direct refactoring | Tests validate behavior, no external consumers |
| **internal/models/ fate** | Delete after migration | Consolidate everything into domain/ |

### Alternatives Considered

1. **Keep internal/models/ and move orphans there**
   - Rejected: Doesn't align with DDD strategy (clarification answer)
   - Would perpetuate flat model structure

2. **Use type aliases for compatibility**
   - Rejected: Unnecessary complexity for internal codebase
   - Would delay full cleanup

3. **Gradual strangler fig migration**
   - Rejected: Small codebase (~100 files), direct refactoring is tractable
   - Tests provide safety net

## 9. Success Criteria

### Compilation
- [ ] `go build ./...` succeeds with zero errors
- [ ] `go vet ./...` passes with zero warnings
- [ ] `go mod tidy` shows no unused dependencies

### Testing
- [ ] All 51 existing tests pass
- [ ] Contract tests validate API JSON compatibility
- [ ] Integration tests validate auth flows work end-to-end

### Code Quality
- [ ] No imports to `todo-app/internal/models`
- [ ] No imports to `todo-app/models`
- [ ] All imports follow `todo-app/domain/[context]/[layer]` pattern
- [ ] `backend/models/` directory removed
- [ ] `backend/internal/models/` directory removed (or contains only non-domain helpers)

### Documentation
- [ ] quickstart.md documents new import patterns
- [ ] CLAUDE.md updated with DDD structure
- [ ] Migration notes added to README (if needed)

## Appendices

### A. Full List of Model Files

**internal/models/** (4 files):
1. health.go
2. task.go (deprecated, prefer domain/task/entities/task.go)
3. user.go (deprecated, prefer domain/user/entities/user.go)
4. google_identity.go

**models/** (2 files):
1. session.go (orphaned, needs migration)
2. oauth_state.go (orphaned, needs migration)

**domain/** (10 entity files):
- domain/user/entities/user.go
- domain/task/entities/task.go
- 8 value object files

### B. Import Pattern Reference

```go
// OLD (broken/deprecated)
import "todo-app/internal/models"
import legacymodels "todo-app/internal/models"

// NEW (DDD structure)
import "todo-app/domain/auth/entities"
import "todo-app/domain/auth/valueobjects"
import "todo-app/domain/auth/repositories"
import "todo-app/domain/auth/services"
import "todo-app/domain/user/entities"
import "todo-app/domain/task/entities"
import "todo-app/domain/health/entities"
```

### C. Recommended Task Order (Preview for Phase 2)

1. Create domain/auth/ structure
2. Move models/session.go → domain/auth/entities/authentication_session.go
3. Move models/oauth_state.go → domain/auth/entities/oauth_state.go
4. Update internal/config/database.go imports (remove legacymodels)
5. Update services/auth/*.go imports
6. Update jobs/*.go imports
7. Update all test imports
8. Update remaining 34 files with internal/models imports
9. Remove backend/models/ directory
10. Remove backend/internal/models/ or deprecate remaining files
11. Run full test suite
12. Verify build success

---

**Research complete**. All NEEDS CLARIFICATION items resolved. Ready for Phase 1 design.
