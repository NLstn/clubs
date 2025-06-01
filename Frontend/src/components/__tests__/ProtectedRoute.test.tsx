import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import ProtectedRoute from '../auth/ProtectedRoute'
import { AuthProvider } from '../../context/AuthProvider'

// Mock the useAuth hook
const mockUseAuth = vi.fn()

vi.mock('../../context/AuthContext', () => ({
  AuthProvider: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  useAuth: () => mockUseAuth()
}))

// Mock Navigate component from react-router-dom
const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    Navigate: ({ to, replace }: { to: string; replace?: boolean }) => {
      mockNavigate(to, replace)
      return <div data-testid="navigate">Redirecting to {to}</div>
    }
  }
})

const TestChild = () => <div data-testid="protected-content">Protected Content</div>

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      <AuthProvider>
        {component}
      </AuthProvider>
    </BrowserRouter>
  )
}

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders children when user is authenticated', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      accessToken: 'valid-token',
      refreshToken: 'valid-refresh-token',
      login: vi.fn(),
      logout: vi.fn(),
      api: {}
    })

    renderWithRouter(
      <ProtectedRoute>
        <TestChild />
      </ProtectedRoute>
    )

    expect(screen.getByTestId('protected-content')).toBeInTheDocument()
    expect(screen.queryByTestId('navigate')).not.toBeInTheDocument()
  })

  it('redirects to login when user is not authenticated', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      accessToken: null,
      refreshToken: null,
      login: vi.fn(),
      logout: vi.fn(),
      api: {}
    })

    renderWithRouter(
      <ProtectedRoute>
        <TestChild />
      </ProtectedRoute>
    )

    expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument()
    expect(screen.getByTestId('navigate')).toBeInTheDocument()
    expect(mockNavigate).toHaveBeenCalledWith('/login', true)
  })
})