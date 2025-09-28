
# Implementation Plan: Backend Domain-Driven Design Implementation

**Branch**: `004-` | **Date**: 2025-09-28 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-/spec.md`

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
Implement Domain-Driven Design architecture in the backend through a complete restructure. This involves creating two bounded contexts (Task Management and User Management), implementing four architectural layers (Domain, Application, Infrastructure, Presentation), and applying core DDD patterns (Entities, Value Objects, Aggregates, Domain Services) with pure DDD repository patterns.

## Technical Context
**Language/Version**: Go 1.23+
**Primary Dependencies**: Gin web framework, GORM ORM, testify testing framework
**Storage**: SQLite database (development), existing schema to be migrated
**Testing**: Go testing package with testify, existing contract and unit tests
**Target Platform**: Linux/macOS server environment
**Project Type**: Web application (backend + frontend separation)
**Performance Goals**: Maintain existing API response times during migration
**Constraints**: Must preserve all existing functionality, complete architectural rewrite approach
**Scale/Scope**: Existing todo application with task management and user functionality

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Key Constitutional Requirements**:
- ✅ **Test-First Development**: All DDD patterns must be validated through tests before implementation
- ✅ **Clear Library Boundaries**: Domain, Application, Infrastructure layers must be independently testable
- ✅ **Simplicity Principle**: Core DDD patterns only (Entities, Value Objects, Aggregates, Domain Services) - no overengineering
- ✅ **Backward Compatibility**: All existing functionality preserved during architectural migration
- ⚠️ **Complexity Justification**: Big bang rewrite approach requires justification (documented in Complexity Tracking)

**Initial Gate Status**: CONDITIONAL PASS - pending complexity justification

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
├── domain/
│   ├── task/           # Task Management bounded context
│   │   ├── entities/
│   │   ├── valueobjects/
│   │   ├── aggregates/
│   │   ├── services/
│   │   └── repositories/
│   └── user/           # User Management bounded context
│       ├── entities/
│       ├── valueobjects/
│       ├── aggregates/
│       ├── services/
│       └── repositories/
├── application/
│   ├── task/           # Task use cases and application services
│   └── user/           # User use cases and application services
├── infrastructure/
│   ├── persistence/    # Repository implementations
│   ├── external/       # External service integrations
│   └── config/         # Configuration
├── presentation/
│   ├── http/           # REST API controllers
│   └── middleware/     # HTTP middleware
├── cmd/
│   └── server/
└── tests/
    ├── contract/
    ├── integration/
    └── unit/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/
```

**Structure Decision**: Web application with DDD layered architecture. Backend restructured into four distinct layers (Domain, Application, Infrastructure, Presentation) with two bounded contexts (Task Management, User Management). Frontend remains unchanged but may require updates to consume new DDD-structured APIs.

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
- Load `.specify/templates/tasks-template.md` as base
- Generate DDD-specific tasks from Phase 1 design docs
- Domain layer tasks: Entity, Value Object, Aggregate, Domain Service creation [P]
- Application layer tasks: Use case and Application Service implementation [P]
- Infrastructure layer tasks: Repository implementations, database adapters
- Presentation layer tasks: HTTP handlers, request/response models
- Contract validation tasks: API compliance tests for each endpoint
- Migration tasks: Data transformation from current to DDD structure

**DDD-Specific Ordering Strategy**:
1. **Domain Layer First**: Entities, Value Objects, Aggregates (pure business logic)
2. **Application Layer**: Use cases that orchestrate domain operations
3. **Infrastructure Layer**: Repository implementations, database migrations
4. **Presentation Layer**: HTTP handlers using application services
5. **Integration Tests**: End-to-end validation of DDD architecture
6. **Migration Validation**: Ensure backward compatibility with existing data

**Task Categories**:
- **[P] Domain Tasks**: Can be implemented in parallel (no dependencies)
- **[S] Sequential Tasks**: Must follow dependency order (Application → Infrastructure → Presentation)
- **[M] Migration Tasks**: Data and schema transformation tasks
- **[T] Test Tasks**: Contract tests, integration tests, validation scenarios

**Estimated Output**: 35-40 numbered, dependency-ordered tasks covering complete DDD transformation

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
| Big bang rewrite approach | Clear architectural boundaries and proper DDD implementation requires complete restructure | Gradual refactoring would create hybrid architecture with unclear boundaries, making it harder to maintain DDD principles |
| Four-layer architecture | Domain-driven design requires clear separation of concerns between business logic and infrastructure | Simpler layering would mix domain logic with infrastructure concerns, violating core DDD principles |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command) - ✅ research.md created
- [x] Phase 1: Design complete (/plan command) - ✅ data-model.md, contracts/, quickstart.md, CLAUDE.md updated
- [x] Phase 2: Task planning complete (/plan command - describe approach only) - ✅ DDD-specific task strategy defined
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: CONDITIONAL PASS - complexity justified
- [x] Post-Design Constitution Check: PASS - DDD architecture aligns with constitutional principles
- [x] All NEEDS CLARIFICATION resolved - clarified via /clarify session
- [x] Complexity deviations documented - big bang rewrite and four-layer architecture justified

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
