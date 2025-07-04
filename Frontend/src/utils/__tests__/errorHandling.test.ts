import { describe, it, expect, vi } from 'vitest';
import { getErrorMessage, logError, createErrorHandler } from '../errorHandling';

describe('Error Handling Utilities', () => {
  describe('getErrorMessage', () => {
    it('should return error message for Error instances', () => {
      const error = new Error('Test error');
      expect(getErrorMessage(error)).toBe('Test error');
    });

    it('should return string errors as-is', () => {
      const error = 'String error';
      expect(getErrorMessage(error)).toBe('String error');
    });

    it('should return fallback message for unknown errors', () => {
      const error = { someProperty: 'value' };
      expect(getErrorMessage(error)).toBe('An unknown error occurred');
    });

    it('should return custom fallback message', () => {
      const error = null;
      expect(getErrorMessage(error, 'Custom fallback')).toBe('Custom fallback');
    });
  });

  describe('logError', () => {
    it('should log error with context', () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      const error = new Error('Test error');
      
      logError('TestComponent', error);
      
      expect(consoleSpy).toHaveBeenCalledWith('Error in TestComponent:', 'Test error');
      consoleSpy.mockRestore();
    });
  });

  describe('createErrorHandler', () => {
    it('should create error handler that logs and sets error', () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      const setError = vi.fn();
      const error = new Error('Test error');
      
      const handler = createErrorHandler('TestComponent', setError);
      handler(error);
      
      expect(consoleSpy).toHaveBeenCalledWith('Error in TestComponent:', 'Test error');
      expect(setError).toHaveBeenCalledWith('Test error');
      consoleSpy.mockRestore();
    });

    it('should use custom fallback message', () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      const setError = vi.fn();
      const error = null;
      
      const handler = createErrorHandler('TestComponent', setError, 'Custom fallback');
      handler(error);
      
      expect(setError).toHaveBeenCalledWith('Custom fallback');
      consoleSpy.mockRestore();
    });
  });
});