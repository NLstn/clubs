import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import '@testing-library/jest-dom'
import AdminClubTeamList from '../AdminClubTeamList'

// Mock the api module
vi.mock('../../../../../utils/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    patch: vi.fn()
  }
}))

// Mock react-router-dom useParams
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useParams: () => ({ id: 'test-club-id' })
  }
})

// Mock the translation hook
vi.mock('../../../../../hooks/useTranslation', () => ({
  useT: () => ({
    t: (key: string, params?: { [key: string]: unknown }) => {
      const translations: Record<string, string> = {
        'teams.title': 'Teams',
        'teams.createTeam': 'Create Team',
        'teams.addMember': 'Add Member',
        'teams.membersOf': `Members of ${params?.teamName || 'Team'}`,
        'teams.noMembers': 'No team members yet.',
        'common.name': 'Name',
        'common.role': 'Role',
        'common.actions': 'Actions',
        'teams.joinedAt': 'Joined',
        'teams.roles.admin': 'Admin',
        'teams.roles.member': 'Member',
        'teams.promoteToAdmin': 'Promote to Admin',
        'teams.demoteToMember': 'Demote to Member',
        'common.remove': 'Remove',
        'common.edit': 'Edit',
        'common.delete': 'Delete'
      }
      return translations[key] || key
    }
  })
}))

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  )
}

const mockTeams = [
  {
    id: 'team-1',
    name: 'Empty Team',
    description: 'A team with no members',
    createdAt: '2024-01-01T00:00:00Z',
    clubId: 'test-club-id'
  },
  {
    id: 'team-2',
    name: 'Team with Members',
    description: 'A team with members',
    createdAt: '2024-01-01T00:00:00Z',
    clubId: 'test-club-id'
  }
]

const mockClubMembers = [
  {
    id: 'member-1',
    userId: 'user-1',
    name: 'John Doe',
    role: 'member'
  },
  {
    id: 'member-2',
    userId: 'user-2',
    name: 'Jane Smith',
    role: 'admin'
  }
]

const mockTeamMembers = {
  value: [
    {
      ID: 'team-member-1',
      UserID: 'user-1',
      Role: 'member',
      CreatedAt: '2024-01-01T00:00:00Z',
      User: {
        FirstName: 'John',
        LastName: 'Doe'
      }
    }
  ]
}

