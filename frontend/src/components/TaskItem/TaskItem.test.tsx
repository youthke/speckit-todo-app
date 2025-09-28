import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import TaskItem from './TaskItem';
import { Task, TaskItemProps } from '../../types';

// Mock the API service
jest.mock('../../services/api', () => ({
  updateTask: jest.fn(),
  deleteTask: jest.fn(),
}));

const mockApi = require('../../services/api');

const mockTask: Task = {
  id: 1,
  title: 'Buy groceries',
  completed: false,
  created_at: '2025-09-27T10:00:00Z',
  updated_at: '2025-09-27T10:00:00Z',
};

const mockOnTaskChange = jest.fn<void, []>();

describe('TaskItem Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders task information correctly', () => {
    const props: TaskItemProps = { task: mockTask, onTaskChange: mockOnTaskChange };
    render(<TaskItem {...props} />);

    expect(screen.getByText('Buy groceries')).toBeInTheDocument();
    expect(screen.getByRole('checkbox')).not.toBeChecked();
  });

  test('renders completed task with different styling', () => {
    const completedTask: Task = { ...mockTask, completed: true };
    const props: TaskItemProps = { task: completedTask, onTaskChange: mockOnTaskChange };

    render(<TaskItem {...props} />);

    expect(screen.getByRole('checkbox')).toBeChecked();
    expect(screen.getByText('Buy groceries')).toHaveClass('completed');
  });

  test('toggles task completion when checkbox is clicked', async () => {
    const updatedTask: Task = { ...mockTask, completed: true };
    mockApi.updateTask.mockResolvedValue(updatedTask);

    const props: TaskItemProps = { task: mockTask, onTaskChange: mockOnTaskChange };
    render(<TaskItem {...props} />);

    const checkbox = screen.getByRole('checkbox') as HTMLInputElement;
    fireEvent.click(checkbox);

    await waitFor(() => {
      expect(mockApi.updateTask).toHaveBeenCalledWith(1, { completed: true });
      expect(mockOnTaskChange).toHaveBeenCalled();
    });
  });

  test('enters edit mode when title is clicked', () => {
    const props: TaskItemProps = { task: mockTask, onTaskChange: mockOnTaskChange };
    render(<TaskItem {...props} />);

    const title = screen.getByText('Buy groceries');
    fireEvent.click(title);

    expect(screen.getByDisplayValue('Buy groceries')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument();
  });

  test('saves edited title when save button is clicked', async () => {
    const updatedTask: Task = { ...mockTask, title: 'Buy groceries and cook' };
    mockApi.updateTask.mockResolvedValue(updatedTask);

    const props: TaskItemProps = { task: mockTask, onTaskChange: mockOnTaskChange };
    render(<TaskItem {...props} />);

    // Enter edit mode
    const title = screen.getByText('Buy groceries');
    fireEvent.click(title);

    // Edit the title
    const input = screen.getByDisplayValue('Buy groceries') as HTMLInputElement;
    fireEvent.change(input, { target: { value: 'Buy groceries and cook' } });

    // Save changes
    const saveButton = screen.getByRole('button', { name: /save/i });
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(mockApi.updateTask).toHaveBeenCalledWith(1, {
        title: 'Buy groceries and cook',
      });
      expect(mockOnTaskChange).toHaveBeenCalled();
    });
  });

  test('cancels edit mode when cancel button is clicked', () => {
    const props: TaskItemProps = { task: mockTask, onTaskChange: mockOnTaskChange };
    render(<TaskItem {...props} />);

    // Enter edit mode
    const title = screen.getByText('Buy groceries');
    fireEvent.click(title);

    // Change the input
    const input = screen.getByDisplayValue('Buy groceries') as HTMLInputElement;
    fireEvent.change(input, { target: { value: 'Changed title' } });

    // Cancel changes
    const cancelButton = screen.getByRole('button', { name: /cancel/i });
    fireEvent.click(cancelButton);

    // Should show original title and exit edit mode
    expect(screen.getByText('Buy groceries')).toBeInTheDocument();
    expect(screen.queryByRole('textbox')).not.toBeInTheDocument();
  });

  test('deletes task when delete button is clicked', async () => {
    mockApi.deleteTask.mockResolvedValue(undefined);

    // Mock window.confirm to return true
    Object.defineProperty(window, 'confirm', {
      writable: true,
      value: jest.fn(() => true),
    });

    const props: TaskItemProps = { task: mockTask, onTaskChange: mockOnTaskChange };
    render(<TaskItem {...props} />);

    const deleteButton = screen.getByRole('button', { name: /delete/i });
    fireEvent.click(deleteButton);

    await waitFor(() => {
      expect(mockApi.deleteTask).toHaveBeenCalledWith(1);
      expect(mockOnTaskChange).toHaveBeenCalled();
    });
  });

  test('prevents saving empty title', async () => {
    const props: TaskItemProps = { task: mockTask, onTaskChange: mockOnTaskChange };
    render(<TaskItem {...props} />);

    // Enter edit mode
    const title = screen.getByText('Buy groceries');
    fireEvent.click(title);

    // Clear the title
    const input = screen.getByDisplayValue('Buy groceries') as HTMLInputElement;
    fireEvent.change(input, { target: { value: '' } });

    // Try to save
    const saveButton = screen.getByRole('button', { name: /save/i });
    fireEvent.click(saveButton);

    // Should not call API and should show error
    expect(mockApi.updateTask).not.toHaveBeenCalled();
    expect(screen.getByText(/title cannot be empty/i)).toBeInTheDocument();
  });

  test('handles API errors gracefully', async () => {
    mockApi.updateTask.mockRejectedValue(new Error('API Error'));

    const props: TaskItemProps = { task: mockTask, onTaskChange: mockOnTaskChange };
    render(<TaskItem {...props} />);

    const checkbox = screen.getByRole('checkbox') as HTMLInputElement;
    fireEvent.click(checkbox);

    await waitFor(() => {
      expect(screen.getByText(/error updating task/i)).toBeInTheDocument();
    });
  });

  test('displays task creation date', () => {
    const props: TaskItemProps = { task: mockTask, onTaskChange: mockOnTaskChange };
    render(<TaskItem {...props} />);

    // Should show some form of date display
    expect(screen.getByText(/2025-09-27/)).toBeInTheDocument();
  });
});