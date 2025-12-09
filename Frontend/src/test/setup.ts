import '@testing-library/jest-dom'
import { setupDefaultApiMocks, cleanupApiMocks } from './api-mocks'
import { beforeEach, afterEach } from 'vitest'

// Mock environment variables for testing
import.meta.env.VITE_API_HOST = 'http://localhost:8080'

// Mock window.matchMedia for ThemeProvider
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {}, // deprecated
    removeListener: () => {}, // deprecated
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => true,
  }),
})

// Configure React Testing Library to not show act warnings in console
// This is handled by our test utilities which properly wrap updates in act()
import { configure } from '@testing-library/react'

configure({
  // Reduce the timeout for async operations in tests
  asyncUtilTimeout: 1000,
})

// Suppress React act warnings in tests since we handle them properly
// This is especially important for React 19 where act() behavior has changed
const originalError = console.error
beforeEach(() => {
  console.error = (...args: unknown[]) => {
    if (
      typeof args[0] === 'string' &&
      (args[0].includes('act(...)') || 
       args[0].includes('Warning: An update to') ||
       args[0].includes('When testing, code that causes React state updates'))
    ) {
      return
    }
    originalError.call(console, ...args)
  }
  
  setupDefaultApiMocks()
})

// Clean up after each test
afterEach(() => {
  console.error = originalError
  cleanupApiMocks()
})