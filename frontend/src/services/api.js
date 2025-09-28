import axios from 'axios';

// Base API configuration
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000, // 10 seconds
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for logging
api.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    console.log(`API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  (error) => {
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
 * @param {Object} params - Query parameters
 * @param {boolean} params.completed - Filter by completion status
 * @returns {Promise<Object>} Response with tasks array and count
 */
export const getTasks = async (params = {}) => {
  const response = await api.get('/tasks', { params });
  return response.data;
};

/**
 * Get a specific task by ID
 * @param {number} id - Task ID
 * @returns {Promise<Object>} Task object
 */
export const getTask = async (id) => {
  const response = await api.get(`/tasks/${id}`);
  return response.data;
};

/**
 * Create a new task
 * @param {Object} taskData - Task data
 * @param {string} taskData.title - Task title (required)
 * @returns {Promise<Object>} Created task object
 */
export const createTask = async (taskData) => {
  const response = await api.post('/tasks', taskData);
  return response.data;
};

/**
 * Update an existing task
 * @param {number} id - Task ID
 * @param {Object} updates - Updates to apply
 * @param {string} updates.title - New task title
 * @param {boolean} updates.completed - New completion status
 * @returns {Promise<Object>} Updated task object
 */
export const updateTask = async (id, updates) => {
  const response = await api.put(`/tasks/${id}`, updates);
  return response.data;
};

/**
 * Delete a task
 * @param {number} id - Task ID
 * @returns {Promise<void>}
 */
export const deleteTask = async (id) => {
  await api.delete(`/tasks/${id}`);
};

/**
 * Check API health with enhanced monitoring information
 * @returns {Promise<Object>} Enhanced health status
 * @returns {string} returns.status - Overall health status ("healthy", "degraded", "unhealthy")
 * @returns {string} returns.database - Database connectivity status ("connected", "disconnected", "error")
 * @returns {string} returns.timestamp - ISO 8601 timestamp when check was performed
 * @returns {string} returns.version - Application version (optional)
 * @returns {number} returns.uptime - Service uptime in seconds (optional)
 */
export const checkHealth = async () => {
  const response = await api.get('/../health'); // Go up one level from /api/v1
  return response.data;
};

/**
 * Check if the service is currently healthy
 * @returns {Promise<boolean>} True if service is healthy, false otherwise
 */
export const isServiceHealthy = async () => {
  try {
    const health = await checkHealth();
    return health.status === 'healthy' && health.database === 'connected';
  } catch (error) {
    console.error('Health check failed:', error);
    return false;
  }
};

/**
 * Get detailed service status information
 * @returns {Promise<Object>} Detailed status information for monitoring
 */
export const getServiceStatus = async () => {
  try {
    const health = await checkHealth();

    return {
      isHealthy: health.status === 'healthy',
      isDegraded: health.status === 'degraded',
      isUnhealthy: health.status === 'unhealthy',
      databaseConnected: health.database === 'connected',
      timestamp: new Date(health.timestamp),
      version: health.version || 'unknown',
      uptime: health.uptime || 0,
      uptimeFormatted: formatUptime(health.uptime || 0),
      responseTime: null, // Can be measured by caller
    };
  } catch (error) {
    console.error('Failed to get service status:', error);
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
      error: error.message,
    };
  }
};

/**
 * Format uptime seconds into human-readable string
 * @param {number} seconds - Uptime in seconds
 * @returns {string} Formatted uptime string
 */
const formatUptime = (seconds) => {
  if (!seconds || seconds < 0) return '0 seconds';

  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);

  const parts = [];
  if (days > 0) parts.push(`${days} day${days !== 1 ? 's' : ''}`);
  if (hours > 0) parts.push(`${hours} hour${hours !== 1 ? 's' : ''}`);
  if (minutes > 0) parts.push(`${minutes} minute${minutes !== 1 ? 's' : ''}`);
  if (secs > 0 || parts.length === 0) parts.push(`${secs} second${secs !== 1 ? 's' : ''}`);

  return parts.join(', ');
};

export default api;