import { describe, it, expect } from 'vitest';
import { screen } from '@testing-library/react';
import { renderWithRouter } from '../../test/testUtils';
import { Route, Routes } from 'react-router-dom';
import ProtectedRoute from './ProtectedRoute';

describe('ProtectedRoute', () => {
  it('redirects to login when not authenticated', () => {
    renderWithRouter(
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<div>Protected Content</div>} />
        </Route>
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      {
        authContextValue: { isAuthenticated: false, isLoading: false },
      }
    );

    expect(screen.getByText('Login Page')).toBeInTheDocument();
    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
  });

  it('renders child routes when authenticated', () => {
    renderWithRouter(
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<div>Protected Content</div>} />
        </Route>
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      {
        authContextValue: {
          isAuthenticated: true,
          isLoading: false,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        },
      }
    );

    expect(screen.getByText('Protected Content')).toBeInTheDocument();
    expect(screen.queryByText('Login Page')).not.toBeInTheDocument();
  });

  it('shows loading state while checking auth', () => {
    renderWithRouter(
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<div>Protected Content</div>} />
        </Route>
      </Routes>,
      {
        authContextValue: { isLoading: true, isAuthenticated: false },
      }
    );

    expect(screen.getByText('Loading...')).toBeInTheDocument();
    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
  });

  it('uses custom redirect path when provided', () => {
    renderWithRouter(
      <Routes>
        <Route element={<ProtectedRoute redirectTo="/custom-login" />}>
          <Route path="/" element={<div>Protected Content</div>} />
        </Route>
        <Route path="/custom-login" element={<div>Custom Login</div>} />
      </Routes>,
      {
        authContextValue: { isAuthenticated: false, isLoading: false },
      }
    );

    expect(screen.getByText('Custom Login')).toBeInTheDocument();
  });
});