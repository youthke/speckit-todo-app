# Tasks: Backend Domain-Driven Design Implementation

**Input**: Design documents from `/specs/004-/`
**Prerequisites**: plan.md (✓), research.md (✓), data-model.md (✓), contracts/ (✓)

## Execution Flow (main)
```
1. Load plan.md from feature directory ✓
   → Tech stack: Go 1.23+, Gin, GORM, testify
   → Structure: Web app (backend/), DDD 4-layer architecture
2. Load design documents: ✓
   → data-model.md: Task/User entities, Value Objects, Domain Services
   → contracts/: task-api.yaml, user-api.yaml (8 endpoints total)
   → research.md: Go DDD patterns, interface-based dependency injection
3. Generate tasks by category: ✓
   → Setup: DDD directory structure, dependencies
   → Tests: Contract tests for 8 endpoints, domain tests
   → Core: Domain entities, Application services, Infrastructure
   → Integration: Repository implementations, HTTP handlers
   → Polish: Migration scripts, validation, performance tests
4. Apply task rules: ✓
   → [P] = Different files/packages (parallel execution)
   → TDD: All tests before implementation
   → Dependencies: Domain → Application → Infrastructure → Presentation
5. Total tasks: 42 (T001-T042)
6. Dependencies mapped with blocking relationships
7. Parallel execution examples for [P] tasks
8. Validation: All contracts tested, entities modeled, endpoints implemented ✓
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- File paths relative to repository root

## Phase 3.1: Setup & Prerequisites
- [x] T001 Create DDD directory structure in backend/ per implementation plan
- [x] T002 Initialize Go modules for layered architecture with proper dependencies
- [x] T003 [P] Configure golangci-lint with DDD-specific rules and formatting
- [x] T004 [P] Create backup of existing implementation for rollback capability

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests for Task API
- [x] T005 [P] Contract test GET /api/v1/tasks in backend/tests/contract/tasks_get_test.go
- [x] T006 [P] Contract test POST /api/v1/tasks in backend/tests/contract/tasks_post_test.go
- [x] T007 [P] Contract test GET /api/v1/tasks/{id} in backend/tests/contract/tasks_get_by_id_test.go
- [x] T008 [P] Contract test PUT /api/v1/tasks/{id} in backend/tests/contract/tasks_put_update_test.go
- [x] T009 [P] Contract test DELETE /api/v1/tasks/{id} in backend/tests/contract/tasks_delete_new_test.go

### Contract Tests for User API
- [x] T010 [P] Contract test POST /api/v1/users/register in backend/tests/contract/users_register_test.go
- [x] T011 [P] Contract test GET /api/v1/users/profile in backend/tests/contract/users_profile_test.go
- [x] T012 [P] Contract test PUT /api/v1/users/profile in backend/tests/contract/users_profile_test.go
- [x] T013 [P] Contract test GET /api/v1/users/preferences in backend/tests/contract/users_preferences_test.go
- [x] T014 [P] Contract test PUT /api/v1/users/preferences in backend/tests/contract/users_preferences_test.go

### Domain Tests
- [x] T015 [P] Unit tests for Task entity behavior in backend/domain/task/entities/task_test.go
- [x] T016 [P] Unit tests for User entity behavior in backend/domain/user/entities/user_test.go
- [x] T017 [P] Unit tests for TaskTitle value object in backend/domain/task/valueobjects/task_title_test.go
- [x] T018 [P] Unit tests for Email value object in backend/domain/user/valueobjects/email_test.go

## Phase 3.3: Domain Layer (ONLY after tests are failing)

### Task Management Domain
- [x] T019 [P] TaskID value object in backend/domain/task/valueobjects/task_id.go
- [x] T020 [P] TaskTitle value object in backend/domain/task/valueobjects/task_title.go
- [x] T021 [P] TaskDescription value object in backend/domain/task/valueobjects/task_description.go
- [x] T022 [P] TaskStatus value object in backend/domain/task/valueobjects/task_status.go
- [x] T023 [P] TaskPriority value object in backend/domain/task/valueobjects/task_priority.go
- [x] T024 Task entity with behavior methods in backend/domain/task/entities/task.go
- [x] T025 TaskRepository interface in backend/domain/task/repositories/task_repository.go
- [x] T026 [P] TaskValidationService in backend/domain/task/services/task_validation_service.go
- [x] T027 [P] TaskSearchService in backend/domain/task/services/task_search_service.go

### User Management Domain
- [x] T028 [P] UserID value object in backend/domain/user/valueobjects/user_id.go
- [x] T029 [P] Email value object in backend/domain/user/valueobjects/email.go
- [x] T030 [P] UserProfile value object in backend/domain/user/valueobjects/user_profile.go
- [x] T031 [P] UserPreferences value object in backend/domain/user/valueobjects/user_preferences.go
- [x] T032 User entity with behavior methods in backend/domain/user/entities/user.go
- [x] T033 UserRepository interface in backend/domain/user/repositories/user_repository.go
- [x] T034 [P] UserAuthenticationService in backend/domain/user/services/user_authentication_service.go
- [x] T035 [P] UserProfileService in backend/domain/user/services/user_profile_service.go

## Phase 3.4: Application Layer

### Application Services
- [x] T036 Task application service in backend/application/task/task_application_service.go
- [x] T037 User application service in backend/application/user/user_application_service.go

## Phase 3.5: Infrastructure Layer

### Repository Implementations
- [x] T038 GORM TaskRepository implementation in backend/infrastructure/persistence/gorm_task_repository.go
- [x] T039 GORM UserRepository implementation in backend/infrastructure/persistence/gorm_user_repository.go
- [x] T040 Database migration scripts in backend/infrastructure/persistence/migrations/

## Phase 3.6: Presentation Layer

### HTTP Handlers
- [x] T041 Task HTTP handlers in backend/presentation/http/task_handlers.go
- [x] T042 User HTTP handlers in backend/presentation/http/user_handlers.go

## Dependencies
**Critical Path: Domain → Application → Infrastructure → Presentation**

### Blocking Dependencies
- T005-T018 (all tests) MUST complete before T019-T042 (any implementation)
- T019-T023 (Task value objects) block T024 (Task entity)
- T028-T031 (User value objects) block T032 (User entity)
- T024,T025 (Task entity/repo interface) block T036 (Task app service)
- T032,T033 (User entity/repo interface) block T037 (User app service)
- T025 (Task repo interface) blocks T038 (Task repo implementation)
- T033 (User repo interface) blocks T039 (User repo implementation)
- T036,T038 (Task services) block T041 (Task handlers)
- T037,T039 (User services) block T042 (User handlers)

### Sequential Chains
1. T019→T020→T021→T022→T023→T024→T025→T026,T027→T036→T038→T041
2. T028→T029→T030→T031→T032→T033→T034,T035→T037→T039→T042

## Parallel Execution Examples

### Phase 3.2: All Contract Tests (T005-T014)
```bash
# Launch all contract tests simultaneously:
Task: "Contract test GET /api/v1/tasks in backend/tests/contract/tasks_get_test.go"
Task: "Contract test POST /api/v1/tasks in backend/tests/contract/tasks_post_test.go"
Task: "Contract test GET /api/v1/tasks/{id} in backend/tests/contract/tasks_get_by_id_test.go"
Task: "Contract test PUT /api/v1/tasks/{id} in backend/tests/contract/tasks_put_test.go"
Task: "Contract test DELETE /api/v1/tasks/{id} in backend/tests/contract/tasks_delete_test.go"
Task: "Contract test POST /api/v1/users/register in backend/tests/contract/users_register_test.go"
Task: "Contract test GET /api/v1/users/profile in backend/tests/contract/users_profile_get_test.go"
Task: "Contract test PUT /api/v1/users/profile in backend/tests/contract/users_profile_put_test.go"
Task: "Contract test GET /api/v1/users/preferences in backend/tests/contract/users_preferences_get_test.go"
Task: "Contract test PUT /api/v1/users/preferences in backend/tests/contract/users_preferences_put_test.go"
```

### Phase 3.2: Domain Unit Tests (T015-T018)
```bash
# Launch domain tests simultaneously:
Task: "Unit tests for Task entity behavior in backend/domain/task/entities/task_test.go"
Task: "Unit tests for User entity behavior in backend/domain/user/entities/user_test.go"
Task: "Unit tests for TaskTitle value object in backend/domain/task/valueobjects/task_title_test.go"
Task: "Unit tests for Email value object in backend/domain/user/valueobjects/email_test.go"
```

### Phase 3.3: Task Value Objects (T019-T023)
```bash
# Launch Task value objects simultaneously:
Task: "TaskID value object in backend/domain/task/valueobjects/task_id.go"
Task: "TaskTitle value object in backend/domain/task/valueobjects/task_title.go"
Task: "TaskDescription value object in backend/domain/task/valueobjects/task_description.go"
Task: "TaskStatus value object in backend/domain/task/valueobjects/task_status.go"
Task: "TaskPriority value object in backend/domain/task/valueobjects/task_priority.go"
```

### Phase 3.3: User Value Objects (T028-T031)
```bash
# Launch User value objects simultaneously:
Task: "UserID value object in backend/domain/user/valueobjects/user_id.go"
Task: "Email value object in backend/domain/user/valueobjects/email.go"
Task: "UserProfile value object in backend/domain/user/valueobjects/user_profile.go"
Task: "UserPreferences value object in backend/domain/user/valueobjects/user_preferences.go"
```

### Phase 3.3: Domain Services (T026-T027, T034-T035)
```bash
# Launch domain services simultaneously (after entities are complete):
Task: "TaskValidationService in backend/domain/task/services/task_validation_service.go"
Task: "TaskSearchService in backend/domain/task/services/task_search_service.go"
Task: "UserAuthenticationService in backend/domain/user/services/user_authentication_service.go"
Task: "UserProfileService in backend/domain/user/services/user_profile_service.go"
```

## Validation Checklist
*GATE: Verified during main() execution*

- [✓] All contracts have corresponding tests (T005-T014 cover all 10 API endpoints)
- [✓] All entities have model tasks (Task: T024, User: T032)
- [✓] All tests come before implementation (T005-T018 before T019-T042)
- [✓] Parallel tasks truly independent (different packages/files)
- [✓] Each task specifies exact file path
- [✓] No task modifies same file as another [P] task
- [✓] DDD layer dependencies respected (Domain → App → Infra → Presentation)

## Notes
- **[P] tasks**: Different files/packages, no shared dependencies
- **TDD Critical**: ALL tests (T005-T018) must fail before ANY implementation
- **Layer Isolation**: Domain layer has zero external dependencies
- **Migration**: T040 handles data transformation from existing schema
- **Rollback**: T004 provides safety net for architectural changes
- **Performance**: Maintain existing API response times per plan requirements

## DDD Architecture Validation
Each task enforces specific DDD principles:
- **Domain Purity**: T019-T035 create dependency-free business logic
- **Dependency Inversion**: T025,T033 define interfaces, T038,T039 implement
- **Bounded Contexts**: Clear separation between Task and User domains
- **Aggregate Boundaries**: T024,T032 maintain consistency within entities
- **Repository Pattern**: T025→T038, T033→T039 follow pure DDD patterns