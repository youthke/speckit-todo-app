import React, { useState } from 'react';
import { createTask } from '../../services/api';
import { TaskFormProps, Task } from '../../types';
import './TaskForm.css';

const TaskForm: React.FC<TaskFormProps> = ({ onTaskCreated }) => {
  const [title, setTitle] = useState<string>('');
  const [isCreating, setIsCreating] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    await createNewTask();
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>): void => {
    if (e.key === 'Enter') {
      e.preventDefault();
      createNewTask();
    }
  };

  const createNewTask = async (): Promise<void> => {
    const trimmedTitle = title.trim();

    // Validation
    if (!trimmedTitle) {
      setError('Title cannot be empty');
      return;
    }

    if (trimmedTitle.length > 500) {
      setError('Title must be 500 characters or less');
      return;
    }

    try {
      setIsCreating(true);
      setError(null);

      const newTask: Task = await createTask({ title: trimmedTitle });
      setTitle(''); // Clear the form

      if (onTaskCreated) {
        onTaskCreated(newTask);
      }
    } catch (err) {
      setError('Error creating task. Please try again.');
      console.error('Error creating task:', err);
    } finally {
      setIsCreating(false);
    }
  };

  const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    setTitle(e.target.value);
    // Clear error when user starts typing
    if (error) {
      setError(null);
    }
  };

  return (
    <div className="task-form">
      <form onSubmit={handleSubmit} className="task-form-container">
        <div className="task-form-input-group">
          <input
            type="text"
            value={title}
            onChange={handleTitleChange}
            onKeyPress={handleKeyPress}
            placeholder="Enter task title..."
            className={`task-form-input ${error ? 'error' : ''}`}
            disabled={isCreating}
            maxLength={500}
          />
          <button
            type="submit"
            disabled={isCreating || !title.trim()}
            className="task-form-submit"
          >
            {isCreating ? 'Creating...' : 'Add Task'}
          </button>
        </div>

        {error && (
          <div className="task-form-error">
            {error}
          </div>
        )}

        <div className="task-form-meta">
          <span className="character-count">
            {title.length}/500 characters
          </span>
          <span className="hint">
            Press Enter to add task
          </span>
        </div>
      </form>
    </div>
  );
};

export default TaskForm;