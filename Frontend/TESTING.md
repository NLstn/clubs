# Frontend Testing Setup

This project includes a comprehensive testing setup using Vitest and React Testing Library.

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
    __tests__/
      TypeAheadDropdown.test.tsx
      Header.test.tsx
      ProtectedRoute.test.tsx
  context/
    __tests__/
      AuthContext.test.tsx
  utils/
    __tests__/
      api.test.ts
```

## Example Tests

### Component Testing

#### TypeAheadDropdown Component
Tests user interactions, dropdown behavior, option selection, and search functionality.

#### Header Component  
Tests navigation, dropdown menus, user interactions, and authentication actions.

#### ProtectedRoute Component
Tests authentication-based route protection and redirects.

### Context Testing

#### AuthContext
Tests authentication state management, login/logout functionality, localStorage integration, and API calls.

### Utility Testing

#### API Module
Tests axios configuration, interceptor setup, and module exports.

## GitHub Actions Integration

Frontend tests run automatically on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- When Frontend code changes

The workflow includes:
- Dependency installation
- ESLint linting
- Build verification
- Test execution

## Configuration

- **Vite config**: `vite.config.ts` includes test configuration
- **Test setup**: `src/test/setup.ts` configures testing environment
- **Global test utilities**: Available via vitest globals configuration