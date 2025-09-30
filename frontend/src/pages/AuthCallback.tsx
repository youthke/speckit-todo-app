import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

interface AuthCallbackProps {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
  defaultRedirect?: string;
}

const AuthCallback: React.FC<AuthCallbackProps> = ({
  onSuccess,
  onError,
  defaultRedirect = '/dashboard',
}) => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(true);

  useEffect(() => {
    const processCallback = async () => {
      try {
        // Extract query parameters from URL
        const code = searchParams.get('code');
        const state = searchParams.get('state');
        const errorParam = searchParams.get('error');

        // Handle OAuth error from Google
        if (errorParam) {
          const errorMessage = getErrorMessage(errorParam);
          throw new Error(errorMessage);
        }

        // Validate required parameters
        if (!code || !state) {
          throw new Error('Missing authorization code or state parameter');
        }

        // Send callback request to backend
        const response = await fetch(
          `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/auth/google/callback?code=${encodeURIComponent(code)}&state=${encodeURIComponent(state)}`,
          {
            method: 'GET',
            credentials: 'include', // Include cookies for session
            headers: {
              'Accept': 'application/json',
            },
          }
        );

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || 'Authentication failed');
        }

        // Authentication successful
        onSuccess?.();

        // Redirect to intended destination
        const redirectTo = sessionStorage.getItem('oauth_redirect') || defaultRedirect;
        sessionStorage.removeItem('oauth_redirect');
        navigate(redirectTo, { replace: true });
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Authentication failed';
        setError(errorMessage);
        onError?.(err instanceof Error ? err : new Error(errorMessage));
        setIsProcessing(false);

        // Redirect to login page after 3 seconds
        setTimeout(() => {
          navigate('/login', { replace: true });
        }, 3000);
      }
    };

    processCallback();
  }, [searchParams, navigate, onSuccess, onError, defaultRedirect]);

  const getErrorMessage = (error: string): string => {
    const errorMessages: Record<string, string> = {
      access_denied: 'You denied access to your Google account. Please try again.',
      server_error: 'Google encountered a server error. Please try again later.',
      temporarily_unavailable: 'Google service is temporarily unavailable. Please try again later.',
    };

    return errorMessages[error] || `Authentication error: ${error}`;
  };

  return (
    <div className="auth-callback-container">
      {isProcessing ? (
        <div className="auth-callback-loading">
          <div className="spinner"></div>
          <p>Completing authentication...</p>
        </div>
      ) : (
        <div className="auth-callback-error">
          <div className="error-icon">⚠️</div>
          <h2>Authentication Failed</h2>
          <p>{error}</p>
          <p className="redirect-message">Redirecting to login page...</p>
        </div>
      )}
    </div>
  );
};

export default AuthCallback;