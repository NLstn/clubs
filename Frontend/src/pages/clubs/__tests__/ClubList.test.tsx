import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import ClubList from '../ClubList';
import api from '../../../utils/api';

// Mock the API
vi.mock('../../../utils/api');
const mockedApi = vi.mocked(api);

// Mock the translation hook
vi.mock('../../../hooks/useTranslation', () => ({
    useT: () => ({
        t: (key: string) => {
            const translations: { [key: string]: string } = {
                'clubs.roles.owner': 'Owner',
                'clubs.roles.admin': 'Admin',
                'clubs.roles.member': 'Member',
            };
            return translations[key] || key;
        }
    })
}));

// Mock Layout component to avoid AuthProvider dependency
vi.mock('../../../components/layout/Layout', () => ({
    default: ({ children, title }: { children: React.ReactNode; title: string }) => (
        <div data-testid="layout" data-title={title}>
            {children}
        </div>
    ),
}));

// Mock navigation
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => mockNavigate,
    };
});

const renderWithRouter = (component: React.ReactElement) => {
    return render(
        <BrowserRouter>
            {component}
        </BrowserRouter>
    );
};

const mockClubs = [
    {
        id: '1',
        name: 'Admin Club',
        description: 'A club where I am admin',
        user_role: 'admin',
        created_at: '2024-01-01T00:00:00Z',
    },
    {
        id: '2',
        name: 'Owner Club',
        description: 'A club where I am owner',
        user_role: 'owner',
        created_at: '2024-01-02T00:00:00Z',
    },
    {
        id: '3',
        name: 'Member Club',
        description: 'A club where I am just a member',
        user_role: 'member',
        created_at: '2024-01-03T00:00:00Z',
    },
    {
        id: '4',
        name: 'Deleted Club',
        description: 'A deleted club where I am owner',
        user_role: 'owner',
        created_at: '2024-01-04T00:00:00Z',
        deleted: true,
    },
];

describe('ClubList', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mockNavigate.mockClear();
    });

    it('renders loading state initially', () => {
        mockedApi.get.mockImplementation(() => new Promise(() => {})); // Never resolves
        renderWithRouter(<ClubList />);
        expect(screen.getByText('Loading clubs...')).toBeInTheDocument();
    });

    it('renders error state when API fails', async () => {
        mockedApi.get.mockRejectedValue(new Error('API Error'));
        renderWithRouter(<ClubList />);
        
        await waitFor(() => {
            expect(screen.getByText('Failed to fetch clubs')).toBeInTheDocument();
        });
    });

    it('renders empty state when no clubs', async () => {
        mockedApi.get.mockResolvedValue({ data: [] });
        renderWithRouter(<ClubList />);
        
        await waitFor(() => {
            expect(screen.getByText('No Clubs Yet')).toBeInTheDocument();
            expect(screen.getByText("You're not a member of any clubs yet.")).toBeInTheDocument();
        });
    });

    it('renders clubs separated by role sections', async () => {
        mockedApi.get.mockResolvedValue({ data: mockClubs });
        renderWithRouter(<ClubList />);
        
        await waitFor(() => {
            // Check section headers
            expect(screen.getByText('Clubs I Manage')).toBeInTheDocument();
            expect(screen.getByText("Clubs I'm a Member Of")).toBeInTheDocument();
            
            // Check admin/owner clubs
            expect(screen.getByText('Admin Club')).toBeInTheDocument();
            expect(screen.getByText('Owner Club')).toBeInTheDocument();
            
            // Check member clubs
            expect(screen.getByText('Member Club')).toBeInTheDocument();
            
            // Check deleted club is shown (owner can see it)
            expect(screen.getByText('Deleted Club')).toBeInTheDocument();
        });
    });

    it('navigates to club details when club card is clicked', async () => {
        mockedApi.get.mockResolvedValue({ data: mockClubs });
        renderWithRouter(<ClubList />);
        
        await waitFor(() => {
            const adminClub = screen.getByText('Admin Club');
            fireEvent.click(adminClub.closest('.club-card')!);
            expect(mockNavigate).toHaveBeenCalledWith('/clubs/1');
        });
    });

    it('navigates to create club page when button is clicked in empty state', async () => {
        mockedApi.get.mockResolvedValue({ data: [] });
        renderWithRouter(<ClubList />);
        
        await waitFor(() => {
            const createButton = screen.getByText('Create Your First Club');
            fireEvent.click(createButton);
            expect(mockNavigate).toHaveBeenCalledWith('/createClub');
        });
    });

    it('displays role badges correctly', async () => {
        mockedApi.get.mockResolvedValue({ data: mockClubs });
        renderWithRouter(<ClubList />);
        
        await waitFor(() => {
            expect(screen.getByText('Admin')).toBeInTheDocument();
            expect(screen.getAllByText('Owner')).toHaveLength(2); // Two owner clubs
            expect(screen.getByText('Member')).toBeInTheDocument();
        });
    });

    it('displays deleted badge for deleted clubs', async () => {
        mockedApi.get.mockResolvedValue({ data: mockClubs });
        renderWithRouter(<ClubList />);
        
        await waitFor(() => {
            expect(screen.getByText('Deleted')).toBeInTheDocument();
        });
    });
});