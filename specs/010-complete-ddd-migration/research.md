# Research: Complete DDD Migration

**Feature**: 010-complete-ddd-migration
**Date**: 2025-10-04
**Status**: Complete

## Overview

This document consolidates research findings for implementing the Mapper/Adapter pattern to complete the DDD migration in the todo-app backend. All technical unknowns have been resolved through analysis of the existing codebase and DDD best practices.

---

## Research Areas

### 1. Mapper Pattern in Go DDD Applications

**Decision**: Implement stateless mapper structs in `application/mappers/` package

**Rationale**:
- **Separation of Concerns**: DTOs handle serialization/persistence, entities handle business logic
- **Testability**: Mappers can be unit tested independently
- **Maintainability**: Centralized conversion logic, no scattered mapping code
- **Performance**: Stateless mappers have minimal overhead (< 1ms per conversion)
- **Golang Idioms**: Struct-based mappers align with Go's composition patterns

**Alternatives Considered**:
1. **Embedded conversion methods on DTOs**
   - ❌ Violates single responsibility (DTOs should only handle data structure)
   - ❌ Creates circular dependencies between DTO and entity packages

2. **Embedded conversion methods on Entities**
   - ❌ Pollutes domain layer with infrastructure concerns
   - ❌ Entities shouldn't know about DTOs (dependency direction violation)

3. **Selected: Dedicated mapper layer**
   - ✅ Clean separation between layers
   - ✅ Unidirectional dependencies (mapper knows both DTO and entity)
   - ✅ Easy to mock for testing
   - ✅ Follows hexagonal architecture principles

**Implementation Pattern**:
```go
type UserMapper struct{}

func (m *UserMapper) ToEntity(dto *dtos.UserDTO) (*entities.User, error) {
    // DTO → Entity with validation
}

func (m *UserMapper) ToDTO(entity *entities.User) *dtos.UserDTO {
    // Entity → DTO (no validation needed)
}
```

**References**:
- Domain-Driven Design by Eric Evans (Chapter 5: Model-Driven Design)
- Clean Architecture patterns for Go applications
- Existing pattern in feature 009 for Auth/Health domains

---

### 2. DTO vs Entity Separation Strategy

**Decision**: Keep DTOs for GORM persistence, use entities for business logic

**Rationale**:
- **GORM Compatibility**: GORM requires public fields and struct tags, which violates encapsulation
- **API Stability**: DTOs provide stable JSON contracts independent of domain model changes
- **Database Decoupling**: Schema changes don't ripple into domain layer
- **DDD Principles**: Entities remain pure domain objects without infrastructure pollution

**Alternatives Considered**:
1. **Use entities directly with GORM**
   - ❌ Requires public fields (breaks encapsulation)
   - ❌ GORM tags pollute domain models
   - ❌ Database schema couples to business logic

2. **Use DTOs everywhere (abandon DDD)**
   - ❌ Loses feature 009 progress (75% complete)
   - ❌ No place for business invariants
   - ❌ Anemic domain model anti-pattern

3. **Selected: Dual model approach (DTO + Entity)**
   - ✅ GORM works with DTOs (public fields, tags)
   - ✅ Entities encapsulate business logic
   - ✅ Clear layer boundaries
   - ✅ Gradual migration path

**Layer Responsibilities**:
- **DTOs** (`internal/dtos/`): Database schema, JSON serialization, GORM tags
- **Entities** (`domain/*/entities/`): Business logic, invariants, value objects
- **Mappers** (`application/mappers/`): Bidirectional transformation

---

### 3. Error Handling in Mappers

**Decision**: Return errors from `ToEntity()`, no errors from `ToDTO()`

**Rationale**:
- **DTO → Entity**: External data may be invalid, validation required
- **Entity → DTO**: Entities are always valid (invariants enforced), conversion cannot fail
- **Fail-Fast**: Detect invalid data at system boundary (handlers)
- **Type Safety**: Errors propagate to caller for appropriate HTTP status codes

**Pattern**:
```go
func (m *UserMapper) ToEntity(dto *dtos.UserDTO) (*entities.User, error) {
    email, err := valueobjects.NewEmail(dto.Email)
    if err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }
    return entities.NewUser(...), nil
}

func (m *UserMapper) ToDTO(entity *entities.User) *dtos.UserDTO {
    // No error possible - entity is always valid
    return &dtos.UserDTO{...}
}
```

