import React, { useState, useEffect } from 'react';
import { checkHealth } from '../services/api';
import { Link } from 'react-router-dom';
import { ROUTES } from '../routes/routeConfig';

const Dashboard: React.FC = () => {
  const [serverStatus, setServerStatus] = useState('checking');

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

  return (
    <div className="dashboard-page">
      <header className="dashboard-header">
        <h1>Dashboard</h1>
        <div className="server-status">
          <span className={`status-indicator ${serverStatus}`}>
            {serverStatus === 'checking' && '🔄 Checking server...'}
            {serverStatus === 'connected' && '✅ Server Connected'}
            {serverStatus === 'disconnected' && '❌ Server Unavailable'}
          </span>
        </div>
      </header>

      <div className="dashboard-content">
        <div className="dashboard-card">
          <h2>Welcome to Todo App</h2>
          <p>Manage your tasks efficiently with our todo application.</p>

          <div className="quick-actions">
            <Link to={ROUTES.TASKS} className="action-button primary">
              Go to Tasks
            </Link>
            <Link to={ROUTES.PROFILE} className="action-button secondary">
              View Profile
            </Link>
          </div>
        </div>

        {serverStatus === 'disconnected' && (
          <div className="server-error">
            <h3>Server Connection Error</h3>
            <p>Cannot connect to the backend server. Please ensure it's running on port 8080.</p>
            <button onClick={() => window.location.reload()} className="retry-button">
              Retry Connection
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;