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

vi.mock('../../../hooks/useTranslation', () => ({
  useT: () => ({
    t: (key: string) => key
  })
}));

vi.mock('../../../components/layout/Layout', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div data-testid="layout">{children}</div>
}));

vi.mock('../ProfileSidebar', () => ({
  default: () => <div data-testid="profile-sidebar">Sidebar</div>
}));

describe('ProfileShifts', () => {
  it('renders the ProfileShifts component', () => {
    render(
      <BrowserRouter>
        <ProfileShifts />
      </BrowserRouter>
    );

    expect(screen.getByTestId('layout')).toBeInTheDocument();
    expect(screen.getByTestId('profile-sidebar')).toBeInTheDocument();
    expect(screen.getByText('My Future Shifts')).toBeInTheDocument();
  });

  it('shows loading state initially', () => {
    render(
      <BrowserRouter>
        <ProfileShifts />
      </BrowserRouter>
    );

    expect(screen.getByText('Loading shifts...')).toBeInTheDocument();
  });
});