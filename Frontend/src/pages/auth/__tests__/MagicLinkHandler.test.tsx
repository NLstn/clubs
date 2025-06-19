import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import MagicLinkHandler from '../MagicLinkHandler'
import { AuthProvider } from '../../../context/AuthProvider'
import { useAuth } from '../../../hooks/useAuth'

// Mock the api module
vi.mock('../../../utils/api', () => ({
  default: {}
}))

// Mock environment variable
vi.stubEnv('VITE_API_HOST', 'http://localhost:3000')

// Mock fetch for the verifyMagicLink API call
const mockFetch = vi.fn()
global.fetch = mockFetch

// Mock document.cookie with a more realistic implementation
let mockCookieStore = ''
Object.defineProperty(document, 'cookie', {
  get: () => mockCookieStore,
  set: (value) => {
    // Parse and add the cookie to our mock store
    const [cookieDef] = value.split(';')
    const [name, cookieValue] = cookieDef.split('=')
    const existingCookies = mockCookieStore ? mockCookieStore.split('; ') : []
    
    // Remove existing cookie with same name
    const filteredCookies = existingCookies.filter(c => !c.startsWith(name + '='))
    
    // Add new cookie
    if (cookieValue && cookieValue !== '') {
      filteredCookies.push(`${name}=${cookieValue}`)
    }
    
    mockCookieStore = filteredCookies.join('; ')
    console.log('Cookie set:', name, '=', cookieValue, 'Full cookie string:', mockCookieStore)
  },
  configurable: true,
})

// Mock CookieConsent component
vi.mock('../../../components/CookieConsent', () => ({
  default: () => <div data-testid="cookie-consent">Cookie Consent</div>
}))

// Helper component to test authentication state
const AuthStateIndicator = () => {
  const { isAuthenticated, accessToken } = useAuth()
  return (
    <div>
      <div data-testid="auth-status">{isAuthenticated ? 'authenticated' : 'not-authenticated'}</div>
      <div data-testid="access-token">{accessToken || 'no-token'}</div>
    </div>
  )
}

describe('MagicLinkHandler Integration', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockCookieStore = ''
    // Mock console.log to reduce noise in tests
    vi.spyOn(console, 'log').mockImplementation(() => {})
    vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  afterEach(() => {
    vi.clearAllMocks()
    vi.restoreAllMocks()
  })

  it('should authenticate user after successful magic link verification with immediate cookie availability', async () => {
    // Mock successful API response
    mockFetch.mockResolvedValueOnce({
      ok: true,
      text: async () => 'success'
    })

    // Simulate immediate cookie availability
    Object.defineProperty(document, 'cookie', {
      get: () => 'access_token=mock-access-token; refresh_token=mock-refresh-token',
      set: () => {},
      configurable: true,
    })

    render(
      <AuthProvider>
        <MemoryRouter initialEntries={['/auth/magic?token=valid-token']}>
          <div>
            <MagicLinkHandler />
            <AuthStateIndicator />
          </div>
        </MemoryRouter>
      </AuthProvider>
    )

    // Wait for success message and authentication
    await waitFor(() => {
      expect(screen.getByText('Login Successful!')).toBeInTheDocument()
    })

    await waitFor(() => {
      expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated')
    })
  })

  it('should authenticate user after successful magic link verification with delayed cookie availability', async () => {
    // Mock successful API response
    mockFetch.mockResolvedValueOnce({
      ok: true,
      text: async () => 'success'
    })

    // Simulate delayed cookie availability (real browser behavior)
    let cookiesAvailable = false
    const originalCookieGetter = Object.getOwnPropertyDescriptor(document, 'cookie')?.get
    
    Object.defineProperty(document, 'cookie', {
      get: () => {
        if (cookiesAvailable) {
          return 'access_token=mock-access-token; refresh_token=mock-refresh-token'
        }
        return ''
      },
      set: () => {
        // Simulate cookies becoming available after a delay
        setTimeout(() => {
          cookiesAvailable = true
        }, 100)
      },
      configurable: true,
    })

    render(
      <AuthProvider>
        <MemoryRouter initialEntries={['/auth/magic?token=valid-token']}>
          <div>
            <MagicLinkHandler />
            <AuthStateIndicator />
          </div>
        </MemoryRouter>
      </AuthProvider>
    )

    // Initially should be verifying
    expect(screen.getByText('Verifying your login...')).toBeInTheDocument()
    expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated')

    // Wait for the API call
    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/auth/verifyMagicLink?token=valid-token',
        {
          method: 'GET',
          credentials: 'include'
        }
      )
    })

    // Wait for success message
    await waitFor(() => {
      expect(screen.getByText('Login Successful!')).toBeInTheDocument()
    })

    // The auth state should eventually become authenticated after cookies are available
    await waitFor(() => {
      expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated')
    }, { timeout: 1000 })
  })

  it('should handle missing token gracefully', async () => {
    render(
      <AuthProvider>
        <MemoryRouter initialEntries={['/auth/magic']}>
          <div>
            <MagicLinkHandler />
            <AuthStateIndicator />
          </div>
        </MemoryRouter>
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.getByText('Login Failed')).toBeInTheDocument()
      expect(screen.getByText('Invalid link. No token found.')).toBeInTheDocument()
    })

    expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated')
  })

  it('should handle API errors gracefully', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      text: async () => 'Invalid token'
    })

    render(
      <AuthProvider>
        <MemoryRouter initialEntries={['/auth/magic?token=invalid-token']}>
          <div>
            <MagicLinkHandler />
            <AuthStateIndicator />
          </div>
        </MemoryRouter>
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.getByText('Login Failed')).toBeInTheDocument()
      expect(screen.getByText('Authentication failed: Invalid token')).toBeInTheDocument()
    })

    expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated')
  })
})