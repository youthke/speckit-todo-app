# Data Model: Complete DDD Migration

**Feature**: 010-complete-ddd-migration
**Date**: 2025-10-04
**Status**: Complete

## Overview

This document defines the data model for the mapper layer that bridges DTOs (database/API representation) and DDD entities (domain representation). The migration introduces a new mapper layer while maintaining existing DTO and entity structures.

---

## Model Layers

### Layer 1: DTOs (Data Transfer Objects)
**Location**: `internal/dtos/`
**Purpose**: Database persistence via GORM, JSON API serialization
**Characteristics**: Public fields, GORM tags, no business logic

### Layer 2: Entities (Domain Models)
**Location**: `domain/*/entities/`
**Purpose**: Business logic, invariants, domain rules
**Characteristics**: Private fields, value objects, rich behavior

### Layer 3: Mappers (Transformation Layer)
**Location**: `application/mappers/`
**Purpose**: Bidirectional conversion between DTOs and Entities
**Characteristics**: Stateless, pure functions, validation

---

## DTO Models

### UserDTO

**File**: `internal/dtos/user_dto.go`

**Fields**:
```go
type UserDTO struct {
    ID             uint       `json:"id" gorm:"primaryKey"`
    Email          string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
    Name           string     `json:"name" gorm:"type:varchar(255);not null"`
    PasswordHash   string     `json:"-" gorm:"type:varchar(255)"`
    AuthMethod     string     `json:"auth_method" gorm:"type:varchar(50);not null;default:'password'"`
    GoogleID       string     `json:"google_id,omitempty" gorm:"type:varchar(255);uniqueIndex"`
    OAuthProvider  string     `json:"oauth_provider,omitempty" gorm:"type:varchar(50)"`
    OAuthCreatedAt *time.Time `json:"oauth_created_at,omitempty"`
    IsActive       bool       `json:"is_active" gorm:"default:true"`
    CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
```

**GORM Configuration**:
- Table: `users`
- Primary Key: `id`
- Unique Indexes: `email`, `google_id`
- Auto-managed: `created_at`, `updated_at`

**Validation Rules**:
- Email: Required, unique, valid email format
- Name: Required, max 255 characters
- Either `password_hash` OR `google_id` must be present
- If `google_id` present, `oauth_provider` must be "google"

---

### TaskDTO

**File**: `internal/dtos/task_dto.go`

**Fields**:
```go
type TaskDTO struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title" gorm:"type:varchar(500);not null"`
    Completed bool      `json:"completed" gorm:"default:false"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

**GORM Configuration**:
- Table: `tasks`
- Primary Key: `id`
- Auto-managed: `created_at`, `updated_at`

**Validation Rules**:
- Title: Required, max 500 characters
- Completed: Boolean, defaults to false

---

## Entity Models

### User Entity

**File**: `domain/user/entities/user.go`

**Structure**:
```go
type User struct {
    id          valueobjects.UserID
    email       valueobjects.Email
    profile     valueobjects.UserProfile
    preferences valueobjects.UserPreferences
    createdAt   time.Time
    updatedAt   time.Time
}
```

**Value Objects**:
- `UserID`: Unique identifier (uint)
- `Email`: Validated email address with domain methods
- `UserProfile`: Display name and profile data
- `UserPreferences`: Theme, notifications, default task priority

**Business Methods**:
- `UpdateProfile(profile)`: Change user profile
- `UpdatePreferences(prefs)`: Modify user preferences
- `ChangeEmail(email)`: Update email with validation
- `EnableEmailNotifications()`: Toggle notifications
- `UpdateThemePreference(theme)`: Change theme

**Invariants**:
- Email cannot be empty
- User ID cannot be zero
- Profile display name required

---

### Task Entity

**File**: `domain/task/entities/task.go`

**Structure**:
```go
type Task struct {
    id          valueobjects.TaskID
    title       valueobjects.TaskTitle
    description valueobjects.TaskDescription
    status      valueobjects.TaskStatus
    priority    valueobjects.TaskPriority
    userID      uservo.UserID
    createdAt   time.Time
    updatedAt   time.Time
}
```

**Value Objects**:
- `TaskID`: Unique identifier (uint)
- `TaskTitle`: Validated title (max 500 chars)
- `TaskDescription`: Optional description
- `TaskStatus`: Pending/Completed/Archived
- `TaskPriority`: Low/Medium/High
- `UserID`: Owner reference

**Business Methods**:
- `MarkAsCompleted()`: Transition to completed status
- `UpdateTitle(title)`: Change title (if not archived)
- `UpdateDescription(desc)`: Change description
- `ChangePriority(priority)`: Update priority (if pending)
- `Archive()`: Mark as archived
- `IsOwnedBy(userID)`: Check ownership

**Invariants**:
- Task ID cannot be zero
- User ID cannot be zero
- Cannot modify archived tasks
- Can only change priority on pending tasks

---

## Mapper Models

### UserMapper

**File**: `application/mappers/user_mapper.go`

**Structure**:
```go
type UserMapper struct{}

