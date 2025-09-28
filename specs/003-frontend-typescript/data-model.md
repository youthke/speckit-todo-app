# Data Model: Frontend TypeScript Migration

## Core TypeScript Interfaces

### Task Entity
```typescript
interface Task {
  id: number;
  title: string;
  completed: boolean;
  created_at: string; // ISO 8601 timestamp
  updated_at: string; // ISO 8601 timestamp
}
```

**Properties**:
- `id`: Unique identifier for the task
- `title`: Task description (1-500 characters)
- `completed`: Boolean completion status
- `created_at`: Timestamp when task was created
- `updated_at`: Timestamp when task was last modified

**Validation Rules**:
- Title must be non-empty and ≤ 500 characters
- Timestamps are ISO 8601 format strings
- ID is a positive integer

### API Response Types
```typescript
interface TasksResponse {
  tasks: Task[];
  count: number;
}

interface HealthResponse {
  status: 'healthy' | 'degraded' | 'unhealthy';
  database: 'connected' | 'disconnected' | 'error';
  timestamp: string; // ISO 8601
  version?: string;
  uptime?: number; // seconds
}

interface ServiceStatus {
  isHealthy: boolean;
  isDegraded: boolean;
  isUnhealthy: boolean;
  databaseConnected: boolean;
  timestamp: Date;
  version: string;
  uptime: number;
  uptimeFormatted: string;
  responseTime: number | null;
  error?: string;
}
```

### Component Props Types
```typescript
interface TaskFormProps {
  onTaskCreated?: (task: Task) => void;
}

interface TaskItemProps {
  task: Task;
  onTaskChange: () => void;
}

interface TaskListProps {
  showCompleted?: boolean;
  onTaskChange: () => void;
  key?: number; // React key prop
}
```

### API Function Parameters
```typescript
interface CreateTaskData {
  title: string;
}

interface UpdateTaskData {
  title?: string;
  completed?: boolean;
}

interface GetTasksParams {
  completed?: boolean;
}
```

### Component State Types
```typescript
// App component state
interface AppState {
  filter: 'all' | 'pending' | 'completed';
  refreshKey: number;
  serverStatus: 'checking' | 'connected' | 'disconnected';
}

// TaskForm component state
interface TaskFormState {
  title: string;
  isCreating: boolean;
  error: string | null;
}

// TaskItem component state
interface TaskItemState {
  isEditing: boolean;
  editTitle: string;
  isUpdating: boolean;
  error: string | null;
}
```

### Event Handler Types
```typescript
type TaskChangeHandler = () => void;
type TaskCreatedHandler = (task: Task) => void;
type FormSubmitHandler = (e: React.FormEvent<HTMLFormElement>) => void;
type KeyPressHandler = (e: React.KeyboardEvent<HTMLInputElement>) => void;
type ButtonClickHandler = (e: React.MouseEvent<HTMLButtonElement>) => void;
```

## Type Relationships

### State Transitions
- **Task Status**: pending → completed → pending (toggleable)
- **Server Status**: checking → connected/disconnected
- **Form State**: idle → creating → idle
- **Edit State**: viewing → editing → viewing

### Data Flow
1. **Task Creation**: CreateTaskData → API → Task → TasksResponse
2. **Task Update**: UpdateTaskData → API → Task
3. **Task Deletion**: ID → API → void
4. **Health Check**: void → API → HealthResponse → ServiceStatus

### Component Hierarchy Type Dependencies
```
App (AppState)
├── TaskForm (TaskFormProps, TaskFormState)
├── TaskList (TaskListProps)
│   └── TaskItem (TaskItemProps, TaskItemState)
└── API Service (all API types)
```

## Migration-Specific Types

### File Conversion Map
```typescript
interface FileConversion {
  source: string; // .js/.jsx file path
  target: string; // .ts/.tsx file path
  hasJSX: boolean;
  dependencies: string[]; // other files that import this
}
```

### TypeScript Configuration
```typescript
interface TSConfig {
  compilerOptions: {
    target: string;
    lib: string[];
    allowJs: boolean;
    skipLibCheck: boolean;
    esModuleInterop: boolean;
    allowSyntheticDefaultImports: boolean;
    strict: boolean;
    forceConsistentCasingInFileNames: boolean;
    moduleResolution: string;
    resolveJsonModule: boolean;
    isolatedModules: boolean;
    noEmit: boolean;
    jsx: string;
  };
  include: string[];
  exclude: string[];
}
```