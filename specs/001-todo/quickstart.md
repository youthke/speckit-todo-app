# Quickstart Guide: TODO App

This guide validates the TODO app implementation by walking through all core user scenarios.

## Prerequisites

- Go 1.21+ installed
- Node.js 18+ installed
- Git repository cloned

## Setup

### 1. Backend Setup
```bash
cd backend
go mod init todo-app
go mod tidy
go run cmd/server/main.go
```
Expected output: `Server starting on :8080`

### 2. Frontend Setup
```bash
cd frontend
npm install
npm start
```
Expected output: `Development server started on http://localhost:3000`

### 3. Verify API Connection
```bash
curl http://localhost:8080/api/v1/tasks
```
Expected response: `{"tasks": [], "count": 0}`

## Core User Scenarios Validation

### Scenario 1: Creating Tasks
**Given**: Empty TODO list
**When**: User adds a new task

**Frontend Actions**:
1. Open http://localhost:3000
2. Enter "Buy groceries" in the task input field
3. Click "Add Task" button

**Expected Results**:
- ✅ Task appears in the list with title "Buy groceries"
- ✅ Task shows as pending (not completed)
- ✅ Task has creation timestamp
- ✅ Input field clears after creation

**API Test**:
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries"}'
```
Expected: HTTP 201 with task object including ID

### Scenario 2: Marking Tasks Complete
**Given**: Task "Buy groceries" exists in list
**When**: User marks task as completed

**Frontend Actions**:
1. Click the checkbox next to "Buy groceries"
2. Observe visual change

**Expected Results**:
- ✅ Task shows visual indication of completion (strikethrough/different color)
- ✅ Task moves to completed section (if applicable)
- ✅ Completion status persists on page refresh

**API Test**:
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```
Expected: HTTP 200 with updated task (completed: true)

### Scenario 3: Editing Task Titles
**Given**: Task exists in the list
**When**: User edits the task title

**Frontend Actions**:
1. Click edit button/icon on an existing task
2. Change title from "Buy groceries" to "Buy groceries and cook dinner"
3. Save changes

**Expected Results**:
- ✅ Task title updates to new text
- ✅ Updated timestamp changes
- ✅ Other task properties remain unchanged

**API Test**:
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries and cook dinner"}'
```
Expected: HTTP 200 with updated task title

### Scenario 4: Deleting Tasks
**Given**: Multiple tasks exist
**When**: User deletes a specific task

**Frontend Actions**:
1. Click delete button/icon on a task
2. Confirm deletion (if confirmation dialog exists)

**Expected Results**:
- ✅ Task immediately disappears from list
- ✅ Remaining tasks are unaffected
- ✅ Task deletion persists on page refresh

**API Test**:
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1
```
Expected: HTTP 204 (No Content)

### Scenario 5: Viewing Task Status
**Given**: Mix of completed and pending tasks
**When**: User views the task list

**Expected Results**:
- ✅ Clear visual distinction between completed and pending tasks
- ✅ Tasks display in logical order (newest first recommended)
- ✅ All task information visible (title, status, timestamps)

## Edge Cases Validation

### Empty Title Prevention
**Test**: Try to create task with empty title
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": ""}'
```
Expected: HTTP 400 with validation error

**Frontend**: Submit form with empty input
Expected: Form validation prevents submission

### Long Title Handling
**Test**: Create task with 500+ character title
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "'"$(python3 -c "print('a' * 501)")"'"}'
```
Expected: HTTP 400 with validation error

### Non-existent Task Operations
**Test**: Try to update/delete non-existent task
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/999 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated title"}'
```
Expected: HTTP 404 with error message

## Data Persistence Validation

### Session Persistence
1. Create several tasks
2. Mark some as completed
3. Close browser/stop frontend
4. Restart frontend
5. Verify all tasks and their states are preserved

### Server Restart Persistence
1. Create tasks via API
2. Stop backend server
3. Restart backend server
4. Verify tasks are still available via API

## Performance Validation

### Response Time Check
```bash
time curl http://localhost:8080/api/v1/tasks
```
Expected: Response time < 500ms (as per technical requirements)

### UI Responsiveness
- Task creation should be instantaneous
- Status toggles should respond immediately
- No noticeable lag in UI interactions

## Success Criteria

✅ All API endpoints respond correctly to valid requests
✅ All API endpoints reject invalid requests with appropriate errors
✅ Frontend displays tasks correctly
✅ Frontend allows all CRUD operations
✅ Data persists between sessions
✅ Performance meets stated requirements
✅ Edge cases handled gracefully

## Troubleshooting

### Backend Issues
- **Port 8080 in use**: Change port in server configuration
- **Database connection**: Check SQLite file permissions
- **CORS errors**: Verify frontend origin in CORS config

### Frontend Issues
- **API connection failed**: Check backend is running on correct port
- **Build errors**: Run `npm install` to ensure dependencies
- **State not updating**: Check React component state management

### Integration Issues
- **API calls failing**: Verify API base URL configuration
- **Data not persisting**: Check database file creation and permissions