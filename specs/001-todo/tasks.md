# Tasks: TODO App

**Input**: Design documents from `/specs/001-todo/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Extract: Go 1.21+, Node.js 18+, Gin backend, React frontend
   → Structure: backend/ and frontend/ directories
2. Load design documents:
   → data-model.md: Task entity with validation
   → contracts/api.yaml: 4 REST endpoints (GET, POST, PUT, DELETE)
   → quickstart.md: 5 core user scenarios + edge cases
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, handlers, React components
   → Integration: DB, CORS, React-API connection
   → Polish: unit tests, validation, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Tests before implementation (TDD)
   → Dependencies block execution
5. Validate all contracts/entities/scenarios covered
6. SUCCESS: 24 tasks ready for execution
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- All file paths are absolute from repository root

## Phase 3.1: Setup

- [x] T001 Create backend project structure: `backend/cmd/server/`, `backend/internal/{models,handlers,services,storage}/`, `backend/pkg/api/`, `backend/tests/{contract,integration,unit}/`
- [x] T002 Initialize Go module in `backend/go.mod` with dependencies: gin-gonic/gin, gorm.io/gorm, gorm.io/driver/sqlite, stretchr/testify
- [x] T003 [P] Create React app in `frontend/` with dependencies: axios, @testing-library/react, @testing-library/jest-dom
- [x] T004 [P] Configure Go linting with golangci-lint in `backend/.golangci.yml`
- [x] T005 [P] Configure ESLint and Prettier in `frontend/.eslintrc.js` and `frontend/.prettierrc`

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

- [x] T006 [P] Contract test GET /api/v1/tasks in `backend/tests/contract/tasks_get_test.go`
- [x] T007 [P] Contract test POST /api/v1/tasks in `backend/tests/contract/tasks_post_test.go`
- [x] T008 [P] Contract test PUT /api/v1/tasks/{id} in `backend/tests/contract/tasks_put_test.go`
- [x] T009 [P] Contract test DELETE /api/v1/tasks/{id} in `backend/tests/contract/tasks_delete_test.go`
- [x] T010 [P] Integration test task creation scenario in `backend/tests/integration/task_creation_test.go`
- [x] T011 [P] Integration test task completion scenario in `backend/tests/integration/task_completion_test.go`
- [x] T012 [P] Integration test task editing scenario in `backend/tests/integration/task_editing_test.go`
- [x] T013 [P] Integration test task deletion scenario in `backend/tests/integration/task_deletion_test.go`
- [x] T014 [P] React component test for TaskList in `frontend/src/components/TaskList/TaskList.test.js`
- [x] T015 [P] React component test for TaskItem in `frontend/src/components/TaskItem/TaskItem.test.js`
- [x] T016 [P] React component test for TaskForm in `frontend/src/components/TaskForm/TaskForm.test.js`

## Phase 3.3: Core Implementation (ONLY after tests are failing)

- [x] T017 Task model with validation in `backend/internal/models/task.go`
- [x] T018 Database setup and migrations in `backend/internal/storage/database.go`
- [x] T019 Task service with CRUD operations in `backend/internal/services/task_service.go`
- [x] T020 GET /api/v1/tasks handler in `backend/internal/handlers/task_handlers.go`
- [x] T021 POST /api/v1/tasks handler in same file as T020
- [x] T022 PUT /api/v1/tasks/{id} handler in same file as T020
- [x] T023 DELETE /api/v1/tasks/{id} handler in same file as T020
- [x] T024 Main server setup with routes in `backend/cmd/server/main.go`

## Phase 3.4: Frontend Implementation

- [x] T025 [P] TaskList component in `frontend/src/components/TaskList/TaskList.js`
- [x] T026 [P] TaskItem component in `frontend/src/components/TaskItem/TaskItem.js`
- [x] T027 [P] TaskForm component in `frontend/src/components/TaskForm/TaskForm.js`
- [x] T028 API service for backend communication in `frontend/src/services/api.js`
- [x] T029 Main App component integrating all components in `frontend/src/App.js`
- [x] T030 CSS styling for task components in `frontend/src/App.css`

## Phase 3.5: Integration & Polish

- [x] T031 CORS configuration for frontend-backend communication in `backend/cmd/server/main.go`
- [x] T032 Error handling and validation middleware in `backend/internal/handlers/middleware.go`
- [x] T033 [P] Unit tests for task validation in `backend/tests/unit/task_validation_test.go`
- [x] T034 [P] Edge case tests (empty title, long title) in `backend/tests/unit/edge_cases_test.go`
- [x] T035 Performance validation: API response times <500ms
- [x] T036 Manual testing following quickstart.md scenarios
- [x] T037 [P] Update README.md with setup and development instructions

## Dependencies

### Critical Dependencies (TDD)
- **Tests (T006-T016) MUST complete and FAIL before implementation (T017-T024)**
- **T017 (model) blocks T018, T019**
- **T018 (database) blocks T019, T020-T023**
- **T019 (service) blocks T020-T023**

### Implementation Dependencies
- **T020-T023 are sequential (same file)**
- **T024 requires T020-T023 (routes need handlers)**
- **T028 blocks T029 (App needs API service)**
- **T031-T032 can run after T024**

### Frontend Dependencies
- **T025-T027 can run in parallel (different files)**
- **T029 requires T025-T027 (integration component)**

## Parallel Execution Examples

### Phase 3.1 Setup (T003, T004, T005)
```bash
# Launch these together:
Task: "Create React app in frontend/ with dependencies"
Task: "Configure Go linting with golangci-lint"
Task: "Configure ESLint and Prettier in frontend/"
```

### Phase 3.2 Contract Tests (T006-T009)
```bash
# Launch API contract tests together:
Task: "Contract test GET /api/v1/tasks in backend/tests/contract/tasks_get_test.go"
Task: "Contract test POST /api/v1/tasks in backend/tests/contract/tasks_post_test.go"
Task: "Contract test PUT /api/v1/tasks/{id} in backend/tests/contract/tasks_put_test.go"
Task: "Contract test DELETE /api/v1/tasks/{id} in backend/tests/contract/tasks_delete_test.go"
```

### Phase 3.2 Integration Tests (T010-T013)
```bash
# Launch integration tests together:
Task: "Integration test task creation scenario in backend/tests/integration/task_creation_test.go"
Task: "Integration test task completion scenario in backend/tests/integration/task_completion_test.go"
Task: "Integration test task editing scenario in backend/tests/integration/task_editing_test.go"
Task: "Integration test task deletion scenario in backend/tests/integration/task_deletion_test.go"
```

### Phase 3.2 React Component Tests (T014-T016)
```bash
# Launch React tests together:
Task: "React component test for TaskList in frontend/src/components/TaskList/TaskList.test.js"
Task: "React component test for TaskItem in frontend/src/components/TaskItem/TaskItem.test.js"
Task: "React component test for TaskForm in frontend/src/components/TaskForm/TaskForm.test.js"
```

### Phase 3.4 React Components (T025-T027)
```bash
# Launch component implementation together:
Task: "TaskList component in frontend/src/components/TaskList/TaskList.js"
Task: "TaskItem component in frontend/src/components/TaskItem/TaskItem.js"
Task: "TaskForm component in frontend/src/components/TaskForm/TaskForm.js"
```

## Task Generation Rules Applied

### From Contracts (api.yaml)
✅ **GET /tasks** → T006 (contract test) + T020 (implementation)
✅ **POST /tasks** → T007 (contract test) + T021 (implementation)
✅ **PUT /tasks/{id}** → T008 (contract test) + T022 (implementation)
✅ **DELETE /tasks/{id}** → T009 (contract test) + T023 (implementation)

### From Data Model
✅ **Task entity** → T017 (model) + T018 (database) + T019 (service)

### From Quickstart Scenarios
✅ **Creating tasks** → T010 (integration test) + T025-T027 (components)
✅ **Marking complete** → T011 (integration test)
✅ **Editing tasks** → T012 (integration test)
✅ **Deleting tasks** → T013 (integration test)
✅ **Edge cases** → T034 (unit tests)

## Validation Checklist ✅

- [x] All contracts have corresponding tests (T006-T009)
- [x] Task entity has model task (T017)
- [x] All tests come before implementation (T006-T016 before T017+)
- [x] Parallel tasks are truly independent (different files)
- [x] Each task specifies exact file path
- [x] No [P] task modifies same file as another [P] task
- [x] All quickstart scenarios covered (T010-T013 + frontend tests)
- [x] Setup → Tests → Implementation → Polish ordering maintained

## Notes

- **TDD Enforcement**: Tests T006-T016 must be written first and MUST FAIL
- **File Conflicts**: T020-T023 are sequential (same handler file)
- **Dependencies**: Database (T018) required before handlers (T020-T023)
- **Performance**: T035 validates <500ms requirement from technical context
- **Edge Cases**: T034 covers empty title and 500+ character validation
- **Manual Validation**: T036 executes full quickstart.md scenarios