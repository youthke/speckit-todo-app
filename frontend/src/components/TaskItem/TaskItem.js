import React, { useState } from 'react';
import { updateTask, deleteTask } from '../../services/api';
import './TaskItem.css';

const TaskItem = ({ task, onTaskChange }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [editTitle, setEditTitle] = useState(task.title);
  const [isUpdating, setIsUpdating] = useState(false);
  const [error, setError] = useState(null);

  const handleToggleComplete = async () => {
    try {
      setIsUpdating(true);
      setError(null);
      await updateTask(task.id, { completed: !task.completed });
      onTaskChange();
    } catch (err) {
      setError('Error updating task. Please try again.');
      console.error('Error toggling task completion:', err);
    } finally {
      setIsUpdating(false);
    }
  };

  const handleStartEdit = () => {
    setIsEditing(true);
    setEditTitle(task.title);
    setError(null);
  };

  const handleSaveEdit = async () => {
    const trimmedTitle = editTitle.trim();

    if (!trimmedTitle) {
      setError('Title cannot be empty');
      return;
    }

    if (trimmedTitle.length > 500) {
      setError('Title must be 500 characters or less');
      return;
    }

    try {
      setIsUpdating(true);
      setError(null);
      await updateTask(task.id, { title: trimmedTitle });
      setIsEditing(false);
      onTaskChange();
    } catch (err) {
      setError('Error updating task. Please try again.');
      console.error('Error updating task title:', err);
    } finally {
      setIsUpdating(false);
    }
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
    setEditTitle(task.title);
    setError(null);
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this task?')) {
      try {
        setIsUpdating(true);
        setError(null);
        await deleteTask(task.id);
        onTaskChange();
      } catch (err) {
        setError('Error deleting task. Please try again.');
        console.error('Error deleting task:', err);
        setIsUpdating(false);
      }
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleSaveEdit();
    } else if (e.key === 'Escape') {
      handleCancelEdit();
    }
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className={`task-item ${task.completed ? 'completed' : ''} ${isUpdating ? 'updating' : ''}`}>
      <div className="task-item-main">
        <input
          type="checkbox"
          checked={task.completed}
          onChange={handleToggleComplete}
          disabled={isUpdating || isEditing}
          className="task-checkbox"
        />

        {isEditing ? (
          <div className="task-edit-section">
            <input
              type="text"
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
              onKeyDown={handleKeyPress}
              className="task-edit-input"
              autoFocus
              maxLength={500}
            />
            <div className="task-edit-buttons">
              <button
                onClick={handleSaveEdit}
                disabled={isUpdating}
                className="save-button"
              >
                Save
              </button>
              <button
                onClick={handleCancelEdit}
                disabled={isUpdating}
                className="cancel-button"
              >
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <div className="task-content">
            <span
              className={`task-title ${task.completed ? 'completed' : ''}`}
              onClick={handleStartEdit}
              title="Click to edit"
            >
              {task.title}
            </span>
            <div className="task-meta">
              <span className="task-date">
                Created: {formatDate(task.created_at)}
              </span>
              {task.updated_at !== task.created_at && (
                <span className="task-date">
                  Updated: {formatDate(task.updated_at)}
                </span>
              )}
            </div>
          </div>
        )}

        <button
          onClick={handleDelete}
          disabled={isUpdating || isEditing}
          className="delete-button"
          title="Delete task"
        >
          ×
        </button>
      </div>

      {error && (
        <div className="task-error">
          {error}
        </div>
      )}
    </div>
  );
};

export default TaskItem;