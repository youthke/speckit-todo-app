import React, { useState } from 'react';

interface GoogleLoginButtonProps {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
  redirectUri?: string;
  className?: string;
  disabled?: boolean;
}

const GoogleLoginButton: React.FC<GoogleLoginButtonProps> = ({
  onSuccess,
  onError,
  redirectUri,
  className = '',
  disabled = false,
}) => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGoogleLogin = async () => {
    setIsLoading(true);
    setError(null);

    try {
      // Build query parameters
      const params = new URLSearchParams();
      if (redirectUri) {
        params.append('redirect_uri', redirectUri);
      }

      // Call backend to initiate OAuth flow
      const response = await fetch(
        `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/auth/google/login?${params.toString()}`,
        {
          method: 'GET',
          credentials: 'include', // Include cookies
          headers: {
            'Accept': 'application/json',
          },
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to initiate Google login');
      }

      const data = await response.json();

      // Redirect to Google OAuth authorization URL
      if (data.auth_url) {
        window.location.href = data.auth_url;
        onSuccess?.();
      } else {
        throw new Error('No authorization URL received');
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
      setError(errorMessage);
      onError?.(err instanceof Error ? err : new Error(errorMessage));
      setIsLoading(false);
    }
  };

  return (
    <div className="google-login-container">
      <button
        onClick={handleGoogleLogin}
        disabled={disabled || isLoading}
        className={`google-login-button ${className}`}
        type="button"
      >
        {isLoading ? (
          <span className="loading-spinner">Loading...</span>
        ) : (
          <>
            <svg className="google-icon" viewBox="0 0 24 24" width="18" height="18">
              <path
                fill="#4285F4"
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
              />
              <path
                fill="#34A853"
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
              />
              <path
                fill="#FBBC05"
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
              />
              <path
                fill="#EA4335"
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
              />
            </svg>
            <span>Sign in with Google</span>
          </>
        )}
      </button>
      {error && (
        <div className="google-login-error" role="alert">
          {error}
        </div>
      )}
    </div>
  );
};

export default GoogleLoginButton;