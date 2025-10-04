# API Contracts: Complete DDD Migration

**Feature**: 010-complete-ddd-migration
**Date**: 2025-10-04
**Contract Version**: 1.0
**Breaking Changes**: NONE (backward compatibility required)

## Overview

This document specifies the API contracts that MUST remain unchanged during the DDD migration. All endpoints must preserve their current request/response formats to ensure zero breaking changes.

---

## Contract Constraints

### Critical Requirements
1. **JSON Structure**: Field names, types, nesting must remain identical
2. **HTTP Status Codes**: Same codes for success/error scenarios
3. **Content-Type**: `application/json` for all requests/responses
4. **Authentication**: Same session/token mechanisms
5. **Error Format**: Consistent error response structure

### Testing Requirements
- All contract tests in `backend/tests/contract/` must pass
- Response schema validation before/after migration
- No changes to test expectations (only import paths)

---

## User Endpoints

### GET /api/users/:id

**Description**: Retrieve user profile by ID

**Request**:
```http
GET /api/users/1 HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

**Success Response** (200 OK):
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "auth_method": "password",
  "is_active": true,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-02T00:00:00Z"
}
```

**Error Response** (404 Not Found):
```json
{
  "error": "User not found"
}
```

**Contract Requirements**:
- ✅ All fields present in response
- ✅ Field types match (id: number, email: string, etc.)
- ✅ Timestamps in RFC3339 format
- ✅ OAuth fields omitted if not OAuth user

**Mapper Involvement**:
- Handler receives request → extracts user ID
- Repository returns User entity
- **Mapper converts User entity → UserDTO**
- Handler serializes UserDTO → JSON response

---

### POST /api/users/register

**Description**: Register new user with password

**Request**:
```http
POST /api/users/register HTTP/1.1
Content-Type: application/json

{
  "email": "newuser@example.com",
  "name": "Jane Smith",
  "password": "securepassword123"
}
```

**Success Response** (201 Created):
```json
{
  "id": 2,
  "email": "newuser@example.com",
  "name": "Jane Smith",
  "auth_method": "password",
  "is_active": true,
  "created_at": "2025-01-03T00:00:00Z",
  "updated_at": "2025-01-03T00:00:00Z"
}
```

**Error Response** (400 Bad Request):
```json
{
  "error": "Email already exists"
}
```

**Contract Requirements**:
- ✅ Password field not returned in response
- ✅ Default `auth_method` = "password"
- ✅ Default `is_active` = true
- ✅ Auto-generated timestamps

**Mapper Involvement**:
- Handler receives JSON → parses to CreateUserRequest DTO
- Service creates User entity (via domain factory)
- **Mapper converts User entity → UserDTO**
- Repository persists UserDTO via GORM
- Handler returns UserDTO as JSON

---

### PUT /api/users/:id/profile

**Description**: Update user profile (name)

**Request**:
```http
PUT /api/users/1/profile HTTP/1.1
Content-Type: application/json

{
  "name": "John Updated"
}
```

**Success Response** (200 OK):
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Updated",
  "auth_method": "password",
  "is_active": true,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z"
}
```

**Contract Requirements**:
- ✅ Updated `name` field reflected
- ✅ Updated `updated_at` timestamp
- ✅ Other fields unchanged

**Mapper Involvement**:
- Handler receives JSON → UpdateProfileRequest DTO
- Repository fetches User entity
- Service updates User entity profile
- **Mapper converts User entity → UserDTO**
- Handler returns updated UserDTO

---

## Task Endpoints

### GET /api/tasks

**Description**: Retrieve all tasks for authenticated user

**Request**:
```http
GET /api/tasks HTTP/1.1
Host: localhost:8080
Authorization: Bearer {session_token}
```

**Success Response** (200 OK):
```json
{
  "tasks": [
    {
      "id": 1,
      "title": "Complete migration",
      "completed": false,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "title": "Write tests",
      "completed": true,
      "created_at": "2025-01-02T00:00:00Z",
      "updated_at": "2025-01-02T12:00:00Z"
    }
  ],
  "count": 2
}
```

**Contract Requirements**:
- ✅ Array of task objects
- ✅ Count field = number of tasks
- ✅ Empty array if no tasks
- ✅ Tasks sorted by creation date (newest first)

**Mapper Involvement**:
- Repository fetches Task entities for user
- **Mapper converts []Task entities → []TaskDTO**
- Handler wraps in TaskResponse structure
- Handler serializes to JSON

---

### POST /api/tasks

**Description**: Create new task

**Request**:
```http
POST /api/tasks HTTP/1.1
Content-Type: application/json
Authorization: Bearer {session_token}

