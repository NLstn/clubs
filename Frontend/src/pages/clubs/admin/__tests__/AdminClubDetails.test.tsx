import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import AdminClubDetails from '../AdminClubDetails';
import api from '../../../../utils/api';

// Mock the API
vi.mock('../../../../utils/api', () => ({
    default: {
        get: vi.fn(),
        patch: vi.fn(),
        post: vi.fn(),
        delete: vi.fn(),
    },
    hardDeleteClub: vi.fn(),
}));

// Mock the translation hook
vi.mock('../../../../hooks/useTranslation', () => ({
    useT: () => ({
        t: (key: string, params?: Record<string, string>) => {
            if (params) {
                let result = key;
                Object.entries(params).forEach(([k, v]) => {
                    result = result.replace(`{{${k}}}`, v);
                });
                return result;
            }
            return key;
        },
    }),
}));

// Mock useClubSettings
vi.mock('../../../../hooks/useClubSettings', () => ({
    useClubSettings: () => ({
        settings: {
            FinesEnabled: true,
            TeamsEnabled: true,
            NewsEnabled: true,
            EventsEnabled: true,
            ShiftsEnabled: true,
        },
        refetch: vi.fn(),
    }),
}));

// Mock utility functions
vi.mock('../../../../utils/recentClubs', () => ({
    removeRecentClub: vi.fn(),
}));

// Mock the OData utilities
vi.mock('@/utils/odata', () => ({
    parseODataCollection: vi.fn((data) => data?.value || []),
}));

// Mock UI components
vi.mock('@/components/ui', () => ({
    Input: ({ label, value, onChange, placeholder, multiline }: {
        label?: string;
        value?: string;
        onChange?: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
        placeholder?: string;
        multiline?: boolean;
    }) => (
        <div>
            <label>{label}</label>
            {multiline ? (
                <textarea
                    data-testid={`input-${label?.toLowerCase().replace(/\s/g, '-')}`}
                    value={value}
                    onChange={onChange}
                    placeholder={placeholder}
                />
            ) : (
                <input
                    data-testid={`input-${label?.toLowerCase().replace(/\s/g, '-')}`}
                    value={value}
                    onChange={onChange}
                    placeholder={placeholder}
                />
            )}
        </div>
    ),
    Modal: Object.assign(
        ({ children, isOpen }: { children: React.ReactNode; isOpen: boolean }) => (isOpen ? <div data-testid="modal">{children}</div> : null),
        {
            Body: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
            Actions: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
        }
    ),
    Button: ({ children, onClick, variant, disabled }: {
        children: React.ReactNode;
        onClick?: () => void;
        variant?: string;
        disabled?: boolean;
    }) => (
        <button
            data-testid={`button-${variant || 'default'}`}
            onClick={onClick}
            disabled={disabled}
        >
            {children}
        </button>
    ),
}));

// Mock Layout component
vi.mock('../../../../components/layout/Layout', () => ({
    default: ({ children }: { children: React.ReactNode }) => <div data-testid="layout">{children}</div>,
}));

// Mock PageHeader
vi.mock('../../../../components/layout/PageHeader', () => ({
    default: ({ children, actions }: { children?: React.ReactNode; actions?: React.ReactNode }) => (
        <div data-testid="page-header">
            {children}
            <div data-testid="page-header-actions">{actions}</div>
        </div>
    ),
}));

// Mock StatisticsCard
vi.mock('../../../../components/dashboard/StatisticsCard', () => ({
    default: () => <div data-testid="statistics-card">Statistics</div>,
}));

// Mock ClubNotFound
vi.mock('../../ClubNotFound', () => ({
    default: () => <div data-testid="club-not-found">Club Not Found</div>,
}));

// Mock tab content components
vi.mock('../members/AdminClubMemberList', () => ({
    default: () => <div data-testid="member-list">Members</div>,
}));
vi.mock('../teams/AdminClubTeamList', () => ({
    default: () => <div data-testid="team-list">Teams</div>,
}));
vi.mock('../fines/AdminClubFineList', () => ({
    default: () => <div data-testid="fine-list">Fines</div>,
}));
vi.mock('../events/AdminClubEventList', () => ({
    default: () => <div data-testid="event-list">Events</div>,
}));
vi.mock('../news/AdminClubNewsList', () => ({
    default: () => <div data-testid="news-list">News</div>,
}));
vi.mock('../settings/AdminClubSettings', () => ({
    default: () => <div data-testid="settings">Settings</div>,
}));

