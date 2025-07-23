import { describe, it, expect, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom';
import { renderWithProviders, screen, act } from '../../test/test-utils';
import Dashboard from '../Dashboard';
import { useAuth } from '../../hooks/useAuth';
import { useDashboardData } from '../../hooks/useDashboardData';
import type { AuthContextType } from '../../context/AuthContext';

// Mock the hooks
vi.mock('../../hooks/useAuth');
vi.mock('../../hooks/useDashboardData');
vi.mock('../../hooks/useCurrentUser', () => ({
  useCurrentUser: () => ({
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
  }),
}));

const mockUseAuth = vi.mocked(useAuth);
const mockUseDashboardData = vi.mocked(useDashboardData);

// Mock Layout component
vi.mock('../../components/layout/Layout', () => ({
  default: ({ children, title }: { children: React.ReactNode; title: string }) => (
    <div data-testid="layout" data-title={title}>
      {children}
    </div>
  ),
}));

describe('Dashboard', () => {
  beforeEach(() => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      accessToken: 'mock-token',
      refreshToken: 'mock-refresh-token',
      login: vi.fn(),
      logout: vi.fn(),
      api: {} as AuthContextType['api'],
    });
  });

  it('renders activity feed when activities are available', async () => {
    const mockActivities = [
      {
        id: '1',
        type: 'news',
        title: 'Test News',
        content: 'Test news content',
        club_name: 'Test Club',
        club_id: '1',
        created_at: '2024-01-01T10:00:00Z',
        updated_at: '2024-01-01T10:00:00Z',
        actor: 'user1',
        actor_name: 'John Doe',
      },
      {
        id: '2',
        type: 'event',
        title: 'Test Event',
        content: 'Test event content',
        club_name: 'Test Club',
        club_id: '1',
        created_at: '2024-01-01T11:00:00Z',
        updated_at: '2024-01-01T11:00:00Z',
        actor: 'user2',
        actor_name: 'Jane Smith',
        metadata: {
          start_time: '2024-01-01T15:00:00Z',
          end_time: '2024-01-01T17:00:00Z',
        },
      },
    ];

    mockUseDashboardData.mockReturnValue({
      news: [],
      events: [],
      activities: mockActivities,
      loading: false,
      error: null,
      refetch: vi.fn(),
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    // Check if activity feed is rendered
    expect(screen.getByText('Activity Feed')).toBeInTheDocument();
    expect(screen.getByText('Test News')).toBeInTheDocument();
    expect(screen.getByText('Test Event')).toBeInTheDocument();
    expect(screen.getByText('news')).toBeInTheDocument();
    expect(screen.getByText('event')).toBeInTheDocument();
    
    // Check if creator information is displayed (only for non-event activities)
    expect(screen.getByText(/Created by John Doe/)).toBeInTheDocument();
    // Events don't show creator information according to the current logic
    expect(screen.queryByText(/Created by Jane Smith/)).not.toBeInTheDocument();
  });

  it('renders empty state when no activities are available', async () => {
    mockUseDashboardData.mockReturnValue({
      news: [],
      events: [],
      activities: [],
      loading: false,
      error: null,
      refetch: vi.fn(),
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    expect(screen.getByText('Activity Feed')).toBeInTheDocument();
    expect(screen.getByText('No recent activities from your clubs.')).toBeInTheDocument();
  });

  it('renders loading state', async () => {
    mockUseDashboardData.mockReturnValue({
      news: [],
      events: [],
      activities: [],
      loading: true,
      error: null,
      refetch: vi.fn(),
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    expect(screen.getByText('Loading dashboard...')).toBeInTheDocument();
  });

  it('renders error state', async () => {
    mockUseDashboardData.mockReturnValue({
      news: [],
      events: [],
      activities: [],
      loading: false,
      error: 'Failed to load dashboard data',
      refetch: vi.fn(),
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    expect(screen.getByText('Failed to load dashboard data')).toBeInTheDocument();
  });
});