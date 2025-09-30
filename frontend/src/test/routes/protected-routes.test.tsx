import { describe, it, expect } from 'vitest';
import { screen } from '@testing-library/react';
import { renderWithRouter } from '../testUtils';
import { Routes, Route } from 'react-router-dom';
import ProtectedRoute from '../../components/routes/ProtectedRoute';

describe('Protected Routes Integration', () => {
  it('unauthenticated access to /tasks redirects to /login', () => {
    renderWithRouter(
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/tasks" element={<div>Tasks Page</div>} />
        </Route>
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/tasks'] },
        authContextValue: { isAuthenticated: false, isLoading: false }
      }
    );

    expect(screen.getByText('Login Page')).toBeInTheDocument();
    expect(screen.queryByText('Tasks Page')).not.toBeInTheDocument();
  });

  it('unauthenticated access to /dashboard redirects to /login', () => {
    renderWithRouter(
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/dashboard" element={<div>Dashboard Page</div>} />
        </Route>
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/dashboard'] },
        authContextValue: { isAuthenticated: false, isLoading: false }
      }
    );

    expect(screen.getByText('Login Page')).toBeInTheDocument();
    expect(screen.queryByText('Dashboard Page')).not.toBeInTheDocument();
  });

  it('unauthenticated access to /profile redirects to /login', () => {
    renderWithRouter(
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/profile" element={<div>Profile Page</div>} />
        </Route>
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/profile'] },
        authContextValue: { isAuthenticated: false, isLoading: false }
      }
    );

    expect(screen.getByText('Login Page')).toBeInTheDocument();
    expect(screen.queryByText('Profile Page')).not.toBeInTheDocument();
  });

  it('authenticated access to protected routes succeeds', () => {
    renderWithRouter(
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/tasks" element={<div>Tasks Page</div>} />
        </Route>
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/tasks'] },
        authContextValue: {
          isAuthenticated: true,
          isLoading: false,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        }
      }
    );

    expect(screen.getByText('Tasks Page')).toBeInTheDocument();
    expect(screen.queryByText('Login Page')).not.toBeInTheDocument();
  });

  it('loading state prevents premature redirect', () => {
    renderWithRouter(
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/tasks" element={<div>Tasks Page</div>} />
        </Route>
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/tasks'] },
        authContextValue: { isLoading: true, isAuthenticated: false }
      }
    );

    expect(screen.getByText('Loading...')).toBeInTheDocument();
    expect(screen.queryByText('Tasks Page')).not.toBeInTheDocument();
    expect(screen.queryByText('Login Page')).not.toBeInTheDocument();
  });
});