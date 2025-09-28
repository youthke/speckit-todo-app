/**
 * Contract tests for API function types and health status
 * These tests verify the API-related TypeScript type definitions work as expected
 */

import {
  GetTasksParams,
  HealthResponse,
  ServiceStatus,
  GetTasksFunction,
  CreateTaskFunction,
  UpdateTaskFunction,
  DeleteTaskFunction,
  CheckHealthFunction,
} from '../index';

describe('API Types Contract', () => {
  it('should define GetTasksParams interface correctly', () => {
    const allTasks: GetTasksParams = {};
    const completedTasks: GetTasksParams = { completed: true };
    const pendingTasks: GetTasksParams = { completed: false };

    expect(allTasks.completed).toBeUndefined();
    expect(completedTasks.completed).toBe(true);
    expect(pendingTasks.completed).toBe(false);
  });

  it('should define HealthResponse interface correctly', () => {
    const healthyResponse: HealthResponse = {
      status: 'healthy',
      database: 'connected',
      timestamp: '2025-09-28T12:00:00Z',
      version: '1.0.0',
      uptime: 3600,
    };

    const unhealthyResponse: HealthResponse = {
      status: 'unhealthy',
      database: 'disconnected',
      timestamp: '2025-09-28T12:00:00Z',
    };

    expect(healthyResponse.status).toBe('healthy');
    expect(healthyResponse.database).toBe('connected');
    expect(healthyResponse.version).toBe('1.0.0');
    expect(healthyResponse.uptime).toBe(3600);

    expect(unhealthyResponse.status).toBe('unhealthy');
    expect(unhealthyResponse.database).toBe('disconnected');
    expect(unhealthyResponse.version).toBeUndefined();
    expect(unhealthyResponse.uptime).toBeUndefined();
  });

  it('should define ServiceStatus interface correctly', () => {
    const serviceStatus: ServiceStatus = {
      isHealthy: true,
      isDegraded: false,
      isUnhealthy: false,
      databaseConnected: true,
      timestamp: new Date('2025-09-28T12:00:00Z'),
      version: '1.0.0',
      uptime: 3600,
      uptimeFormatted: '1 hour',
      responseTime: 150,
      error: undefined,
    };

    expect(serviceStatus.isHealthy).toBe(true);
    expect(serviceStatus.timestamp).toBeInstanceOf(Date);
    expect(serviceStatus.responseTime).toBe(150);
  });

  it('should define API function types correctly', () => {
    // Mock functions that should match the type signatures
    const mockGetTasks: GetTasksFunction = async (params) => {
      return {
        tasks: [],
        count: 0,
      };
    };

    const mockCreateTask: CreateTaskFunction = async (taskData) => {
      return {
        id: 1,
        title: taskData.title,
        completed: false,
        created_at: '2025-09-28T12:00:00Z',
        updated_at: '2025-09-28T12:00:00Z',
      };
    };

    const mockUpdateTask: UpdateTaskFunction = async (id, updates) => {
      return {
        id,
        title: updates.title || 'Default Title',
        completed: updates.completed || false,
        created_at: '2025-09-28T12:00:00Z',
        updated_at: '2025-09-28T12:00:00Z',
      };
    };

    const mockDeleteTask: DeleteTaskFunction = async (id) => {
      // Returns void
    };

    const mockCheckHealth: CheckHealthFunction = async () => {
      return {
        status: 'healthy',
        database: 'connected',
        timestamp: '2025-09-28T12:00:00Z',
      };
    };

    // These should compile without TypeScript errors
    expect(typeof mockGetTasks).toBe('function');
    expect(typeof mockCreateTask).toBe('function');
    expect(typeof mockUpdateTask).toBe('function');
    expect(typeof mockDeleteTask).toBe('function');
    expect(typeof mockCheckHealth).toBe('function');
  });

  it('should enforce status literal types', () => {
    // Valid status values
    const healthyStatus: HealthResponse['status'] = 'healthy';
    const degradedStatus: HealthResponse['status'] = 'degraded';
    const unhealthyStatus: HealthResponse['status'] = 'unhealthy';

    const connectedDb: HealthResponse['database'] = 'connected';
    const disconnectedDb: HealthResponse['database'] = 'disconnected';
    const errorDb: HealthResponse['database'] = 'error';

    expect(['healthy', 'degraded', 'unhealthy']).toContain(healthyStatus);
    expect(['connected', 'disconnected', 'error']).toContain(connectedDb);

    // These would cause TypeScript compilation errors if uncommented:
    /*
    const invalidStatus: HealthResponse['status'] = 'invalid';
    const invalidDatabase: HealthResponse['database'] = 'unknown';
    */

    // This test passes if TypeScript compilation succeeds
    expect(true).toBe(true);
  });
});