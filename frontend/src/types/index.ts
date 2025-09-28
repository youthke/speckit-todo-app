/**
 * TypeScript Type Definitions Contract
 *
 * This file defines the expected TypeScript types and interfaces
 * that must be implemented during the migration. These types serve
 * as contracts between components and the API layer.
 */

// ============================================================================
// CORE ENTITY TYPES
// ============================================================================

/**
 * Task entity as returned by the API
 */
export interface Task {
  id: number;
  title: string;
  completed: boolean;
  created_at: string; // ISO 8601 timestamp
  updated_at: string; // ISO 8601 timestamp
}

/**
 * Task creation payload (subset of Task)
 */
export interface CreateTaskData {
  title: string;
}

/**
 * Task update payload (partial Task)
 */
export interface UpdateTaskData {
  title?: string;
  completed?: boolean;
}

// ============================================================================
// API RESPONSE TYPES
// ============================================================================

/**
 * Response from GET /api/v1/tasks
 */
export interface TasksResponse {
  tasks: Task[];
  count: number;
}

/**
 * Response from GET /health
 */
export interface HealthResponse {
  status: 'healthy' | 'degraded' | 'unhealthy';
  database: 'connected' | 'disconnected' | 'error';
  timestamp: string; // ISO 8601
  version?: string;
  uptime?: number; // seconds
}

/**
 * Processed service status for UI consumption
 */
export interface ServiceStatus {
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

// ============================================================================
// COMPONENT PROPS TYPES
// ============================================================================

/**
 * Props for TaskForm component
 */
export interface TaskFormProps {
  onTaskCreated?: (task: Task) => void;
}

/**
 * Props for TaskItem component
 */
export interface TaskItemProps {
  task: Task;
  onTaskChange: () => void;
}

/**
 * Props for TaskList component
 */
export interface TaskListProps {
  showCompleted?: boolean;
  onTaskChange: () => void;
}

// ============================================================================
// API FUNCTION TYPES
// ============================================================================

/**
 * Parameters for getTasks API function
 */
export interface GetTasksParams {
  completed?: boolean;
}

/**
 * API function type definitions
 */
export type GetTasksFunction = (params?: GetTasksParams) => Promise<TasksResponse>;
export type GetTaskFunction = (id: number) => Promise<Task>;
export type CreateTaskFunction = (taskData: CreateTaskData) => Promise<Task>;
export type UpdateTaskFunction = (id: number, updates: UpdateTaskData) => Promise<Task>;
export type DeleteTaskFunction = (id: number) => Promise<void>;
export type CheckHealthFunction = () => Promise<HealthResponse>;

// ============================================================================
// EVENT HANDLER TYPES
// ============================================================================

export type TaskChangeHandler = () => void;
export type TaskCreatedHandler = (task: Task) => void;
export type FormSubmitHandler = (e: React.FormEvent<HTMLFormElement>) => void;
export type InputChangeHandler = (e: React.ChangeEvent<HTMLInputElement>) => void;
export type KeyPressHandler = (e: React.KeyboardEvent<HTMLInputElement>) => void;
export type ButtonClickHandler = (e: React.MouseEvent<HTMLButtonElement>) => void;

// ============================================================================
// STATE TYPES
// ============================================================================

/**
 * App component state structure
 */
export interface AppState {
  filter: 'all' | 'pending' | 'completed';
  refreshKey: number;
  serverStatus: 'checking' | 'connected' | 'disconnected';
}

/**
 * TaskForm component state structure
 */
export interface TaskFormState {
  title: string;
  isCreating: boolean;
  error: string | null;
}

/**
 * TaskItem component state structure
 */
export interface TaskItemState {
  isEditing: boolean;
  editTitle: string;
  isUpdating: boolean;
  error: string | null;
}

// ============================================================================
// VALIDATION TYPES
// ============================================================================

/**
 * Validation rules for task title
 */
export interface TaskTitleValidation {
  readonly MIN_LENGTH: 1;
  readonly MAX_LENGTH: 500;
  readonly REQUIRED: true;
}

/**
 * Error types that can occur in the application
 */
export type TaskErrorType =
  | 'TITLE_EMPTY'
  | 'TITLE_TOO_LONG'
  | 'NETWORK_ERROR'
  | 'SERVER_ERROR'
  | 'VALIDATION_ERROR'
  | 'UNKNOWN_ERROR';

/**
 * Error object structure
 */
export interface TaskError {
  type: TaskErrorType;
  message: string;
  field?: string;
}