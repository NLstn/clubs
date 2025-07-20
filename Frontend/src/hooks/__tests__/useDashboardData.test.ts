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

const mockNews = [
  { 
    id: '1', 
    title: 'News 1', 
    content: 'Content 1', 
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    club_name: 'Club 1',
    club_id: 'club1'
  }
];

const mockEvents = [
  { 
    id: '1', 
    name: 'Event 1', 
    start_time: '2024-01-01T00:00:00Z',
    end_time: '2024-01-01T02:00:00Z',
    club_name: 'Club 1',
    club_id: 'club1'
  }
];

const mockActivities = [
  { 
    id: '1', 
    type: 'news', 
    title: 'User joined',
    content: 'A new user has joined the club',
    club_name: 'Club 1',
    club_id: 'club1',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z'
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
    // Mock successful API responses
    mockApiGet.mockImplementation((url: string) => {
      if (url === '/api/v1/dashboard/news') {
        return Promise.resolve({ data: mockNews });
      } else if (url === '/api/v1/dashboard/events') {
        return Promise.resolve({ data: mockEvents });
      } else if (url === '/api/v1/dashboard/activities') {
        return Promise.resolve({ data: mockActivities });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    const { result } = renderHook(() => useDashboardData());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.news).toEqual(mockNews);
    expect(result.current.events).toEqual(mockEvents);
    expect(result.current.activities).toEqual(mockActivities);
    expect(result.current.error).toBeNull();
    
    // Verify all three API endpoints were called
    expect(mockApiGet).toHaveBeenCalledWith('/api/v1/dashboard/news');
    expect(mockApiGet).toHaveBeenCalledWith('/api/v1/dashboard/events');
    expect(mockApiGet).toHaveBeenCalledWith('/api/v1/dashboard/activities');
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
    // Simple test that focuses on the refetch functionality without complex mock switching
    mockApiGet.mockImplementation((url: string) => {
      if (url === '/api/v1/dashboard/news') {
        return Promise.resolve({ data: mockNews });
      } else if (url === '/api/v1/dashboard/events') {
        return Promise.resolve({ data: mockEvents });
      } else if (url === '/api/v1/dashboard/activities') {
        return Promise.resolve({ data: mockActivities });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    const { result } = renderHook(() => useDashboardData());

    // Wait for initial data to load
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.news).toEqual(mockNews);
    expect(result.current.events).toEqual(mockEvents);
    expect(result.current.activities).toEqual(mockActivities);

    // Clear the call count to track refetch calls
    mockApiGet.mockClear();

    // Call refetch
    await act(async () => {
      await result.current.refetch();
    });

    // Verify refetch called all endpoints
    expect(mockApiGet).toHaveBeenCalledWith('/api/v1/dashboard/news');
    expect(mockApiGet).toHaveBeenCalledWith('/api/v1/dashboard/events');
    expect(mockApiGet).toHaveBeenCalledWith('/api/v1/dashboard/activities');
    expect(mockApiGet).toHaveBeenCalledTimes(3);
  });
});
