import { describe, it, expect } from 'vitest';
import { screen } from '@testing-library/react';
import { renderWithRouter } from '../../test/testUtils';
import { Routes, Route } from 'react-router-dom';
import MainLayout from './MainLayout';

describe('MainLayout', () => {
  it('renders navigation component', () => {
    renderWithRouter(
      <Routes>
        <Route element={<MainLayout />}>
          <Route path="/" element={<div>Child Content</div>} />
        </Route>
      </Routes>,
      {
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        }
      }
    );

    // Navigation component should render when authenticated
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
  });

  it('renders child routes in outlet', () => {
    renderWithRouter(
      <Routes>
        <Route element={<MainLayout />}>
          <Route path="/" element={<div>Child Content</div>} />
        </Route>
      </Routes>,
      {
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        }
      }
    );

    expect(screen.getByText('Child Content')).toBeInTheDocument();
  });

  it('applies correct CSS classes', () => {
    const { container } = renderWithRouter(
      <Routes>
        <Route element={<MainLayout />}>
          <Route path="/" element={<div>Child Content</div>} />
        </Route>
      </Routes>,
      {
        authContextValue: {
          isAuthenticated: true,
          user: { id: 1, email: 'test@example.com', name: 'Test User' }
        }
      }
    );

    expect(container.querySelector('.main-layout')).toBeInTheDocument();
    expect(container.querySelector('.main-content')).toBeInTheDocument();
  });
});