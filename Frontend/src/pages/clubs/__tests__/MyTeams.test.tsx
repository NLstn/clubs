import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import '@testing-library/jest-dom';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import MyTeams from '../MyTeams';

// Mock API
vi.mock('../../../utils/api', () => ({
  default: {
    get: vi.fn()
  }
}));

// Import the mocked api to use in tests
import api from '../../../utils/api';
const mockGet = vi.fn();
vi.mocked(api).get = mockGet;

// Mock the hooks
vi.mock('../../../hooks/useCurrentUser', () => ({
  useCurrentUser: () => ({
    user: { ID: 'test-user-id', Email: 'test@example.com' },
    loading: false,
    error: null
  })
}));

vi.mock('../../../hooks/useTranslation', () => ({
  useT: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'teams.myTeams': 'My Teams'
      };
      return translations[key] || key;
    }
  })
}));

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useParams: () => ({ id: 'test-club-id' })
  };
});

// Mock API
vi.mock('../../../utils/api', () => ({
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

describe('MyTeams Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

  it('renders teams when user has teams', async () => {
    const mockTeams = [
      {
        id: 'team-1',
        name: 'Development Team',
        description: 'Team for developers',
        createdAt: '2024-01-01T00:00:00Z',
        clubId: 'test-club-id'
      },
      {
        id: 'team-2',
        name: 'Marketing Team',
        description: 'Team for marketing',
        createdAt: '2024-01-01T00:00:00Z',
        clubId: 'test-club-id'
      }
    ];

    mockGet.mockResolvedValueOnce({ data: mockTeams });

    renderWithRouter(<MyTeams />);

    // Wait for teams to load
    await waitFor(() => {
      expect(screen.getByText('My Teams')).toBeInTheDocument();
    });

    expect(screen.getByText('Development Team')).toBeInTheDocument();
    expect(screen.getByText('Team for developers')).toBeInTheDocument();
    expect(screen.getByText('Marketing Team')).toBeInTheDocument();
    expect(screen.getByText('Team for marketing')).toBeInTheDocument();

    expect(mockGet).toHaveBeenCalledWith('/api/v1/clubs/test-club-id/teams?user');
  });

  it('does not render anything when user has no teams', async () => {
    mockGet.mockResolvedValueOnce({ data: [] });

    const { container } = renderWithRouter(<MyTeams />);

    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled();
    });

    // Component should not render anything when no teams
    expect(container.firstChild).toBeNull();
  });

  it('handles API errors gracefully', async () => {
    mockGet.mockRejectedValueOnce(new Error('API Error'));

    renderWithRouter(<MyTeams />);

    await waitFor(() => {
      expect(screen.getByText('Failed to fetch teams')).toBeInTheDocument();
    });
  });

  it('handles 403 errors by not showing error message', async () => {
    const error = {
      response: { status: 403 }
    };
    mockGet.mockRejectedValueOnce(error);

    const { container } = renderWithRouter(<MyTeams />);

    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled();
    });

    // Should not show error for 403, should just not render
    expect(container.firstChild).toBeNull();
  });

  it('shows loading state initially', () => {
    // Don't resolve the promise to keep it in loading state
    mockGet.mockImplementation(() => new Promise(() => {}));

    renderWithRouter(<MyTeams />);

    expect(screen.getByText('Loading teams...')).toBeInTheDocument();
  });
});
