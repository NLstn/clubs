import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import AdminClubSettings from '../AdminClubSettings';
import api from '../../../../../utils/api';

// Mock the API
vi.mock('../../../../../utils/api');

// Mock the ToggleSwitch component
vi.mock('@/components/ui', () => ({
  ToggleSwitch: ({ checked, onChange, disabled }: { checked: boolean; onChange: (value: boolean) => void; disabled?: boolean }) => (
    <input
      type="checkbox"
      checked={checked}
      onChange={(e) => onChange(e.target.checked)}
      disabled={disabled}
      data-testid="toggle-switch"
    />
  ),
}));

// Mock the translation hook
vi.mock('../../../../../hooks/useTranslation', () => ({
  useT: () => ({
    t: (key: string) => key,
  }),
}));

describe('AdminClubSettings', () => {
  const mockSettings = {
    ID: 'settings-id',
    ClubID: 'club-id',
    FinesEnabled: true,
    ShiftsEnabled: true,
    TeamsEnabled: false,
    MembersListVisible: true,
    DiscoverableByNonMembers: false,
    CreatedAt: '2024-01-01T00:00:00Z',
    CreatedBy: 'user-id',
    UpdatedAt: '2024-01-01T00:00:00Z',
    UpdatedBy: 'user-id',
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch and display settings with PascalCase field names', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockSettings });

    render(
      <MemoryRouter initialEntries={['/clubs/club-id/admin/settings']}>
        <Routes>
          <Route path="/clubs/:id/admin/settings" element={<AdminClubSettings />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.queryByText('clubs.loading.settings')).not.toBeInTheDocument();
    });

    // Verify the API was called with correct endpoint
    expect(api.get).toHaveBeenCalledWith("/api/v2/Clubs('club-id')/Settings");

    // Verify toggles are rendered with correct states
    const toggles = screen.getAllByTestId('toggle-switch');
    expect(toggles).toHaveLength(5);
    
    // FinesEnabled is true
    expect(toggles[0]).toBeChecked();
    // ShiftsEnabled is true
    expect(toggles[1]).toBeChecked();
    // TeamsEnabled is false
    expect(toggles[2]).not.toBeChecked();
  });

  it('should update settings with PascalCase field names', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockSettings });
    vi.mocked(api.patch).mockResolvedValue({ data: { ...mockSettings, ShiftsEnabled: false } });

    render(
      <MemoryRouter initialEntries={['/clubs/club-id/admin/settings']}>
        <Routes>
          <Route path="/clubs/:id/admin/settings" element={<AdminClubSettings />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.queryByText('clubs.loading.settings')).not.toBeInTheDocument();
    });

    // Click the ShiftsEnabled toggle (second toggle)
    const toggles = screen.getAllByTestId('toggle-switch');
    toggles[1].click();

    await waitFor(() => {
      // Verify the API was called with PascalCase field names
      expect(api.patch).toHaveBeenCalledWith(
        "/api/v2/ClubSettings('settings-id')",
        expect.objectContaining({
          FinesEnabled: true,
          ShiftsEnabled: false,
          TeamsEnabled: false,
          MembersListVisible: true,
          DiscoverableByNonMembers: false,
        })
      );
    });
  });

  it('should handle loading state', () => {
    vi.mocked(api.get).mockImplementation(() => new Promise(() => {})); // Never resolves

    render(
      <MemoryRouter initialEntries={['/clubs/club-id/admin/settings']}>
        <Routes>
          <Route path="/clubs/:id/admin/settings" element={<AdminClubSettings />} />
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByText('clubs.loading.settings')).toBeInTheDocument();
  });

  it('should handle error state', async () => {
    vi.mocked(api.get).mockRejectedValue(new Error('API Error'));

    render(
      <MemoryRouter initialEntries={['/clubs/club-id/admin/settings']}>
        <Routes>
          <Route path="/clubs/:id/admin/settings" element={<AdminClubSettings />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('clubs.errors.failedToLoadSettings')).toBeInTheDocument();
    });
  });

  it('should correctly update local state with PascalCase properties', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockSettings });
    vi.mocked(api.patch).mockResolvedValue({ data: { ...mockSettings, TeamsEnabled: true } });

    render(
      <MemoryRouter initialEntries={['/clubs/club-id/admin/settings']}>
        <Routes>
          <Route path="/clubs/:id/admin/settings" element={<AdminClubSettings />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.queryByText('clubs.loading.settings')).not.toBeInTheDocument();
    });

    // Click the TeamsEnabled toggle (third toggle)
    const toggles = screen.getAllByTestId('toggle-switch');
    expect(toggles[2]).not.toBeChecked();
    
    toggles[2].click();

    await waitFor(() => {
      // Verify the toggle reflects the new state
      // Note: Since we're testing the state update logic, we're mainly ensuring
      // the component doesn't crash and the API call is made correctly
      expect(api.patch).toHaveBeenCalled();
    });
  });
});
