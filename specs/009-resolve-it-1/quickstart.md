# Quickstart: Import Path Cleanup Verification

**Feature**: 009-resolve-it-1
**Purpose**: Validate that import path refactoring completed successfully

## Prerequisites

- Go 1.24+ installed
- Repository cloned at `/Users/youthke/practice/speckit/todo-app`
- Branch `009-resolve-it-1` checked out

## Quick Verification (5 minutes)

### Step 1: Check Branch
```bash
cd /Users/youthke/practice/speckit/todo-app
git branch --show-current
# Expected output: 009-resolve-it-1
```

### Step 2: Verify Build Success
```bash
cd backend
go build ./...
```

**Success Criteria**: ✅ Zero compilation errors

**Expected Output**:
```
# Should complete silently with exit code 0
echo $?
# Output: 0
```

**Failure Signs**:
- ❌ `undefined: models.AuthenticationSession`
- ❌ `undefined: models.OAuthState`
- ❌ `undefined: legacymodels.AuthenticationSession`
- ❌ Any package not found errors

### Step 3: Verify No Forbidden Imports
```bash
# Check for legacy imports (should find ZERO matches)
grep -r "todo-app/internal/models" backend/ --include="*.go" | grep -v "^Binary"

# Check for orphaned models directory (should find ZERO matches)
grep -r '"todo-app/models"' backend/ --include="*.go"

# Check for legacymodels alias (should find ZERO matches)
grep -r "legacymodels" backend/ --include="*.go"
```

**Success Criteria**: All three commands return no output

### Step 4: Run Test Suite
```bash
cd backend
go test ./... -v
```

**Success Criteria**: ✅ All tests pass (PASS), zero FAIL

**Expected Output Sample**:
```
=== RUN   TestUserEntity_Creation
--- PASS: TestUserEntity_Creation (0.00s)
=== RUN   TestTaskEntity_Validation
--- PASS: TestTaskEntity_Validation (0.00s)
...
PASS
ok      todo-app/domain/user/entities    0.123s
ok      todo-app/domain/task/entities    0.089s
ok      todo-app/domain/auth/entities    0.145s
```

### Step 5: Verify DDD Structure
```bash
# List domain structure
ls -la backend/domain/
```

**Expected Output**:
```
drwxr-xr-x  auth/       # NEW - Authentication domain
drwxr-xr-x  health/     # NEW - Health check domain
drwxr-xr-x  task/       # EXISTING - Task domain
drwxr-xr-x  user/       # EXISTING - User domain
```

```bash
# Verify auth domain structure
ls -la backend/domain/auth/
```

**Expected Output**:
```
drwxr-xr-x  entities/
drwxr-xr-x  valueobjects/
drwxr-xr-x  repositories/
drwxr-xr-x  services/
```

### Step 6: Verify Legacy Cleanup
```bash
# These directories should NOT exist or be empty
ls backend/models/ 2>&1
# Expected: "No such file or directory"

ls backend/internal/models/ 2>&1
# Expected: "No such file or directory" OR only non-domain helper files
```

## Detailed Verification (15 minutes)

### Test Case 1: Session Management
```bash
cd backend
go test ./domain/auth/entities -v -run TestAuthenticationSession
```

**Success Criteria**: Tests for AuthenticationSession entity pass

### Test Case 2: OAuth Flow
```bash
go test ./domain/auth/entities -v -run TestOAuthState
```

**Success Criteria**: Tests for OAuthState entity pass

### Test Case 3: Integration Tests
```bash
go test ./tests/integration/session_management_test.go -v
go test ./tests/integration/oauth_* -v
```

**Success Criteria**: Session and OAuth integration tests pass

### Test Case 4: Contract Tests
```bash
go test ./tests/contract/ -v
```

**Success Criteria**: All 23 contract tests pass (API JSON contracts unchanged)

### Test Case 5: Import Path Compliance
Create and run this verification script:

```bash
cat > /tmp/check_imports.sh << 'EOF'
#!/bin/bash
echo "Checking for forbidden import patterns..."

FORBIDDEN_FOUND=0

# Check for internal/models imports
if grep -r "todo-app/internal/models" backend/ --include="*.go" | grep -v "^Binary" > /dev/null; then
    echo "❌ FAIL: Found forbidden import 'todo-app/internal/models'"
    grep -r "todo-app/internal/models" backend/ --include="*.go" | head -5
    FORBIDDEN_FOUND=1
fi

# Check for models directory imports
if grep -r '"todo-app/models"' backend/ --include="*.go" > /dev/null; then
    echo "❌ FAIL: Found forbidden import 'todo-app/models'"
    grep -r '"todo-app/models"' backend/ --include="*.go" | head -5
    FORBIDDEN_FOUND=1
fi

# Check for legacymodels alias
if grep -r "legacymodels" backend/ --include="*.go" > /dev/null; then
    echo "❌ FAIL: Found forbidden 'legacymodels' alias"
    grep -r "legacymodels" backend/ --include="*.go" | head -5
    FORBIDDEN_FOUND=1
fi

if [ $FORBIDDEN_FOUND -eq 0 ]; then
    echo "✅ PASS: No forbidden import patterns found"
fi

# Verify DDD imports exist
DDD_COUNT=$(grep -r "todo-app/domain/" backend/ --include="*.go" | wc -l)
echo "✅ INFO: Found $DDD_COUNT DDD-style imports"

if [ $DDD_COUNT -lt 50 ]; then
    echo "⚠️  WARNING: Expected at least 50 DDD imports, found $DDD_COUNT"
    echo "   This might indicate incomplete migration"
fi

exit $FORBIDDEN_FOUND
EOF

chmod +x /tmp/check_imports.sh
/tmp/check_imports.sh
```

