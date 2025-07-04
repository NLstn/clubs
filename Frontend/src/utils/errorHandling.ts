/**
 * Utility functions for consistent error handling across the application
 */

/**
 * Formats an error message consistently
 * @param error - The error object or unknown value
 * @param fallbackMessage - Default message if error is not an Error instance
 * @returns Formatted error message
 */
export const getErrorMessage = (error: unknown, fallbackMessage = 'An unknown error occurred'): string => {
  if (error instanceof Error) {
    return error.message;
  }
  if (typeof error === 'string') {
    return error;
  }
  return fallbackMessage;
};

/**
 * Logs an error consistently with context
 * @param context - Context where the error occurred
 * @param error - The error to log
 */
export const logError = (context: string, error: unknown): void => {
  const errorMessage = getErrorMessage(error);
  console.error(`Error in ${context}:`, errorMessage);
};

/**
 * Creates a consistent error handler for async operations
 * @param context - Context where the error occurred
 * @param setError - State setter for error message
 * @param fallbackMessage - Default error message
 * @returns Error handler function
 */
export const createErrorHandler = (
  context: string,
  setError: (error: string | null) => void,
  fallbackMessage?: string
) => {
  return (error: unknown) => {
    const errorMessage = getErrorMessage(error, fallbackMessage);
    logError(context, error);
    setError(errorMessage);
  };
};