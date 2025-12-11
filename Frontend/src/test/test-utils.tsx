import React from 'react';
import { render, RenderOptions, act } from '@testing-library/react';
import { TestI18nProvider } from './i18n-test-utils';
import { BrowserRouter } from 'react-router-dom';
import { vi } from 'vitest';

interface MockAuthContextValue {
  isAuthenticated?: boolean;
  isLoading?: boolean;
  logout?: () => void;
}

interface TestProvidersProps {
  children: React.ReactNode;
  authValue?: MockAuthContextValue;
}

// Custom render function that wraps components with necessary providers
interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  withRouter?: boolean;
  withAuth?: boolean;
  withI18n?: boolean;
  authValue?: MockAuthContextValue;
}

// Test wrapper component that provides all necessary context
export const TestProviders: React.FC<TestProvidersProps> = ({ 
  children, 
  authValue = { isAuthenticated: true, isLoading: false, logout: () => {} }
}) => {
  return (
    <TestI18nProvider>
      <BrowserRouter>
        <MockAuthProvider value={authValue}>
          {children}
        </MockAuthProvider>
      </BrowserRouter>
    </TestI18nProvider>
  );
};

// Mock AuthProvider component for testing
const MockAuthProvider: React.FC<{ children: React.ReactNode; value: MockAuthContextValue }> = ({ 
  children 
}) => {
  // The actual auth context mocking should be done in individual test files using vi.mock
  return <>{children}</>;
};

// Wrapper component for tests that includes common providers
const createWrapper = (
  withRouter: boolean = true, 
  withAuth: boolean = true, 
  withI18n: boolean = true,
  authValue: MockAuthContextValue = { isAuthenticated: true, isLoading: false, logout: () => {} }
) => {
  return ({ children }: { children: React.ReactNode }) => {
    let content = children;

    if (withI18n) {
      content = <TestI18nProvider>{content}</TestI18nProvider>;
    }

    if (withRouter) {
      content = <BrowserRouter>{content}</BrowserRouter>;
    }

    if (withAuth) {
      content = <MockAuthProvider value={authValue}>{content}</MockAuthProvider>;
    }

    return <>{content}</>;
  };
};

export const renderWithProviders = (
  ui: React.ReactElement,
  { 
    withRouter = true, 
    withAuth = true, 
    withI18n = true,
    authValue = { isAuthenticated: true, isLoading: false, logout: () => {} },
    ...options 
  }: CustomRenderOptions = {}
) => {
  const Wrapper = createWrapper(withRouter, withAuth, withI18n, authValue);
  
  return render(ui, {
    wrapper: Wrapper,
    ...options,
  });
};

// Utility to wrap state updates in act
export const actAsync = async (fn: () => Promise<void> | void) => {
  await act(async () => {
    await fn();
  });
};

// Utility to handle async renders with act
export const renderWithActAsync = async (
  ui: React.ReactElement,
  options: CustomRenderOptions = {}
) => {
  let result: ReturnType<typeof render>;
  
  await act(async () => {
    result = renderWithProviders(ui, options);
  });
  
  return result!;
};

// Mock implementations for common hooks
export const mockUseAuth = () => ({
  isAuthenticated: true,
  accessToken: 'mock-token',
  refreshToken: 'mock-refresh-token',
  login: vi.fn(),
  logout: vi.fn(),
  user: {
    ID: 'user-123',
    Username: 'testuser',
    Email: 'test@example.com',
    FirstName: 'Test',
    LastName: 'User',
  },
  api: {
    get: vi.fn().mockResolvedValue({ data: [] }),
    post: vi.fn().mockResolvedValue({ data: {} }),
    put: vi.fn().mockResolvedValue({ data: {} }),
    delete: vi.fn().mockResolvedValue({ data: {} }),
  },
});

export const mockUseCurrentUser = () => ({
  user: {
    ID: 'user-123',
    Username: 'testuser',
    Email: 'test@example.com',
    FirstName: 'Test',
    LastName: 'User',
  },
  loading: false,
  error: null,
  refetch: vi.fn(),
});

// Utility to wait for async operations in tests
export const waitForAsyncUpdates = async () => {
  await act(async () => {
    await new Promise(resolve => setTimeout(resolve, 0));
  });
};

// Re-export everything from testing library for convenience
export * from '@testing-library/react';
export { act };
