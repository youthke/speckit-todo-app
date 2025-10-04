# Quickstart: Complete DDD Migration

**Feature**: 010-complete-ddd-migration
**Date**: 2025-10-04
**Estimated Time**: 15 minutes

## Overview

This quickstart guide provides step-by-step instructions to verify the DDD migration completed successfully. It tests the mapper layer, repository integration, and API contract preservation.

---

## Prerequisites

- Go 1.24.7 installed
- Backend already running or can be started
- SQLite database initialized
- Feature 009 completed (Auth/Health domains migrated)

---

## Setup

### 1. Navigate to Backend Directory
```bash
cd /Users/youthke/practice/speckit/todo-app/backend
```

### 2. Verify Dependencies
```bash
go mod verify
```

**Expected Output**:
```
all modules verified
```

### 3. Build Backend
```bash
go build ./...
```

**Expected Output**:
```
# No errors, successful compilation
```

**Success Criteria**: Build completes in < 10 seconds with zero errors

---

## Verification Tests

### Test 1: Mapper Unit Tests

**Purpose**: Verify mappers correctly convert between DTOs and entities

**Command**:
```bash
go test -v ./application/mappers/...
```

**Expected Output**:
```
=== RUN   TestUserMapper_ToEntity
--- PASS: TestUserMapper_ToEntity (0.00s)
=== RUN   TestUserMapper_ToDTO
--- PASS: TestUserMapper_ToDTO (0.00s)
=== RUN   TestUserMapper_ToEntity_InvalidEmail
--- PASS: TestUserMapper_ToEntity_InvalidEmail (0.00s)
=== RUN   TestTaskMapper_ToEntity
--- PASS: TestTaskMapper_ToEntity (0.00s)
=== RUN   TestTaskMapper_ToDTO
--- PASS: TestTaskMapper_ToDTO (0.00s)
=== RUN   TestTaskMapper_ToEntity_EmptyTitle
--- PASS: TestTaskMapper_ToEntity_EmptyTitle (0.00s)
PASS
ok      application/mappers     0.XXXs
```

**Success Criteria**:
- ✅ All mapper tests pass
- ✅ Coverage ≥ 90%
- ✅ No panics or errors

---

### Test 2: Repository Integration Tests

**Purpose**: Verify repositories use mappers correctly and persist entities

**Command**:
```bash
go test -v ./tests/integration/... -run "TestUserRepository|TestTaskRepository"
```

**Expected Output**:
```
=== RUN   TestUserRepository_Save
--- PASS: TestUserRepository_Save (0.XX s)
=== RUN   TestUserRepository_FindByID
--- PASS: TestUserRepository_FindByID (0.XX s)
=== RUN   TestTaskRepository_Save
--- PASS: TestTaskRepository_Save (0.XX s)
=== RUN   TestTaskRepository_FindByID
--- PASS: TestTaskRepository_FindByID (0.XX s)
PASS
```

**Success Criteria**:
- ✅ Repositories return entities (not DTOs)
- ✅ Data persisted correctly to database
- ✅ Mapper conversions successful

---

### Test 3: Contract Tests (API Compatibility)

**Purpose**: Verify API responses unchanged after migration

**Command**:
```bash
go test -v ./tests/contract/users_*.go
go test -v ./tests/contract/tasks_*.go
```

**Expected Output**:
```
=== RUN   TestUsersProfile_Get
--- PASS: TestUsersProfile_Get (0.XX s)
=== RUN   TestUsersRegister_Post
--- PASS: TestUsersRegister_Post (0.XX s)
=== RUN   TestTasksGet_Success
--- PASS: TestTasksGet_Success (0.XX s)
=== RUN   TestTasksPost_Success
--- PASS: TestTasksPost_Success (0.XX s)
=== RUN   TestTasksPut_Success
--- PASS: TestTasksPut_Success (0.XX s)
=== RUN   TestTasksDelete_Success
--- PASS: TestTasksDelete_Success (0.XX s)
PASS
```

**Success Criteria**:
- ✅ All contract tests pass
- ✅ JSON responses match expected schema
- ✅ No breaking changes detected

---

### Test 4: Full Test Suite

**Purpose**: Verify all tests pass (unit, integration, contract, domain)

**Command**:
```bash
go test ./... -v -count=1
```

**Expected Output**:
```
?       todo-app/cmd/server     [no test files]
ok      todo-app/application/mappers          0.XXXs
ok      todo-app/domain/user/entities         0.XXXs
ok      todo-app/domain/user/valueobjects     0.XXXs
ok      todo-app/domain/task/entities         0.XXXs
ok      todo-app/domain/task/valueobjects     0.XXXs
ok      todo-app/tests/unit                   0.XXXs
ok      todo-app/tests/integration            X.XXXs
ok      todo-app/tests/contract               X.XXXs
PASS
```

