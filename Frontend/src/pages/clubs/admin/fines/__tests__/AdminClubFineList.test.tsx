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
    
    // Mock API to return empty OData response with count
    mockGet.mockResolvedValue({ data: { value: [], '@odata.count': 0 } })

    renderWithRouter(<AdminClubFineList />)

    // Check that basic UI elements are rendered
    expect(screen.getByText('Fines')).toBeInTheDocument()
    expect(screen.getByText('Show all fines')).toBeInTheDocument()
    expect(screen.getByText('Manage Templates')).toBeInTheDocument()
    expect(screen.getByText('Add Fine')).toBeInTheDocument()

    // Wait for API call to complete - ODataTable now makes calls with pagination params
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
      const callArg = mockGet.mock.calls[0][0] as string
      expect(callArg).toContain('/api/v2/Fines')
      expect(callArg).toContain('$filter=ClubID eq \'test-club-id\' and Paid eq false')
      expect(callArg).toContain('$expand=User')
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
      expect(mockGet).toHaveBeenCalled()
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
      expect(mockGet).toHaveBeenCalled()
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
    // ODataTable shows "Error loading data" by default
    await waitFor(() => {
      expect(screen.getByText('Error loading data')).toBeInTheDocument()
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
          User: { FirstName: 'John', LastName: 'Doe' }
        }
      ],
      '@odata.count': 1
    }

    mockGet.mockResolvedValue({ data: mockFines })

    renderWithRouter(<AdminClubFineList />)

    // Wait for API call to complete - ODataTable filters to Paid eq false by default
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
      const callArg = mockGet.mock.calls[0][0] as string
      expect(callArg).toContain('/api/v2/Fines')
      expect(callArg).toContain('$filter=ClubID eq \'test-club-id\' and Paid eq false')
    })

    // Check that fines are rendered (only unpaid fines by default)
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument()
      expect(screen.getByText('25.5')).toBeInTheDocument()
      expect(screen.getByText('Late arrival')).toBeInTheDocument()
    })
  })

  it('shows all fines when "Show all fines" is checked', async () => {
    const { default: api } = await import('../../../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    const mockUnpaidFines = {
      value: [
        {
          ID: '1',
          Amount: 25.50,
          Reason: 'Late arrival',
          CreatedAt: '2024-01-01T10:00:00Z',
          UpdatedAt: '2024-01-01T10:00:00Z',
          Paid: false,
          User: { FirstName: 'John', LastName: 'Doe' }
        }
      ],
      '@odata.count': 1
    }

    const mockAllFines = {
      value: [
        {
          ID: '1',
          Amount: 25.50,
          Reason: 'Late arrival',
          CreatedAt: '2024-01-01T10:00:00Z',
          UpdatedAt: '2024-01-01T10:00:00Z',
          Paid: false,
          User: { FirstName: 'John', LastName: 'Doe' }
        },
        {
          ID: '2',
          Amount: 10.00,
          Reason: 'Missed meeting',
          CreatedAt: '2024-01-02T10:00:00Z',
          UpdatedAt: '2024-01-02T10:00:00Z',
          Paid: true,
          User: { FirstName: 'Jane', LastName: 'Smith' }
        }
      ],
      '@odata.count': 2
    }

    // First call returns unpaid fines
    mockGet.mockResolvedValueOnce({ data: mockUnpaidFines })

    renderWithRouter(<AdminClubFineList />)

    // Wait for initial API call (unpaid fines only)
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })

    // Initially only unpaid fine should be visible
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument()
    })

    // Second call after checkbox returns all fines
    mockGet.mockResolvedValueOnce({ data: mockAllFines })

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