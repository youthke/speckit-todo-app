import React, { useState, useEffect } from 'react';
import { getTasks } from '../../services/api';
import TaskItem from '../TaskItem/TaskItem';
import './TaskList.css';

const TaskList = ({ showCompleted, onTaskChange }) => {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadTasks = async () => {
    try {
      setLoading(true);
      setError(null);

      const params = {};
      if (showCompleted !== undefined) {
        params.completed = showCompleted;
      }

      const response = await getTasks(params);
      setTasks(response.tasks);
    } catch (err) {
      setError('Failed to load tasks. Please try again.');
      console.error('Error loading tasks:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTasks();
  }, [showCompleted]);

  const handleTaskChange = () => {
    loadTasks();
    if (onTaskChange) {
      onTaskChange();
    }
  };

  if (loading) {
    return <div className="task-list-loading">Loading tasks...</div>;
  }

  if (error) {
    return (
      <div className="task-list-error">
        <p>Error: {error}</p>
        <button onClick={loadTasks} className="retry-button">
          Retry
        </button>
      </div>
    );
  }

  if (tasks.length === 0) {
    return (
      <div className="task-list-empty">
        <p>No tasks found. Add a new task to get started!</p>
      </div>
    );
  }

  return (
    <div className="task-list">
      <div className="task-list-header">
        <h2>
          {showCompleted === true && 'Completed Tasks'}
          {showCompleted === false && 'Pending Tasks'}
          {showCompleted === undefined && 'All Tasks'}
          ({tasks.length})
        </h2>
      </div>
      <div className="task-list-items">
        {tasks.map((task) => (
          <TaskItem
            key={task.id}
            task={task}
            onTaskChange={handleTaskChange}
          />
        ))}
      </div>
    </div>
  );
};

export default TaskList;