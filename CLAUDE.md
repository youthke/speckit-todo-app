# todo-app Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-09-27

## Active Technologies
- Go 1.21+, Node.js 18+ + Gin/Echo (Go backend), React 18, SQLite for local developmen (001-todo)
- Go 1.23+ + Gin web framework, GORM ORM, SQLite database (002-api-health)
- SQLite (development), existing database connection from TODO app (002-api-health)
- TypeScript 5.x (migration from JavaScript ES6+) + React 19.1.1, React Scripts 5.0.1, Testing Library, Axios (003-frontend-typescript)
- N/A (frontend-only migration) (003-frontend-typescript)
- Go 1.23+ + Gin web framework, GORM ORM, testify testing framework (004-)
- SQLite database (development), existing schema to be migrated (004-)
- Go 1.23+, React 19.1.1 with TypeScript 5.x + Gin web framework, GORM ORM, Google OAuth 2.0 libraries (005-google)
- SQLite (development), existing database schema (005-google)
- TypeScript 5.9.2, React 19.1.1 + React Router (version TBD - v6 or v7), React DOM 19.1.1, Vite 6.0.11 (006-react-router)
- N/A (frontend routing only, no data persistence) (006-react-router)
- Backend: Go 1.24.0, Frontend: TypeScript 5.9.2 + React 19.1.1 + Backend: Gin web framework, GORM ORM, Google OAuth 2.0 libraries; Frontend: React Router 6.30.1, Axios 1.12.2, Vite 6.0.11 (007-google)
- SQLite database (development), existing schema to be extended with Google OAuth entities (007-google)
- Backend: Go 1.24.0, Frontend: TypeScript 5.9.2 + React 19.1.1 + Backend: Gin web framework, GORM ORM, Google OAuth 2.0 libraries (`golang.org/x/oauth2`); Frontend: React Router 6.30.1, Axios 1.12.2, Vite 6.0.11 (008-signup)
- SQLite (development), existing database schema with users table and Google OAuth suppor (008-signup)
- Go 1.24.7 + GORM ORM, Gin web framework, golang.org/x/oauth2 (009-resolve-it-1)
- Go 1.24.7 + Gin web framework v1.11.0, GORM ORM v1.31.0, testify v1.11.1, golang.org/x/oauth2 v0.31.0 (010-complete-ddd-migration)
- SQLite (GORM driver v1.6.0, development database) (010-complete-ddd-migration)

## Project Structure (DDD Architecture)
```
backend/
  domain/                   # Domain layer (entities, value objects, services)
    auth/                   # Auth domain (Google OAuth)
    health/                 # Health check domain
    task/                   # Task domain
    user/                   # User domain
  application/              # Application layer (use cases, commands)
    mappers/                # Entity↔DTO conversion
    task/                   # Task application services
    user/                   # User application services
  infrastructure/           # Infrastructure layer (repositories, persistence)
    persistence/            # GORM repositories
  internal/                 # Internal adapters (handlers, services, middleware)
    dtos/                   # Data Transfer Objects (API contracts)
    handlers/               # HTTP handlers
    services/               # Legacy services (being phased out)
    storage/                # Database connection
  tests/                    # Test suites
frontend/
tests/
```

## Commands
# Add commands for Go 1.21+, Node.js 18+

## Code Style
Go 1.21+, Node.js 18+: Follow standard conventions

## Recent Changes
- 010-complete-ddd-migration: **IN PROGRESS (85% complete)** - Migrated to full DDD architecture with domain entities, value objects, repositories, mappers, and application services. Core architecture complete, handlers and integration remaining.
  - ✅ Domain layer: Entities with rich behavior, value objects with validation
  - ✅ Application layer: Use case services, entity↔DTO mappers (84% coverage, <4μs performance)
  - ✅ Infrastructure layer: Repositories with mapper integration
  - ⚠️ Presentation layer: Handlers need mapper integration
  - ✅ Models renamed to DTOs (Data Transfer Objects)
  - ⚠️ Some test imports need domain layer updates
- 009-resolve-it-1: Added Go 1.24.7 + GORM ORM, Gin web framework, golang.org/x/oauth2
- 008-signup: Added Backend: Go 1.24.0, Frontend: TypeScript 5.9.2 + React 19.1.1 + Backend: Gin web framework, GORM ORM, Google OAuth 2.0 libraries (`golang.org/x/oauth2`); Frontend: React Router 6.30.1, Axios 1.12.2, Vite 6.0.11

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
