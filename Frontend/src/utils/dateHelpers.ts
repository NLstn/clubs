/**
 * Utility functions for consistent date formatting across the application
 */

/**
 * Formats a date string to a localized string
 * @param dateString - ISO date string
 * @param options - Intl.DateTimeFormatOptions
 * @returns Formatted date string
 */
export const formatDate = (
  dateString: string,
  options: Intl.DateTimeFormatOptions = {}
): string => {
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      return 'Invalid date';
    }
    return date.toLocaleDateString(undefined, options);
  } catch {
    return 'Invalid date';
  }
};

/**
 * Formats a date string to a localized date and time string
 * @param dateString - ISO date string
 * @param options - Intl.DateTimeFormatOptions
 * @returns Formatted date and time string
 */
export const formatDateTime = (
  dateString: string,
  options: Intl.DateTimeFormatOptions = {}
): string => {
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      return 'Invalid date';
    }
    return date.toLocaleString(undefined, options);
  } catch {
    return 'Invalid date';
  }
};

/**
 * Formats a date string for HTML datetime-local input
 * @param dateString - ISO date string
 * @returns Formatted date string for datetime-local input
 */
export const formatDateTimeLocal = (dateString: string): string => {
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      return '';
    }
    return date.toISOString().slice(0, 16);
  } catch {
    return '';
  }
};

/**
 * Validates if a date string represents a valid date
 * @param dateString - Date string to validate
 * @returns True if valid, false otherwise
 */
export const isValidDate = (dateString: string): boolean => {
  try {
    const date = new Date(dateString);
    return !isNaN(date.getTime());
  } catch {
    return false;
  }
};