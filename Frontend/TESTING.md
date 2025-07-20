# Frontend Testing Setup

This project includes a comprehensive testing setup using Vitest and React Testing Library with automatic API mocking.

## Testing Stack

- **Vitest**: Fast unit test framework designed for Vite projects
- **React Testing Library**: Testing utilities for React components
- **jsdom**: DOM environment for Node.js testing
- **@testing-library/jest-dom**: Custom Jest matchers for DOM assertions
- **axios-mock-adapter**: HTTP request mocking for API calls

## Running Tests

```bash
# Run tests once
npm test

# Run tests in watch mode (development)
npm run test:watch

# Run tests with coverage
npm run test:coverage
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
  test/
    api-mocks.ts
    mock-values.ts
    test-utils.tsx
    i18n-test-utils.tsx
    setup.ts
```

## API Mocking

All HTTP requests are automatically mocked during tests to prevent real API calls:

### Automatic Setup
- API mocks are configured in `src/test/api-mocks.ts`
- Mocks are automatically set up before each test in `src/test/setup.ts`
- Both the main `api` instance and direct `axios` calls are mocked

### Mock Data Available
- **Notifications**: Mock notifications and counts
- **Clubs**: Mock club data with members, teams, fines
- **Users**: Mock user profiles and authentication
- **Dashboard**: Mock activity data
- **Authentication**: Mock login/logout and token refresh

### Customizing Mocks in Tests

```typescript
import { mockApi, setupErrorApiMocks } from '../../test/api-mocks';

// Set up specific response for a test
beforeEach(() => {
  mockApi.onGet('/api/v1/custom-endpoint').reply(200, { data: 'test' });
});

// Test error conditions
it('handles API errors', () => {
  setupErrorApiMocks(); // All API calls will fail
  // ... test error handling
});
```

### Mock Hook Values

```typescript
import { defaultMockNotificationValue, mockNotificationsWithData } from '../../test/mock-values';

// Mock useNotifications hook
vi.mock('../../hooks/useNotifications', () => ({
  useNotifications: () => mockNotificationsWithData
}));
```

## Example Tests

### Component Testing

#### TypeAheadDropdown Component
Tests user interactions, dropdown behavior, option selection, and search functionality.

#### Header Component  
Tests navigation, dropdown menus, user interactions, and authentication actions. Automatically mocks notification API calls.

#### ProtectedRoute Component
Tests authentication-based route protection and redirects.

### Context Testing

#### AuthContext
Tests authentication state management, login/logout functionality, localStorage integration, and API calls.

### Utility Testing

#### API Module
Tests axios configuration, interceptor setup, and module exports.

## Common Patterns

### Mocking Hooks
```typescript
const mockUseAuth = vi.fn();
const mockUseNotifications = vi.fn();

vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth()
}));

vi.mock('../../hooks/useNotifications', () => ({
  useNotifications: () => mockUseNotifications()
}));

beforeEach(() => {
  mockUseAuth.mockReturnValue(defaultMockAuthValue);
  mockUseNotifications.mockReturnValue(defaultMockNotificationValue);
});
```

### Testing Components with Providers
```typescript
import { TestProviders } from '../../test/test-utils';

const renderWithProviders = (component) => {
  return render(
    <TestProviders>
      {component}
    </TestProviders>
  );
};
```

## GitHub Actions Integration

Frontend tests run automatically on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- When Frontend code changes

The workflow includes:
- Dependency installation
- ESLint linting
- Build verification  
- Test execution with API mocking

## Configuration

- **Vite config**: `vite.config.ts` includes test configuration
- **Test setup**: `src/test/setup.ts` configures testing environment and API mocks
- **Global test utilities**: Available via vitest globals configuration
- **TypeScript types**: jest-dom types are included in `src/vite-env.d.ts`

## Troubleshooting

### HTTP Connection Errors
If you see "ECONNREFUSED" errors, it means API calls aren't being mocked properly:

1. Ensure your component tests mock all required hooks (`useAuth`, `useNotifications`, etc.)
2. Check that `src/test/setup.ts` is being loaded in your test configuration
3. Verify that new API endpoints are added to `src/test/api-mocks.ts`

### Mock Not Working
- Ensure mocks are set up in `beforeEach()` blocks
- Check that mock functions are cleared with `vi.clearAllMocks()`
- Verify mock adapters are reset between tests