**Success Criteria**: Script exits with code 0, reports "PASS"

## Acceptance Scenarios (from spec.md)

### Scenario 1: Codebase Scan
**Given**: The backend codebase with unknown import issues
**When**: I perform a codebase scan
**Then**: All files with undefined import references are identified and reported

**Verification**:
```bash
go build ./... 2>&1 | grep "undefined:"
# Expected: No output (no undefined references)
```

### Scenario 2: Compilation Success
**Given**: The backend codebase with broken import paths
**When**: I run the build command after fixes
**Then**: The compilation completes successfully without import-related errors

**Verification**:
```bash
cd backend && go build ./... && echo "✅ Compilation succeeded"
```

### Scenario 3: Test Execution
**Given**: All import paths have been corrected
**When**: I run the test suite
**Then**: All tests can execute without import resolution failures

**Verification**:
```bash
cd backend && go test ./... -count=1 && echo "✅ All tests passed"
```

### Scenario 4: CI/CD Build
**Given**: The codebase is ready for deployment
**When**: The CI/CD pipeline runs
**Then**: The build step succeeds and produces deployable artifacts

**Verification**:
```bash
cd backend && go build -o /tmp/todo-api cmd/server/main.go && \
test -x /tmp/todo-api && \
echo "✅ Deployable binary created"
```

## Troubleshooting

### Issue: "undefined: models.AuthenticationSession"
**Cause**: Import still references old location
**Fix**: Update import to `"todo-app/domain/auth/entities"`

### Issue: "package todo-app/models is not in GOROOT"
**Cause**: Import references deleted directory
**Fix**: Update import to appropriate domain path

### Issue: Tests fail with "no such table: authentication_sessions"
**Cause**: Database migration needed (should not happen - tables preserved)
**Fix**: Check GORM `TableName()` methods are correct

### Issue: "too many arguments in call to function"
**Cause**: Function signature changed during refactoring
**Fix**: Review data-model.md for correct entity constructor signatures

## Success Checklist

Copy this checklist to verify completion:

```
Compilation
- [ ] go build ./... succeeds with zero errors
- [ ] go vet ./... reports zero warnings
- [ ] go mod tidy completes without changes

Import Paths
- [ ] Zero matches for grep "todo-app/internal/models"
- [ ] Zero matches for grep "todo-app/models\"" (orphaned)
- [ ] Zero matches for grep "legacymodels"
- [ ] At least 50 matches for grep "todo-app/domain/"

Directory Structure
- [ ] backend/domain/auth/ exists with 4 subdirectories
- [ ] backend/domain/health/ exists
- [ ] backend/models/ does not exist
- [ ] backend/internal/models/ does not exist (or empty)

Testing
- [ ] Unit tests pass: go test ./domain/... -v
- [ ] Integration tests pass: go test ./tests/integration/... -v
- [ ] Contract tests pass: go test ./tests/contract/... -v
- [ ] All 51+ tests pass: go test ./... -v

Functionality
- [ ] Application starts: go run cmd/server/main.go
- [ ] Health endpoint responds: curl http://localhost:8080/health
- [ ] Login flow works (manual test via frontend)
- [ ] Task CRUD works (manual test via frontend)

Documentation
- [ ] specs/009-resolve-it-1/research.md complete
- [ ] specs/009-resolve-it-1/data-model.md complete
- [ ] specs/009-resolve-it-1/contracts/import-paths.md complete
- [ ] CLAUDE.md updated with new DDD structure
```

## Rollback Procedure (if needed)

If verification fails and rollback is needed:

```bash
# Stash any uncommitted changes
git stash

# Return to main branch
git checkout main

# Delete feature branch (if desired)
git branch -D 009-resolve-it-1

# Rebuild to confirm stable state
cd backend && go build ./...
```

## Next Steps After Success

1. Create pull request from `009-resolve-it-1` to `main`
2. Ensure CI/CD pipeline passes
3. Request code review focusing on import path changes
4. Merge to main after approval
5. Deploy to staging environment
6. Verify production deployment checklist
7. Document lesson learned in project wiki

---

**Quickstart complete**. Use this guide to validate the feature implementation.
