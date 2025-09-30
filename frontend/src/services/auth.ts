const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface User {
  id: number;
  email: string;
  name: string;
  oauth_provider?: string;
  created_at: string;
}

export interface SessionInfo {
  session_id: string;
  expires_at: string;
  last_activity: string;
}

export interface SessionResponse {
  user: User;
  session: SessionInfo;
}

export interface AuthError {
  error: string;
  message: string;
  details?: Record<string, unknown>;
}

class AuthService {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  /**
   * Initiate Google OAuth login flow
   */
  async initiateGoogleLogin(redirectUri?: string): Promise<{ auth_url: string }> {
    const params = new URLSearchParams();
    if (redirectUri) {
      params.append('redirect_uri', redirectUri);
    }

    const response = await fetch(`${this.baseUrl}/api/v1/auth/google/login?${params.toString()}`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData: AuthError = await response.json();
      throw new Error(errorData.message || 'Failed to initiate Google login');
    }

    return response.json();
  }

  /**
   * Validate current session
   */
  async validateSession(): Promise<SessionResponse> {
    const response = await fetch(`${this.baseUrl}/api/v1/auth/session/validate`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      if (response.status === 401) {
        throw new Error('Session expired or invalid');
      }
      const errorData: AuthError = await response.json();
      throw new Error(errorData.message || 'Session validation failed');
    }

    return response.json();
  }

  /**
   * Refresh OAuth tokens
   */
  async refreshSession(): Promise<{ status: string; expires_at: string }> {
    const response = await fetch(`${this.baseUrl}/api/v1/auth/session/refresh`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      if (response.status === 401) {
        throw new Error('Session expired or invalid');
      }
      const errorData: AuthError = await response.json();
      throw new Error(errorData.message || 'Token refresh failed');
    }

    return response.json();
  }

  /**
   * Logout user and terminate session
   */
  async logout(): Promise<{ status: string; message: string }> {
    const response = await fetch(`${this.baseUrl}/api/v1/auth/logout`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData: AuthError = await response.json();
      throw new Error(errorData.message || 'Logout failed');
    }

    return response.json();
  }

  /**
   * Check if user is authenticated
   */
  async isAuthenticated(): Promise<boolean> {
    try {
      await this.validateSession();
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Get current user information
   */
  async getCurrentUser(): Promise<User | null> {
    try {
      const { user } = await this.validateSession();
      return user;
    } catch {
      return null;
    }
  }

  /**
   * Set redirect URI in session storage for OAuth flow
   */
  setOAuthRedirect(redirectUri: string): void {
    sessionStorage.setItem('oauth_redirect', redirectUri);
  }

  /**
   * Get and clear OAuth redirect URI from session storage
   */
  getOAuthRedirect(): string | null {
    const redirect = sessionStorage.getItem('oauth_redirect');
    if (redirect) {
      sessionStorage.removeItem('oauth_redirect');
    }
    return redirect;
  }

  /**
   * Clear OAuth redirect URI from session storage
   */
  clearOAuthRedirect(): void {
    sessionStorage.removeItem('oauth_redirect');
  }
}

// Export singleton instance
const authService = new AuthService();
export default authService;