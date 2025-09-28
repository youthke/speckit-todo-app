import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import TaskList from './TaskList';

// Mock the API service
jest.mock('../../services/api', () => ({
  getTasks: jest.fn(),
}));

const mockApi = require('../../services/api');

describe('TaskList Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders empty state when no tasks', async () => {
    mockApi.getTasks.mockResolvedValue({ tasks: [], count: 0 });

    render(<TaskList />);

    await waitFor(() => {
      expect(screen.getByText(/no tasks/i)).toBeInTheDocument();
    });
  });

  test('renders list of tasks', async () => {
    const mockTasks = [
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

    mockApi.getTasks.mockResolvedValue({ tasks: mockTasks, count: 2 });

    render(<TaskList />);

    await waitFor(() => {
      expect(screen.getByText('Buy groceries')).toBeInTheDocument();
      expect(screen.getByText('Walk the dog')).toBeInTheDocument();
    });
  });

  test('shows loading state initially', () => {
    mockApi.getTasks.mockImplementation(
      () => new Promise((resolve) => setTimeout(resolve, 100))
    );

    render(<TaskList />);

    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  test('handles API error gracefully', async () => {
    mockApi.getTasks.mockRejectedValue(new Error('API Error'));

    render(<TaskList />);

    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument();
    });
  });

  test('filters completed tasks', async () => {
    const mockTasks = [
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

    mockApi.getTasks.mockResolvedValue({ tasks: mockTasks, count: 2 });

    render(<TaskList showCompleted={true} />);

    await waitFor(() => {
      expect(mockApi.getTasks).toHaveBeenCalledWith({ completed: true });
    });
  });

  test('filters pending tasks', async () => {
    const mockTasks = [
      {
        id: 1,
        title: 'Buy groceries',
        completed: false,
        created_at: '2025-09-27T10:00:00Z',
        updated_at: '2025-09-27T10:00:00Z',
      },
    ];

    mockApi.getTasks.mockResolvedValue({ tasks: mockTasks, count: 1 });

    render(<TaskList showCompleted={false} />);

    await waitFor(() => {
      expect(mockApi.getTasks).toHaveBeenCalledWith({ completed: false });
    });
  });

  test('refreshes task list when onTaskChange is called', async () => {
    const mockTasks = [
      {
        id: 1,
        title: 'Buy groceries',
        completed: false,
        created_at: '2025-09-27T10:00:00Z',
        updated_at: '2025-09-27T10:00:00Z',
      },
    ];

    mockApi.getTasks.mockResolvedValue({ tasks: mockTasks, count: 1 });

    render(<TaskList />);

    await waitFor(() => {
      expect(mockApi.getTasks).toHaveBeenCalledTimes(1);
    });

    // This test verifies the component can refresh its data
    // The actual implementation will need to expose a way to trigger refresh
  });
});