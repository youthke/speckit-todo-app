# Feature Specification: Complete DDD Migration

**Feature ID**: 010-complete-ddd-migration
**Status**: Draft
**Priority**: Medium
**Created**: 2025-10-04
**Dependencies**: 009-resolve-it-1 (partial completion required)

## Overview

Complete the Domain-Driven Design (DDD) migration started in feature 009 by resolving the architectural incompatibility between legacy GORM DTOs (`internal/models`) and rich DDD entities (`domain/`), enabling full migration of User and Task models.

## Context

Feature 009-resolve-it-1 successfully migrated Auth and Health domains to DDD structure, achieving 75% completion and fixing all compilation errors. However, 15 tasks remain blocked due to an architectural incompatibility:

**Current State**:
- ✅ Auth domain: Fully migrated to DDD (`domain/auth/`)
- ✅ Health domain: Fully migrated to DDD (`domain/health/`)
- ⚠️ User domain: Split between `internal/models/user.go` (DTO) and `domain/user/entities/user.go` (DDD)
- ⚠️ Task domain: Split between `internal/models/task.go` (DTO) and `domain/task/entities/task.go` (DDD)

**Problem**:
- `internal/models`: Simple GORM DTOs with public fields, direct database mapping
- `domain/entities`: Rich DDD entities with private fields, value objects, business logic
- 30 files use `internal/models` imports (11 source files + 19 test files)
- Direct replacement breaks functionality due to incompatible interfaces

**Impact**:
- Test suite fails (19 test files need updates)
- Services mixing DTO and DDD patterns
- Inconsistent architecture across codebase

## Problem Statement

**As a** backend developer
**I want** a unified DDD architecture across all domains
**So that** the codebase follows consistent patterns, tests pass, and future features align with domain-driven design principles

### Current Pain Points

1. **Architectural Inconsistency**: Some domains use DTOs, others use DDD entities
2. **Test Failures**: 19 test files fail due to import mismatches
3. **Maintenance Burden**: Developers must understand two different model patterns
4. **Technical Debt**: Legacy `internal/models` directory should be deprecated
5. **Migration Incompleteness**: Feature 009 left 15 tasks incomplete

## Goals

### Primary Goals

1. ✅ **Unify Architecture**: All domains use DDD structure consistently
2. ✅ **Fix Test Suite**: All tests pass after migration
3. ✅ **Remove Legacy Code**: Deprecate `internal/models/user.go` and `task.go`
4. ✅ **Maintain Compatibility**: No breaking changes to APIs or database schema

### Secondary Goals

- Document mapper/adapter pattern for future migrations
- Establish clear guidelines for DTO ↔ DDD entity conversion
- Enable smooth transition path for similar refactoring efforts

### Non-Goals

- Changing database schema or table structures
- Modifying public API contracts (JSON responses)
- Refactoring application logic beyond model layer
- Adding new features or functionality

## Proposed Solution

### High-Level Approach

Implement a **Mapper/Adapter Pattern** to bridge DTOs and DDD entities:

```
┌─────────────────────────────────────────────────────────┐
│  Presentation Layer (Handlers)                          │
│  - Receives HTTP requests                               │
│  - Uses DTOs for JSON serialization                     │
└────────────────┬────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────┐
│  Mapper Layer (NEW)                                     │
│  - Converts DTO ↔ DDD Entity                            │
│  - Handles validation and transformation                │
│  - Location: application/mappers/                       │
└────────────────┬────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────┐
│  Domain Layer (DDD Entities)                            │
│  - Business logic and invariants                        │
│  - Rich domain models with value objects                │
│  - Location: domain/user/, domain/task/                 │
└────────────────┬────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────┐
│  Infrastructure Layer (GORM Repositories)               │
│  - Converts DDD Entity ↔ Database Record                │
│  - Persistence logic                                    │
└─────────────────────────────────────────────────────────┘
```

### Solution Components

#### 1. DTO Layer (Keep for API Compatibility)

**Purpose**: JSON serialization/deserialization for HTTP APIs

**Location**: `internal/dtos/` (renamed from `internal/models/`)

```go
// internal/dtos/user_dto.go
type UserDTO struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    Email     string    `json:"email" gorm:"index"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (UserDTO) TableName() string {
    return "users" // Preserve DB table name
}
```

**Characteristics**:
- Public fields for GORM and JSON marshaling
- Struct tags for JSON/GORM
- No business logic
- Direct database mapping

#### 2. Mapper Layer (NEW)

**Purpose**: Transform between DTOs and DDD entities

**Location**: `application/mappers/`

```go
// application/mappers/user_mapper.go
type UserMapper struct{}

