# Import Path Contracts

**Feature**: 009-resolve-it-1
**Purpose**: Define the contract for import paths after refactoring

## Import Path Standards

All Go imports must follow the Domain-Driven Design structure pattern:

```
todo-app/domain/[bounded-context]/[layer]/[package]
```

### Bounded Contexts
- `user` - User identity and profile
- `task` - Task management
- `auth` - Authentication and session management
- `health` - System health monitoring

### Layers
- `entities/` - Domain entities (aggregate roots)
- `valueobjects/` - Immutable value objects
- `repositories/` - Persistence interfaces
- `services/` - Domain services

## Allowed Import Patterns

### ✅ Correct Imports

```go
// User domain
import "todo-app/domain/user/entities"
import "todo-app/domain/user/valueobjects"
import "todo-app/domain/user/repositories"
import "todo-app/domain/user/services"

// Task domain
import "todo-app/domain/task/entities"
import "todo-app/domain/task/valueobjects"
import "todo-app/domain/task/repositories"
import "todo-app/domain/task/services"

// Auth domain
import "todo-app/domain/auth/entities"
import "todo-app/domain/auth/valueobjects"
import "todo-app/domain/auth/repositories"
import "todo-app/domain/auth/services"

// Health domain
import "todo-app/domain/health/entities"

// Infrastructure
import "todo-app/infrastructure/persistence"

// Application services
import "todo-app/application/user"
import "todo-app/application/task"

// Presentation
import "todo-app/presentation/http"

// Configuration
import "todo-app/internal/config"
```

### ❌ Forbidden Imports (MUST NOT EXIST after refactoring)

```go
// FORBIDDEN: Legacy flat models
import "todo-app/internal/models"

// FORBIDDEN: Orphaned models directory
import "todo-app/models"

// FORBIDDEN: Legacy services directory
import "todo-app/services/auth"

// FORBIDDEN: Any import with "legacy" in the name
import legacymodels "todo-app/internal/models"
```

## Type Reference Patterns

### Session Types

```go
// OLD (broken)
import "todo-app/internal/models"
session := models.AuthenticationSession{}

// NEW (correct)
import "todo-app/domain/auth/entities"
session := entities.AuthenticationSession{}
```

### OAuth Types

```go
// OLD (broken)
import "todo-app/internal/models"
state := models.OAuthState{}

// NEW (correct)
import "todo-app/domain/auth/entities"
state := entities.OAuthState{}
```

### User Types

```go
// OLD (deprecated flat model)
import "todo-app/internal/models"
user := models.User{}

// NEW (rich domain entity)
import "todo-app/domain/user/entities"
user := entities.User{}
```

### Task Types

```go
// OLD (deprecated flat model)
import "todo-app/internal/models"
task := models.Task{}

// NEW (rich domain entity)
import "todo-app/domain/task/entities"
task := entities.Task{}
```

## Verification Commands

### Check for Forbidden Imports

```bash
# Should return NO matches
grep -r "todo-app/internal/models" backend/ --include="*.go"
grep -r "todo-app/models" backend/ --include="*.go" | grep -v "domain"
grep -r "legacymodels" backend/ --include="*.go"
```

### Check for Correct DDD Imports

```bash
# Should match all import statements
grep -r "todo-app/domain/" backend/ --include="*.go"
```

### Compilation Check

```bash
# Must succeed with zero errors
cd backend && go build ./...
```

## Contract Test

A contract test will verify import path compliance:

```go
// Test: All imports follow DDD structure
func TestImportPathCompliance(t *testing.T) {
    forbiddenPatterns := []string{
        "todo-app/internal/models",
        "todo-app/models\"",  // Exclude domain/models
        "legacymodels",
    }

    files := getAllGoFiles("../../../backend")

    for _, file := range files {
        content := readFile(file)
        for _, pattern := range forbiddenPatterns {
            if strings.Contains(content, pattern) {
                t.Errorf("File %s contains forbidden import pattern: %s", file, pattern)
            }
        }
    }
}
```

## Migration Checklist

- [ ] All imports updated to DDD structure
- [ ] Zero references to `todo-app/internal/models`
- [ ] Zero references to `todo-app/models` (except domain subdirs)
- [ ] Zero references to `legacymodels` alias
- [ ] `go build ./...` succeeds
- [ ] Contract test passes
- [ ] All 51 existing tests pass

---

**Contract defined**. This document serves as the acceptance criteria for import path correctness.
