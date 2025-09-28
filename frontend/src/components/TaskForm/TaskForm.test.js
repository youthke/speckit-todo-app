import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import TaskForm from './TaskForm';

// Mock the API service
jest.mock('../../services/api', () => ({
  createTask: jest.fn(),
}));

const mockApi = require('../../services/api');

const mockOnTaskCreated = jest.fn();

describe('TaskForm Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders form elements correctly', () => {
    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    expect(screen.getByPlaceholderText(/enter task title/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /add task/i })).toBeInTheDocument();
  });

  test('creates task when form is submitted with valid title', async () => {
    const newTask = {
      id: 1,
      title: 'New task',
      completed: false,
      created_at: '2025-09-27T10:00:00Z',
      updated_at: '2025-09-27T10:00:00Z',
    };

    mockApi.createTask.mockResolvedValue(newTask);

    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    const input = screen.getByPlaceholderText(/enter task title/i);
    const submitButton = screen.getByRole('button', { name: /add task/i });

    // Enter task title
    fireEvent.change(input, { target: { value: 'New task' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockApi.createTask).toHaveBeenCalledWith({ title: 'New task' });
      expect(mockOnTaskCreated).toHaveBeenCalledWith(newTask);
    });

    // Input should be cleared after successful creation
    expect(input.value).toBe('');
  });

  test('creates task when Enter key is pressed', async () => {
    const newTask = {
      id: 1,
      title: 'New task',
      completed: false,
      created_at: '2025-09-27T10:00:00Z',
      updated_at: '2025-09-27T10:00:00Z',
    };

    mockApi.createTask.mockResolvedValue(newTask);

    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    const input = screen.getByPlaceholderText(/enter task title/i);

    // Enter task title and press Enter
    fireEvent.change(input, { target: { value: 'New task' } });
    fireEvent.keyPress(input, { key: 'Enter', code: 'Enter', charCode: 13 });

    await waitFor(() => {
      expect(mockApi.createTask).toHaveBeenCalledWith({ title: 'New task' });
      expect(mockOnTaskCreated).toHaveBeenCalledWith(newTask);
    });
  });

  test('prevents creating task with empty title', async () => {
    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    const submitButton = screen.getByRole('button', { name: /add task/i });

    // Try to submit with empty title
    fireEvent.click(submitButton);

    // Should not call API
    expect(mockApi.createTask).not.toHaveBeenCalled();
    expect(mockOnTaskCreated).not.toHaveBeenCalled();

    // Should show validation error
    expect(screen.getByText(/title cannot be empty/i)).toBeInTheDocument();
  });

  test('prevents creating task with only whitespace title', async () => {
    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    const input = screen.getByPlaceholderText(/enter task title/i);
    const submitButton = screen.getByRole('button', { name: /add task/i });

    // Enter only whitespace
    fireEvent.change(input, { target: { value: '   ' } });
    fireEvent.click(submitButton);

    // Should not call API
    expect(mockApi.createTask).not.toHaveBeenCalled();
    expect(mockOnTaskCreated).not.toHaveBeenCalled();

    // Should show validation error
    expect(screen.getByText(/title cannot be empty/i)).toBeInTheDocument();
  });

  test('prevents creating task with title longer than 500 characters', async () => {
    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    const input = screen.getByPlaceholderText(/enter task title/i);
    const submitButton = screen.getByRole('button', { name: /add task/i });

    // Enter very long title
    const longTitle = 'a'.repeat(501);
    fireEvent.change(input, { target: { value: longTitle } });
    fireEvent.click(submitButton);

    // Should not call API
    expect(mockApi.createTask).not.toHaveBeenCalled();
    expect(mockOnTaskCreated).not.toHaveBeenCalled();

    // Should show validation error
    expect(screen.getByText(/title must be 500 characters or less/i)).toBeInTheDocument();
  });

  test('shows loading state while creating task', async () => {
    mockApi.createTask.mockImplementation(
      () => new Promise((resolve) => setTimeout(resolve, 100))
    );

    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    const input = screen.getByPlaceholderText(/enter task title/i);
    const submitButton = screen.getByRole('button', { name: /add task/i });

    fireEvent.change(input, { target: { value: 'New task' } });
    fireEvent.click(submitButton);

    // Button should be disabled and show loading state
    expect(submitButton).toBeDisabled();
    expect(screen.getByText(/creating/i)).toBeInTheDocument();
  });

  test('handles API errors gracefully', async () => {
    mockApi.createTask.mockRejectedValue(new Error('API Error'));

    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    const input = screen.getByPlaceholderText(/enter task title/i);
    const submitButton = screen.getByRole('button', { name: /add task/i });

    fireEvent.change(input, { target: { value: 'New task' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/error creating task/i)).toBeInTheDocument();
    });

    // Form should be re-enabled
    expect(submitButton).not.toBeDisabled();
    // Input value should be preserved
    expect(input.value).toBe('New task');
  });

  test('trims whitespace from task title', async () => {
    const newTask = {
      id: 1,
      title: 'Trimmed task',
      completed: false,
      created_at: '2025-09-27T10:00:00Z',
      updated_at: '2025-09-27T10:00:00Z',
    };

    mockApi.createTask.mockResolvedValue(newTask);

    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    const input = screen.getByPlaceholderText(/enter task title/i);
    const submitButton = screen.getByRole('button', { name: /add task/i });

    // Enter task title with surrounding whitespace
    fireEvent.change(input, { target: { value: '  Trimmed task  ' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockApi.createTask).toHaveBeenCalledWith({ title: 'Trimmed task' });
    });
  });

  test('clears validation errors when user starts typing', async () => {
    render(<TaskForm onTaskCreated={mockOnTaskCreated} />);

    const input = screen.getByPlaceholderText(/enter task title/i);
    const submitButton = screen.getByRole('button', { name: /add task/i });

    // Trigger validation error
    fireEvent.click(submitButton);
    expect(screen.getByText(/title cannot be empty/i)).toBeInTheDocument();

    // Start typing
    fireEvent.change(input, { target: { value: 'N' } });

    // Error should be cleared
    expect(screen.queryByText(/title cannot be empty/i)).not.toBeInTheDocument();
  });
});