**Validation Placement**:
- **Mapper Layer**: Format validation (email syntax, ID format)
- **Entity Layer**: Business validation (email uniqueness, task assignment rules)
- **Handler Layer**: HTTP input validation (required fields, JSON structure)

---

### 4. Repository Pattern with Mappers

**Decision**: Inject mappers into repositories, return entities from repository methods

**Rationale**:
- **Domain Interface**: Repositories in `domain/*/repositories/` define entity-based interfaces
- **Infrastructure Implementation**: GORM repositories in `infrastructure/persistence/` use DTOs internally
- **Dependency Injection**: Mappers injected via constructor for testability
- **Consistency**: All repository methods return domain types

**Implementation Pattern**:
```go
type GormUserRepository struct {
    db     *gorm.DB
    mapper *mappers.UserMapper
}

func NewGormUserRepository(db *gorm.DB, mapper *mappers.UserMapper) *GormUserRepository {
    return &GormUserRepository{db: db, mapper: mapper}
}

func (r *GormUserRepository) FindByID(id valueobjects.UserID) (*entities.User, error) {
    var dto dtos.UserDTO
    if err := r.db.First(&dto, "id = ?", id.Value()).Error; err != nil {
        return nil, err
    }
    return r.mapper.ToEntity(&dto)
}
```

**Benefits**:
- Repositories remain infrastructure-agnostic (interface in domain)
- Easy to swap GORM for another ORM
- Testable via mock mappers

---

### 5. Testing Strategy

**Decision**: 3-layer test coverage (mapper tests, repository tests, integration tests)

**Test Types**:

1. **Mapper Unit Tests** (`application/mappers/*_test.go`):
   - Test DTO → Entity conversion (valid/invalid cases)
   - Test Entity → DTO conversion
   - Test error handling for malformed data
   - Coverage target: 100% (mappers are pure functions)

2. **Repository Integration Tests** (`tests/integration/*_test.go`):
   - Test CRUD operations return entities
   - Test mapper integration with GORM
   - Use in-memory SQLite for speed
   - Coverage target: All repository methods

3. **Contract Tests** (`tests/contract/*_test.go`):
   - Verify JSON API contracts unchanged
   - Test DTO serialization matches expected format
   - Ensure backward compatibility
   - Coverage target: All endpoints affected by migration

**Existing Test Updates**:
- 19 test files currently import `internal/models`
- Update imports to `internal/dtos` or use entities via mappers
- No functional changes, only import path updates

---

### 6. Migration Path for Legacy Imports

**Decision**: Incremental migration with automated verification

**Strategy**:

**Phase 1: Rename Package**
```bash
# Rename internal/models → internal/dtos
mv backend/internal/models backend/internal/dtos
# Rename files for clarity
mv backend/internal/dtos/user.go backend/internal/dtos/user_dto.go
mv backend/internal/dtos/task.go backend/internal/dtos/task_dto.go
```

**Phase 2: Create Mappers**
- Generate `application/mappers/user_mapper.go`
- Generate `application/mappers/task_mapper.go`
- Write comprehensive mapper tests (100% coverage)

**Phase 3: Update Repositories**
- Inject mappers into GORM repositories
- Convert repository methods to use entities
- Test with existing integration tests

**Phase 4: Update Services & Handlers**
- Services receive/return entities
- Handlers use mappers at boundaries (JSON ↔ Entity)
- Maintain API contract compatibility

**Phase 5: Update Tests**
- Update import paths in 19 test files
- Add mapper-specific tests
- Verify all 51+ tests pass

**Phase 6: Cleanup**
- Verify no `internal/models` references remain
- Remove deprecated code paths
- Update documentation

**Verification Commands**:
```bash
# Check for legacy imports
grep -r "internal/models" backend/

# Verify tests pass
go test ./... -v

# Verify build succeeds
go build ./...
```

---

### 7. Performance Considerations

**Decision**: Accept minimal mapper overhead as acceptable trade-off for clean architecture

**Analysis**:

**Mapper Overhead**:
- Estimated: 0.1-0.5 ms per conversion (struct field copying)
- Typical request: 2-4 mapper calls (DTO→Entity, Entity→DTO)
- Total overhead: < 2ms per request (< 4% of 50ms target)

