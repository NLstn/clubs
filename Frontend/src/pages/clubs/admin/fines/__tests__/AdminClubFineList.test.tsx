import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, act } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import '@testing-library/jest-dom'
import AdminClubFineList from '../AdminClubFineList'

// Mock the api module
vi.mock('../../../../../utils/api', () => ({
  default: {
    get: vi.fn()
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

// Mock child components to avoid their complexity in unit tests
vi.mock('../AddFine', () => ({
  default: ({ isOpen }: { isOpen: boolean }) => (
    <div data-testid="add-fine-modal">{isOpen ? 'Modal Open' : 'Modal Closed'}</div>
  )
}))

vi.mock('../AdminClubFineTemplateList', () => ({
  default: () => <div data-testid="fine-template-list">Fine Template List</div>
}))

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  )
}

describe('AdminClubFineList', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders with empty fines list without crashing', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock API to return empty OData response
    mockGet.mockResolvedValue({ data: { value: [] } })

    renderWithRouter(<AdminClubFineList />)

    // Check that basic UI elements are rendered
    expect(screen.getByText('Fines')).toBeInTheDocument()
    expect(screen.getByText('Show all fines')).toBeInTheDocument()
    expect(screen.getByText('Manage Templates')).toBeInTheDocument()
    expect(screen.getByText('Add Fine')).toBeInTheDocument()

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/api/v2/Fines?$filter=ClubID eq \'test-club-id\'&$expand=User')
    })

    // After loading is complete, check that table headers are rendered
    await waitFor(() => {
      expect(screen.getByText('User')).toBeInTheDocument()
      expect(screen.getByText('Amount')).toBeInTheDocument()
      expect(screen.getByText('Reason')).toBeInTheDocument()
      expect(screen.getByText('Created At')).toBeInTheDocument()
      expect(screen.getByText('Updated At')).toBeInTheDocument()
      expect(screen.getByText('Paid')).toBeInTheDocument()
    })

    // Verify no error is displayed
    expect(screen.queryByText(/Failed to load fines/)).not.toBeInTheDocument()
    // Verify empty message is shown
    expect(screen.getByText('No fines available')).toBeInTheDocument()
  })

  it('handles null API response without crashing', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock API to return null data
    mockGet.mockResolvedValue({ data: null })

    renderWithRouter(<AdminClubFineList />)

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/api/v2/Fines?$filter=ClubID eq \'test-club-id\'&$expand=User')
    })

    // Component should still render basic UI elements
    expect(screen.getByText('Fines')).toBeInTheDocument()
    expect(screen.getByText('Add Fine')).toBeInTheDocument()
  })

  it('handles undefined API response without crashing', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock API to return undefined data
    mockGet.mockResolvedValue({ data: undefined })

    renderWithRouter(<AdminClubFineList />)

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/api/v2/Fines?$filter=ClubID eq \'test-club-id\'&$expand=User')
    })

    // Component should still render basic UI elements
    expect(screen.getByText('Fines')).toBeInTheDocument()
    expect(screen.getByText('Add Fine')).toBeInTheDocument()
  })

  it('displays error message when API call fails', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock API to reject
    const mockError = new Error('Network error')
    mockGet.mockRejectedValue(mockError)

    renderWithRouter(<AdminClubFineList />)

    // Wait for API call to complete and error to be displayed
    await waitFor(() => {
      expect(screen.getByText('Failed to load fines')).toBeInTheDocument()
    })

    // Component should still render basic UI elements
    expect(screen.getByText('Fines')).toBeInTheDocument()
    expect(screen.getByText('Add Fine')).toBeInTheDocument()
  })

  it('renders fines when API returns valid data', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    const mockFines = {
      value: [
        {
          ID: '1',
          Amount: 25.50,
          Reason: 'Late arrival',
          CreatedAt: '2024-01-01T10:00:00Z',
          UpdatedAt: '2024-01-01T10:00:00Z',
          Paid: false,
          User: { Name: 'John Doe' }
        },
        {
          ID: '2',
          Amount: 10.00,
          Reason: 'Missed meeting',
          CreatedAt: '2024-01-02T10:00:00Z',
          UpdatedAt: '2024-01-02T10:00:00Z',
          Paid: true,
          User: { Name: 'Jane Smith' }
        }
      ]
    }

    mockGet.mockResolvedValue({ data: mockFines })

    renderWithRouter(<AdminClubFineList />)

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/api/v2/Fines?$filter=ClubID eq \'test-club-id\'&$expand=User')
    })

    // Check that fines are rendered (by default only unpaid fines)
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument()
      expect(screen.getByText('25.5')).toBeInTheDocument()
      expect(screen.getByText('Late arrival')).toBeInTheDocument()
    })

    // Jane Smith's paid fine should not be visible by default
    expect(screen.queryByText('Jane Smith')).not.toBeInTheDocument()
  })

  it('shows all fines when "Show all fines" is checked', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    const mockFines = {
      value: [
        {
          ID: '1',
          Amount: 25.50,
          Reason: 'Late arrival',
          CreatedAt: '2024-01-01T10:00:00Z',
          UpdatedAt: '2024-01-01T10:00:00Z',
          Paid: false,
          User: { Name: 'John Doe' }
        },
        {
          ID: '2',
          Amount: 10.00,
          Reason: 'Missed meeting',
          CreatedAt: '2024-01-02T10:00:00Z',
          UpdatedAt: '2024-01-02T10:00:00Z',
          Paid: true,
          User: { Name: 'Jane Smith' }
        }
      ]
    }

    mockGet.mockResolvedValue({ data: mockFines })

    renderWithRouter(<AdminClubFineList />)

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/api/v2/Fines?$filter=ClubID eq \'test-club-id\'&$expand=User')
    })

    // Initially only unpaid fine should be visible
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument()
    })
    expect(screen.queryByText('Jane Smith')).not.toBeInTheDocument()

    // Check the "Show all fines" checkbox
    const showAllCheckbox = screen.getByRole('checkbox', { name: /show all fines/i })
    act(() => {
      showAllCheckbox.click()
    })

    // Now both fines should be visible
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument()
      expect(screen.getByText('Jane Smith')).toBeInTheDocument()
    })
  })
})