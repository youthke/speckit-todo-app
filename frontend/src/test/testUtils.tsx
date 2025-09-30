import { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { MemoryRouter, MemoryRouterProps } from 'react-router-dom';
import { AuthContext } from '../hooks/useAuth';
import { vi } from 'vitest';

interface AuthContextType {
  user: { id: number; email: string; name: string } | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (redirectUri?: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshSession: () => Promise<void>;
  validateSession: () => Promise<void>;
}

interface ExtendedRenderOptions extends RenderOptions {
  routerProps?: MemoryRouterProps;
  authContextValue?: Partial<AuthContextType>;
}

export const createMockAuthContext = (
  overrides?: Partial<AuthContextType>
): AuthContextType => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
  login: vi.fn(),
  logout: vi.fn(),
  refreshSession: vi.fn(),
  validateSession: vi.fn(),
  ...overrides,
});

export const renderWithRouter = (
  ui: ReactElement,
  {
    routerProps = { initialEntries: ['/'] },
    authContextValue,
    ...renderOptions
  }: ExtendedRenderOptions = {}
) => {
  const mockAuthContext = createMockAuthContext(authContextValue);

  const Wrapper = ({ children }: { children: React.ReactNode }) => (
    <MemoryRouter {...routerProps}>
      <AuthContext.Provider value={mockAuthContext}>
        {children}
      </AuthContext.Provider>
    </MemoryRouter>
  );

  return {
    ...render(ui, { wrapper: Wrapper, ...renderOptions }),
    mockAuthContext,
  };
};