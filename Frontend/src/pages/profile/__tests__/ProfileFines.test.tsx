import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import '@testing-library/jest-dom'
import ProfileFines from '../ProfileFines'

// Mock the api module
vi.mock('../../../utils/api', () => ({
  default: {
    get: vi.fn()
  }
}))

// Mock Layout component to avoid AuthProvider dependency
vi.mock('../../../components/layout/Layout', () => ({
  default: ({ children, title }: { children: React.ReactNode; title: string }) => (
    <div data-testid="layout" data-title={title}>
      {children}
    </div>
  ),
}))

// Mock ProfileContentLayout to simplify testing
vi.mock('../../../components/layout/ProfileContentLayout', () => ({
  default: ({ title, children }: { title: string; children: React.ReactNode }) => (
    <div data-testid="profile-content-layout">
      <h1>{title}</h1>
      {children}
    </div>
  )
}))

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  )
}

describe('ProfileFines', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders with empty fines list without crashing', async () => {
    const { default: api } = await import('../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock API to return empty array
    mockGet.mockResolvedValue({ status: 200, data: [] })

    renderWithRouter(<ProfileFines />)

    // Check that basic UI elements are rendered
    expect(screen.getByText('Fines')).toBeInTheDocument()

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/api/v1/me/fines')
    })

    // Verify empty message is shown
    expect(screen.getByText('No fines found')).toBeInTheDocument()
    
    // Should not display footer with empty fines
    expect(screen.queryByText(/Total:/)).not.toBeInTheDocument()
  })

  it('renders fines and displays total amount in footer', async () => {
    const { default: api } = await import('../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    const mockFines = [
      {
        id: '1',
        clubName: 'Test Club A',
        amount: 25.50,
        reason: 'Late arrival',
        createdAt: '2024-01-01T10:00:00Z',
        updatedAt: '2024-01-01T10:00:00Z',
        paid: false
      },
      {
        id: '2',
        clubName: 'Test Club B',
        amount: 10.00,
        reason: 'Missed meeting',
        createdAt: '2024-01-02T10:00:00Z',
        updatedAt: '2024-01-02T10:00:00Z',
        paid: true
      },
      {
        id: '3',
        clubName: 'Test Club C',
        amount: 15.75,
        reason: 'Equipment damage',
        createdAt: '2024-01-03T10:00:00Z',
        updatedAt: '2024-01-03T10:00:00Z',
        paid: false
      }
    ]

    mockGet.mockResolvedValue({ status: 200, data: mockFines })

    renderWithRouter(<ProfileFines />)

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/api/v1/me/fines')
    })

    // Check that fines are rendered
    await waitFor(() => {
      expect(screen.getByText('Test Club A')).toBeInTheDocument()
      expect(screen.getByText('Test Club B')).toBeInTheDocument()
      expect(screen.getByText('Test Club C')).toBeInTheDocument()
      expect(screen.getByText('€25.50')).toBeInTheDocument()
      expect(screen.getByText('€10.00')).toBeInTheDocument()
      expect(screen.getByText('€15.75')).toBeInTheDocument()
    })

    // Check that total amount is displayed in footer
    // Total should be 25.50 + 10.00 + 15.75 = 51.25
    await waitFor(() => {
      expect(screen.getByText('Total: €51.25')).toBeInTheDocument()
    })
  })

  it('calculates total correctly with decimal amounts', async () => {
    const { default: api } = await import('../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    const mockFines = [
      {
        id: '1',
        clubName: 'Test Club',
        amount: 12.33,
        reason: 'Test reason',
        createdAt: '2024-01-01T10:00:00Z',
        updatedAt: '2024-01-01T10:00:00Z',
        paid: false
      },
      {
        id: '2',
        clubName: 'Test Club',
        amount: 7.67,
        reason: 'Test reason',
        createdAt: '2024-01-02T10:00:00Z',
        updatedAt: '2024-01-02T10:00:00Z',
        paid: true
      }
    ]

    mockGet.mockResolvedValue({ status: 200, data: mockFines })

    renderWithRouter(<ProfileFines />)

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/api/v1/me/fines')
    })

    // Check that total amount is calculated correctly
    // Total should be 12.33 + 7.67 = 20.00
    await waitFor(() => {
      expect(screen.getByText('Total: €20.00')).toBeInTheDocument()
    })
  })

  it('displays error message when API call fails', async () => {
    const { default: api } = await import('../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    // Mock API to reject
    const mockError = new Error('Network error')
    mockGet.mockRejectedValue(mockError)

    renderWithRouter(<ProfileFines />)

    // Wait for API call to complete and error to be displayed
    await waitFor(() => {
      expect(screen.getByText('Failed to load fines')).toBeInTheDocument()
    })

    // Should not display footer when there's an error
    expect(screen.queryByText(/Total:/)).not.toBeInTheDocument()
  })

  it('handles single fine correctly', async () => {
    const { default: api } = await import('../../../utils/api')
    const mockGet = vi.mocked(api.get)
    
    const mockFines = [
      {
        id: '1',
        clubName: 'Test Club',
        amount: 42.99,
        reason: 'Single fine',
        createdAt: '2024-01-01T10:00:00Z',
        updatedAt: '2024-01-01T10:00:00Z',
        paid: false
      }
    ]

    mockGet.mockResolvedValue({ status: 200, data: mockFines })

    renderWithRouter(<ProfileFines />)

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/api/v1/me/fines')
    })

    // Check that fine is rendered
    await waitFor(() => {
      expect(screen.getByText('Test Club')).toBeInTheDocument()
      expect(screen.getByText('€42.99')).toBeInTheDocument()
    })

    // Check that total matches the single fine amount
    await waitFor(() => {
      expect(screen.getByText('Total: €42.99')).toBeInTheDocument()
    })
  })
})