import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { vi, describe, it, expect, beforeEach, Mock } from 'vitest';
import GlobalSearch from '../GlobalSearch';
import { useAuth } from '../../../hooks/useAuth';

// Mock the useAuth hook
vi.mock('../../../hooks/useAuth');

// Mock react-router-dom navigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

const mockApi = {
  get: vi.fn(),
};

const MockedUseAuth = useAuth as Mock;

const renderWithRouter = (component: React.ReactElement) => {
  return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('GlobalSearch', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    MockedUseAuth.mockReturnValue({
      api: mockApi,
      isAuthenticated: true,
      accessToken: 'mock-token',
      refreshToken: 'mock-refresh-token',
      login: vi.fn(),
      logout: vi.fn(),
    });
  });

  it('renders search input', () => {
    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    expect(searchInput).toBeDefined();
  });

  it('shows loading indicator during search', async () => {
    // Mock a delayed API response
    mockApi.get.mockImplementation(() => 
      new Promise(resolve => 
        setTimeout(() => resolve({ data: { clubs: [], events: [] } }), 100)
      )
    );

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    fireEvent.change(searchInput, { target: { value: 'test' } });

    // Should show loading indicator briefly
    await waitFor(() => {
      const loadingIcon = screen.getByText('⟳');
      expect(loadingIcon).toBeDefined();
    });
  });

  it('performs search and displays results', async () => {
    const mockResults = {
      clubs: [
        {
          type: 'club',
          id: '1',
          name: 'Test Club',
          description: 'A test club',
        },
      ],
      events: [
        {
          type: 'event',
          id: '2',
          name: 'Test Event',
          club_id: '1',
          club_name: 'Test Club',
          start_time: '2024-06-01T10:00:00Z',
          end_time: '2024-06-01T12:00:00Z',
        },
      ],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    fireEvent.change(searchInput, { target: { value: 'test' } });

    // Wait for search to complete
    await waitFor(() => {
      expect(mockApi.get).toHaveBeenCalledWith("/api/v2/SearchGlobal(query='test')");
    });

    // Check if results are displayed
    await waitFor(() => {
      expect(screen.getByText('Clubs (1)')).toBeDefined();
      expect(screen.getByText('Events (1)')).toBeDefined();
      expect(screen.getAllByText('Test Club')).toHaveLength(2); // One in club result, one in event's club name
      expect(screen.getByText('Test Event')).toBeDefined();
    });
  });

  it('shows no results message when no results found', async () => {
    mockApi.get.mockResolvedValue({ data: { clubs: [], events: [] } });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    fireEvent.change(searchInput, { target: { value: 'nonexistent' } });

    await waitFor(() => {
      expect(screen.getByText('No results found for "nonexistent"')).toBeDefined();
    });
  });

  it('navigates to club when club result is clicked', async () => {
    const mockResults = {
      clubs: [
        {
          type: 'club',
          id: 'club-1',
          name: 'Test Club',
          description: 'A test club',
        },
      ],
      events: [],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    fireEvent.change(searchInput, { target: { value: 'test' } });

    // Wait for results and click on club
    await waitFor(() => {
      const clubResult = screen.getByText('Test Club');
      expect(clubResult).toBeDefined();
      fireEvent.click(clubResult);
    });

    expect(mockNavigate).toHaveBeenCalledWith('/clubs/club-1');
  });

  it('navigates to club when event result is clicked', async () => {
    const mockResults = {
      clubs: [],
      events: [
        {
          type: 'event',
          id: 'event-1',
          name: 'Test Event',
          club_id: 'club-1',
          club_name: 'Test Club',
          start_time: '2024-06-01T10:00:00Z',
          end_time: '2024-06-01T12:00:00Z',
        },
      ],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    fireEvent.change(searchInput, { target: { value: 'event' } });

    // Wait for results and click on event
    await waitFor(() => {
      const eventResult = screen.getByText('Test Event');
      expect(eventResult).toBeDefined();
      fireEvent.click(eventResult);
    });

    expect(mockNavigate).toHaveBeenCalledWith('/clubs/club-1');
  });

  it('closes dropdown when clicking outside', async () => {
    const mockResults = {
      clubs: [{ type: 'club', id: '1', name: 'Test Club', description: '' }],
      events: [],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    fireEvent.change(searchInput, { target: { value: 'test' } });

    // Wait for dropdown to appear
    await waitFor(() => {
      expect(screen.getByText('Test Club')).toBeDefined();
    });

    // Click outside the dropdown
    fireEvent.mouseDown(document.body);

    // Dropdown should close
    await waitFor(() => {
      expect(screen.queryByText('Test Club')).toBeNull();
    });
  });

  it('clears results when search input is cleared', async () => {
    const mockResults = {
      clubs: [{ type: 'club', id: '1', name: 'Test Club', description: '' }],
      events: [],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    
    // First search
    fireEvent.change(searchInput, { target: { value: 'test' } });
    
    await waitFor(() => {
      expect(screen.getByText('Test Club')).toBeDefined();
    });

    // Clear search
    fireEvent.change(searchInput, { target: { value: '' } });

    // Results should disappear
    await waitFor(() => {
      expect(screen.queryByText('Test Club')).toBeNull();
    });
  });

  it('debounces search requests', async () => {
    mockApi.get.mockResolvedValue({ data: { clubs: [], events: [] } });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    
    // Type multiple characters quickly
    fireEvent.change(searchInput, { target: { value: 't' } });
    fireEvent.change(searchInput, { target: { value: 'te' } });
    fireEvent.change(searchInput, { target: { value: 'tes' } });
    fireEvent.change(searchInput, { target: { value: 'test' } });

    // Wait for debounce
    await waitFor(() => {
      expect(mockApi.get).toHaveBeenCalledTimes(1);
      expect(mockApi.get).toHaveBeenCalledWith("/api/v2/SearchGlobal(query='test')");
    }, { timeout: 1000 });
  });

  it('handles search errors gracefully', async () => {
    mockApi.get.mockRejectedValue(new Error('API Error'));
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    fireEvent.change(searchInput, { target: { value: 'test' } });

    await waitFor(() => {
      expect(consoleSpy).toHaveBeenCalledWith('Search failed:', expect.any(Error));
    });

    consoleSpy.mockRestore();
  });

  it('formats dates correctly', async () => {
    const mockResults = {
      clubs: [],
      events: [
        {
          type: 'event',
          id: '1',
          name: 'Test Event',
          club_id: 'club-1',
          club_name: 'Test Club',
          start_time: '2024-06-01T14:30:00Z',
          end_time: '2024-06-01T16:30:00Z',
        },
      ],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    fireEvent.change(searchInput, { target: { value: 'event' } });

    await waitFor(() => {
      // Check if date is formatted (exact format depends on locale)
      const dateElement = screen.getByText(/Jun|6/);
      expect(dateElement).toBeDefined();
    });
  });

  it('truncates long descriptions', async () => {
    const longDescription = 'This is a very long description that should be truncated because it exceeds the maximum length allowed in the search results dropdown component';
    
    const mockResults = {
      clubs: [
        {
          type: 'club',
          id: '1',
          name: 'Test Club',
          description: longDescription,
        },
      ],
      events: [],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search clubs and events...');
    fireEvent.change(searchInput, { target: { value: 'test' } });

    await waitFor(() => {
      const truncatedText = screen.getByText(/This is a very long description.*\.\.\./);
      expect(truncatedText).toBeDefined();
    });
  });
});
