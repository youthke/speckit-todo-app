# Data Model: Backend Domain-Driven Design Implementation

**Feature**: Backend Domain-Driven Design Implementation
**Date**: 2025-09-28
**Status**: Design Phase

## Domain Models

### Task Management Domain

#### Task Aggregate Root
```go
// Domain Entity
type Task struct {
    id          TaskID
    title       TaskTitle        // Value Object
    description TaskDescription  // Value Object
    status      TaskStatus       // Value Object
    priority    TaskPriority     // Value Object
    createdAt   time.Time
    updatedAt   time.Time
    userID      UserID           // Foreign key to User domain
}

type TaskID struct {
    value uint
}

type TaskTitle struct {
    value string // max 500 characters, non-empty
}

type TaskDescription struct {
    value string // optional, max 2000 characters
}

type TaskStatus struct {
    value string // "pending", "completed", "archived"
}

type TaskPriority struct {
    value string // "low", "medium", "high"
}
```

#### Task Domain Services
```go
type TaskValidationService interface {
    ValidateTaskCreation(title TaskTitle, userID UserID) error
    ValidateTaskUpdate(task *Task, updates TaskUpdates) error
}

type TaskSearchService interface {
    FindTasksByStatus(userID UserID, status TaskStatus) ([]Task, error)
    FindTasksByPriority(userID UserID, priority TaskPriority) ([]Task, error)
}
```

#### Task Repository Interface
```go
type TaskRepository interface {
    Save(task *Task) error
    FindByID(id TaskID) (*Task, error)
    FindByUserID(userID UserID) ([]Task, error)
    Delete(id TaskID) error
    Update(task *Task) error
}
```

### User Management Domain

#### User Aggregate Root
```go
// Domain Entity
type User struct {
    id          UserID
    email       Email            // Value Object
    profile     UserProfile      // Value Object
    preferences UserPreferences  // Value Object
    createdAt   time.Time
    updatedAt   time.Time
}

type UserID struct {
    value uint
}

type Email struct {
    value string // validated email format
}

type UserProfile struct {
    firstName string
    lastName  string
    timezone  string
}

type UserPreferences struct {
    defaultTaskPriority TaskPriority
    emailNotifications  bool
    themePreference     string
}
```

#### User Domain Services
```go
type UserAuthenticationService interface {
    ValidateEmailUniqueness(email Email) error
    GenerateUserCredentials(email Email) (*UserCredentials, error)
}

type UserProfileService interface {
    UpdateProfile(userID UserID, profile UserProfile) error
    ValidateProfileData(profile UserProfile) error
}
```

#### User Repository Interface
```go
type UserRepository interface {
    Save(user *User) error
    FindByID(id UserID) (*User, error)
    FindByEmail(email Email) (*User, error)
    Update(user *User) error
    Delete(id UserID) error
}
```

## Cross-Domain Relationships

### Domain Events
```go
type TaskCreatedEvent struct {
    TaskID    TaskID
    UserID    UserID
    Title     TaskTitle
    CreatedAt time.Time
}

type TaskCompletedEvent struct {
    TaskID      TaskID
    UserID      UserID
    CompletedAt time.Time
}

type UserRegisteredEvent struct {
    UserID    UserID
    Email     Email
    CreatedAt time.Time
}
```

## Value Object Validation Rules

### Task Domain
- **TaskTitle**: Required, 1-500 characters, no special characters except alphanumeric, spaces, and basic punctuation
- **TaskDescription**: Optional, max 2000 characters
- **TaskStatus**: Must be one of: "pending", "completed", "archived"
- **TaskPriority**: Must be one of: "low", "medium", "high"

### User Domain
- **Email**: Valid email format, unique across system, max 255 characters
- **UserProfile.firstName**: Required, 1-50 characters, letters only
- **UserProfile.lastName**: Required, 1-50 characters, letters only
- **UserProfile.timezone**: Valid IANA timezone identifier

## Entity Behavior

### Task Entity Methods
```go
func (t *Task) MarkAsCompleted() error
func (t *Task) UpdateTitle(title TaskTitle) error
func (t *Task) UpdateDescription(description TaskDescription) error
func (t *Task) ChangePriority(priority TaskPriority) error
func (t *Task) Archive() error
func (t *Task) IsOwnedBy(userID UserID) bool
```

### User Entity Methods
```go
func (u *User) UpdateProfile(profile UserProfile) error
func (u *User) UpdatePreferences(preferences UserPreferences) error
func (u *User) ChangeEmail(email Email) error
func (u *User) GetDisplayName() string
```

## Aggregate Consistency Rules

### Task Aggregate
- Only the task owner (UserID) can modify the task
- Task title cannot be empty when status is "pending" or "completed"
- Archived tasks cannot be modified except to unarchive
- Task priority can only be changed on pending tasks

### User Aggregate
- Email address must be unique across all users
- User profile changes require all mandatory fields (firstName, lastName)
- User preferences are optional and have sensible defaults
- User cannot be deleted if they have active (non-archived) tasks

## Migration from Current Model

### Existing Task Model Mapping
```go
// Current models.Task -> New domain.Task
ID        -> TaskID{value: ID}
Title     -> TaskTitle{value: Title}
Completed -> TaskStatus{value: completed ? "completed" : "pending"}
CreatedAt -> createdAt
UpdatedAt -> updatedAt
// New fields with defaults:
Description -> TaskDescription{value: ""}
Priority    -> TaskPriority{value: "medium"}
UserID      -> UserID{value: 1} // default user for migration
```

### Data Migration Strategy
1. Create default user for existing tasks
2. Convert boolean `Completed` to `TaskStatus` enum
3. Add default priority "medium" for existing tasks
4. Add empty description for existing tasks
5. Maintain existing ID values for backward compatibility