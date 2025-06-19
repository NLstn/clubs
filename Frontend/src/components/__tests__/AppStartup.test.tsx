import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AuthProvider } from '../../context/AuthProvider'
import ProtectedRoute from '../auth/ProtectedRoute'

// Mock the api module
vi.mock('../../utils/api', () => ({
  default: {}
}))

// Mock a simple protected component
const MockProtectedComponent = () => <div data-testid="protected-content">Protected Content</div>

describe('App Startup with Existing Cookies', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('should not redirect to login when valid cookies are present on app startup', async () => {
    // Mock cookies being present on startup
    Object.defineProperty(document, 'cookie', {
      get: () => 'access_token=valid-token; refresh_token=valid-refresh-token',
      set: () => {},
      configurable: true,
    })

    const { container } = render(
      <AuthProvider>
        <MemoryRouter initialEntries={['/']}>
          <ProtectedRoute>
            <MockProtectedComponent />
          </ProtectedRoute>
        </MemoryRouter>
      </AuthProvider>
    )

    // Should render the protected content, not redirect to login
    await waitFor(() => {
      expect(screen.getByTestId('protected-content')).toBeInTheDocument()
    })

    // Should not contain any navigation to login
    expect(container).not.toHaveTextContent('Login')
  })

  it('should redirect to login when no cookies are present on app startup', async () => {
    // Mock no cookies
    Object.defineProperty(document, 'cookie', {
      get: () => '',
      set: () => {},
      configurable: true,
    })

    render(
      <AuthProvider>
        <MemoryRouter initialEntries={['/']}>
          <ProtectedRoute>
            <MockProtectedComponent />
          </ProtectedRoute>
        </MemoryRouter>
      </AuthProvider>
    )

    // Should redirect to login (MemoryRouter will show the Navigate component behavior)
    // In a real app, this would cause a redirect to /login, but in our test environment
    // we can't easily test the redirect, so we just verify the protected content is not shown
    expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument()
  })

  it('should handle cookie availability timing issues during app startup', async () => {
    // Simulate cookies that become available after a short delay (browser timing)
    let cookiesAvailable = false
    
    Object.defineProperty(document, 'cookie', {
      get: () => {
        if (cookiesAvailable) {
          return 'access_token=delayed-token; refresh_token=delayed-refresh-token'
        }
        return ''
      },
      set: () => {},
      configurable: true,
    })

    // Make cookies available after a short delay
    setTimeout(() => {
      cookiesAvailable = true
    }, 100)

    render(
      <AuthProvider>
        <MemoryRouter initialEntries={['/']}>
          <ProtectedRoute>
            <MockProtectedComponent />
          </ProtectedRoute>
        </MemoryRouter>
      </AuthProvider>
    )

    // Initially should not show protected content
    expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument()

    // After cookies become available, should show protected content
    await waitFor(() => {
      expect(screen.getByTestId('protected-content')).toBeInTheDocument()
    }, { timeout: 1000 })
  })
})