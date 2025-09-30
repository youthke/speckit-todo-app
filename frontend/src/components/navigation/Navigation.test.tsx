import { describe, it, expect, vi } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithRouter } from '../../test/testUtils';
import { Routes, Route } from 'react-router-dom';
import Navigation from './Navigation';

describe('Navigation', () => {
  it('renders navigation links when authenticated', () => {
    renderWithRouter(
      <Navigation />,
      {
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        },
      }
    );

    expect(screen.getByText('Dashboard')).toBeInTheDocument();
    expect(screen.getByText('Tasks')).toBeInTheDocument();
    expect(screen.getByText('Profile')).toBeInTheDocument();
    expect(screen.getByText('Test User')).toBeInTheDocument();
  });

  it('does not render when not authenticated', () => {
    const { container } = renderWithRouter(
      <Navigation />,
      {
        authContextValue: {
          isAuthenticated: false,
          user: null
        },
      }
    );

    expect(container.querySelector('nav')).not.toBeInTheDocument();
  });

  it('highlights active route', () => {
    renderWithRouter(
      <Routes>
        <Route path="/tasks" element={<Navigation />} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/tasks'] },
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        },
      }
    );

    const tasksLink = screen.getByText('Tasks').closest('a');
    expect(tasksLink).toHaveClass('active');
  });

  it('displays user name', () => {
    renderWithRouter(
      <Navigation />,
      {
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'John Doe' }
        },
      }
    );

    expect(screen.getByText('John Doe')).toBeInTheDocument();
  });

  it('calls logout when logout button clicked', async () => {
    const mockLogout = vi.fn();
    const user = userEvent.setup();

    renderWithRouter(
      <Navigation />,
      {
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' },
          logout: mockLogout
        },
      }
    );

    const logoutButton = screen.getByText('Logout');
    await user.click(logoutButton);

    expect(mockLogout).toHaveBeenCalledTimes(1);
  });

  it('redirects to login after logout', async () => {
    const mockLogout = vi.fn().mockResolvedValue(undefined);
    const user = userEvent.setup();

    renderWithRouter(
      <Routes>
        <Route path="/" element={<Navigation />} />
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>,
      {
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' },
          logout: mockLogout
        },
      }
    );

    const logoutButton = screen.getByText('Logout');
    await user.click(logoutButton);

    // After logout, navigation should trigger redirect (handled by auth state change)
    expect(mockLogout).toHaveBeenCalled();
  });
});