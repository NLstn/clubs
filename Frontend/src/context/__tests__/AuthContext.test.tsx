import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import '@testing-library/jest-dom'
import { AuthProvider } from '../AuthProvider'
import { useAuth } from '../../hooks/useAuth'

// Mock the api module
vi.mock('../../utils/api', () => ({
  default: {}
}))

// Mock fetch for logout API call
const mockFetch = vi.fn()
// Setup global fetch mock
globalThis.fetch = mockFetch

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
      <button onClick={() => logout()}>Logout</button>
    </div>
  )
}

describe('AuthContext', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.getItem.mockReturnValue(null)
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

  it('initializes with authentication when tokens exist in localStorage', () => {
    localStorageMock.getItem.mockImplementation((key) => {
      if (key === 'auth_token') return 'existing-access-token'
      if (key === 'refresh_token') return 'existing-refresh-token'
      return null
    })

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('true')
    expect(screen.getByTestId('accessToken')).toHaveTextContent('existing-access-token')
    expect(screen.getByTestId('refreshToken')).toHaveTextContent('existing-refresh-token')
  })

  it('updates state and localStorage when login is called', async () => {
    localStorageMock.getItem.mockReturnValue(null)

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    // Initially not authenticated
    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('false')

    // Trigger login
    await act(async () => {
      screen.getByText('Login').click()
    })

    // Should be authenticated now
    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('true')
    expect(screen.getByTestId('accessToken')).toHaveTextContent('new-access-token')
    expect(screen.getByTestId('refreshToken')).toHaveTextContent('new-refresh-token')

    // Should update localStorage
    expect(localStorageMock.setItem).toHaveBeenCalledWith('auth_token', 'new-access-token')
    expect(localStorageMock.setItem).toHaveBeenCalledWith('refresh_token', 'new-refresh-token')
  })

  it('clears state and localStorage when logout is called without refresh token', async () => {
    localStorageMock.getItem.mockImplementation((key) => {
      if (key === 'auth_token') return 'existing-access-token'
      if (key === 'refresh_token') return null
      return null
    })

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

    // Should clear localStorage
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_token')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')

    // Should not call logout API
    expect(mockFetch).not.toHaveBeenCalled()
  })

  it('calls logout API and clears state when logout is called with refresh token', async () => {
    localStorageMock.getItem.mockImplementation((key) => {
      if (key === 'auth_token') return 'existing-access-token'
      if (key === 'refresh_token') return 'existing-refresh-token'
      return null
    })

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
      'http://localhost:3000/api/v1/auth/keycloak/logout',
      {
        method: 'POST',
        headers: {
          'Authorization': 'existing-refresh-token',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          post_logout_redirect_uri: 'http://localhost:3000/login',
          id_token: null
        })
      }
    )

    // Should be logged out
    expect(screen.getByTestId('isAuthenticated')).toHaveTextContent('false')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_token')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')
  })

  it('handles logout API error gracefully', async () => {
    localStorageMock.getItem.mockImplementation((key) => {
      if (key === 'auth_token') return 'existing-access-token'
      if (key === 'refresh_token') return 'existing-refresh-token'
      return null
    })

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
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_token')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')
    expect(consoleSpy).toHaveBeenCalledWith('Error during logout:', expect.any(Error))

    consoleSpy.mockRestore()
  })
})