**Success Criteria**:
- ✅ All 51+ tests pass
- ✅ Zero test failures
- ✅ Total run time < 30 seconds

---

## Manual API Testing

### Test 5: User Profile Endpoint

**Purpose**: Verify User entity → UserDTO conversion in real API

**Start Server**:
```bash
# Terminal 1
cd backend
go run cmd/server/main.go
```

**Create Test User** (Terminal 2):
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
  }' | jq
```

**Expected Response**:
```json
{
  "id": 1,
  "email": "test@example.com",
  "name": "Test User",
  "auth_method": "password",
  "is_active": true,
  "created_at": "2025-10-04T...",
  "updated_at": "2025-10-04T..."
}
```

**Get User Profile**:
```bash
curl http://localhost:8080/api/users/1 | jq
```

**Expected Response**: Same structure as above

**Success Criteria**:
- ✅ User created successfully
- ✅ Response matches UserDTO schema
- ✅ No password_hash in response
- ✅ Timestamps formatted correctly

---

### Test 6: Task CRUD Endpoint

**Purpose**: Verify Task entity → TaskDTO conversion in real API

**Create Task**:
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {session_token}" \
  -d '{
    "title": "Test migration"
  }' | jq
```

**Expected Response**:
```json
{
  "id": 1,
  "title": "Test migration",
  "completed": false,
  "created_at": "2025-10-04T...",
  "updated_at": "2025-10-04T..."
}
```

**List Tasks**:
```bash
curl http://localhost:8080/api/tasks \
  -H "Authorization: Bearer {session_token}" | jq
```

**Expected Response**:
```json
{
  "tasks": [
    {
      "id": 1,
      "title": "Test migration",
      "completed": false,
      "created_at": "2025-10-04T...",
      "updated_at": "2025-10-04T..."
    }
  ],
  "count": 1
}
```

**Update Task**:
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {session_token}" \
  -d '{
    "completed": true
  }' | jq
