import React from 'react';
import { TestI18nProvider } from './i18n-test-utils';
import { BrowserRouter } from 'react-router-dom';

interface MockAuthContextValue {
  isAuthenticated?: boolean;
  isLoading?: boolean;
  logout?: () => void;
}

interface TestProvidersProps {
  children: React.ReactNode;
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
