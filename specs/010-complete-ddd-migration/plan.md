
# Implementation Plan: Complete DDD Migration

**Branch**: `010-complete-ddd-migration` | **Date**: 2025-10-04 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/010-complete-ddd-migration/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Complete the Domain-Driven Design (DDD) migration by resolving architectural incompatibility between legacy GORM DTOs (`internal/models`) and rich DDD entities (`domain/`). Implement a Mapper/Adapter pattern to bridge DTOs (for API/persistence) and DDD entities (for business logic), enabling full migration of User and Task models while maintaining API and database compatibility.

## Technical Context
**Language/Version**: Go 1.24.7
**Primary Dependencies**: Gin web framework v1.11.0, GORM ORM v1.31.0, testify v1.11.1, golang.org/x/oauth2 v0.31.0
**Storage**: SQLite (GORM driver v1.6.0, development database)
**Testing**: Go testing framework, testify assertions, 51+ existing tests (8 unit, 19 contract, 13 integration, 11 domain)
**Target Platform**: Backend REST API server (Linux/macOS)
**Project Type**: web (backend + frontend)
**Performance Goals**: Mapper overhead < 1ms per conversion, API response time < 50ms, build time < 10 seconds
**Constraints**: Zero API breaking changes, zero database schema changes, 100% test pass rate, backward compatibility maintained
**Scale/Scope**: 4 domains (Auth, Health, User, Task), 30 files with legacy imports (11 source + 19 tests), 47 migration tasks

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Status**: PASS (constitution file is template-only, no specific gates defined)

No constitutional violations detected for this refactoring feature:
- ✅ Internal refactoring only (no new external interfaces)
- ✅ Maintains existing API contracts
- ✅ No new dependencies or libraries
- ✅ Uses existing test framework (testify)
- ✅ Follows established DDD pattern from feature 009

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
backend/
├── domain/                      # DDD entities (existing)
│   ├── auth/entities/
│   ├── health/entities/
│   ├── user/
│   │   ├── entities/user.go
│   │   ├── valueobjects/email.go
│   │   └── repositories/user_repository.go
│   └── task/
│       ├── entities/task.go
│       ├── valueobjects/task_status.go
│       └── repositories/task_repository.go
├── application/                 # NEW - Mapper layer
│   └── mappers/
│       ├── user_mapper.go
│       ├── user_mapper_test.go
│       ├── task_mapper.go
│       └── task_mapper_test.go
├── internal/
│   ├── models/                  # RENAME TO dtos/
│   │   ├── user.go             # Becomes user_dto.go
│   │   └── task.go             # Becomes task_dto.go
│   ├── handlers/               # UPDATE - Use mappers
│   └── services/               # UPDATE - Use DDD entities
├── infrastructure/
│   └── persistence/            # UPDATE - Inject mappers
│       ├── gorm_user_repository.go
│       └── gorm_task_repository.go
└── tests/
    ├── unit/                   # 8 tests
    ├── contract/               # 19 tests (subset need updates)
    ├── integration/            # 13 tests
    └── domain/                 # 11 tests

frontend/
└── [No changes - frontend unaffected]
```

**Structure Decision**: Web application structure (backend + frontend). This feature focuses exclusively on backend refactoring. The existing DDD structure from feature 009 will be enhanced with a new mapper layer (`application/mappers/`) to bridge DTOs and entities. Legacy `internal/models/` will be renamed to `internal/dtos/` and modified files will use the mapper pattern.

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh claude`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:

The /tasks command will generate tasks following this structure:

1. **Infrastructure Setup** (5 tasks)
   - Create `application/mappers/` directory structure
   - Rename `internal/models/` to `internal/dtos/`
   - Rename model files (user.go → user_dto.go, task.go → task_dto.go)
   - Update package declarations in DTO files
   - Create mapper file skeletons (user_mapper.go, task_mapper.go)

2. **Mapper Implementation** (6 tasks) [TDD]
   - Write UserMapper tests (ToEntity/ToDTO, error cases)
   - Implement UserMapper (DTO ↔ Entity conversion)
   - Write TaskMapper tests (ToEntity/ToDTO, error cases)
   - Implement TaskMapper (DTO ↔ Entity conversion)
   - Benchmark mapper performance (< 1ms target)
   - Verify 100% mapper test coverage

