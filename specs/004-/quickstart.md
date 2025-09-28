# Quickstart: Backend Domain-Driven Design Implementation

**Feature**: Backend Domain-Driven Design Implementation
**Date**: 2025-09-28
**Status**: Implementation Ready

## Prerequisites

- Go 1.23+ installed
- SQLite3 available
- Existing todo-app project cloned
- Backend and frontend directories present

## Setup Instructions

### 1. Backup Current Implementation
```bash
# Create backup branch
git checkout -b backup-pre-ddd
git add -A
git commit -m "Backup before DDD implementation"
git checkout 004-
```

### 2. Verify Current System Works
```bash
# Start current backend
cd backend
go mod tidy
go run cmd/server/main.go

# In another terminal, test current API
curl http://localhost:8080/api/health
curl http://localhost:8080/api/v1/tasks
```

### 3. Run Existing Tests
```bash
# Run current test suite to establish baseline
cd backend
go test ./tests/contract/...
go test ./tests/unit/...
```

## Implementation Validation Steps

### Phase 1: Domain Layer Validation
```bash
# After implementing domain models
cd backend
go test ./domain/task/... -v
go test ./domain/user/... -v

# Verify no external dependencies in domain layer
go list -deps ./domain/... | grep -v "todo-app/domain"
```

### Phase 2: Application Layer Validation
```bash
# After implementing application services
go test ./application/task/... -v
go test ./application/user/... -v

# Verify application layer only depends on domain
go list -deps ./application/... | grep -v "todo-app/\(domain\|application\)"
```

### Phase 3: Infrastructure Layer Validation
```bash
# After implementing repositories and external services
go test ./infrastructure/... -v

# Test database connectivity and migrations
go run cmd/migrate/main.go
```

### Phase 4: Presentation Layer Validation
```bash
# After implementing HTTP handlers
go test ./presentation/... -v

# Test API contract compliance
go run cmd/server/main.go &
SERVER_PID=$!

# Test all endpoints match OpenAPI spec
curl -X GET http://localhost:8080/api/v1/tasks
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Task", "priority": "medium"}'

kill $SERVER_PID
```

## End-to-End Validation Scenarios

### Scenario 1: Task Management Workflow
```bash
# Start the DDD-restructured server
cd backend
go run cmd/server/main.go &
SERVER_PID=$!

# 1. Create a new task
TASK_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Test DDD Task", "description": "Testing the new DDD structure", "priority": "high"}')

echo "Created task: $TASK_RESPONSE"
TASK_ID=$(echo $TASK_RESPONSE | jq -r '.id')

# 2. Retrieve the task
curl -s http://localhost:8080/api/v1/tasks/$TASK_ID | jq '.'

# 3. Update the task
curl -s -X PUT http://localhost:8080/api/v1/tasks/$TASK_ID \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}' | jq '.'

# 4. List all tasks
curl -s http://localhost:8080/api/v1/tasks | jq '.'

# 5. Delete the task
curl -s -X DELETE http://localhost:8080/api/v1/tasks/$TASK_ID

# 6. Verify deletion
curl -s http://localhost:8080/api/v1/tasks/$TASK_ID

kill $SERVER_PID
```

### Scenario 2: User Management Workflow
```bash
# Start server
cd backend
go run cmd/server/main.go &
SERVER_PID=$!

# 1. Register a new user
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "profile": {
      "first_name": "John",
      "last_name": "Doe",
      "timezone": "America/New_York"
    },
    "preferences": {
      "default_task_priority": "high",
      "email_notifications": true,
      "theme_preference": "dark"
    }
  }')

echo "Registered user: $USER_RESPONSE"

# 2. Get user profile
curl -s http://localhost:8080/api/v1/users/profile | jq '.'

# 3. Update user preferences
curl -s -X PUT http://localhost:8080/api/v1/users/preferences \
  -H "Content-Type: application/json" \
  -d '{"default_task_priority": "medium", "theme_preference": "light"}' | jq '.'

kill $SERVER_PID
```

## Performance Validation

### Response Time Benchmarks
```bash
# Measure API response times (should be similar to pre-DDD)
cd backend
go run cmd/server/main.go &
SERVER_PID=$!

# Warm up
for i in {1..10}; do
  curl -s http://localhost:8080/api/v1/tasks > /dev/null
done

# Benchmark task listing
time curl -s http://localhost:8080/api/v1/tasks > /dev/null

# Benchmark task creation
time curl -s -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Benchmark Task"}' > /dev/null

kill $SERVER_PID
```

### Memory Usage Validation
```bash
# Monitor memory usage during operation
cd backend
go run cmd/server/main.go &
SERVER_PID=$!

# Create multiple tasks to test memory usage
for i in {1..100}; do
  curl -s -X POST http://localhost:8080/api/v1/tasks \
    -H "Content-Type: application/json" \
    -d "{\"title\": \"Task $i\"}" > /dev/null
done

# Check memory usage
ps aux | grep "server/main"

kill $SERVER_PID
```

## Success Criteria

### ✅ Functional Requirements Validation
- [ ] All existing API endpoints work with same responses
- [ ] Task CRUD operations maintain backward compatibility
- [ ] Database queries produce same results
- [ ] Error responses match existing format

### ✅ DDD Architecture Validation
- [ ] Domain layer has no external dependencies
- [ ] Application layer depends only on domain
- [ ] Infrastructure implements domain interfaces
- [ ] Presentation layer orchestrates use cases

### ✅ Code Quality Validation
- [ ] All tests pass with new structure
- [ ] Code coverage maintained or improved
- [ ] No circular dependencies between layers
- [ ] Repository pattern properly implemented

### ✅ Performance Validation
- [ ] API response times within 10% of baseline
- [ ] Memory usage not significantly increased
- [ ] Database query performance maintained

## Troubleshooting

### Common Issues

**Import Cycle Error**
```bash
# Check for circular dependencies
go list -deps ./... | sort | uniq -c | sort -nr
```

**Interface Not Satisfied**
```bash
# Verify interface implementations
go build ./...
```

**Database Migration Issues**
```bash
# Reset database if needed
rm -f todo.db
go run cmd/migrate/main.go
```

### Rollback Procedure
If validation fails:
```bash
git checkout backup-pre-ddd
git checkout -b rollback-from-ddd
# Continue with original implementation
```

## Next Steps After Validation

1. Run complete test suite: `go test ./...`
2. Update documentation with new architecture
3. Train team on DDD patterns used
4. Monitor production metrics for performance impact
5. Plan future DDD enhancements (events, CQRS if needed)