import { describe, it, expect } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithRouter } from '../testUtils';
import { Routes, Route } from 'react-router-dom';
import NotFound from '../../pages/NotFound';

describe('404 Handling Integration', () => {
  it('invalid route shows 404 page', () => {
    renderWithRouter(
      <Routes>
        <Route path="/dashboard" element={<div>Dashboard Page</div>} />
        <Route path="*" element={<NotFound />} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/invalid-route'] }
      }
    );

    expect(screen.getByText('404')).toBeInTheDocument();
    expect(screen.getByText('Page Not Found')).toBeInTheDocument();
  });

  it('404 page displays attempted path', () => {
    renderWithRouter(
      <Routes>
        <Route path="*" element={<NotFound />} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/nonexistent-page'] }
      }
    );

    expect(screen.getByText(/\/nonexistent-page/)).toBeInTheDocument();
  });

  it('clicking "Go to Dashboard" navigates correctly', async () => {
    const user = userEvent.setup();

    renderWithRouter(
      <Routes>
        <Route path="/dashboard" element={<div>Dashboard Page</div>} />
        <Route path="*" element={<NotFound />} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/invalid'] }
      }
    );

    const dashboardLink = screen.getByText('Go to Dashboard');
    await user.click(dashboardLink);

    expect(screen.getByText('Dashboard Page')).toBeInTheDocument();
    expect(screen.queryByText('404')).not.toBeInTheDocument();
  });

  it('clicking "Go to Home" navigates correctly', async () => {
    const user = userEvent.setup();

    renderWithRouter(
      <Routes>
        <Route path="/" element={<div>Home Page</div>} />
        <Route path="*" element={<NotFound />} />
      </Routes>,
      {
        routerProps: { initialEntries: ['/invalid'] }
      }
    );

    const homeLink = screen.getByText('Go to Home');
    await user.click(homeLink);

    expect(screen.getByText('Home Page')).toBeInTheDocument();
    expect(screen.queryByText('404')).not.toBeInTheDocument();
  });
});