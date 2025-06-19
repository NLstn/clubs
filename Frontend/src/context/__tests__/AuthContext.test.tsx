import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import { AuthProvider } from '../AuthProvider'
import { useAuth } from '../../hooks/useAuth'

// Mock the api module
vi.mock('../../utils/api', () => ({
  default: {}
}))

// Mock fetch for logout API call
const mockFetch = vi.fn()
global.fetch = mockFetch

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

// Mock environment variable
vi.stubEnv('VITE_API_HOST', 'http://localhost:3000')

// Mock document.cookie
const mockCookie = {
  get: vi.fn(() => ''),
  set: vi.fn(),
}
Object.defineProperty(document, 'cookie', {
  get: () => mockCookie.get(),
  set: (value) => mockCookie.set(value),
  configurable: true,
})

// Test component that uses the AuthContext
const TestComponent = () => {
  const { isAuthenticated, accessToken, refreshToken, login, logout } = useAuth()
  
  return (
    <div>
      <div data-testid="isAuthenticated">{isAuthenticated.toString()}</div>
      <div data-testid="accessToken">{accessToken || 'null'}</div>
      <div data-testid="refreshToken">{refreshToken || 'null'}</div>
      <button onClick={() => login('new-access-token', 'new-refresh-token')}>
        Login
      </button>
      <button onClick={logout}>Logout</button>
    </div>
  )
}

describe('AuthContext', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.getItem.mockReturnValue(null)
    mockCookie.get.mockReturnValue('')
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('throws error when useAuth is used outside AuthProvider', () => {
    // Suppress console.error for this test
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    
    expect(() => {
      render(<TestComponent />)
    }).toThrow('useAuth must be used within an AuthProvider')
    
    consoleSpy.mockRestore()
  })

  it('initializes with no authentication when no tokens in localStorage', () => {
    localStorageMock.getItem.mockReturnValue(null)

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('false')
    expect(screen.getByTestId('accessToken')).toHaveTextContent('null')
    expect(screen.getByTestId('refreshToken')).toHaveTextContent('null')
  })

  it('initializes with authentication when tokens exist in cookies', () => {
    mockCookie.get.mockReturnValue('access_token=existing-access-token; refresh_token=existing-refresh-token')

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('true')
    expect(screen.getByTestId('accessToken')).toHaveTextContent('existing-access-token')
    expect(screen.getByTestId('refreshToken')).toHaveTextContent('existing-refresh-token')
  })

  it('updates state when login is called', () => {
    mockCookie.get.mockReturnValue('')

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    // Initially not authenticated
    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('false')

    // Trigger login
    act(() => {
      screen.getByText('Login').click()
    })

    // Should be authenticated now
    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('true')
    expect(screen.getByTestId('accessToken')).toHaveTextContent('new-access-token')
    expect(screen.getByTestId('refreshToken')).toHaveTextContent('new-refresh-token')

    // No localStorage calls should be made since we only use cookies
    expect(localStorageMock.setItem).not.toHaveBeenCalled()
  })

  it('clears state when logout is called without refresh token', async () => {
    mockCookie.get.mockReturnValue('access_token=existing-access-token')

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    // Initially authenticated
    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('true')

    // Trigger logout
    await act(async () => {
      screen.getByText('Logout').click()
    })

    // Should be logged out
    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('false')
    expect(screen.getByTestId('accessToken')).toHaveTextContent('null')
    expect(screen.getByTestId('refreshToken')).toHaveTextContent('null')

    // No localStorage calls should be made since we only use cookies
    expect(localStorageMock.removeItem).not.toHaveBeenCalled()

    // Should not call logout API
    expect(mockFetch).not.toHaveBeenCalled()
  })

  it('calls logout API and clears state when logout is called with refresh token', async () => {
    mockCookie.get.mockReturnValue('access_token=existing-access-token; refresh_token=existing-refresh-token')

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({})
    })

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    // Trigger logout
    await act(async () => {
      screen.getByText('Logout').click()
    })

    // Should call logout API
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:3000/api/v1/auth/logout',
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include'
      }
    )

    // Should be logged out
    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('false')
    // No localStorage calls should be made since we only use cookies
    expect(localStorageMock.removeItem).not.toHaveBeenCalled()
  })

  it('handles logout API error gracefully', async () => {
    mockCookie.get.mockReturnValue('access_token=existing-access-token; refresh_token=existing-refresh-token')

    mockFetch.mockRejectedValueOnce(new Error('Network error'))
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    // Trigger logout
    await act(async () => {
      screen.getByText('Logout').click()
    })

    // Should still clear state even if API call fails
    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('false')
    // No localStorage calls should be made since we only use cookies
    expect(localStorageMock.removeItem).not.toHaveBeenCalled()
    expect(consoleSpy).toHaveBeenCalledWith('Error during logout:', expect.any(Error))

    consoleSpy.mockRestore()
  })
})