import { describe, it, expect } from 'vitest'
import { formatCurrency, isValidEmail } from '../helpers'

describe('Utility Functions', () => {
  describe('formatCurrency', () => {
    it('formats positive numbers correctly', () => {
      expect(formatCurrency(25.99)).toBe('$25.99')
      expect(formatCurrency(100)).toBe('$100.00')
    })

    it('formats zero correctly', () => {
      expect(formatCurrency(0)).toBe('$0.00')
    })

    it('formats negative numbers correctly', () => {
      expect(formatCurrency(-10.50)).toBe('-$10.50')
    })
  })

  describe('isValidEmail', () => {
    it('validates correct email addresses', () => {
      expect(isValidEmail('test@example.com')).toBe(true)
      expect(isValidEmail('user.name@domain.co.uk')).toBe(true)
      expect(isValidEmail('user+tag@example.org')).toBe(true)
    })

    it('rejects invalid email addresses', () => {
      expect(isValidEmail('invalid-email')).toBe(false)
      expect(isValidEmail('test@')).toBe(false)
      expect(isValidEmail('@example.com')).toBe(false)
      expect(isValidEmail('')).toBe(false)
    })
  })
})