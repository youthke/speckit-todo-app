# Data Model: TODO App

## Core Entities

### Task
**Purpose**: Represents a single TODO item that users can create, modify, and complete.

**Attributes**:
- `id` (uint): Unique identifier, auto-generated primary key
- `title` (string): The task description/title, required, max 500 characters
- `completed` (bool): Task completion status, defaults to false
- `created_at` (timestamp): When the task was created, auto-generated
- `updated_at` (timestamp): When the task was last modified, auto-updated

**Validation Rules**:
- Title must not be empty (FR-007)
- Title length must be ≤ 500 characters
- completed field must be boolean
- id must be unique and positive

**State Transitions**:
```
[New] --create--> [Pending] --complete--> [Completed]
                     ↑           ↓
                     ←--uncomplete--
```

**Database Schema** (SQLite/PostgreSQL):
```sql
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(500) NOT NULL CHECK(LENGTH(title) > 0),
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for common queries
CREATE INDEX idx_tasks_completed ON tasks(completed);
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);
```

**Go Struct**:
```go
type Task struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title" gorm:"type:varchar(500);not null" validate:"required,max=500"`
    Completed bool      `json:"completed" gorm:"default:false"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

**TypeScript Interface**:
```typescript
interface Task {
  id: number;
  title: string;
  completed: boolean;
  created_at: string; // ISO 8601 format
  updated_at: string; // ISO 8601 format
}
```

## Data Access Patterns

### Common Queries
1. **Get all tasks**: `SELECT * FROM tasks ORDER BY created_at DESC`
2. **Get pending tasks**: `SELECT * FROM tasks WHERE completed = false ORDER BY created_at DESC`
3. **Get completed tasks**: `SELECT * FROM tasks WHERE completed = true ORDER BY updated_at DESC`
4. **Get task by ID**: `SELECT * FROM tasks WHERE id = ?`

### Performance Considerations
- Tasks table expected to grow to ~50-100 entries for typical personal use
- No complex joins required
- Simple indexes on completed and created_at sufficient
- No pagination needed for expected data volume

## Data Flow

### Task Creation
```
User Input → Validation → Database Insert → Return Task with ID
```

### Task Update
```
Task ID + Updates → Validation → Database Update → Return Updated Task
```

### Task Completion Toggle
```
Task ID → Find Task → Toggle completed → Update timestamp → Save
```

### Task Deletion
```
Task ID → Find Task → Delete from Database → Confirm Deletion
```

## Error Handling

### Validation Errors
- Empty title: "Task title cannot be empty"
- Title too long: "Task title must be 500 characters or less"
- Invalid ID: "Task not found"

### Database Errors
- Constraint violations: Return 400 Bad Request
- Not found: Return 404 Not Found
- Server errors: Return 500 Internal Server Error

## Migration Strategy

### Initial Schema
- Create tasks table with all required fields
- Add indexes for performance
- Seed with sample data for development

### Future Considerations
- Add migration system for schema changes
- Consider soft deletion (deleted_at field) instead of hard deletion
- Potential fields: priority, due_date, category (based on clarifications)