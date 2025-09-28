# Tasks: Frontend TypeScript Migration

**Input**: Design documents from `/specs/003-frontend-typescript/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Web app**: `frontend/src/` for TypeScript migration
- All paths relative to repository root

## Phase 3.1: Setup
- [x] T001 Install TypeScript dependencies in frontend/package.json: typescript, @types/node, @types/react, @types/react-dom, @types/jest
- [x] T002 Create TypeScript configuration file frontend/tsconfig.json with React and Create React App settings
- [x] T003 [P] Create TypeScript type contracts file frontend/src/types/index.ts from contracts/typescript-types.contract.ts

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T004 [P] Create TypeScript contract test for Task interface in frontend/src/types/__tests__/task.test.ts
- [x] T005 [P] Create TypeScript contract test for API functions in frontend/src/types/__tests__/api.test.ts
- [x] T006 [P] Create TypeScript contract test for component props in frontend/src/types/__tests__/components.test.ts
- [x] T007 [P] Convert existing TaskForm test to TypeScript: frontend/src/components/TaskForm/TaskForm.test.js → TaskForm.test.tsx

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [x] T008 [P] Create main types file frontend/src/types/index.ts with all interfaces from data model
- [x] T009 Convert API service to TypeScript: frontend/src/services/api.js → api.ts with proper type annotations
- [x] T010 Convert reportWebVitals utility: frontend/src/reportWebVitals.js → reportWebVitals.ts
- [x] T011 Convert main index file: frontend/src/index.js → index.tsx with proper React types
- [x] T012 Convert App component: frontend/src/App.js → App.tsx with TypeScript interfaces
- [x] T013 Convert TaskForm component: frontend/src/components/TaskForm/TaskForm.js → TaskForm.tsx with props interface
- [x] T014 Convert TaskItem component: frontend/src/components/TaskItem/TaskItem.js → TaskItem.tsx with props interface
- [x] T015 Convert remaining component files to TypeScript with proper type annotations

## Phase 3.4: Integration
- [x] T016 Update all import statements to use TypeScript file extensions where needed
- [x] T017 Fix any TypeScript compilation errors in converted files
- [x] T018 Update component prop passing to match TypeScript interfaces
- [x] T019 Verify type checking works in development mode (npm start)

## Phase 3.5: Polish
- [x] T020 [P] Convert remaining test files to TypeScript: frontend/src/components/TaskItem/TaskItem.test.js → TaskItem.test.tsx
- [x] T021 [P] Add type assertions and better error handling in API service
- [x] T022 Run TypeScript compiler check: npx tsc --noEmit to verify no type errors
- [x] T023 Verify build process works: npm run build with TypeScript files
- [x] T024 Execute quickstart verification steps to ensure all functionality preserved

## Dependencies
- Setup (T001-T003) before all other phases
- Tests (T004-T007) before implementation (T008-T015)
- T008 (types) blocks T009-T015 (all implementation files need types)
- T009 (API service) blocks T012-T014 (components depend on API types)
- Implementation (T008-T015) before integration (T016-T019)
- Integration before polish (T020-T024)

## Parallel Example
```
# Launch T004-T007 together (test creation):
Task: "Create TypeScript contract test for Task interface in frontend/src/types/__tests__/task.test.ts"
Task: "Create TypeScript contract test for API functions in frontend/src/types/__tests__/api.test.ts"
Task: "Create TypeScript contract test for component props in frontend/src/types/__tests__/components.test.ts"
Task: "Convert existing TaskForm test to TypeScript: frontend/src/components/TaskForm/TaskForm.test.js → TaskForm.test.tsx"

# Launch T020-T021 together (polish tasks):
Task: "Convert remaining test files to TypeScript: frontend/src/components/TaskItem/TaskItem.test.js → TaskItem.test.tsx"
Task: "Add type assertions and better error handling in API service"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing TypeScript conversion
- Maintain existing functionality throughout migration
- Test each converted file individually before proceeding
- Use incremental approach - one file at a time

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - typescript-types.contract.ts → type definition tasks and tests

2. **From Data Model**:
   - Each interface → type creation task
   - Component props → component conversion tasks

3. **From Quickstart**:
   - Each verification step → validation task
   - Build and test scenarios → polish tasks

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All contracts have corresponding tests (T004-T007)
- [x] All entities have type definition tasks (T008)
- [x] All tests come before implementation (T004-T007 before T008-T015)
- [x] Parallel tasks truly independent (different files)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task