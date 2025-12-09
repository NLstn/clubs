import { describe, it, expect } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useTheme } from '../useTheme';

describe('useTheme', () => {
  it('should throw an error when used outside ThemeProvider', () => {
    // Suppress console.error for this test since we expect an error
    const originalError = console.error;
    console.error = () => {};

    expect(() => {
      renderHook(() => useTheme());
    }).toThrow('useTheme must be used within a ThemeProvider');

    // Restore console.error
    console.error = originalError;
  });
});
