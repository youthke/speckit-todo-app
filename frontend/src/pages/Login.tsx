import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import GoogleLoginButton from '../components/auth/GoogleLoginButton';
import { useAuth } from '../hooks/useAuth';

const Login: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated, isLoading } = useAuth();

  // Redirect to dashboard if already authenticated
  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      navigate('/dashboard', { replace: true });
    }
  }, [isAuthenticated, isLoading, navigate]);

  const handleLoginSuccess = () => {
    // GoogleLoginButton will redirect to Google OAuth
    // After callback, AuthCallback component will handle redirect
  };

  const handleLoginError = (error: Error) => {
    console.error('Login error:', error);
    // Could show toast notification here
  };

  if (isLoading) {
    return (
      <div className="login-page loading">
        <div className="spinner"></div>
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="login-page">
      <div className="login-container">
        <div className="login-header">
          <h1>Welcome to Todo App</h1>
          <p>Sign in to manage your tasks</p>
        </div>

        <div className="login-content">
          <GoogleLoginButton
            onSuccess={handleLoginSuccess}
            onError={handleLoginError}
            className="login-google-button"
          />

          <div className="login-divider">
            <span>or</span>
          </div>

          <form className="login-form">
            <div className="form-group">
              <label htmlFor="email">Email</label>
              <input
                type="email"
                id="email"
                placeholder="Enter your email"
                autoComplete="email"
              />
            </div>

            <div className="form-group">
              <label htmlFor="password">Password</label>
              <input
                type="password"
                id="password"
                placeholder="Enter your password"
                autoComplete="current-password"
              />
            </div>

            <button type="submit" className="login-button">
              Sign In
            </button>
          </form>

          <div className="login-footer">
            <p>Don't have an account? <a href="/signup">Sign up</a></p>
          </div>
        </div>
      </div>

      <style>{`
        .login-page {
          display: flex;
          justify-content: center;
          align-items: center;
          min-height: 100vh;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          padding: 20px;
        }

        .login-page.loading {
          flex-direction: column;
          color: white;
        }

        .login-container {
          background: white;
          border-radius: 12px;
          box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
          padding: 40px;
          max-width: 400px;
          width: 100%;
        }

        .login-header {
          text-align: center;
          margin-bottom: 32px;
        }

        .login-header h1 {
          font-size: 28px;
          font-weight: 700;
          color: #1a202c;
          margin: 0 0 8px 0;
        }

        .login-header p {
          font-size: 14px;
          color: #718096;
          margin: 0;
        }

        .login-content {
          display: flex;
          flex-direction: column;
          gap: 24px;
        }

        .google-login-button {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 12px;
          width: 100%;
          padding: 12px 24px;
          background: white;
          border: 2px solid #e2e8f0;
          border-radius: 8px;
          font-size: 16px;
          font-weight: 500;
          color: #1a202c;
          cursor: pointer;
          transition: all 0.2s;
        }

        .google-login-button:hover {
          background: #f7fafc;
          border-color: #cbd5e0;
        }

        .google-login-button:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }

        .google-icon {
          width: 18px;
          height: 18px;
        }

        .login-divider {
          display: flex;
          align-items: center;
          text-align: center;
          color: #a0aec0;
          font-size: 14px;
        }

        .login-divider::before,
        .login-divider::after {
          content: '';
          flex: 1;
          border-bottom: 1px solid #e2e8f0;
        }

        .login-divider span {
          padding: 0 16px;
        }

        .login-form {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .form-group {
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .form-group label {
          font-size: 14px;
          font-weight: 500;
          color: #4a5568;
        }

        .form-group input {
          padding: 12px;
          border: 2px solid #e2e8f0;
          border-radius: 8px;
          font-size: 14px;
          transition: border-color 0.2s;
        }

        .form-group input:focus {
          outline: none;
          border-color: #667eea;
        }

        .login-button {
          padding: 12px 24px;
          background: #667eea;
          color: white;
          border: none;
          border-radius: 8px;
          font-size: 16px;
          font-weight: 500;
          cursor: pointer;
          transition: background 0.2s;
        }

        .login-button:hover {
          background: #5568d3;
        }

        .login-footer {
          text-align: center;
          font-size: 14px;
          color: #718096;
        }

        .login-footer a {
          color: #667eea;
          text-decoration: none;
          font-weight: 500;
        }

        .login-footer a:hover {
          text-decoration: underline;
        }

        .spinner {
          border: 4px solid rgba(255, 255, 255, 0.3);
          border-radius: 50%;
          border-top: 4px solid white;
          width: 40px;
          height: 40px;
          animation: spin 1s linear infinite;
        }

        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }

        .google-login-error {
          color: #e53e3e;
          font-size: 14px;
          margin-top: 8px;
          text-align: center;
        }
      `}</style>
    </div>
  );
};

export default Login;