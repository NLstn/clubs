import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import type { AxiosRequestConfig } from 'axios'

vi.stubEnv('VITE_API_HOST', 'http://localhost:3000')

let requestInterceptor: (config: AxiosRequestConfig) => Promise<AxiosRequestConfig>
const axiosPostMock = vi.fn()

// Mock axios module
vi.mock('axios', () => {
  const instance = {
    interceptors: {
      request: {
        use: vi.fn((handler) => {
          requestInterceptor = handler
        })
      },
      response: { use: vi.fn() }
    }
  }

  return {
    default: {
      create: vi.fn(() => instance),
      post: axiosPostMock
    }
  }
})

// Mock jwt-decode
const jwtDecodeMock = vi.fn()
vi.mock('jwt-decode', () => ({
  jwtDecode: jwtDecodeMock
}))

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}
const mockWindow = { location: { href: '' } } as Window & typeof globalThis
let locationHrefSpy = vi.fn()

describe('API request interceptor', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.resetModules()
    vi.stubGlobal('window', mockWindow)
    vi.stubGlobal('localStorage', localStorageMock)
    locationHrefSpy = vi.fn()
    Object.defineProperty(window, 'localStorage', { value: localStorageMock })
    Object.defineProperty(window.location, 'href', {
      get: () => mockWindow.location.href,
      set: locationHrefSpy,
      configurable: true
    })
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('refreshes expired token and sets authorization header', async () => {
    localStorageMock.getItem.mockImplementation((key: string) => {
      if (key === 'auth_token') return 'expired-token'
      if (key === 'refresh_token') return 'refresh-token'
      return null
    })

    jwtDecodeMock.mockReturnValue({ exp: Math.floor(Date.now() / 1000) - 60 })

    axiosPostMock.mockResolvedValueOnce({
      data: { access: 'new-access', refresh: 'new-refresh' }
    })

    await import('../api')

    const cfg = { headers: {} as Record<string, string> }
    const result = await requestInterceptor(cfg)

    expect(axiosPostMock).toHaveBeenCalledWith(
      'http://localhost:3000/api/v1/auth/refreshToken',
      {},
      { headers: { Authorization: 'refresh-token' } }
    )
    expect(localStorageMock.setItem).toHaveBeenCalledWith('auth_token', 'new-access')
    expect(localStorageMock.setItem).toHaveBeenCalledWith('refresh_token', 'new-refresh')
    expect(result.headers?.Authorization).toBe('Bearer new-access')
  })

  it('clears storage and redirects to login when no refresh token', async () => {
    localStorageMock.getItem.mockImplementation((key: string) => {
      if (key === 'auth_token') return 'expired-token'
      if (key === 'refresh_token') return null
      return null
    })

    jwtDecodeMock.mockReturnValue({ exp: Math.floor(Date.now() / 1000) - 60 })

    await import('../api')

    const cfg = { headers: {} as Record<string, string> }
    await expect(requestInterceptor(cfg)).rejects.toThrow('No refresh token available')

    expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_token')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')
    expect(locationHrefSpy).toHaveBeenCalledWith('/login')
  })
})