func (m *UserMapper) ToEntity(dto *dtos.UserDTO) (*entities.User, error) {
    email, err := valueobjects.NewEmail(dto.Email)
    if err != nil {
        return nil, err
    }

    userID := valueobjects.NewUserID(dto.ID)

    return entities.NewUser(userID, email, dto.Name), nil
}

func (m *UserMapper) ToDTO(user *entities.User) *dtos.UserDTO {
    return &dtos.UserDTO{
        ID:        user.ID().Value(),
        Email:     user.Email().Value(),
        Name:      user.Name(),
        CreatedAt: user.CreatedAt(),
        UpdatedAt: user.UpdatedAt(),
    }
}
```

**Characteristics**:
- Bidirectional conversion
- Validation during DTO → Entity conversion
- Error handling for invalid data
- Stateless transformation logic

#### 3. Domain Layer (Enhanced DDD)

**Purpose**: Business logic and invariants

**Location**: `domain/user/entities/`, `domain/task/entities/`

```go
// domain/user/entities/user.go
type User struct {
    id        valueobjects.UserID
    email     valueobjects.Email
    name      string
    createdAt time.Time
    updatedAt time.Time
}

func NewUser(id valueobjects.UserID, email valueobjects.Email, name string) *User {
    return &User{
        id:        id,
        email:     email,
        name:      name,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }
}

// Getters with encapsulation
func (u *User) ID() valueobjects.UserID { return u.id }
func (u *User) Email() valueobjects.Email { return u.email }
func (u *User) Name() string { return u.name }

// Business methods
func (u *User) UpdateEmail(newEmail valueobjects.Email) error {
    if !newEmail.IsValid() {
        return errors.New("invalid email")
    }
    u.email = newEmail
    u.updatedAt = time.Now()
    return nil
}
```

**Characteristics**:
- Private fields (encapsulation)
- Value objects for type safety
- Business logic methods
- Immutability where appropriate

#### 4. Repository Pattern Update

**Purpose**: Persist DDD entities using GORM

**Location**: `infrastructure/persistence/`

```go
// infrastructure/persistence/gorm_user_repository.go
type GormUserRepository struct {
    db     *gorm.DB
    mapper *mappers.UserMapper
}

func (r *GormUserRepository) Save(user *entities.User) error {
    dto := r.mapper.ToDTO(user)
    return r.db.Save(dto).Error
}

func (r *GormUserRepository) FindByID(id valueobjects.UserID) (*entities.User, error) {
    var dto dtos.UserDTO
    if err := r.db.First(&dto, "id = ?", id.Value()).Error; err != nil {
        return nil, err
    }
    return r.mapper.ToEntity(&dto)
}
```

**Characteristics**:
- Uses mapper for DTO ↔ Entity conversion
- Maintains GORM compatibility
- Returns domain entities
- Database-agnostic interface

### Migration Strategy

#### Phase 1: Setup Mapper Infrastructure

1. Create `application/mappers/` directory
2. Rename `internal/models/` → `internal/dtos/` (preserve DTOs for GORM)
3. Implement `UserMapper` and `TaskMapper`
4. Add mapper tests

#### Phase 2: Update Repositories

1. Inject mappers into GORM repositories
2. Update `FindByID`, `Save`, `Update`, `Delete` methods
3. Ensure repositories return DDD entities
4. Test repository layer

#### Phase 3: Update Services

1. Update services to use DDD entities internally
2. Convert at service boundaries (DTO → Entity → DTO)
3. Update business logic to use entity methods
4. Test service layer

#### Phase 4: Update Handlers

1. Handlers receive/send DTOs (JSON compatibility)
2. Call mappers to convert DTO → Entity
3. Pass entities to services
4. Convert Entity → DTO for responses
5. Test handler layer

#### Phase 5: Update Tests

1. Update 19 test files to use new imports
2. Update mock objects to use DDD entities
3. Add mapper tests
4. Ensure all tests pass

#### Phase 6: Cleanup

1. Verify no references to old `internal/models` paths
2. Update documentation
3. Remove deprecated code paths
4. Final verification

## Technical Specifications

### Technology Stack

- **Language**: Go 1.24.7
- **Framework**: Gin web framework
- **ORM**: GORM (existing)
- **Testing**: Go testing framework, testify
- **Database**: SQLite (development)

### File Structure

```
backend/
├── domain/
│   ├── user/
│   │   ├── entities/user.go          # DDD entity (exists, enhance)
│   │   ├── valueobjects/              # Email, UserID (exists)
│   │   └── repositories/user_repository.go  # Interface (exists)
│   └── task/
│       ├── entities/task.go          # DDD entity (exists, enhance)
│       ├── valueobjects/              # TaskID, TaskStatus (exists)
│       └── repositories/task_repository.go  # Interface (exists)
├── application/
│   └── mappers/                      # NEW
│       ├── user_mapper.go            # DTO ↔ Entity conversion
│       ├── user_mapper_test.go
│       ├── task_mapper.go
│       └── task_mapper_test.go
├── internal/
│   └── dtos/                         # RENAMED from models/
│       ├── user_dto.go               # Simple GORM DTO
│       └── task_dto.go               # Simple GORM DTO
├── infrastructure/
│   └── persistence/
│       ├── gorm_user_repository.go   # UPDATE: Add mapper
│       └── gorm_task_repository.go   # UPDATE: Add mapper
├── internal/
│   ├── handlers/                     # UPDATE: Use mappers
│   └── services/                     # UPDATE: Use DDD entities
└── tests/
    ├── unit/                         # UPDATE: 8 test files
    ├── integration/                  # UPDATE: 13 test files
    └── contract/                     # UPDATE: 23 test files (subset)