3. **Repository Updates** (8 tasks)
   - Update GormUserRepository constructor (inject mapper)
   - Update GormUserRepository methods (use mapper for conversions)
   - Add repository integration tests for User
   - Update GormTaskRepository constructor (inject mapper)
   - Update GormTaskRepository methods (use mapper for conversions)
   - Add repository integration tests for Task
   - Update repository factory/initialization code
   - Verify all repository tests pass

4. **Service Layer Updates** (6 tasks)
   - Update UserService to use User entities internally
   - Update UserService tests (mock mappers if needed)
   - Update TaskService to use Task entities internally
   - Update TaskService tests (mock mappers if needed)
   - Update service initialization/dependency injection
   - Verify all service tests pass

5. **Handler Layer Updates** (8 tasks)
   - Update UserHandler (inject mapper, use at boundaries)
   - Update TaskHandler (inject mapper, use at boundaries)
   - Update handler initialization in main.go/routes
   - Update handler tests (verify JSON contracts)
   - Verify user endpoint contract tests pass
   - Verify task endpoint contract tests pass
   - Run manual API smoke tests
   - Capture before/after response baselines

6. **Test Updates** (12 tasks)
   - Update unit tests imports (models → dtos)
   - Update integration tests imports
   - Update contract tests imports
   - Fix any test compilation errors
   - Run full test suite (unit tests)
   - Run full test suite (integration tests)
   - Run full test suite (contract tests)
   - Run full test suite (domain tests)
   - Verify 100% test pass rate
   - Check test coverage metrics
   - Benchmark test suite runtime
   - Document any test changes needed

7. **Cleanup & Verification** (8 tasks)
   - Search for remaining `internal/models` imports
   - Update any missed import references
   - Remove deprecated code/comments
   - Update go.mod if needed
   - Run `go build ./...` (verify compilation)
   - Run `go vet ./...` (static analysis)
   - Run `go test ./...` (full suite)
   - Execute quickstart.md validation steps

8. **Documentation** (3 tasks)
   - Add mapper pattern documentation to CLAUDE.md
   - Create migration guide for future refactoring
   - Update architecture diagrams if they exist

**Ordering Strategy**:
- **Sequential phases**: 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8
- **Parallel within phase**: Tasks marked [P] can run concurrently
- **TDD enforcement**: Test tasks before implementation tasks
- **Incremental verification**: Test after each phase

**Dependencies**:
- Phase 3 depends on Phase 2 (mappers must exist before repository updates)
- Phase 4 depends on Phase 3 (services need updated repositories)
- Phase 5 depends on Phase 4 (handlers need updated services)
- Phase 6 can partially overlap with Phases 3-5 (update tests as code changes)
- Phase 7 depends on all previous phases

**Estimated Task Count**: 56 tasks total (aligned with spec estimate of 47+ from feature 009)

**Parallel Execution Opportunities**:
- Phase 2: UserMapper and TaskMapper can be developed in parallel [P]
- Phase 3: User and Task repositories can be updated in parallel [P]
- Phase 4: User and Task services can be updated in parallel [P]
- Phase 6: Different test categories can be updated in parallel [P]

**Test-First Approach**:
- Each implementation task (mapper, repository, service, handler) preceded by test task
- Tests written first → verify they fail → implement → verify they pass
- Maintains TDD discipline throughout migration

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command) ✅ 2025-10-04
- [x] Phase 1: Design complete (/plan command) ✅ 2025-10-04
- [x] Phase 2: Task planning complete (/plan command - describe approach only) ✅ 2025-10-04
- [ ] Phase 3: Tasks generated (/tasks command) - Pending
- [ ] Phase 4: Implementation complete - Pending
- [ ] Phase 5: Validation passed - Pending

**Gate Status**:
- [x] Initial Constitution Check: PASS ✅
- [x] Post-Design Constitution Check: PASS ✅
- [x] All NEEDS CLARIFICATION resolved: N/A (Technical Context fully specified) ✅
- [x] Complexity deviations documented: None (no constitutional violations) ✅

**Artifacts Generated**:
- [x] research.md - 9 research areas, all decisions documented
- [x] data-model.md - 2 DTOs, 2 Entities, 2 Mappers specified
- [x] contracts/api-contracts.md - 8 API endpoints documented
- [x] quickstart.md - 11 verification tests defined
- [x] CLAUDE.md - Updated with feature 010 context
- [x] plan.md - This document (complete)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
