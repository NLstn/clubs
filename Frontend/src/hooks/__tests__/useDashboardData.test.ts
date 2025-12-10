import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useDashboardData } from '../useDashboardData';

// Create a mock API object with mocked methods
const mockApiGet = vi.fn();

// Mock useAuth to return our mocked API instance
vi.mock('../useAuth', () => ({
  useAuth: () => ({
    api: {
      get: mockApiGet
    },
    isAuthenticated: true,
    user: { id: 'user1' }
  })
}));

// Mock timeline items from the unified Timeline entity
const mockTimelineItems = [
  {
    ID: 'news-1',
    ClubID: 'club1',
    ClubName: 'Club 1',
    Type: 'news',
    Title: 'News 1',
    Content: 'Content 1',
    Timestamp: '2024-01-01T00:00:00Z',
    CreatedAt: '2024-01-01T00:00:00Z',
    UpdatedAt: '2024-01-01T00:00:00Z',
    Metadata: {}
  },
  {
    ID: 'event-1',
    ClubID: 'club1',
    ClubName: 'Club 1',
    Type: 'event',
    Title: 'Event 1',
    Content: '',
    Timestamp: '2024-01-01T00:00:00Z',
    CreatedAt: '2024-01-01T00:00:00Z',
    UpdatedAt: '2024-01-01T00:00:00Z',
    StartTime: '2024-01-01T00:00:00Z',
    EndTime: '2024-01-01T02:00:00Z',
    Metadata: {}
  },
  {
    ID: 'activity-1',
    ClubID: 'club1',
    ClubName: 'Club 1',
    Type: 'activity',
    Title: 'User joined',
    Content: 'A new user has joined the club',
    Timestamp: '2024-01-01T00:00:00Z',
    CreatedAt: '2024-01-01T00:00:00Z',
    UpdatedAt: '2024-01-01T00:00:00Z',
    Metadata: {}
  }
];

describe('useDashboardData', () => {
  beforeEach(() => {
    mockApiGet.mockClear();
  });

  it('should return default values initially', () => {
    // Mock API calls to return pending promises
    mockApiGet.mockImplementation(() => new Promise(() => {}));
    
    const { result } = renderHook(() => useDashboardData());

    expect(result.current.news).toEqual([]);
    expect(result.current.events).toEqual([]);
    expect(result.current.activities).toEqual([]);
    expect(result.current.loading).toBe(true);
    expect(result.current.error).toBeNull();
    expect(typeof result.current.refetch).toBe('function');
  });

  it('should fetch dashboard data successfully', async () => {
    // Mock successful API response with OData v2 format using Timeline entity
    mockApiGet.mockResolvedValue({ data: { value: mockTimelineItems } });

    const { result } = renderHook(() => useDashboardData());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // Check that timeline data is available
    expect(result.current.timeline).toHaveLength(3);
    expect(result.current.timeline[0].Type).toBe('news');
    expect(result.current.timeline[1].Type).toBe('event');
    expect(result.current.timeline[2].Type).toBe('activity');
    
    // Check that legacy format is still available
    expect(result.current.news).toHaveLength(1);
    expect(result.current.news[0].id).toBe('news-1');
    expect(result.current.events).toHaveLength(1);
    expect(result.current.events[0].id).toBe('event-1');
    expect(result.current.activities).toHaveLength(1);
    expect(result.current.activities[0].ID).toBe('activity-1');
    
    expect(result.current.error).toBeNull();
    
    // Verify the unified Timeline endpoint was called
    expect(mockApiGet).toHaveBeenCalledWith('/api/v2/TimelineItems');
    expect(mockApiGet).toHaveBeenCalledTimes(1);
  });

  it('should handle API errors gracefully', async () => {
    // Mock API calls to throw errors
    mockApiGet.mockRejectedValue(new Error('API Error'));

    const { result } = renderHook(() => useDashboardData());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.news).toEqual([]);
    expect(result.current.events).toEqual([]);
    expect(result.current.activities).toEqual([]);
    expect(result.current.error).toBe('Failed to load dashboard data');
  });

  it('should refetch data when refetch is called', async () => {
    // Mock successful API response
    mockApiGet.mockResolvedValue({ data: { value: mockTimelineItems } });

    const { result } = renderHook(() => useDashboardData());

    // Wait for initial data to load
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.timeline).toHaveLength(3);
    expect(result.current.news).toHaveLength(1);
    expect(result.current.events).toHaveLength(1);
    expect(result.current.activities).toHaveLength(1);

    // Clear the call count to track refetch calls
    mockApiGet.mockClear();

    // Call refetch
    await act(async () => {
      await result.current.refetch();
    });

    // Verify refetch called the Timeline endpoint
    expect(mockApiGet).toHaveBeenCalledWith('/api/v2/TimelineItems');
    expect(mockApiGet).toHaveBeenCalledTimes(1);
  });
});