describe('AdminClubTeamList', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Mock console methods to avoid cluttering test output
    vi.spyOn(console, 'error').mockImplementation(() => {})
    vi.spyOn(console, 'warn').mockImplementation(() => {})
  })

  it('renders without crashing when loading', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock APIs to never resolve (simulate loading state)
    mockGet.mockImplementation(() => new Promise(() => {}))

    renderWithRouter(<AdminClubTeamList />)

    expect(screen.getByText('Loading teams...')).toBeInTheDocument()
  })

  it('renders teams list successfully', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock successful API responses
    mockGet
      .mockResolvedValueOnce({ data: mockTeams }) // fetchTeams
      .mockResolvedValueOnce({ data: mockClubMembers }) // fetchClubMembers

    renderWithRouter(<AdminClubTeamList />)

    await waitFor(() => {
      expect(screen.getByText('Teams')).toBeInTheDocument()
      expect(screen.getByText('Create Team')).toBeInTheDocument()
      expect(screen.getByText('Empty Team')).toBeInTheDocument()
      expect(screen.getByText('Team with Members')).toBeInTheDocument()
    })
  })

  it('handles team with no members without crashing', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock successful API responses
    mockGet
      .mockResolvedValueOnce({ data: mockTeams }) // fetchTeams
      .mockResolvedValueOnce({ data: mockClubMembers }) // fetchClubMembers
      .mockResolvedValueOnce({ data: [] }) // fetchTeamMembers - empty array

    renderWithRouter(<AdminClubTeamList />)

    // Wait for initial load
    await waitFor(() => {
      expect(screen.getByText('Empty Team')).toBeInTheDocument()
    })

    // Click on the empty team
    fireEvent.click(screen.getByText('Empty Team'))

    // Wait for team members section to appear
    await waitFor(() => {
      expect(screen.getByText('Members of Empty Team')).toBeInTheDocument()
      expect(screen.getByText('No team members yet.')).toBeInTheDocument()
      expect(screen.getByText('Add Member')).toBeInTheDocument()
    })

    // Ensure table headers are present
    expect(screen.getByText('Name')).toBeInTheDocument()
    expect(screen.getByText('Role')).toBeInTheDocument()
    expect(screen.getByText('Joined')).toBeInTheDocument()
    expect(screen.getByText('Actions')).toBeInTheDocument()
  })

  it('handles team with null members response without crashing', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock successful API responses
    mockGet
      .mockResolvedValueOnce({ data: mockTeams }) // fetchTeams
      .mockResolvedValueOnce({ data: mockClubMembers }) // fetchClubMembers
      .mockResolvedValueOnce({ data: null }) // fetchTeamMembers - null response

    renderWithRouter(<AdminClubTeamList />)

    // Wait for initial load
    await waitFor(() => {
      expect(screen.getByText('Empty Team')).toBeInTheDocument()
    })

    // Click on the empty team
    fireEvent.click(screen.getByText('Empty Team'))

    // Wait for team members section to appear
    await waitFor(() => {
      expect(screen.getByText('Members of Empty Team')).toBeInTheDocument()
      expect(screen.getByText('No team members yet.')).toBeInTheDocument()
    })

    // The component should not crash and should show empty state
    expect(screen.queryByText('TypeError')).not.toBeInTheDocument()
  })

  it('handles API error when fetching team members', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock successful API responses for initial load
    mockGet
      .mockResolvedValueOnce({ data: mockTeams }) // fetchTeams
      .mockResolvedValueOnce({ data: mockClubMembers }) // fetchClubMembers
      .mockRejectedValueOnce(new Error('API Error')) // fetchTeamMembers - error

    renderWithRouter(<AdminClubTeamList />)

    // Wait for initial load
    await waitFor(() => {
      expect(screen.getByText('Empty Team')).toBeInTheDocument()
    })

    // Click on the empty team
    fireEvent.click(screen.getByText('Empty Team'))

    // Should show error but not crash
    await waitFor(() => {
      expect(screen.getByText('Failed to fetch team members')).toBeInTheDocument()
    })
  })

  it('renders team with members correctly', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock successful API responses
    mockGet
      .mockResolvedValueOnce({ data: mockTeams }) // fetchTeams
      .mockResolvedValueOnce({ data: mockClubMembers }) // fetchClubMembers
      .mockResolvedValueOnce({ data: mockTeamMembers }) // fetchTeamMembers

    renderWithRouter(<AdminClubTeamList />)

    // Wait for initial load
    await waitFor(() => {
      expect(screen.getByText('Team with Members')).toBeInTheDocument()
    })

    // Click on the team with members
    fireEvent.click(screen.getByText('Team with Members'))

    // Wait for team members section to appear
    await waitFor(() => {
      expect(screen.getByText('Members of Team with Members')).toBeInTheDocument()
      expect(screen.getByText('John Doe')).toBeInTheDocument()
      expect(screen.getByText('Member')).toBeInTheDocument()
    })

    // Should not show empty state message
    expect(screen.queryByText('No team members yet.')).not.toBeInTheDocument()
  })

  it('shows add member modal with available members when team has no members', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock successful API responses
    mockGet
      .mockResolvedValueOnce({ data: mockTeams }) // fetchTeams
      .mockResolvedValueOnce({ data: mockClubMembers }) // fetchClubMembers
      .mockResolvedValueOnce({ data: [] }) // fetchTeamMembers - empty

    renderWithRouter(<AdminClubTeamList />)

    // Wait for initial load and select empty team
    await waitFor(() => {
      expect(screen.getByText('Empty Team')).toBeInTheDocument()
    })

    fireEvent.click(screen.getByText('Empty Team'))

    await waitFor(() => {
      expect(screen.getByText('Add Member')).toBeInTheDocument()
    })

    // Click Add Member button
    fireEvent.click(screen.getByText('Add Member'))

    // Should show all club members as available (since team is empty)
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument()
      expect(screen.getByText('Jane Smith')).toBeInTheDocument()
    })
  })

  it('prevents regression: teamMembers.map crash when teamMembers is null', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock APIs returning data that could cause the original crash
    mockGet
      .mockResolvedValueOnce({ data: mockTeams })
      .mockResolvedValueOnce({ data: mockClubMembers })
      .mockResolvedValueOnce({ data: null }) // This would have caused the crash

    // This test specifically ensures the component doesn't crash
    expect(() => {
      renderWithRouter(<AdminClubTeamList />)
    }).not.toThrow()

    await waitFor(() => {
      expect(screen.getByText('Empty Team')).toBeInTheDocument()
    })

    // Clicking on team should not crash
    expect(() => {
      fireEvent.click(screen.getByText('Empty Team'))
    }).not.toThrow()

    // Component should handle null gracefully
    await waitFor(() => {
      expect(screen.getByText('No team members yet.')).toBeInTheDocument()
    })
  })
})
