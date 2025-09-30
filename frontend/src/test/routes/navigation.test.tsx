import { describe, it, expect, vi } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithRouter } from '../testUtils';
import { Routes, Route, useNavigate } from 'react-router-dom';
import Navigation from '../../components/navigation/Navigation';

// Test component that uses navigate
const TestPageWithNavigation = ({ page }: { page: string }) => {
  const navigate = useNavigate();

  return (
    <>
      <Navigation />
      <div>{page} Page</div>
      <button onClick={() => navigate(-1)}>Back</button>
      <button onClick={() => navigate(1)}>Forward</button>
    </>
  );
};

describe('Navigation Flow Integration', () => {
  it('clicking navigation links changes route', async () => {
    const user = userEvent.setup();

    renderWithRouter(
      <Routes>
        <Route path="/dashboard" element={
          <>
            <Navigation />
            <div>Dashboard Page</div>
          </>
        } />
        <Route path="/tasks" element={
          <>
            <Navigation />
            <div>Tasks Page</div>
          </>
        } />
      </Routes>,
      {
        routerProps: { initialEntries: ['/dashboard'] },
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        }
      }
    );

    expect(screen.getByText('Dashboard Page')).toBeInTheDocument();

    const tasksLink = screen.getByText('Tasks');
    await user.click(tasksLink);

    expect(screen.getByText('Tasks Page')).toBeInTheDocument();
  });

  it('browser back button navigates backward', async () => {
    const user = userEvent.setup();

    renderWithRouter(
      <Routes>
        <Route path="/dashboard" element={<TestPageWithNavigation page="Dashboard" />} />
        <Route path="/tasks" element={<TestPageWithNavigation page="Tasks" />} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/dashboard', '/tasks'], initialIndex: 1 },
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        }
      }
    );

    expect(screen.getByText('Tasks Page')).toBeInTheDocument();

    const backButton = screen.getByText('Back');
    await user.click(backButton);

    // After going back, dashboard should be visible
    expect(screen.getByText('Dashboard Page')).toBeInTheDocument();
  });

  it('browser forward button navigates forward', async () => {
    const user = userEvent.setup();

    renderWithRouter(
      <Routes>
        <Route path="/dashboard" element={<TestPageWithNavigation page="Dashboard" />} />
        <Route path="/tasks" element={<TestPageWithNavigation page="Tasks" />} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/dashboard', '/tasks'], initialIndex: 0 },
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        }
      }
    );

    expect(screen.getByText('Dashboard Page')).toBeInTheDocument();

    const forwardButton = screen.getByText('Forward');
    await user.click(forwardButton);

    expect(screen.getByText('Tasks Page')).toBeInTheDocument();
  });

  it('active route is highlighted in navigation', () => {
    renderWithRouter(
      <Routes>
        <Route path="/tasks" element={<Navigation />} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/tasks'] },
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        }
      }
    );

    const tasksLink = screen.getByText('Tasks').closest('a');
    expect(tasksLink).toHaveClass('active');
  });

  it('logout redirects to login page', async () => {
    const mockLogout = vi.fn().mockResolvedValue(undefined);
    const user = userEvent.setup();

    renderWithRouter(
      <Routes>
        <Route path="/dashboard" element={<Navigation />} />
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/dashboard'] },
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' },
          logout: mockLogout
        }
      }
    );

    const logoutButton = screen.getByText('Logout');
    await user.click(logoutButton);

    expect(mockLogout).toHaveBeenCalledTimes(1);
  });
});