func (m *UserMapper) ToEntity(dto *dtos.UserDTO) (*entities.User, error)
func (m *UserMapper) ToDTO(entity *entities.User) *dtos.UserDTO
```

**ToEntity Mapping** (DTO → Entity):
```
UserDTO.ID        → UserID value object
UserDTO.Email     → Email value object (with validation)
UserDTO.Name      → UserProfile value object
[defaults]        → UserPreferences value object
UserDTO.CreatedAt → createdAt
UserDTO.UpdatedAt → updatedAt
```

**ToDTO Mapping** (Entity → DTO):
```
User.ID()          → UserDTO.ID
User.Email()       → UserDTO.Email
User.Profile()     → UserDTO.Name (extract display name)
[defaults]         → UserDTO.AuthMethod, IsActive
User.CreatedAt()   → UserDTO.CreatedAt
User.UpdatedAt()   → UserDTO.UpdatedAt
```

**Error Conditions**:
- Invalid email format (ToEntity)
- Empty email (ToEntity)
- Zero user ID (ToEntity)

**Data Loss**:
- DTO fields not in Entity: `PasswordHash`, `GoogleID`, `OAuthProvider`, `AuthMethod`, `IsActive`
- These fields are managed separately by Auth domain (not User domain concern)
- Mapper focuses on core user profile data only

---

### TaskMapper

**File**: `application/mappers/task_mapper.go`

**Structure**:
```go
type TaskMapper struct{}

func (m *TaskMapper) ToEntity(dto *dtos.TaskDTO) (*entities.Task, error)
func (m *TaskMapper) ToDTO(entity *entities.Task) *dtos.TaskDTO
```

**ToEntity Mapping** (DTO → Task):
```
TaskDTO.ID        → TaskID value object
TaskDTO.Title     → TaskTitle value object (with validation)
[empty]           → TaskDescription value object
TaskDTO.Completed → TaskStatus value object (completed or pending)
[default]         → TaskPriority value object (medium)
[TBD]             → UserID value object (from context/session)
TaskDTO.CreatedAt → createdAt
TaskDTO.UpdatedAt → updatedAt
```

**ToDTO Mapping** (Task → TaskDTO):
```
Task.ID()         → TaskDTO.ID
Task.Title()      → TaskDTO.Title
Task.Status()     → TaskDTO.Completed (convert status to boolean)
Task.CreatedAt()  → TaskDTO.CreatedAt
Task.UpdatedAt()  → TaskDTO.UpdatedAt
```

**Error Conditions**:
- Empty title (ToEntity)
- Title > 500 characters (ToEntity)
- Zero task ID (ToEntity)

**Data Transformation**:
- `TaskDTO.Completed` (boolean) ↔ `TaskStatus` (enum: Pending/Completed/Archived)
  - `false` → Pending
  - `true` → Completed
  - Archived status not represented in DTO (requires separate flag or status field)

**Missing Fields**:
- Entity has `description`, `priority`, `userID` not in DTO
- DTO is minimal (ID, title, completed) for legacy API compatibility
- Future: Extend DTO to include all entity fields

---

## Relationships

### User ↔ Task
**Type**: One-to-Many (User has many Tasks)
**Implementation**: Task entity contains `userID` field
**Mapping**: Not handled by mappers (repository layer responsibility)

**Notes**:
- Current DTO models don't include relationships
- Future enhancement: Add `UserID` field to TaskDTO
- Current workaround: User ID passed separately in repository methods

---

## Validation Strategy

### Layer-Specific Validation

**DTO Validation** (GORM hooks):
- Format checks: email syntax, field lengths
- Database constraints: uniqueness, NOT NULL
- Handled by: GORM BeforeCreate/BeforeUpdate hooks

**Mapper Validation** (ToEntity):
- Value object creation: Email format, TaskTitle length
- Required field checks: non-empty, non-zero
- Handled by: Value object constructors

**Entity Validation** (Business logic):
- Business rules: cannot modify archived task
- State transitions: only pending tasks can be completed
- Handled by: Entity methods

---

## Migration Strategy

### Phase 1: Create DTOs from Models
```bash
# Rename models → dtos
mv internal/models internal/dtos
mv internal/dtos/user.go internal/dtos/user_dto.go
mv internal/dtos/task.go internal/dtos/task_dto.go

