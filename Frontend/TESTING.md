# Frontend Testing Setup

This project now includes a complete testing setup using Vitest and React Testing Library.

## Testing Stack

- **Vitest**: Fast unit test framework designed for Vite projects
- **React Testing Library**: Testing utilities for React components
- **jsdom**: DOM environment for Node.js testing
- **@testing-library/jest-dom**: Custom Jest matchers for DOM assertions

## Running Tests

```bash
# Run tests once
npm run test:run

# Run tests in watch mode (development)
npm test

# Run tests with UI (if @vitest/ui is installed)
npm run test:ui
```

## Test Structure

Tests are located in `__tests__` directories next to the code they test:

```
src/
  components/
    common/
      Button.tsx
      __tests__/
        Button.test.tsx
  utils/
    helpers.ts
    __tests__/
      helpers.test.ts
```

## Example Tests

### Component Testing
See `src/components/common/__tests__/Button.test.tsx` for an example of testing React components with user interactions.

### Utility Function Testing
See `src/utils/__tests__/helpers.test.ts` for an example of testing utility functions.

## Configuration

- **Vite config**: `vite.config.ts` includes test configuration
- **Test setup**: `src/test/setup.ts` configures testing environment
- **Global test utilities**: Available via vitest globals configuration