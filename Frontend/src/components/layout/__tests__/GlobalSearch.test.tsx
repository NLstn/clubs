import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { vi, describe, it, expect, beforeEach, Mock } from 'vitest';
import GlobalSearch from '../GlobalSearch';
import { useAuth } from '../../../hooks/useAuth';
import { TestI18nProvider } from '../../../test/i18n-test-utils';

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

// Mock the recentClubs utility
const mockGetRecentClubs = vi.fn();
const mockRemoveRecentClub = vi.fn();
vi.mock('../../../utils/recentClubs', () => ({
  getRecentClubs: () => mockGetRecentClubs(),
  removeRecentClub: (id: string) => mockRemoveRecentClub(id),
}));

const mockApi = {
  get: vi.fn(),
};

const MockedUseAuth = useAuth as Mock;

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <TestI18nProvider>
      <BrowserRouter>{component}</BrowserRouter>
    </TestI18nProvider>
  );
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
    mockGetRecentClubs.mockReturnValue([]);
  });

  it('renders search input', () => {
    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    expect(searchInput).toBeDefined();
  });

  it('shows recent clubs when focused with empty query', () => {
    const mockClubs = [
      { id: '1', name: 'Club A', visitedAt: 1000 },
      { id: '2', name: 'Club B', visitedAt: 2000 },
    ];
    mockGetRecentClubs.mockReturnValue(mockClubs);

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.focus(searchInput);

    // Should show recent clubs dropdown
    expect(screen.getByText('Recent Clubs')).toBeDefined();
    expect(screen.getByText('Club A')).toBeDefined();
    expect(screen.getByText('Club B')).toBeDefined();
  });

  it('shows "No recent clubs" when focused with no recent clubs', () => {
    mockGetRecentClubs.mockReturnValue([]);

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.focus(searchInput);

    expect(screen.getByText('No recent clubs')).toBeDefined();
  });

  it('navigates to club when recent club is clicked', () => {
    const mockClubs = [{ id: 'club-1', name: 'Test Club', visitedAt: 1000 }];
    mockGetRecentClubs.mockReturnValue(mockClubs);

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.focus(searchInput);

    const clubItem = screen.getByText('Test Club');
    fireEvent.click(clubItem);

    expect(mockNavigate).toHaveBeenCalledWith('/clubs/club-1');
  });

  it('navigates to clubs list when "View All Clubs" is clicked', () => {
    mockGetRecentClubs.mockReturnValue([]);

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.focus(searchInput);

    const viewAllButton = screen.getByText('View All Clubs');
    fireEvent.click(viewAllButton);

    expect(mockNavigate).toHaveBeenCalledWith('/clubs');
  });

  it('removes recent club when remove button is clicked', () => {
    const mockClubs = [{ id: 'club-1', name: 'Test Club', visitedAt: 1000 }];
    mockGetRecentClubs.mockReturnValue(mockClubs);

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.focus(searchInput);

    // Find and click the remove button
    const removeButton = screen.getByTitle('Remove from recent clubs');
    fireEvent.click(removeButton);

    expect(mockRemoveRecentClub).toHaveBeenCalledWith('club-1');
  });

  it('shows loading indicator during search', async () => {
    // Mock a delayed API response
    mockApi.get.mockImplementation(() => 
      new Promise(resolve => 
        setTimeout(() => resolve({ data: { Clubs: [], Events: [] } }), 100)
      )
    );

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.change(searchInput, { target: { value: 'test' } });

    // Should show loading indicator briefly
    await waitFor(() => {
      const loadingIcon = screen.getByText('⟳');
      expect(loadingIcon).toBeDefined();
    });
  });

  it('performs search and displays results', async () => {
    const mockResults = {
      Clubs: [
        {
          Type: 'club',
          ID: '1',
          Name: 'Test Club',
          Description: 'A test club',
        },
      ],
      Events: [
        {
          Type: 'event',
          ID: '2',
          Name: 'Test Event',
          ClubID: '1',
          ClubName: 'Test Club',
          StartTime: '2024-06-01T10:00:00Z',
          EndTime: '2024-06-01T12:00:00Z',
        },
      ],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
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
    mockApi.get.mockResolvedValue({ data: { Clubs: [], Events: [] } });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.change(searchInput, { target: { value: 'nonexistent' } });

    await waitFor(() => {
      expect(screen.getByText('No results found for "nonexistent"')).toBeDefined();
    });
  });

  it('navigates to club when club result is clicked', async () => {
    const mockResults = {
      Clubs: [
        {
          Type: 'club',
          ID: 'club-1',
          Name: 'Test Club',
          Description: 'A test club',
        },
      ],
      Events: [],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
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
      Clubs: [],
      Events: [
        {
          Type: 'event',
          ID: 'event-1',
          Name: 'Test Event',
          ClubID: 'club-1',
          ClubName: 'Test Club',
          StartTime: '2024-06-01T10:00:00Z',
          EndTime: '2024-06-01T12:00:00Z',
        },
      ],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
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
    mockGetRecentClubs.mockReturnValue([{ id: '1', name: 'Recent Club', visitedAt: 1000 }]);

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.focus(searchInput);

    // Wait for dropdown to appear
    await waitFor(() => {
      expect(screen.getByText('Recent Club')).toBeDefined();
    });

    // Click outside the dropdown
    fireEvent.mouseDown(document.body);

    // Dropdown should close
    await waitFor(() => {
      expect(screen.queryByText('Recent Club')).toBeNull();
    });
  });

  it('debounces search requests', async () => {
    mockApi.get.mockResolvedValue({ data: { Clubs: [], Events: [] } });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    
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
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.change(searchInput, { target: { value: 'test' } });

    await waitFor(() => {
      expect(consoleSpy).toHaveBeenCalledWith('Search failed:', expect.any(Error));
    });

    consoleSpy.mockRestore();
  });

  it('formats dates correctly', async () => {
    const mockResults = {
      Clubs: [],
      Events: [
        {
          Type: 'event',
          ID: '1',
          Name: 'Test Event',
          ClubID: 'club-1',
          ClubName: 'Test Club',
          StartTime: '2024-06-01T14:30:00Z',
          EndTime: '2024-06-01T16:30:00Z',
        },
      ],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
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
      Clubs: [
        {
          Type: 'club',
          ID: '1',
          Name: 'Test Club',
          Description: longDescription,
        },
      ],
      Events: [],
    };

    mockApi.get.mockResolvedValue({ data: mockResults });

    renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.change(searchInput, { target: { value: 'test' } });

    await waitFor(() => {
      const truncatedText = screen.getByText(/This is a very long description.*\.\.\./);
      expect(truncatedText).toBeDefined();
    });
  });

  it('expands search field when focused', () => {
    const { container } = renderWithRouter(<GlobalSearch />);
    
    const searchInput = screen.getByPlaceholderText('Search');
    const searchContainer = container.querySelector('.global-search');
    
    // Before focus, should not have focused class
    expect(searchContainer?.classList.contains('global-search-focused')).toBe(false);
    
    fireEvent.focus(searchInput);
    
    // After focus, should have focused class
    expect(searchContainer?.classList.contains('global-search-focused')).toBe(true);
  });
});
