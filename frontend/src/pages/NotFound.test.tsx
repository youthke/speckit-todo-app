import { describe, it, expect } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithRouter } from '../test/testUtils';
import { Routes, Route } from 'react-router-dom';
import NotFound from './NotFound';

describe('NotFound', () => {
  it('displays 404 error message', () => {
    renderWithRouter(<NotFound />);

    expect(screen.getByText('404')).toBeInTheDocument();
    expect(screen.getByText('Page Not Found')).toBeInTheDocument();
  });

  it('shows attempted path', () => {
    renderWithRouter(
      <Routes>
        <Route path="*" element={<NotFound />} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/invalid-route'] }
      }
    );

    expect(screen.getByText(/\/invalid-route/)).toBeInTheDocument();
  });

  it('provides link to dashboard', async () => {
    const user = userEvent.setup();

    renderWithRouter(
      <Routes>
        <Route path="*" element={<NotFound />} />
        <Route path="/dashboard" element={<div>Dashboard Page</div>} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/invalid'] }
      }
    );

    const dashboardLink = screen.getByText('Go to Dashboard');
    await user.click(dashboardLink);

    expect(screen.getByText('Dashboard Page')).toBeInTheDocument();
  });

  it('provides link to home', async () => {
    const user = userEvent.setup();

    renderWithRouter(
      <Routes>
        <Route path="*" element={<NotFound />} />
        <Route path="/" element={<div>Home Page</div>} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/invalid'] }
      }
    );

    const homeLink = screen.getByText('Go to Home');
    await user.click(homeLink);

    expect(screen.getByText('Home Page')).toBeInTheDocument();
  });
});