**Optimization Strategies**:
1. **Object Pooling**: Not needed (Go GC handles small allocations efficiently)
2. **Pointer Passing**: Already implemented (avoid struct copies)
3. **Lazy Loading**: Not applicable (mappers are stateless)
4. **Caching**: Not needed (mappers are pure functions)

**Benchmark Target**:
- Single mapper call: < 1ms
- Full request cycle: < 50ms (including DB query)
- Build time: < 10 seconds (no regression)

**Monitoring**:
- Add performance tests for mapper operations
- Track API response times before/after migration
- Ensure no degradation in build/test times

---

### 8. Backward Compatibility Guarantees

**Decision**: Zero-tolerance policy for breaking changes

**Guarantees**:

1. **API Contracts**: JSON structure unchanged
   - Same field names and types
   - Same HTTP status codes
   - Same error response formats

2. **Database Schema**: No migrations required
   - DTOs preserve GORM tags
   - TableName() methods maintained
   - Existing data remains readable

3. **Test Coverage**: 100% pass rate maintained
   - All 51+ existing tests must pass
   - No changes to test expectations
   - Only import path updates

**Validation**:
```bash
# Before migration
curl http://localhost:8080/api/users/123 > before.json

# After migration
curl http://localhost:8080/api/users/123 > after.json

# Verify identical
diff before.json after.json  # Must be empty
```

---

### 9. Open Questions Resolution

**Q1: Mapper Location - `application/mappers/` vs `infrastructure/mappers/`?**

**Decision**: `application/mappers/`

**Rationale**:
- Application layer coordinates between domain and infrastructure
- Mappers orchestrate conversion between layers (not infrastructure concern)
- Aligns with hexagonal architecture (mappers are application services)
- Follows Clean Architecture naming (application layer contains use cases)

---

**Q2: DTO Naming - `internal/dtos/` vs `presentation/dtos/`?**

**Decision**: `internal/dtos/`

**Rationale**:
- DTOs used for both API (presentation) and persistence (infrastructure)
- `internal/` signals "implementation detail, not public API"
- Minimal change from existing `internal/models/` (clearer migration)
- Go convention: `internal/` for unexported packages

---

**Q3: Validation Placement - Mappers vs Entities?**

**Decision**: Both (different validation types)

**Validation Types**:
1. **Format Validation** (Mappers): Email syntax, UUID format, required fields
2. **Business Validation** (Entities): Uniqueness, state transitions, business rules
3. **Input Validation** (Handlers): JSON structure, HTTP headers, request limits

**Example**:
```go
// Mapper: Format validation
func (m *UserMapper) ToEntity(dto *dtos.UserDTO) (*entities.User, error) {
    if dto.Email == "" {
        return nil, errors.New("email required")
    }
    email, err := valueobjects.NewEmail(dto.Email)  // Format check
    if err != nil {
        return nil, err
    }
    return entities.NewUser(...), nil
}

// Entity: Business validation
func (u *User) UpdateEmail(newEmail valueobjects.Email) error {
    if u.IsArchived() {
        return errors.New("cannot update archived user")  // Business rule
    }
    u.email = newEmail
    return nil
}
```

---

**Q4: Error Handling - Return errors vs panic?**

**Decision**: Return errors from `ToEntity()`, no panics

**Rationale**:
- Go idiom: errors are values, handle explicitly
- Panics reserved for programmer errors (nil pointer, index out of bounds)
- Invalid external data is expected, not exceptional
- Allows caller to decide HTTP status code (400 vs 500)

**Pattern**:
```go
// Handler converts mapper errors to HTTP responses
user, err := h.mapper.ToEntity(dto)
if err != nil {
    c.JSON(400, gin.H{"error": err.Error()})  // Bad Request
    return
}
```

---

## Summary

All technical unknowns resolved. Key decisions:

1. **Mapper Pattern**: Stateless structs in `application/mappers/`
2. **Dual Models**: DTOs for persistence, entities for business logic
3. **Error Handling**: Return errors from ToEntity(), no panics
4. **Testing**: 3-layer coverage (mapper, repository, contract)
5. **Migration**: Incremental, automated verification
6. **Performance**: < 1ms mapper overhead acceptable
7. **Compatibility**: Zero breaking changes to API/database

No remaining NEEDS CLARIFICATION items. Ready for Phase 1 (Design & Contracts).

---

**Research Completed**: 2025-10-04
**Next Phase**: Phase 1 - Design & Contracts