```

### API Compatibility

**Critical Constraint**: Zero breaking changes to JSON APIs

**Example - User Endpoint**:

```go
// Handler (external interface - unchanged)
func (h *UserHandler) GetUser(c *gin.Context) {
    // 1. Parse request (DTO for JSON)
    userID := c.Param("id")

    // 2. Convert to domain
    id := valueobjects.NewUserID(userID)

    // 3. Call service (uses DDD entity)
    user, err := h.userService.GetByID(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }

    // 4. Convert to DTO for response
    dto := h.userMapper.ToDTO(user)

    // 5. Return JSON (same structure as before)
    c.JSON(200, dto)
}
```

**JSON Response (unchanged)**:
```json
{
  "id": "user-123",
  "email": "user@example.com",
  "name": "John Doe",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-02T00:00:00Z"
}
```

### Database Compatibility

**Critical Constraint**: Zero schema changes

**Strategy**:
- DTOs preserve `TableName()` methods
- GORM tags remain on DTOs
- Repositories use DTOs for persistence
- Entities are reconstructed from DTOs

**Example**:
```go
// DTO has GORM tags and TableName
type UserDTO struct {
    ID    string `gorm:"primaryKey"`
    Email string `gorm:"index"`
}

func (UserDTO) TableName() string {
    return "users" // Existing table
}