# Update package name
sed -i '' 's/package models/package dtos/g' internal/dtos/*.go
```

### Phase 2: Generate Mappers
- Create `application/mappers/user_mapper.go`
- Create `application/mappers/task_mapper.go`
- Implement ToEntity() and ToDTO() methods
- Handle partial field mapping (User entity doesn't include auth fields)

### Phase 3: Update Repositories
- Inject mappers into GORM repositories
- Convert DTO ↔ Entity at repository boundaries
- Return entities from all repository methods

### Phase 4: Update Handlers/Services
- Handlers receive/return DTOs (JSON)
- Services work with entities internally
- Use mappers at service boundaries

---

## Field Mapping Table

### User: DTO ↔ Entity

| DTO Field         | Entity Field       | Mapping Notes                          |
|-------------------|--------------------|----------------------------------------|
| `ID`              | `id (UserID)`      | uint → UserID value object             |
| `Email`           | `email (Email)`    | string → Email value object, validated |
| `Name`            | `profile`          | Extract display name from UserProfile  |
| `CreatedAt`       | `createdAt`        | Direct copy                            |
| `UpdatedAt`       | `updatedAt`        | Direct copy                            |
| `PasswordHash`    | ❌ Not mapped      | Auth domain concern, not user domain   |
| `AuthMethod`      | ❌ Not mapped      | Auth domain concern                    |
| `GoogleID`        | ❌ Not mapped      | Auth domain concern                    |
| `OAuthProvider`   | ❌ Not mapped      | Auth domain concern                    |
| `OAuthCreatedAt`  | ❌ Not mapped      | Auth domain concern                    |
| `IsActive`        | ❌ Not mapped      | Status managed separately              |
| ❌ No DTO field   | `preferences`      | Defaults to NewDefaultUserPreferences  |

### Task: DTO ↔ Entity

| DTO Field    | Entity Field          | Mapping Notes                         |
|--------------|-----------------------|---------------------------------------|
| `ID`         | `id (TaskID)`         | uint → TaskID value object            |
| `Title`      | `title (TaskTitle)`   | string → TaskTitle value object       |
| `Completed`  | `status (TaskStatus)` | bool → Pending/Completed enum         |
| `CreatedAt`  | `createdAt`           | Direct copy                           |
| `UpdatedAt`  | `updatedAt`           | Direct copy                           |
| ❌ No DTO    | `description`         | Defaults to empty TaskDescription     |
| ❌ No DTO    | `priority`            | Defaults to Medium priority           |
| ❌ No DTO    | `userID`              | Must be provided via context/session  |

---

## Persistence Flow

### Create Flow
```
1. Handler receives JSON → DTO
2. Mapper converts DTO → Entity (with validation)
3. Service applies business logic on Entity
4. Repository converts Entity → DTO (via mapper)
5. GORM persists DTO to database
```

### Read Flow
```
1. Repository fetches DTO from database (GORM)
2. Mapper converts DTO → Entity
3. Service processes Entity
4. Service returns Entity
5. Handler converts Entity → DTO (via mapper)
6. Handler serializes DTO → JSON response
```

---

## Test Data Examples

### UserDTO Example (JSON/Database)
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

### User Entity Example (Domain)
```go
User{
  id: UserID(1),
  email: Email("user@example.com"),
  profile: UserProfile{displayName: "John Doe"},
  preferences: UserPreferences{
    theme: "light",
    emailNotifications: true,
    defaultTaskPriority: "medium",
  },
  createdAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
  updatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
}
```

### TaskDTO Example (JSON/Database)
```json
{
  "id": 100,
  "title": "Complete migration",
  "completed": false,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

### Task Entity Example (Domain)
```go
Task{
  id: TaskID(100),
  title: TaskTitle("Complete migration"),
  description: TaskDescription(""),
  status: TaskStatus(Pending),
  priority: TaskPriority(Medium),
  userID: UserID(1),
  createdAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
  updatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
}
```

---

## Summary

**Model Count**:
- 2 DTOs (UserDTO, TaskDTO)
- 2 Entities (User, Task)
- 2 Mappers (UserMapper, TaskMapper)

**Total Fields Mapped**:
- User: 5 fields mapped, 6 fields ignored (auth-related)
- Task: 5 fields mapped, 3 fields defaulted

**Validation Points**:
- 3 layers: DTO (GORM), Mapper (format), Entity (business)

**Key Decisions**:
- Partial mapping accepted (User entity excludes auth fields)
- Default values used for missing entity fields (Task.userID from context)
- Boolean ↔ Enum transformation (TaskDTO.Completed ↔ TaskStatus)

---

**Data Model Completed**: 2025-10-04
**Next Phase**: Generate contracts and quickstart
