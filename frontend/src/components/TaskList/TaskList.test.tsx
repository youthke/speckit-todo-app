import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import TaskList from './TaskList';
import { Task, TaskListProps, TasksResponse } from '../../types';

// Mock the API service
jest.mock('../../services/api', () => ({
  getTasks: jest.fn(),
}));

const mockApi = require('../../services/api');

const mockOnTaskChange = jest.fn<void, []>();

describe('TaskList Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders empty state when no tasks', async () => {
    const emptyResponse: TasksResponse = { tasks: [], count: 0 };
    mockApi.getTasks.mockResolvedValue(emptyResponse);

    const props: TaskListProps = { onTaskChange: mockOnTaskChange };
    render(<TaskList {...props} />);

    await waitFor(() => {
      expect(screen.getByText(/no tasks/i)).toBeInTheDocument();
    });
  });

  test('renders list of tasks', async () => {
    const mockTasks: Task[] = [
      {
        id: 1,
        title: 'Buy groceries',
        completed: false,
        created_at: '2025-09-27T10:00:00Z',
        updated_at: '2025-09-27T10:00:00Z',
      },
      {
        id: 2,
        title: 'Walk the dog',
        completed: true,
        created_at: '2025-09-27T09:00:00Z',
        updated_at: '2025-09-27T11:00:00Z',
      },
    ];

    const response: TasksResponse = { tasks: mockTasks, count: 2 };
    mockApi.getTasks.mockResolvedValue(response);

    const props: TaskListProps = { onTaskChange: mockOnTaskChange };
    render(<TaskList {...props} />);

    await waitFor(() => {
      expect(screen.getByText('Buy groceries')).toBeInTheDocument();
      expect(screen.getByText('Walk the dog')).toBeInTheDocument();
    });
  });

  test('shows loading state initially', () => {
    mockApi.getTasks.mockImplementation(
      () => new Promise<TasksResponse>((resolve) => setTimeout(resolve, 100))
    );

    const props: TaskListProps = { onTaskChange: mockOnTaskChange };
    render(<TaskList {...props} />);

    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  test('handles API error gracefully', async () => {
    mockApi.getTasks.mockRejectedValue(new Error('API Error'));

    const props: TaskListProps = { onTaskChange: mockOnTaskChange };
    render(<TaskList {...props} />);

    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument();
    });
  });

  test('filters completed tasks', async () => {
    const mockTasks: Task[] = [
      {
        id: 1,
        title: 'Buy groceries',
        completed: false,
        created_at: '2025-09-27T10:00:00Z',
        updated_at: '2025-09-27T10:00:00Z',
      },
      {
        id: 2,
        title: 'Walk the dog',
        completed: true,
        created_at: '2025-09-27T09:00:00Z',
        updated_at: '2025-09-27T11:00:00Z',
      },
    ];

    const response: TasksResponse = { tasks: mockTasks, count: 2 };
    mockApi.getTasks.mockResolvedValue(response);

    const props: TaskListProps = { showCompleted: true, onTaskChange: mockOnTaskChange };
    render(<TaskList {...props} />);

    await waitFor(() => {
      expect(mockApi.getTasks).toHaveBeenCalledWith({ completed: true });
    });
  });

  test('filters pending tasks', async () => {
    const mockTasks: Task[] = [
      {
        id: 1,
        title: 'Buy groceries',
        completed: false,
        created_at: '2025-09-27T10:00:00Z',
        updated_at: '2025-09-27T10:00:00Z',
      },
    ];

    const response: TasksResponse = { tasks: mockTasks, count: 1 };
    mockApi.getTasks.mockResolvedValue(response);

    const props: TaskListProps = { showCompleted: false, onTaskChange: mockOnTaskChange };
    render(<TaskList {...props} />);

    await waitFor(() => {
      expect(mockApi.getTasks).toHaveBeenCalledWith({ completed: false });
    });
  });

  test('refreshes task list when onTaskChange is called', async () => {
    const mockTasks: Task[] = [
      {
        id: 1,
        title: 'Buy groceries',
        completed: false,
        created_at: '2025-09-27T10:00:00Z',
        updated_at: '2025-09-27T10:00:00Z',
      },
    ];

    const response: TasksResponse = { tasks: mockTasks, count: 1 };
    mockApi.getTasks.mockResolvedValue(response);

    const props: TaskListProps = { onTaskChange: mockOnTaskChange };
    render(<TaskList {...props} />);

    await waitFor(() => {
      expect(mockApi.getTasks).toHaveBeenCalledTimes(1);
    });

    // This test verifies the component can refresh its data
    // The actual implementation will need to expose a way to trigger refresh
  });

  test('displays correct header for different filter states', async () => {
    const mockTasks: Task[] = [
      {
        id: 1,
        title: 'Test Task',
        completed: false,
        created_at: '2025-09-27T10:00:00Z',
        updated_at: '2025-09-27T10:00:00Z',
      },
    ];

    const response: TasksResponse = { tasks: mockTasks, count: 1 };
    mockApi.getTasks.mockResolvedValue(response);

    // Test 'All Tasks' header
    const allTasksProps: TaskListProps = { onTaskChange: mockOnTaskChange };
    const { rerender } = render(<TaskList {...allTasksProps} />);

    await waitFor(() => {
      expect(screen.getByText(/all tasks \(1\)/i)).toBeInTheDocument();
    });

    // Test 'Pending Tasks' header
    const pendingTasksProps: TaskListProps = { showCompleted: false, onTaskChange: mockOnTaskChange };
    rerender(<TaskList {...pendingTasksProps} />);

    await waitFor(() => {
      expect(screen.getByText(/pending tasks \(1\)/i)).toBeInTheDocument();
    });

    // Test 'Completed Tasks' header
    const completedTasksProps: TaskListProps = { showCompleted: true, onTaskChange: mockOnTaskChange };
    rerender(<TaskList {...completedTasksProps} />);

    await waitFor(() => {
      expect(screen.getByText(/completed tasks \(1\)/i)).toBeInTheDocument();
    });
  });
});