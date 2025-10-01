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

## Project Structure
```
backend/
frontend/
tests/
```

## Commands
# Add commands for Go 1.21+, Node.js 18+

## Code Style
Go 1.21+, Node.js 18+: Follow standard conventions

## Recent Changes
- 007-google: Added Backend: Go 1.24.0, Frontend: TypeScript 5.9.2 + React 19.1.1 + Backend: Gin web framework, GORM ORM, Google OAuth 2.0 libraries; Frontend: React Router 6.30.1, Axios 1.12.2, Vite 6.0.11
- 006-react-router: Added TypeScript 5.9.2, React 19.1.1 + React Router (version TBD - v6 or v7), React DOM 19.1.1, Vite 6.0.11
- 005-google: Added Go 1.23+, React 19.1.1 with TypeScript 5.x + Gin web framework, GORM ORM, Google OAuth 2.0 libraries

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
