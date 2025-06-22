import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import Header from '../layout/Header'

// Mock the useAuth hook
const mockUseAuth = vi.fn()
const mockLogout = vi.fn()
const mockNavigate = vi.fn()

vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth()
}))

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate
  }
})

// Mock the logo import
vi.mock('../../assets/logo.png', () => ({
  default: 'mock-logo.png'
}))

// Mock the RecentClubsDropdown component
vi.mock('../layout/RecentClubsDropdown', () => ({
  default: () => <div data-testid="recent-clubs-dropdown">Recent Clubs Dropdown</div>
}))

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  )
}

describe('Header', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockUseAuth.mockReturnValue({
      logout: mockLogout
    })
  })

  it('renders with default title', () => {
    renderWithRouter(<Header />)
    expect(screen.getByText('Clubs')).toBeInTheDocument()
  })

  it('renders with custom title', () => {
    renderWithRouter(<Header title="Custom Title" />)
    expect(screen.getByText('Custom Title')).toBeInTheDocument()
  })

  it('renders logo and navigates to home when clicked', () => {
    renderWithRouter(<Header />)
    
    const logo = screen.getByAltText('Logo')
    expect(logo).toBeInTheDocument()
    expect(logo).toHaveAttribute('src', 'mock-logo.png')
    
    fireEvent.click(logo)
    expect(mockNavigate).toHaveBeenCalledWith('/')
  })

  it('opens dropdown when user icon is clicked', () => {
    renderWithRouter(<Header />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    expect(screen.getByText('Profile')).toBeInTheDocument()
    expect(screen.getByText('Create New Club')).toBeInTheDocument()
    expect(screen.getByText('Logout')).toBeInTheDocument()
  })

  it('shows admin panel link when user is club admin', () => {
    renderWithRouter(<Header isClubAdmin={true} clubId="123" />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    expect(screen.getByText('Admin Panel')).toBeInTheDocument()
  })

  it('does not show admin panel link when user is not club admin', () => {
    renderWithRouter(<Header isClubAdmin={false} clubId="123" />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    expect(screen.queryByText('Admin Panel')).not.toBeInTheDocument()
  })

  it('navigates to admin panel when admin panel is clicked', () => {
    renderWithRouter(<Header isClubAdmin={true} clubId="123" />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    const adminButton = screen.getByText('Admin Panel')
    fireEvent.click(adminButton)
    
    expect(mockNavigate).toHaveBeenCalledWith('/clubs/123/admin')
  })

  it('navigates to profile when profile is clicked', () => {
    renderWithRouter(<Header />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    const profileButton = screen.getByText('Profile')
    fireEvent.click(profileButton)
    
    expect(mockNavigate).toHaveBeenCalledWith('/profile')
  })

  it('navigates to create club when create new club is clicked', () => {
    renderWithRouter(<Header />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    const createClubButton = screen.getByText('Create New Club')
    fireEvent.click(createClubButton)
    
    expect(mockNavigate).toHaveBeenCalledWith('/createClub')
  })

  it('calls logout and navigates to login when logout is clicked', async () => {
    renderWithRouter(<Header />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    const logoutButton = screen.getByText('Logout')
    fireEvent.click(logoutButton)
    
    expect(mockLogout).toHaveBeenCalled()
    expect(mockNavigate).toHaveBeenCalledWith('/login')
  })

  it('shows recent clubs dropdown when showRecentClubs is true', () => {
    renderWithRouter(<Header showRecentClubs={true} />)
    
    expect(screen.getByTestId('recent-clubs-dropdown')).toBeInTheDocument()
  })

  it('does not show recent clubs dropdown when showRecentClubs is false', () => {
    renderWithRouter(<Header showRecentClubs={false} />)
    
    expect(screen.queryByTestId('recent-clubs-dropdown')).not.toBeInTheDocument()
  })

  it('shows My Clubs when showRecentClubs is false', () => {
    renderWithRouter(<Header showRecentClubs={false} />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    expect(screen.getByText('My Clubs')).toBeInTheDocument()
  })

  it('does not show My Clubs when showRecentClubs is true', () => {
    renderWithRouter(<Header showRecentClubs={true} />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    expect(screen.queryByText('My Clubs')).not.toBeInTheDocument()
  })

  it('closes dropdown when clicking outside', () => {
    renderWithRouter(<Header />)
    
    const userIcon = screen.getByText('U')
    fireEvent.click(userIcon)
    
    // Dropdown should be open
    expect(screen.getByText('Profile')).toBeInTheDocument()
    
    // Click outside the dropdown
    fireEvent.mouseDown(document.body)
    
    // Dropdown should be closed
    expect(screen.queryByText('Profile')).not.toBeInTheDocument()
  })
})