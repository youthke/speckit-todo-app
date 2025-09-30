import { useState, useEffect, useCallback, useContext, createContext } from 'react';
import authService, { User, SessionResponse } from '../services/auth';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (redirectUri?: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshSession: () => Promise<void>;
  validateSession: () => Promise<void>;
}

// Create Auth Context
export const AuthContext = createContext<AuthContextType | undefined>(undefined);

/**
 * Custom hook to access auth context
 */
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

/**
 * Hook for session management and authentication state
 */
export const useAuthState = () => {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  /**
   * Validate and load current session
   */
  const validateSession = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response: SessionResponse = await authService.validateSession();
      setUser(response.user);
      setIsAuthenticated(true);
    } catch (err) {
      setUser(null);
      setIsAuthenticated(false);
      setError(err instanceof Error ? err.message : 'Session validation failed');
    } finally {
      setIsLoading(false);
    }
  }, []);

  /**
   * Initiate Google OAuth login
   */
  const login = useCallback(async (redirectUri?: string) => {
    try {
      setError(null);
      const { auth_url } = await authService.initiateGoogleLogin(redirectUri);
      // Redirect to Google OAuth
      window.location.href = auth_url;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
      throw err;
    }
  }, []);

  /**
   * Logout and clear session
   */
  const logout = useCallback(async () => {
    try {
      setError(null);
      await authService.logout();
      setUser(null);
      setIsAuthenticated(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Logout failed');
      // Clear local state even if logout request fails
      setUser(null);
      setIsAuthenticated(false);
    }
  }, []);

  /**
   * Refresh OAuth tokens
   */
  const refreshSession = useCallback(async () => {
    try {
      setError(null);
      await authService.refreshSession();
      // Revalidate session after refresh
      await validateSession();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Session refresh failed');
      // If refresh fails, clear session
      setUser(null);
      setIsAuthenticated(false);
      throw err;
    }
  }, [validateSession]);

  /**
   * Auto-validate session on mount
   */
  useEffect(() => {
    validateSession();
  }, [validateSession]);

  /**
   * Auto-refresh session before expiration
   */
  useEffect(() => {
    if (!isAuthenticated || !user) return;

    // Set up interval to check session validity every 5 minutes
    const intervalId = setInterval(() => {
      validateSession();
    }, 5 * 60 * 1000); // 5 minutes

    return () => clearInterval(intervalId);
  }, [isAuthenticated, user, validateSession]);

  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    login,
    logout,
    refreshSession,
    validateSession,
  };
};

/**
 * Hook to check if user is authenticated
 */
export const useIsAuthenticated = (): boolean => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated;
};

/**
 * Hook to get current user
 */
export const useCurrentUser = (): User | null => {
  const { user } = useAuth();
  return user;
};

/**
 * Hook for protected routes
 */
export const useRequireAuth = (redirectTo: string = '/login') => {
  const { isAuthenticated, isLoading } = useAuth();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      // Store intended destination
      authService.setOAuthRedirect(window.location.pathname);
      // Redirect to login
      window.location.href = redirectTo;
    }
  }, [isAuthenticated, isLoading, redirectTo]);

  return { isAuthenticated, isLoading };
};

/**
 * Hook for automatic token refresh
 */
export const useAutoRefresh = (enabled: boolean = true) => {
  const { refreshSession, isAuthenticated } = useAuth();

  useEffect(() => {
    if (!enabled || !isAuthenticated) return;

    // Refresh token every 23 hours (just before 24-hour expiration)
    const refreshInterval = 23 * 60 * 60 * 1000; // 23 hours
    const intervalId = setInterval(() => {
      refreshSession().catch((err) => {
        console.error('Auto-refresh failed:', err);
      });
    }, refreshInterval);

    return () => clearInterval(intervalId);
  }, [enabled, isAuthenticated, refreshSession]);
};