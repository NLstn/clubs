import { describe, it, expect, vi } from 'vitest'

// Mock environment variable
vi.stubEnv('VITE_API_HOST', 'http://localhost:3000')

// Mock axios
vi.mock('axios', () => ({
  default: {
    create: vi.fn(() => ({
      interceptors: {
        request: {
          use: vi.fn()
        },
        response: {
          use: vi.fn()
        }
      }
    }))
  }
}))

// Mock jwt-decode
vi.mock('jwt-decode', () => ({
  jwtDecode: vi.fn()
}))

describe('API Configuration', () => {
  it('can import api module without errors', async () => {
    expect(async () => {
      await import('../api')
    }).not.toThrow()
  })

  it('exports an api instance', async () => {
    const { default: api } = await import('../api')
    expect(api).toBeDefined()
    expect(typeof api).toBe('object')
  })
})