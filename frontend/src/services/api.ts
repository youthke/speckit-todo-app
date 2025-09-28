import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import {
  Task,
  CreateTaskData,
  UpdateTaskData,
  TasksResponse,
  HealthResponse,
  ServiceStatus,
  GetTasksParams,
} from '../types';

// Base API configuration
const API_BASE_URL: string = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000, // 10 seconds
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for logging
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error: any): Promise<never> => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => {
    console.log(`API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  (error: any): Promise<never> => {
    console.error('API Response Error:', error.response?.data || error.message);

    // Handle common error scenarios
    if (error.response?.status === 404) {
      throw new Error('Resource not found');
    } else if (error.response?.status === 400) {
      throw new Error(error.response.data?.message || 'Invalid request');
    } else if (error.response?.status >= 500) {
      throw new Error('Server error. Please try again later.');
    } else if (error.code === 'ECONNABORTED') {
      throw new Error('Request timeout. Please check your connection.');
    } else if (error.code === 'ERR_NETWORK') {
      throw new Error('Network error. Please check your connection and ensure the server is running.');
    }

    throw new Error(error.response?.data?.message || 'An unexpected error occurred');
  }
);

// Task API functions

/**
 * Get all tasks with optional filtering
 */
export const getTasks = async (params: GetTasksParams = {}): Promise<TasksResponse> => {
  const response: AxiosResponse<TasksResponse> = await api.get('/tasks', { params });
  return response.data;
};

/**
 * Get a specific task by ID
 */
export const getTask = async (id: number): Promise<Task> => {
  const response: AxiosResponse<Task> = await api.get(`/tasks/${id}`);
  return response.data;
};

/**
 * Create a new task
 */
export const createTask = async (taskData: CreateTaskData): Promise<Task> => {
  const response: AxiosResponse<Task> = await api.post('/tasks', taskData);
  return response.data;
};

/**
 * Update an existing task
 */
export const updateTask = async (id: number, updates: UpdateTaskData): Promise<Task> => {
  const response: AxiosResponse<Task> = await api.put(`/tasks/${id}`, updates);
  return response.data;
};

/**
 * Delete a task
 */
export const deleteTask = async (id: number): Promise<void> => {
  await api.delete(`/tasks/${id}`);
};

/**
 * Check API health with enhanced monitoring information
 */
export const checkHealth = async (): Promise<HealthResponse> => {
  const response: AxiosResponse<HealthResponse> = await api.get('/../health'); // Go up one level from /api/v1
  return response.data;
};

/**
 * Check if the service is currently healthy
 */
export const isServiceHealthy = async (): Promise<boolean> => {
  try {
    const health: HealthResponse = await checkHealth();
    // Type assertion to ensure we have the expected health response structure
    if (typeof health.status !== 'string' || typeof health.database !== 'string') {
      console.error('Invalid health response structure:', health);
      return false;
    }
    return health.status === 'healthy' && health.database === 'connected';
  } catch (error: unknown) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
    console.error('Health check failed:', errorMessage);
    return false;
  }
};

/**
 * Get detailed service status information
 */
export const getServiceStatus = async (): Promise<ServiceStatus> => {
  try {
    const health: HealthResponse = await checkHealth();

    // Type assertions for better error handling
    if (!health.timestamp || typeof health.timestamp !== 'string') {
      throw new Error('Invalid health response: missing or invalid timestamp');
    }

    const timestamp = new Date(health.timestamp);
    if (isNaN(timestamp.getTime())) {
      throw new Error('Invalid health response: unparseable timestamp');
    }

    return {
      isHealthy: health.status === 'healthy',
      isDegraded: health.status === 'degraded',
      isUnhealthy: health.status === 'unhealthy',
      databaseConnected: health.database === 'connected',
      timestamp,
      version: health.version || 'unknown',
      uptime: health.uptime || 0,
      uptimeFormatted: formatUptime(health.uptime || 0),
      responseTime: null, // Can be measured by caller
    };
  } catch (error: unknown) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
    console.error('Failed to get service status:', errorMessage);
    return {
      isHealthy: false,
      isDegraded: false,
      isUnhealthy: true,
      databaseConnected: false,
      timestamp: new Date(),
      version: 'unknown',
      uptime: 0,
      uptimeFormatted: 'unknown',
      responseTime: null,
      error: errorMessage,
    };
  }
};

/**
 * Format uptime seconds into human-readable string
 */
const formatUptime = (seconds: number): string => {
  if (!seconds || seconds < 0) return '0 seconds';

  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);

  const parts: string[] = [];
  if (days > 0) parts.push(`${days} day${days !== 1 ? 's' : ''}`);
  if (hours > 0) parts.push(`${hours} hour${hours !== 1 ? 's' : ''}`);
  if (minutes > 0) parts.push(`${minutes} minute${minutes !== 1 ? 's' : ''}`);
  if (secs > 0 || parts.length === 0) parts.push(`${secs} second${secs !== 1 ? 's' : ''}`);

  return parts.join(', ');
};

export default api;