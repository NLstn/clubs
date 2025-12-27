import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import ProfileShifts from '../ProfileShifts';

// Mock the hooks
vi.mock('../../../hooks/useAuth', () => ({
  useAuth: () => ({
    api: {
      get: vi.fn().mockRejectedValue(new Error('Mock API error'))
    }
  })
}));

vi.mock('../../../hooks/useCurrentUser', () => ({
  useCurrentUser: () => ({
    user: { ID: 'user-123', Email: 'test@example.com', FirstName: 'Test', LastName: 'User' },
    loading: false,
    error: null
  })
}));

vi.mock('../../../hooks/useTranslation', () => ({
  useT: () => ({
    t: (key: string) => key
  })
}));

vi.mock('../../../components/layout/Layout', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div data-testid="layout">{children}</div>
}));

vi.mock('../../../components/layout/SimpleSettingsLayout', () => ({
  default: ({ title, children }: { title: string; children: React.ReactNode }) => (
    <div data-testid="simple-settings-layout">
      <h1>{title}</h1>
      {children}
    </div>
  )
}));

describe('ProfileShifts', () => {
  it('renders the ProfileShifts component', () => {
    render(
      <BrowserRouter>
        <ProfileShifts />
      </BrowserRouter>
    );

    expect(screen.getByTestId('layout')).toBeInTheDocument();
    expect(screen.getByTestId('simple-settings-layout')).toBeInTheDocument();
    // The mock returns translation keys, so look for the key
    expect(screen.getByText('shifts.myFutureShifts')).toBeInTheDocument();
  });

  it('shows loading state initially', () => {
    render(
      <BrowserRouter>
        <ProfileShifts />
      </BrowserRouter>
    );

    // The mock returns translation keys, so look for the key
    expect(screen.getByText('shifts.loadingShifts')).toBeInTheDocument();
  });
});