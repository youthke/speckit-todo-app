import React, { useState, useEffect } from 'react';
import TaskForm from './components/TaskForm/TaskForm';
import TaskList from './components/TaskList/TaskList';
import { checkHealth } from './services/api';
import './App.css';

function App() {
  const [filter, setFilter] = useState('all'); // 'all', 'pending', 'completed'
  const [refreshKey, setRefreshKey] = useState(0);
  const [serverStatus, setServerStatus] = useState('checking');

  // Check server connection on app load
  useEffect(() => {
    const checkServerConnection = async () => {
      try {
        await checkHealth();
        setServerStatus('connected');
      } catch (error) {
        setServerStatus('disconnected');
        console.error('Server connection failed:', error);
      }
    };

    checkServerConnection();
  }, []);

  const handleTaskChange = () => {
    // Force re-render of task list by updating refresh key
    setRefreshKey(prev => prev + 1);
  };

  const getFilterValue = () => {
    if (filter === 'pending') return false;
    if (filter === 'completed') return true;
    return undefined; // 'all'
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>TODO App</h1>
        <div className="server-status">
          <span className={`status-indicator ${serverStatus}`}>
            {serverStatus === 'checking' && '🔄 Checking server...'}
            {serverStatus === 'connected' && '✅ Connected'}
            {serverStatus === 'disconnected' && '❌ Server unavailable'}
          </span>
        </div>
      </header>

      <main className="App-main">
        {serverStatus === 'disconnected' && (
          <div className="server-error">
            <h2>Server Connection Error</h2>
            <p>Cannot connect to the TODO server. Please ensure the backend server is running on port 8080.</p>
            <button onClick={() => window.location.reload()} className="retry-button">
              Retry Connection
            </button>
          </div>
        )}

        {serverStatus === 'connected' && (
          <>
            <TaskForm onTaskCreated={handleTaskChange} />

            <div className="filter-section">
              <h2>Filter Tasks</h2>
              <div className="filter-buttons">
                <button
                  className={filter === 'all' ? 'active' : ''}
                  onClick={() => setFilter('all')}
                >
                  All Tasks
                </button>
                <button
                  className={filter === 'pending' ? 'active' : ''}
                  onClick={() => setFilter('pending')}
                >
                  Pending
                </button>
                <button
                  className={filter === 'completed' ? 'active' : ''}
                  onClick={() => setFilter('completed')}
                >
                  Completed
                </button>
              </div>
            </div>

            <TaskList
              key={refreshKey}
              showCompleted={getFilterValue()}
              onTaskChange={handleTaskChange}
            />
          </>
        )}
      </main>

      <footer className="App-footer">
        <p>TODO App - Built with React and Go</p>
      </footer>
    </div>
  );
}

export default App;
