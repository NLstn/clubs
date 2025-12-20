import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ProfileNotificationSettings from '../ProfileNotificationSettings';
import * as useNotificationsModule from '../../../hooks/useNotifications';

// Mock the useNotifications hook
vi.mock('../../../hooks/useNotifications');

// Mock the Layout component
vi.mock('../../../components/layout/Layout', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div data-testid="layout">{children}</div>,
}));

// Mock the ProfileContentLayout component
vi.mock('../../../components/layout/ProfileContentLayout', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div data-testid="profile-content-layout">{children}</div>,
}));

// Mock the ToggleSwitch component
vi.mock('@/components/ui', () => ({
  ToggleSwitch: ({ checked, onChange }: { checked: boolean; onChange: (value: boolean) => void }) => (
    <input
      type="checkbox"
      checked={checked}
      onChange={(e) => onChange(e.target.checked)}
      data-testid="toggle-switch"
    />
  ),
}));

// Mock the translation hook
vi.mock('../../../hooks/useTranslation', () => ({
  useT: () => ({
    t: (key: string) => key,
  }),
}));

describe('ProfileNotificationSettings', () => {
  const mockPreferences = {
    ID: 'test-id',
    UserID: 'test-user-id',
    MemberAddedInApp: true,
    MemberAddedEmail: false,
    EventCreatedInApp: true,
    EventCreatedEmail: false,
    FineAssignedInApp: true,
    FineAssignedEmail: true,
    NewsCreatedInApp: false,
    NewsCreatedEmail: false,
    CreatedAt: '2024-01-01T00:00:00Z',
    UpdatedAt: '2024-01-01T00:00:00Z',
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should access preferences with PascalCase field names', async () => {
    const mockUpdatePreferences = vi.fn();
    
    vi.mocked(useNotificationsModule.useNotificationPreferences).mockReturnValue({
      preferences: mockPreferences,
      loading: false,
      error: null,
      fetchPreferences: vi.fn(),
      updatePreferences: mockUpdatePreferences,
    });

    render(
      <MemoryRouter>
        <ProfileNotificationSettings />
      </MemoryRouter>
    );

    await waitFor(() => {
      // Check that the toggles are rendered
      const toggles = screen.getAllByTestId('toggle-switch');
      expect(toggles.length).toBeGreaterThan(0);
    });

    // Verify the component can access PascalCase fields correctly
    // The first toggle should be for MemberAddedInApp, which is true
    const toggles = screen.getAllByTestId('toggle-switch');
    expect(toggles[0]).toBeChecked();
    
    // The second toggle should be for MemberAddedEmail, which is false
    expect(toggles[1]).not.toBeChecked();
  });

  it('should call updatePreferences with PascalCase field names', async () => {
    const mockUpdatePreferences = vi.fn().mockResolvedValue(mockPreferences);
    
    vi.mocked(useNotificationsModule.useNotificationPreferences).mockReturnValue({
      preferences: mockPreferences,
      loading: false,
      error: null,
      fetchPreferences: vi.fn(),
      updatePreferences: mockUpdatePreferences,
    });

    render(
      <MemoryRouter>
        <ProfileNotificationSettings />
      </MemoryRouter>
    );

    await waitFor(() => {
      const toggles = screen.getAllByTestId('toggle-switch');
      expect(toggles.length).toBeGreaterThan(0);
    });

    // Click the first toggle (MemberAddedInApp)
    const toggles = screen.getAllByTestId('toggle-switch');
    toggles[0].click();

    await waitFor(() => {
      // Verify updatePreferences was called with PascalCase field name
      expect(mockUpdatePreferences).toHaveBeenCalledWith(
        expect.objectContaining({
          MemberAddedInApp: expect.any(Boolean),
        })
      );
    });
  });

  it('should handle loading state', () => {
    vi.mocked(useNotificationsModule.useNotificationPreferences).mockReturnValue({
      preferences: null,
      loading: true,
      error: null,
      fetchPreferences: vi.fn(),
      updatePreferences: vi.fn(),
    });

    render(
      <MemoryRouter>
        <ProfileNotificationSettings />
      </MemoryRouter>
    );

    expect(screen.getByText('Loading notification settings...')).toBeInTheDocument();
  });

  it('should handle error state', () => {
    vi.mocked(useNotificationsModule.useNotificationPreferences).mockReturnValue({
      preferences: null,
      loading: false,
      error: 'Failed to load',
      fetchPreferences: vi.fn(),
      updatePreferences: vi.fn(),
    });

    render(
      <MemoryRouter>
        <ProfileNotificationSettings />
      </MemoryRouter>
    );

    expect(screen.getByText('Failed to load notification settings')).toBeInTheDocument();
  });
});
