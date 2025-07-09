import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import RecentClubsDropdown from '../layout/RecentClubsDropdown';

// Mock the recentClubs utility
const mockGetRecentClubs = vi.fn();
const mockRemoveRecentClub = vi.fn();
const mockNavigate = vi.fn();

vi.mock('../../utils/recentClubs', () => ({
  getRecentClubs: () => mockGetRecentClubs(),
  removeRecentClub: () => mockRemoveRecentClub()
}));

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate
  };
});

// Mock the API
vi.mock('../../utils/api', () => ({
  default: {
    get: vi.fn()
  }
}));

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  );
};

describe('RecentClubsDropdown', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders the recent clubs trigger button', () => {
    mockGetRecentClubs.mockReturnValue([]);
    
    renderWithRouter(<RecentClubsDropdown />);
    
    expect(screen.getByTitle('Recent Clubs')).toBeInTheDocument();
    expect(screen.getByText('Recent clubs')).toBeInTheDocument();
  });

  it('opens dropdown when trigger is clicked', () => {
    mockGetRecentClubs.mockReturnValue([]);
    
    renderWithRouter(<RecentClubsDropdown />);
    
    const trigger = screen.getByTitle('Recent Clubs');
    fireEvent.click(trigger);
    
    expect(screen.getByText('Recent Clubs')).toBeInTheDocument();
    expect(screen.getByText('View All Clubs')).toBeInTheDocument();
  });

  it('displays recent clubs when available', () => {
    const mockClubs = [
      { id: '1', name: 'Club A', visitedAt: 1000 },
      { id: '2', name: 'Club B', visitedAt: 2000 },
    ];
    mockGetRecentClubs.mockReturnValue(mockClubs);
    
    renderWithRouter(<RecentClubsDropdown />);
    
    const trigger = screen.getByTitle('Recent Clubs');
    fireEvent.click(trigger);
    
    expect(screen.getByText('Club A')).toBeInTheDocument();
    expect(screen.getByText('Club B')).toBeInTheDocument();
  });

  it('displays "No recent clubs" when no clubs are available', () => {
    mockGetRecentClubs.mockReturnValue([]);
    
    renderWithRouter(<RecentClubsDropdown />);
    
    const trigger = screen.getByTitle('Recent Clubs');
    fireEvent.click(trigger);
    
    expect(screen.getByText('No recent clubs')).toBeInTheDocument();
  });

  it('navigates to club when club item is clicked', async () => {
    const mockClubs = [
      { id: '1', name: 'Club A', visitedAt: 1000 },
    ];
    mockGetRecentClubs.mockReturnValue(mockClubs);
    
    // Mock successful API call
    const mockApi = await import('../../utils/api');
    vi.mocked(mockApi.default.get).mockResolvedValue({});
    
    renderWithRouter(<RecentClubsDropdown />);
    
    const trigger = screen.getByTitle('Recent Clubs');
    fireEvent.click(trigger);
    
    const clubButton = screen.getByText('Club A');
    fireEvent.click(clubButton);
    
    // Wait for the async operation to complete
    await new Promise(resolve => setTimeout(resolve, 0));
    
    expect(mockNavigate).toHaveBeenCalledWith('/clubs/1');
  });

  it('navigates to clubs list when "View All Clubs" is clicked', () => {
    mockGetRecentClubs.mockReturnValue([]);
    
    renderWithRouter(<RecentClubsDropdown />);
    
    const trigger = screen.getByTitle('Recent Clubs');
    fireEvent.click(trigger);
    
    const viewAllButton = screen.getByText('View All Clubs');
    fireEvent.click(viewAllButton);
    
    expect(mockNavigate).toHaveBeenCalledWith('/clubs');
  });

  it('closes dropdown when clicking outside', () => {
    mockGetRecentClubs.mockReturnValue([]);
    
    renderWithRouter(<RecentClubsDropdown />);
    
    const trigger = screen.getByTitle('Recent Clubs');
    fireEvent.click(trigger);
    
    // Dropdown should be open
    expect(screen.getByText('Recent Clubs')).toBeInTheDocument();
    
    // Click outside the dropdown
    fireEvent.mouseDown(document.body);
    
    // Dropdown should be closed
    expect(screen.queryByText('Recent Clubs')).not.toBeInTheDocument();
  });
});