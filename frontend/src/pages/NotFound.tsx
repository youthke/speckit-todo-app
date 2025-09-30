import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { ROUTES } from '../routes/routeConfig';

const NotFound: React.FC = () => {
  const location = useLocation();

  return (
    <div className="not-found-page">
      <div className="not-found-container">
        <h1 className="error-code">404</h1>
        <h2 className="error-title">Page Not Found</h2>
        <p className="error-message">
          The page <code>{location.pathname}</code> does not exist.
        </p>
        <div className="error-actions">
          <Link to={ROUTES.DASHBOARD} className="primary-button">
            Go to Dashboard
          </Link>
          <Link to={ROUTES.HOME} className="secondary-button">
            Go to Home
          </Link>
        </div>
      </div>

      <style>{`
        .not-found-page {
          display: flex;
          justify-content: center;
          align-items: center;
          min-height: 100vh;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          padding: 20px;
        }

        .not-found-container {
          background: white;
          border-radius: 12px;
          box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
          padding: 60px 40px;
          max-width: 500px;
          width: 100%;
          text-align: center;
        }

        .error-code {
          font-size: 96px;
          font-weight: 700;
          color: #667eea;
          margin: 0 0 16px 0;
        }

        .error-title {
          font-size: 32px;
          font-weight: 600;
          color: #1a202c;
          margin: 0 0 16px 0;
        }

        .error-message {
          font-size: 16px;
          color: #718096;
          margin: 0 0 32px 0;
          line-height: 1.6;
        }

        .error-message code {
          background: #f7fafc;
          padding: 2px 8px;
          border-radius: 4px;
          font-family: monospace;
          color: #e53e3e;
        }

        .error-actions {
          display: flex;
          gap: 16px;
          justify-content: center;
        }

        .primary-button {
          padding: 12px 24px;
          background: #667eea;
          color: white;
          border: none;
          border-radius: 8px;
          font-size: 16px;
          font-weight: 500;
          text-decoration: none;
          cursor: pointer;
          transition: background 0.2s;
        }

        .primary-button:hover {
          background: #5568d3;
        }

        .secondary-button {
          padding: 12px 24px;
          background: white;
          color: #667eea;
          border: 2px solid #667eea;
          border-radius: 8px;
          font-size: 16px;
          font-weight: 500;
          text-decoration: none;
          cursor: pointer;
          transition: all 0.2s;
        }

        .secondary-button:hover {
          background: #f7fafc;
        }
      `}</style>
    </div>
  );
};

export default NotFound;