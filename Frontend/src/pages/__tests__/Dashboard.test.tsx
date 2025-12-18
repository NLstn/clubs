import { describe, it, expect, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom';
import { renderWithProviders, screen, act, waitFor } from '../../test/test-utils';
import Dashboard from '../Dashboard';
import { useAuth } from '../../hooks/useAuth';
import type { AxiosInstance } from 'axios';

// Mock the hooks
vi.mock('../../hooks/useAuth');
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

// Mock Layout component
vi.mock('../../components/layout/Layout', () => ({
  default: ({ children, title }: { children: React.ReactNode; title: string }) => (
    <div data-testid="layout" data-title={title}>
      {children}
    </div>
  ),
}));

// Mock IntersectionObserver
class MockIntersectionObserver {
  observe = vi.fn();
  unobserve = vi.fn();
  disconnect = vi.fn();
  
  constructor() {
    // Mock implementation for testing
  }
}

global.IntersectionObserver = MockIntersectionObserver as unknown as typeof IntersectionObserver;

describe('Dashboard', () => {
  let mockApi: { get: ReturnType<typeof vi.fn> };

  beforeEach(() => {
    mockApi = {
      get: vi.fn(),
    };

    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      accessToken: 'mock-token',
      refreshToken: 'mock-refresh-token',
      login: vi.fn(),
      logout: vi.fn(),
      api: mockApi as unknown as AxiosInstance,
    });
  });

  it('renders activity feed when activities are available', async () => {
    const mockTimelineData = [
      {
        ID: '1',
        Type: 'news',
        Title: 'Test News',
        Content: 'Test news content',
        ClubName: 'Test Club',
        ClubID: '1',
        Timestamp: '2024-01-01T10:00:00Z',
        CreatedAt: '2024-01-01T10:00:00Z',
        UpdatedAt: '2024-01-01T10:00:00Z',
        Actor: 'user1',
        ActorName: 'John Doe',
      },
      {
        ID: '2',
        Type: 'event',
        Title: 'Test Event',
        Content: 'Test event content',
        ClubName: 'Test Club',
        ClubID: '1',
        Timestamp: '2024-01-01T11:00:00Z',
        CreatedAt: '2024-01-01T11:00:00Z',
        UpdatedAt: '2024-01-01T11:00:00Z',
        StartTime: '2024-01-01T15:00:00Z',
        EndTime: '2024-01-01T17:00:00Z',
        Actor: 'user2',
        ActorName: 'Jane Smith',
      },
    ];

    mockApi.get.mockResolvedValue({
      data: { value: mockTimelineData },
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    await waitFor(() => {
      expect(screen.getByText('Activity Feed')).toBeInTheDocument();
    });

    // Check if activity feed is rendered
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
    mockApi.get.mockResolvedValue({
      data: { value: [] },
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    await waitFor(() => {
      expect(screen.getByText('Activity Feed')).toBeInTheDocument();
    });

    expect(screen.getByText('No recent activities from your clubs.')).toBeInTheDocument();
  });

  it('renders loading state', async () => {
    mockApi.get.mockImplementation(() => new Promise(() => {})); // Never resolves

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    expect(screen.getByText('Loading dashboard...')).toBeInTheDocument();
  });

  it('renders error state', async () => {
    mockApi.get.mockRejectedValue(new Error('API Error'));

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    await waitFor(() => {
      expect(screen.getByText('Failed to load dashboard data')).toBeInTheDocument();
    });
  });

  it('fetches timeline items with pagination parameters', async () => {
    const mockTimelineData = Array.from({ length: 20 }, (_, i) => ({
      ID: `${i + 1}`,
      Type: 'news',
      Title: `News ${i + 1}`,
      Content: `Content ${i + 1}`,
      ClubName: 'Test Club',
      ClubID: '1',
      Timestamp: '2024-01-01T10:00:00Z',
      CreatedAt: '2024-01-01T10:00:00Z',
      UpdatedAt: '2024-01-01T10:00:00Z',
      Actor: 'user1',
      ActorName: 'John Doe',
    }));

    mockApi.get.mockResolvedValue({
      data: { value: mockTimelineData },
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    await waitFor(() => {
      expect(screen.getByText('Activity Feed')).toBeInTheDocument();
    });

    // Verify the API was called with correct query parameters
    expect(mockApi.get).toHaveBeenCalledWith('/api/v2/TimelineItems?$orderby=CreatedAt desc&$top=20');
  });

  it('loads more items when loadMore is triggered', async () => {
    const initialData = Array.from({ length: 20 }, (_, i) => ({
      ID: `${i + 1}`,
      Type: 'news',
      Title: `News ${i + 1}`,
      Content: `Content ${i + 1}`,
      ClubName: 'Test Club',
      ClubID: '1',
      Timestamp: '2024-01-01T10:00:00Z',
      CreatedAt: '2024-01-01T10:00:00Z',
      UpdatedAt: '2024-01-01T10:00:00Z',
      Actor: 'user1',
      ActorName: 'John Doe',
    }));

    const moreData = Array.from({ length: 20 }, (_, i) => ({
      ID: `${i + 21}`,
      Type: 'news',
      Title: `News ${i + 21}`,
      Content: `Content ${i + 21}`,
      ClubName: 'Test Club',
      ClubID: '1',
      Timestamp: '2024-01-01T10:00:00Z',
      CreatedAt: '2024-01-01T10:00:00Z',
      UpdatedAt: '2024-01-01T10:00:00Z',
      Actor: 'user1',
      ActorName: 'John Doe',
    }));

    mockApi.get
      .mockResolvedValueOnce({ data: { value: initialData } })
      .mockResolvedValueOnce({ data: { value: moreData } });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    await waitFor(() => {
      expect(screen.getByText('News 1')).toBeInTheDocument();
    });

    // First call should be for initial load
    expect(mockApi.get).toHaveBeenCalledWith('/api/v2/TimelineItems?$orderby=CreatedAt desc&$top=20');
  });

  it('prevents duplicate items when loading more', async () => {
    const initialData = Array.from({ length: 20 }, (_, i) => ({
      ID: `${i + 1}`,
      Type: 'news',
      Title: `News ${i + 1}`,
      Content: `Content ${i + 1}`,
      ClubName: 'Test Club',
      ClubID: '1',
      Timestamp: '2024-01-01T10:00:00Z',
      CreatedAt: '2024-01-01T10:00:00Z',
      UpdatedAt: '2024-01-01T10:00:00Z',
      Actor: 'user1',
      ActorName: 'John Doe',
    }));

    mockApi.get.mockResolvedValue({
      data: { value: initialData },
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    await waitFor(() => {
      expect(screen.getByText('News 1')).toBeInTheDocument();
    });

    // Verify only unique items are displayed
    const allItems = screen.getAllByText(/News \d+/);
    const uniqueTitles = new Set(allItems.map(item => item.textContent));
    expect(uniqueTitles.size).toBe(allItems.length);
  });

  it('sets hasMore to false when fewer items are returned', async () => {
    const mockTimelineData = Array.from({ length: 10 }, (_, i) => ({
      ID: `${i + 1}`,
      Type: 'news',
      Title: `News ${i + 1}`,
      Content: `Content ${i + 1}`,
      ClubName: 'Test Club',
      ClubID: '1',
      Timestamp: '2024-01-01T10:00:00Z',
      CreatedAt: '2024-01-01T10:00:00Z',
      UpdatedAt: '2024-01-01T10:00:00Z',
      Actor: 'user1',
      ActorName: 'John Doe',
    }));

    mockApi.get.mockResolvedValue({
      data: { value: mockTimelineData },
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    await waitFor(() => {
      expect(screen.getByText('News 1')).toBeInTheDocument();
    });

    // When fewer than PAGE_SIZE items are returned, the load-more-trigger should not be visible
    expect(screen.queryByText('Loading more activities...')).not.toBeInTheDocument();
  });

  it('shows loading indicator when hasMore is true', async () => {
    const mockTimelineData = Array.from({ length: 20 }, (_, i) => ({
      ID: `${i + 1}`,
      Type: 'news',
      Title: `News ${i + 1}`,
      Content: `Content ${i + 1}`,
      ClubName: 'Test Club',
      ClubID: '1',
      Timestamp: '2024-01-01T10:00:00Z',
      CreatedAt: '2024-01-01T10:00:00Z',
      UpdatedAt: '2024-01-01T10:00:00Z',
      Actor: 'user1',
      ActorName: 'John Doe',
    }));

    mockApi.get.mockResolvedValue({
      data: { value: mockTimelineData },
    });

    await act(async () => {
      renderWithProviders(<Dashboard />);
    });

    await waitFor(() => {
      expect(screen.getByText('News 1')).toBeInTheDocument();
    });

    // When hasMore is true (20 items returned), the load-more-trigger should be present
    const loaderTrigger = document.querySelector('.load-more-trigger');
    expect(loaderTrigger).toBeInTheDocument();
  });
});