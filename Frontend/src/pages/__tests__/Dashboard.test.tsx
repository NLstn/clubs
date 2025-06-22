import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { vi } from 'vitest';
import Dashboard from '../Dashboard';
import { useAuth } from '../../hooks/useAuth';
import { useDashboardData } from '../../hooks/useDashboardData';

// Mock the hooks
vi.mock('../../hooks/useAuth');
vi.mock('../../hooks/useDashboardData');

const mockUseAuth = useAuth as vi.MockedFunction<typeof useAuth>;
const mockUseDashboardData = useDashboardData as vi.MockedFunction<typeof useDashboardData>;

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
      api: {
        get: vi.fn().mockResolvedValue({ data: [] }),
      },
    } as ReturnType<typeof useAuth>);
  });

  it('renders activity feed when activities are available', () => {
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
        created_by: 'user1',
        creator_name: 'John Doe',
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
        created_by: 'user2',
        creator_name: 'Jane Smith',
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

    render(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    // Check if activity feed is rendered
    expect(screen.getByText('Activity Feed')).toBeInTheDocument();
    expect(screen.getByText('Test News')).toBeInTheDocument();
    expect(screen.getByText('Test Event')).toBeInTheDocument();
    expect(screen.getByText('news')).toBeInTheDocument();
    expect(screen.getByText('event')).toBeInTheDocument();
    
    // Check if creator information is displayed
    expect(screen.getByText(/Created by John Doe/)).toBeInTheDocument();
    expect(screen.getByText(/Created by Jane Smith/)).toBeInTheDocument();
  });

  it('renders empty state when no activities are available', () => {
    mockUseDashboardData.mockReturnValue({
      news: [],
      events: [],
      activities: [],
      loading: false,
      error: null,
      refetch: vi.fn(),
    });

    render(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(screen.getByText('Activity Feed')).toBeInTheDocument();
    expect(screen.getByText('No recent activities from your clubs.')).toBeInTheDocument();
  });

  it('renders loading state', () => {
    mockUseDashboardData.mockReturnValue({
      news: [],
      events: [],
      activities: [],
      loading: true,
      error: null,
      refetch: vi.fn(),
    });

    render(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(screen.getByText('Loading dashboard...')).toBeInTheDocument();
  });

  it('renders error state', () => {
    mockUseDashboardData.mockReturnValue({
      news: [],
      events: [],
      activities: [],
      loading: false,
      error: 'Failed to load dashboard data',
      refetch: vi.fn(),
    });

    render(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(screen.getByText('Failed to load dashboard data')).toBeInTheDocument();
  });
});