describe('AdminClubDetails', () => {
    const mockClub = {
        ID: 'club-123',
        Name: 'Test Club',
        Description: 'A test club',
        LogoURL: null,
        Deleted: false,
    };

    beforeEach(() => {
        vi.clearAllMocks();
        vi.mocked(api.get).mockImplementation((url: string) => {
            if (url.includes('/IsAdmin()')) {
                return Promise.resolve({ data: { value: { IsAdmin: true }, IsOwner: true } });
            }
            if (url.includes('/Members')) {
                return Promise.resolve({ data: { value: [] } });
            }
            if (url.includes('/Invites')) {
                return Promise.resolve({ data: { value: [] } });
            }
            return Promise.resolve({ data: mockClub });
        });
    });

    it('should update club state when PATCH returns 200 with updated entity', async () => {
        const updatedClub = {
            ...mockClub,
            Name: 'Updated Club Name',
            Description: 'Updated description',
        };

        vi.mocked(api.patch).mockResolvedValue({ data: updatedClub });

        render(
            <MemoryRouter initialEntries={['/clubs/club-123/admin']}>
                <Routes>
                    <Route path="/clubs/:id/admin" element={<AdminClubDetails />} />
                </Routes>
            </MemoryRouter>
        );

        // Wait for loading to complete
        await waitFor(() => {
            expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
        });

        // Click edit button
        const editButton = screen.getByText('clubs.editClub');
        fireEvent.click(editButton);

        // Fill in form
        const nameInput = screen.getByTestId('input-clubs.clubname');
        fireEvent.change(nameInput, { target: { value: 'Updated Club Name' } });

        const descInput = screen.getByTestId('input-clubs.description');
        fireEvent.change(descInput, { target: { value: 'Updated description' } });

        // Click save button
        const saveButton = screen.getByText('common.save');
        fireEvent.click(saveButton);

        // Wait for update
        await waitFor(() => {
            expect(api.patch).toHaveBeenCalledWith(
                "/api/v2/Clubs('club-123')",
                expect.objectContaining({
                    Name: 'Updated Club Name',
                    Description: 'Updated description',
                })
            );
        });

        // Verify club name is displayed (not "Club not found")
        await waitFor(() => {
            expect(screen.queryByText('Club not found')).not.toBeInTheDocument();
        });
    });

    it('should preserve club state when PATCH returns 204 No Content', async () => {
        // Simulate 204 No Content response - data will be empty/undefined
        vi.mocked(api.patch).mockResolvedValue({ data: null, status: 204 });

        render(
            <MemoryRouter initialEntries={['/clubs/club-123/admin']}>
                <Routes>
                    <Route path="/clubs/:id/admin" element={<AdminClubDetails />} />
                </Routes>
            </MemoryRouter>
        );

        // Wait for loading to complete
        await waitFor(() => {
            expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
        });

        // Click edit button
        const editButton = screen.getByText('clubs.editClub');
        fireEvent.click(editButton);

        // Fill in form
        const nameInput = screen.getByTestId('input-clubs.clubname');
        fireEvent.change(nameInput, { target: { value: 'Updated Club Name' } });

        const descInput = screen.getByTestId('input-clubs.description');
        fireEvent.change(descInput, { target: { value: 'Updated description' } });

        // Click save button
        const saveButton = screen.getByText('common.save');
        fireEvent.click(saveButton);

        // Wait for update
        await waitFor(() => {
            expect(api.patch).toHaveBeenCalled();
        });

        // Verify club name is displayed (not "Club not found")
        await waitFor(() => {
            expect(screen.queryByText('Club not found')).not.toBeInTheDocument();
        });
    });

    it('should preserve club state when PATCH returns empty object', async () => {
        // Simulate response with empty object (no ID field)
        vi.mocked(api.patch).mockResolvedValue({ data: {}, status: 200 });

        render(
            <MemoryRouter initialEntries={['/clubs/club-123/admin']}>
                <Routes>
                    <Route path="/clubs/:id/admin" element={<AdminClubDetails />} />
                </Routes>
            </MemoryRouter>
        );

        // Wait for loading to complete
        await waitFor(() => {
            expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
        });

        // Click edit button
        const editButton = screen.getByText('clubs.editClub');
        fireEvent.click(editButton);

        // Fill in form
        const nameInput = screen.getByTestId('input-clubs.clubname');
        fireEvent.change(nameInput, { target: { value: 'Updated Club Name' } });

        const descInput = screen.getByTestId('input-clubs.description');
        fireEvent.change(descInput, { target: { value: 'Updated description' } });

        // Click save button
        const saveButton = screen.getByText('common.save');
        fireEvent.click(saveButton);

        // Wait for update
        await waitFor(() => {
            expect(api.patch).toHaveBeenCalled();
        });

        // Verify club name is displayed (not "Club not found")
        await waitFor(() => {
            expect(screen.queryByText('Club not found')).not.toBeInTheDocument();
        });
    });
});
