import { describe, it, expect } from 'vitest';
import { formatDate, formatDateTime, formatDateTimeLocal, isValidDate } from '../dateHelpers';

describe('Date Helper Utilities', () => {
  describe('formatDate', () => {
    it('should format valid date string', () => {
      const date = '2023-12-25T10:30:00Z';
      const result = formatDate(date);
      expect(result).toMatch(/12\/25\/2023/); // Basic check for date formatting
    });

    it('should handle invalid date strings', () => {
      const result = formatDate('invalid-date');
      expect(result).toBe('Invalid date');
    });
  });

  describe('formatDateTime', () => {
    it('should format valid datetime string', () => {
      const date = '2023-12-25T10:30:00Z';
      const result = formatDateTime(date);
      expect(result).toContain('2023');
      expect(result).toContain('25');
    });

    it('should handle invalid datetime strings', () => {
      const result = formatDateTime('invalid-date');
      expect(result).toBe('Invalid date');
    });
  });

  describe('formatDateTimeLocal', () => {
    it('should format date for datetime-local input', () => {
      const date = '2023-12-25T10:30:00Z';
      const result = formatDateTimeLocal(date);
      expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/);
    });

    it('should return empty string for invalid dates', () => {
      const result = formatDateTimeLocal('invalid-date');
      expect(result).toBe('');
    });
  });

  describe('isValidDate', () => {
    it('should return true for valid dates', () => {
      expect(isValidDate('2023-12-25T10:30:00Z')).toBe(true);
      expect(isValidDate('2023-12-25')).toBe(true);
    });

    it('should return false for invalid dates', () => {
      expect(isValidDate('invalid-date')).toBe(false);
      expect(isValidDate('')).toBe(false);
    });
  });
});