// Repository uses DTO for database operations
repo.db.First(&dto, "id = ?", id)
entity := mapper.ToEntity(&dto)
```

## Acceptance Criteria

### Must Have

1. ✅ **All 47 tasks from feature 009 completed**
   - T021-T034: Import updates completed
   - T038: Services migrated to domain structure
   - T042: All tests pass
   - T046-T047: Manual verification succeeds

2. ✅ **Test Suite Passes**
   - All 51+ existing tests pass
   - Zero test failures related to imports
   - Mapper unit tests added and passing

3. ✅ **Architecture Unified**
   - Zero files import `internal/models` for User/Task
   - All domains use consistent DDD structure
   - Mapper pattern documented

4. ✅ **Zero Breaking Changes**
   - All API endpoints return same JSON structure
   - Database schema unchanged
   - No functional regressions

5. ✅ **Code Quality**
   - `go build ./...` succeeds
   - `go vet ./...` passes
   - `go test ./...` passes

### Should Have

- Mapper pattern documented in CLAUDE.md
- Migration guide for future similar refactoring
- Performance benchmarks (mapper overhead < 1ms)

### Could Have

- Automated migration script for other DTOs
- Code generation for mappers
- Integration with CI/CD pipeline

## Acceptance Scenarios

### Scenario 1: User Retrieval Flow

**Given**: A User entity exists in the database
**When**: Client calls `GET /api/users/{id}`
**Then**:
- Handler receives request
- Repository fetches DTO from database
- Mapper converts DTO → Entity
- Service processes Entity
- Mapper converts Entity → DTO
- Handler returns JSON (same structure as before)
- Response time < 50ms

### Scenario 2: Task Creation Flow

**Given**: Client sends valid task data
**When**: Client calls `POST /api/tasks`
**Then**:
- Handler parses JSON to DTO
- Mapper converts DTO → Entity with validation
- Service applies business rules using Entity
- Repository persists Entity (via DTO)
- Handler returns created task JSON
- Database contains new record

### Scenario 3: Test Execution

**Given**: All code changes implemented
**When**: Developer runs `go test ./...`
**Then**:
- All unit tests pass (domain entities, mappers)
- All integration tests pass (repositories, services)
- All contract tests pass (API JSON contracts)
- Zero import-related failures

### Scenario 4: Legacy Code Cleanup

**Given**: All migrations completed
**When**: Developer searches for legacy imports
**Then**:
- `grep -r "internal/models" backend/` → 0 matches for User/Task
- `internal/models/` renamed to `internal/dtos/`
- Only DTOs remain in `internal/dtos/`
- Documentation updated

## Risks and Mitigations

### High Risks

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Breaking API contracts | Critical | Extensive contract tests, JSON structure validation |
| Database migration errors | Critical | DTOs preserve GORM tags, no schema changes |
| Performance degradation | Medium | Benchmark mapper operations, optimize if needed |
| Test coverage gaps | Medium | Require tests for all mappers, 100% coverage |

### Medium Risks

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Complex mapper logic | Medium | Keep mappers simple, stateless, well-tested |
| Incomplete migration | Medium | Checklist verification, automated import scanning |
| Developer confusion | Low | Clear documentation, examples, migration guide |

## Dependencies

### Prerequisites

- ✅ Feature 009-resolve-it-1 completed at 75% (T001-T020, T025, T029, T035-T045)
- ✅ Backend compiles successfully
- ✅ DDD structure exists for 4 domains
- ✅ Auth and Health domains fully migrated

### External Dependencies

- None (internal refactoring only)

### Blocking Dependencies

- None

## Metrics and KPIs

### Success Metrics

- **Task Completion**: 47/47 tasks (100%)
- **Test Pass Rate**: 100% (all tests passing)
- **Import Compliance**: 0 forbidden imports
- **Code Coverage**: Maintain or improve current coverage
- **Build Time**: No regression (< 10 seconds)

### Quality Metrics

- **Mapper Performance**: < 1ms per conversion
- **Code Duplication**: < 5% (between DTOs and entities)
- **Cyclomatic Complexity**: Mappers < 10
- **Documentation Coverage**: 100% for mapper pattern

## Timeline Estimate

**Total Effort**: 6-8 hours

| Phase | Effort | Duration |
|-------|--------|----------|
| Phase 1: Mapper Setup | 1 hour | Tasks T1-T5 |
| Phase 2: Repository Updates | 1 hour | Tasks T6-T10 |
| Phase 3: Service Updates | 1.5 hours | Tasks T11-T20 |
| Phase 4: Handler Updates | 1.5 hours | Tasks T21-T30 |
| Phase 5: Test Updates | 2 hours | Tasks T31-T45 |
| Phase 6: Cleanup & Verification | 1 hour | Tasks T46-T50 |

## Open Questions

1. **Mapper Location**: Confirm `application/mappers/` is appropriate (vs `infrastructure/mappers/`)
2. **DTO Naming**: Keep `internal/dtos/` or use `presentation/dtos/`?
3. **Validation Placement**: Should mappers validate, or delegate to entities?
4. **Error Handling**: Return errors from ToEntity() or panic on invalid data?

## Alternatives Considered

### Alternative 1: Replace DTOs Entirely

**Approach**: Use DDD entities everywhere, including GORM

**Pros**:
- Single model type
- True DDD architecture

**Cons**:
- Requires GORM tag pollution on entities
- Breaks encapsulation (public fields required)
- Violates DDD principles

**Decision**: ❌ Rejected - Compromises DDD integrity

### Alternative 2: Keep DTOs, Abandon DDD

**Approach**: Revert to simple DTOs everywhere

**Pros**:
- Simpler architecture
- No mapper layer needed

**Cons**:
- Loses feature 009 progress
- No domain-driven design benefits
- Doesn't align with project goals

**Decision**: ❌ Rejected - Contradicts feature 009 objectives

### Alternative 3: Mapper Pattern (SELECTED)

**Approach**: DTOs for persistence/API, entities for domain logic

**Pros**:
- Clean separation of concerns
- Maintains API compatibility
- Preserves DDD architecture
- Gradual migration path

**Cons**:
- Additional layer complexity
- Mapper maintenance overhead

**Decision**: ✅ Selected - Best balance of DDD and pragmatism

## Notes

- This spec builds on feature 009-resolve-it-1 findings
- Architectural decision made based on real implementation blockers
- Mapper pattern is a proven solution for DDD in Go projects
- Implementation can proceed incrementally (User domain first, then Task)

## References

- Feature 009-resolve-it-1: `/specs/009-resolve-it-1/`
- Research findings: `/specs/009-resolve-it-1/research.md`
- DDD in Go patterns: Standard Go DDD project structures
- Mapper pattern examples: Application layer transformation patterns

---

**Prepared by**: Claude Code
**Last Updated**: 2025-10-04
**Status**: Ready for `/plan` command