```

**Expected Response**:
```json
{
  "id": 1,
  "title": "Test migration",
  "completed": true,
  "created_at": "2025-10-04T...",
  "updated_at": "2025-10-04T..."
}
```

**Success Criteria**:
- ✅ Task created with default completed=false
- ✅ Task listed correctly
- ✅ Task updated successfully
- ✅ Response matches TaskDTO schema

---

## Code Verification

### Test 7: Import Verification

**Purpose**: Ensure no legacy `internal/models` imports remain

**Command**:
```bash
grep -r "internal/models" backend/ --include="*.go" | grep -v "internal/dtos"
```

**Expected Output**:
```
# Empty (no matches)
```

**Alternative Check**:
```bash
# Should only find internal/dtos references
grep -r "internal/dtos" backend/ --include="*.go" | wc -l
```

**Expected Output**: `> 0` (multiple files import dtos)

**Success Criteria**:
- ✅ Zero references to `internal/models` (except in migration scripts)
- ✅ All imports use `internal/dtos` or domain entities

---

### Test 8: Mapper Dependency Check

**Purpose**: Verify repositories inject mappers

**Command**:
```bash
grep -A 5 "type Gorm.*Repository struct" backend/infrastructure/persistence/*.go
```

**Expected Output** (example):
```go
type GormUserRepository struct {
    db     *gorm.DB
    mapper *mappers.UserMapper
}

type GormTaskRepository struct {
    db     *gorm.DB
    mapper *mappers.TaskMapper
}
```

**Success Criteria**:
- ✅ All GORM repositories have mapper field
- ✅ Mappers injected via constructors

---

### Test 9: Entity Return Types

**Purpose**: Verify repositories return entities (not DTOs)

**Command**:
```bash
grep "func.*FindByID" backend/infrastructure/persistence/*.go -A 10
```

**Expected Pattern**:
```go
func (r *GormUserRepository) FindByID(id valueobjects.UserID) (*entities.User, error) {
    var dto dtos.UserDTO
    // ... GORM query ...
    return r.mapper.ToEntity(&dto)
}
```

**Success Criteria**:
- ✅ Return type is `*entities.User` or `*entities.Task`
- ✅ DTO used internally, entity returned
- ✅ Mapper called for conversion

---

## Performance Validation

### Test 10: Mapper Performance Benchmark

**Purpose**: Ensure mapper overhead < 1ms

**Command**:
```bash
go test -bench=BenchmarkMapper -benchmem ./application/mappers/
```

**Expected Output**:
```
BenchmarkUserMapper_ToEntity-8     1000000    800 ns/op    256 B/op    4 allocs/op
BenchmarkUserMapper_ToDTO-8        2000000    500 ns/op    128 B/op    2 allocs/op
BenchmarkTaskMapper_ToEntity-8     1000000    600 ns/op    192 B/op    3 allocs/op
BenchmarkTaskMapper_ToDTO-8        2000000    400 ns/op     96 B/op    1 allocs/op
PASS
```

**Success Criteria**:
- ✅ Single conversion < 1000 ns (1 µs)
- ✅ Memory allocations < 500 B per call
- ✅ No performance regression

---

### Test 11: API Response Time

**Purpose**: Verify no latency regression after mapper introduction

**Command**:
```bash
# Measure response time
time curl -w "\nTime: %{time_total}s\n" http://localhost:8080/api/tasks \
  -H "Authorization: Bearer {session_token}" \
  -o /dev/null -s
```

**Expected Output**:
```
Time: 0.02s  # < 50ms
```

**Success Criteria**:
- ✅ Response time < 50ms
- ✅ No noticeable slowdown compared to pre-migration

---

## Cleanup

### Stop Server
```bash
# Terminal 1: Ctrl+C to stop server
```

### Clean Test Database
```bash
rm backend/todo.db  # Remove test database
```

---

## Success Checklist

**Phase 0: Build & Compile**
- [ ] `go build ./...` succeeds
- [ ] No compilation errors
- [ ] Build time < 10 seconds

**Phase 1: Unit Tests**
- [ ] Mapper tests pass (100% coverage)
- [ ] Domain entity tests pass
- [ ] Value object tests pass

**Phase 2: Integration Tests**
- [ ] Repository tests pass (13 tests)
- [ ] Mappers integrated correctly
- [ ] Entities persisted/retrieved successfully

**Phase 3: Contract Tests**
- [ ] User endpoint tests pass (2 tests)
- [ ] Task endpoint tests pass (7 tests)
- [ ] JSON schema unchanged

**Phase 4: Full Test Suite**
- [ ] All 51+ tests pass
- [ ] Zero test failures
- [ ] No import errors

**Phase 5: Manual API Testing**
- [ ] User registration works
- [ ] Task CRUD operations work
- [ ] API responses match expected schema

**Phase 6: Code Verification**
- [ ] No legacy `internal/models` imports
- [ ] All repositories use mappers
- [ ] All repositories return entities

**Phase 7: Performance**
- [ ] Mapper overhead < 1ms
- [ ] API response time < 50ms
- [ ] No performance regression

---

## Troubleshooting

### Issue: Mapper tests fail with "cannot convert"

**Cause**: Value object constructors may return errors

**Fix**: Check mapper error handling in `ToEntity()` methods

```go
func (m *UserMapper) ToEntity(dto *dtos.UserDTO) (*entities.User, error) {
    email, err := valueobjects.NewEmail(dto.Email)
    if err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }
    // ...
}
```

---

### Issue: Contract tests fail with JSON mismatch

**Cause**: DTO fields missing or renamed

**Fix**: Verify DTO struct tags match legacy models

```go
type UserDTO struct {
    ID    uint   `json:"id"`      // NOT "user_id"
    Email string `json:"email"`   // NOT "Email"
}
```

---

### Issue: Repository tests fail with "mapper is nil"

**Cause**: Mapper not injected in repository constructor

**Fix**: Update repository initialization

```go
repo := persistence.NewGormUserRepository(db, &mappers.UserMapper{})
```

---

### Issue: "internal/models" import errors

**Cause**: Some files not updated to use `internal/dtos`

**Fix**: Find and replace all imports

```bash
find backend -name "*.go" -exec sed -i '' 's|internal/models|internal/dtos|g' {} \;
```

---

## Rollback Plan (If Needed)

If critical issues detected:

1. **Revert commits**:
   ```bash
   git log --oneline  # Find last good commit
   git revert <commit-hash>
   ```

2. **Restore `internal/models`**:
   ```bash
   git checkout HEAD~1 -- backend/internal/models
   ```

3. **Remove mapper layer**:
   ```bash
   rm -rf backend/application/mappers
   ```

4. **Run tests to verify rollback**:
   ```bash
   go test ./... -v
   ```

---

## Next Steps

After successful quickstart completion:

1. **Document learnings**: Update CLAUDE.md with mapper pattern
2. **Performance monitoring**: Add metrics for mapper operations
3. **Code cleanup**: Remove deprecated code comments
4. **Team review**: Schedule walkthrough of mapper architecture

---

**Quickstart Version**: 1.0
**Last Updated**: 2025-10-04
**Estimated Completion Time**: 15 minutes
