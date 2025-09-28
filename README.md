# TODO App

A full-stack TODO application built with Go (backend) and React (frontend).

## Features

- ✅ Create new tasks
- ✅ Mark tasks as completed/pending
- ✅ Edit task titles
- ✅ Delete tasks
- ✅ Filter tasks by status (All/Pending/Completed)
- ✅ Persistent storage with SQLite
- ✅ Real-time server connection status
- ✅ Responsive design
- ✅ Input validation and error handling

## Technology Stack

### Backend
- **Go 1.21+** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **SQLite** - Database (development)
- **PostgreSQL** - Database (production ready)

### Frontend
- **React 18** - UI library
- **Axios** - HTTP client
- **CSS3** - Styling
- **React Testing Library** - Testing

## Project Structure

```
todo-app/
├── backend/                 # Go backend server
│   ├── cmd/server/          # Main application entry point
│   ├── internal/            # Private application code
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── models/          # Data models
│   │   ├── services/        # Business logic
│   │   └── storage/         # Database layer
│   ├── pkg/                 # Public library code
│   ├── tests/               # Test files
│   │   ├── contract/        # API contract tests
│   │   ├── integration/     # Integration tests
│   │   └── unit/            # Unit tests
│   ├── go.mod               # Go module definition
│   └── go.sum               # Go module checksums
├── frontend/                # React frontend application
│   ├── public/              # Static assets
│   ├── src/                 # Source code
│   │   ├── components/      # React components
│   │   │   ├── TaskForm/    # Task creation form
│   │   │   ├── TaskItem/    # Individual task display
│   │   │   └── TaskList/    # Task list container
│   │   ├── services/        # API communication
│   │   ├── App.js           # Main application component
│   │   └── App.css          # Main application styles
│   ├── package.json         # Node.js dependencies
│   └── package-lock.json    # Node.js lock file
└── specs/                   # Feature specifications
    └── 001-todo/            # TODO app specification
        ├── spec.md          # Feature specification
        ├── plan.md          # Implementation plan
        ├── tasks.md         # Implementation tasks
        ├── research.md      # Technical research
        ├── data-model.md    # Data model design
        ├── contracts/       # API contracts
        └── quickstart.md    # Testing guide
```

## Prerequisites

- **Go 1.21+** - [Download and install Go](https://golang.org/dl/)
- **Node.js 18+** - [Download and install Node.js](https://nodejs.org/)
- **Git** - [Download and install Git](https://git-scm.com/)

## Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd todo-app
```

### 2. Start the Backend Server
```bash
cd backend
go mod tidy              # Install Go dependencies
go run cmd/server/main.go # Start the server
```

The backend server will start on `http://localhost:8080`

### 3. Start the Frontend Development Server
```bash
cd frontend
npm install              # Install dependencies
npm start               # Start development server
```

The frontend will start on `http://localhost:3000` and automatically open in your browser.

## Development

### Backend Development

#### Running Tests
```bash
cd backend

# Run all tests
go test ./...

# Run specific test suites
go test ./tests/unit/...
go test ./tests/contract/...
go test ./tests/integration/...

# Run tests with coverage
go test -cover ./...
```

#### Database Operations
```bash
# The application automatically creates and migrates the database
# Database file: todo.db (SQLite)

# To reset the database, simply delete the file:
rm todo.db
```

#### Linting
```bash
cd backend
golangci-lint run
```

### Frontend Development

#### Running Tests
```bash
cd frontend

# Run tests
npm test

# Run tests with coverage
npm test -- --coverage

# Run tests in watch mode
npm test -- --watch
```

#### Building for Production
```bash
cd frontend
npm run build
```

#### Linting
```bash
cd frontend
npm run lint
```

## API Documentation

### Base URL
- Development: `http://localhost:8080/api/v1`

### Endpoints

#### Get All Tasks
```http
GET /tasks
GET /tasks?completed=true   # Filter completed tasks
GET /tasks?completed=false  # Filter pending tasks
```

#### Create Task
```http
POST /tasks
Content-Type: application/json

{
  "title": "Task title"
}
```

#### Get Single Task
```http
GET /tasks/{id}
```

#### Update Task
```http
PUT /tasks/{id}
Content-Type: application/json

{
  "title": "Updated title",     # Optional
  "completed": true             # Optional
}
```

#### Delete Task
```http
DELETE /tasks/{id}
```

#### Health Check
```http
GET /health
```

### Response Format

#### Success Response
```json
{
  "id": 1,
  "title": "Task title",
  "completed": false,
  "created_at": "2025-09-27T10:00:00Z",
  "updated_at": "2025-09-27T10:00:00Z"
}
```

#### Error Response
```json
{
  "error": "validation_error",
  "message": "Task title cannot be empty"
}
```

## Configuration

### Environment Variables

#### Backend
- `PORT` - Server port (default: 8080)
- `DB_PATH` - Database file path (default: todo.db)
- `ENV` - Environment (production/development)

#### Frontend
- `REACT_APP_API_URL` - Backend API URL (default: http://localhost:8080/api/v1)

## Deployment

### Backend Deployment
```bash
cd backend
go build -o todo-server cmd/server/main.go
./todo-server
```

### Frontend Deployment
```bash
cd frontend
npm run build
# Serve the build/ directory with any static file server
```

## Testing

### Manual Testing
Follow the scenarios in `specs/001-todo/quickstart.md` to manually test all features.

### Automated Testing
- Backend: Go testing framework with testify
- Frontend: Jest with React Testing Library
- Contract tests ensure API compliance
- Integration tests validate user scenarios

## Performance

- API response times: < 500ms
- Database: SQLite for development, PostgreSQL recommended for production
- Frontend: Optimized React build with code splitting
- Concurrent request handling with Go goroutines

## Troubleshooting

### Backend Issues
- **Port already in use**: Change the PORT environment variable
- **Database locked**: Ensure no other instances are running
- **Module not found**: Run `go mod tidy`

### Frontend Issues
- **Cannot connect to API**: Ensure backend is running on port 8080
- **Build fails**: Delete node_modules and run `npm install`
- **Tests fail**: Check that all dependencies are installed

### CORS Issues
- The backend is configured to allow requests from `http://localhost:3000`
- For production, update the CORS configuration in `cmd/server/main.go`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License.

## Acknowledgments

- Built with [Gin](https://gin-gonic.com/) web framework
- UI inspired by modern TODO applications
- Generated with [Claude Code](https://claude.ai/code)