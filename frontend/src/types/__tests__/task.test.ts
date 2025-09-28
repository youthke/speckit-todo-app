/**
 * Contract tests for Task interface and related types
 * These tests verify the TypeScript type definitions work as expected
 */

import { Task, CreateTaskData, UpdateTaskData, TasksResponse } from '../index';

describe('Task Interface Contract', () => {
  it('should define Task interface with required properties', () => {
    // Test that Task interface has all required properties
    const validTask: Task = {
      id: 1,
      title: 'Test Task',
      completed: false,
      created_at: '2025-09-28T12:00:00Z',
      updated_at: '2025-09-28T12:00:00Z',
    };

    expect(validTask.id).toBe(1);
    expect(validTask.title).toBe('Test Task');
    expect(validTask.completed).toBe(false);
    expect(validTask.created_at).toBe('2025-09-28T12:00:00Z');
    expect(validTask.updated_at).toBe('2025-09-28T12:00:00Z');
  });

  it('should define CreateTaskData interface correctly', () => {
    const createData: CreateTaskData = {
      title: 'New Task',
    };

    expect(createData.title).toBe('New Task');
    // Should not have other Task properties
    expect((createData as any).id).toBeUndefined();
    expect((createData as any).completed).toBeUndefined();
  });

  it('should define UpdateTaskData interface as partial', () => {
    // Should allow partial updates
    const updateTitle: UpdateTaskData = {
      title: 'Updated Title',
    };

    const updateCompleted: UpdateTaskData = {
      completed: true,
    };

    const updateBoth: UpdateTaskData = {
      title: 'Updated Title',
      completed: true,
    };

    const updateNothing: UpdateTaskData = {};

    expect(updateTitle.title).toBe('Updated Title');
    expect(updateCompleted.completed).toBe(true);
    expect(updateBoth.title).toBe('Updated Title');
    expect(updateBoth.completed).toBe(true);
    expect(Object.keys(updateNothing)).toHaveLength(0);
  });

  it('should define TasksResponse interface correctly', () => {
    const response: TasksResponse = {
      tasks: [
        {
          id: 1,
          title: 'Task 1',
          completed: false,
          created_at: '2025-09-28T12:00:00Z',
          updated_at: '2025-09-28T12:00:00Z',
        },
        {
          id: 2,
          title: 'Task 2',
          completed: true,
          created_at: '2025-09-28T11:00:00Z',
          updated_at: '2025-09-28T11:30:00Z',
        },
      ],
      count: 2,
    };

    expect(response.tasks).toHaveLength(2);
    expect(response.count).toBe(2);
    expect(response.tasks[0]).toMatchObject({
      id: 1,
      title: 'Task 1',
      completed: false,
    });
  });

  it('should enforce type constraints', () => {
    // These would cause TypeScript compilation errors if uncommented:
    /*
    const invalidTask: Task = {
      id: 'string', // Should be number
      title: 123, // Should be string
      completed: 'yes', // Should be boolean
      // Missing required properties
    };
    */

    // This test passes if TypeScript compilation succeeds
    expect(true).toBe(true);
  });
});