{
  "title": "New task"
}
```

**Success Response** (201 Created):
```json
{
  "id": 3,
  "title": "New task",
  "completed": false,
  "created_at": "2025-01-03T00:00:00Z",
  "updated_at": "2025-01-03T00:00:00Z"
}
```

**Error Response** (400 Bad Request):
```json
{
  "error": "Title is required"
}
```

**Contract Requirements**:
- ✅ Default `completed` = false
- ✅ Auto-generated ID
- ✅ Auto-generated timestamps
- ✅ Title max 500 characters

**Mapper Involvement**:
- Handler receives JSON → CreateTaskRequest DTO
- Service creates Task entity with user ID from session
- **Mapper converts Task entity → TaskDTO**
- Repository persists TaskDTO
- Handler returns TaskDTO

---

### GET /api/tasks/:id

**Description**: Retrieve single task by ID

**Request**:
```http
GET /api/tasks/1 HTTP/1.1
Authorization: Bearer {session_token}
```

**Success Response** (200 OK):
```json
{
  "id": 1,
  "title": "Complete migration",
  "completed": false,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**Error Response** (404 Not Found):
```json
{
  "error": "Task not found"
}
```

**Contract Requirements**:
- ✅ Task belongs to authenticated user (ownership check)
- ✅ 404 if task doesn't exist or belongs to another user

**Mapper Involvement**:
- Repository fetches Task entity by ID
- Service verifies ownership via `Task.IsOwnedBy(userID)`
- **Mapper converts Task entity → TaskDTO**
- Handler returns TaskDTO

---

### PUT /api/tasks/:id

**Description**: Update task (title and/or completion status)

**Request**:
```http
PUT /api/tasks/1 HTTP/1.1
Content-Type: application/json
Authorization: Bearer {session_token}

{
  "title": "Updated task title",
  "completed": true
}
```

**Success Response** (200 OK):
```json
{
  "id": 1,
  "title": "Updated task title",
  "completed": true,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-03T15:00:00Z"
}
```

**Contract Requirements**:
- ✅ Partial updates supported (title only, completed only, or both)
- ✅ Updated timestamp reflects change
- ✅ Ownership verified

**Mapper Involvement**:
- Handler receives JSON → UpdateTaskRequest DTO
- Repository fetches Task entity
- Service updates Task entity (UpdateTitle, MarkAsCompleted)
- **Mapper converts Task entity → TaskDTO**
- Handler returns updated TaskDTO

---

### DELETE /api/tasks/:id

**Description**: Delete task

**Request**:
```http
DELETE /api/tasks/1 HTTP/1.1
Authorization: Bearer {session_token}
```

**Success Response** (204 No Content):
```http
HTTP/1.1 204 No Content
```

**Error Response** (404 Not Found):
```json
{
  "error": "Task not found"
}
```

**Contract Requirements**:
- ✅ No response body on success (204)
- ✅ Task permanently deleted from database
- ✅ Ownership verified

**Mapper Involvement**:
- Repository fetches Task entity
- Service verifies ownership
- Repository deletes via GORM (no mapper needed for deletion)

---

## Authentication Endpoints (Out of Scope)

The following endpoints are **NOT affected** by this migration as they belong to the Auth domain (already migrated in feature 009):

- `POST /api/auth/google/login`
- `GET /api/auth/google/callback`
- `POST /api/auth/logout`
- `POST /api/auth/session/refresh`
- `GET /api/auth/session/validate`

These endpoints use Auth domain entities and are not part of the User/Task mapper migration.

---

## Health Endpoint (Out of Scope)

**GET /api/health** - Already migrated in feature 009, no changes needed.

---

## Contract Test Coverage

### Existing Contract Tests (Must Pass)
- `tests/contract/users_profile_test.go` - User profile endpoints
- `tests/contract/users_register_test.go` - User registration
- `tests/contract/tasks_get_test.go` - Get all tasks
- `tests/contract/tasks_post_test.go` - Create task
- `tests/contract/tasks_get_by_id_test.go` - Get task by ID
- `tests/contract/tasks_put_test.go` - Update task
- `tests/contract/tasks_put_update_test.go` - Partial update task
- `tests/contract/tasks_delete_test.go` - Delete task
- `tests/contract/tasks_delete_new_test.go` - Delete newly created task

### Test Modifications Required
- **Import path updates**: Change `internal/models` → `internal/dtos`
- **No functional changes**: Test expectations remain identical
- **Mapper integration**: Tests indirectly verify mapper correctness via API

---

## Validation Strategy

### Pre-Migration Baseline
```bash
# Capture API responses before migration
curl http://localhost:8080/api/users/1 > baseline-user.json
curl http://localhost:8080/api/tasks > baseline-tasks.json
```

### Post-Migration Verification
```bash
# Capture API responses after migration
curl http://localhost:8080/api/users/1 > migrated-user.json
curl http://localhost:8080/api/tasks > migrated-tasks.json

# Verify identical (ignoring timestamps)
diff <(jq -S 'del(.updated_at)' baseline-user.json) \
     <(jq -S 'del(.updated_at)' migrated-user.json)
# Expected: No differences
```

### Contract Test Suite
```bash
# All contract tests must pass
go test ./tests/contract/... -v

# Expected: 0 failures
```

---

## Error Response Format

All endpoints must maintain consistent error format:

**Structure**:
```json
{
  "error": "Error message describing what went wrong"
}
```

**HTTP Status Codes**:
- `400 Bad Request`: Invalid input, validation failure
- `401 Unauthorized`: Missing or invalid authentication
- `404 Not Found`: Resource doesn't exist
- `500 Internal Server Error`: Server-side error

**Examples**:
```json
// Validation error
{"error": "Title is required"}

// Not found
{"error": "Task not found"}

// Authentication
{"error": "Invalid session token"}
```

---

## Summary

**Total Endpoints Affected**: 8 endpoints (5 user + 3 task endpoints)
**Breaking Changes**: NONE
**Mapper Touchpoints**: All read/write operations involve mapper conversion
**Contract Tests**: 9 test files must pass unchanged

**Key Principle**: DTOs serve as the API contract. Entities are internal implementation details. Mappers bridge the two without exposing entity structure to clients.

**Verification Checklist**:
- [ ] All contract tests pass
- [ ] API responses match baseline (before/after comparison)
- [ ] No changes to JSON field names or types
- [ ] HTTP status codes unchanged
- [ ] Error messages consistent
- [ ] Timestamps formatted correctly
- [ ] Default values applied correctly

---

**Contracts Documented**: 2025-10-04
**Next Document**: quickstart.md
