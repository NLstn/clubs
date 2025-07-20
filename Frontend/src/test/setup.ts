import '@testing-library/jest-dom'
import { setupDefaultApiMocks, cleanupApiMocks } from './api-mocks'
import { beforeEach, afterEach } from 'vitest'

// Mock environment variables for testing
import.meta.env.VITE_API_HOST = 'http://localhost:8080'

// Setup API mocks before each test
beforeEach(() => {
  setupDefaultApiMocks()
})

// Clean up after each test
afterEach(() => {
  cleanupApiMocks()
})