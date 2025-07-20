import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { beforeEach, describe, it, expect, vi } from 'vitest';
import Layout from '../Layout';
import { TestI18nProvider } from '../../../test/i18n-test-utils';
import { defaultMockAuthValue, defaultMockNotificationValue } from '../../../test/mock-values';

// Mock the useAuth hook
const mockUseAuth = vi.fn();
const mockLogout = vi.fn();
const mockNavigate = vi.fn();

vi.mock('../../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth()
}));

// Mock the useNotifications hook
const mockUseNotifications = vi.fn();

vi.mock('../../../hooks/useNotifications', () => ({
  useNotifications: () => mockUseNotifications()
}));

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate
  };
});

// Mock the logo import
vi.mock('../../../assets/logo.png', () => ({
  default: 'mock-logo.png'
}));

// Mock the CookieConsent component
vi.mock('../../CookieConsent', () => ({
  default: () => <div data-testid="cookie-consent-mock">Cookie Consent</div>
}));

// Mock the RecentClubsDropdown component
vi.mock('../RecentClubsDropdown', () => ({
  default: () => <div data-testid="recent-clubs-dropdown">Recent Clubs Dropdown</div>
}));

// Mock the NotificationDropdown component
vi.mock('../NotificationDropdown', () => ({
  default: () => <div data-testid="notification-dropdown">Notification Dropdown</div>
}));

describe('Layout', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseAuth.mockReturnValue({
      ...defaultMockAuthValue,
      logout: mockLogout
    });
    mockUseNotifications.mockReturnValue(defaultMockNotificationValue);
  });

  it('renders Layout with CookieConsent component', () => {
    render(
      <TestI18nProvider>
        <BrowserRouter>
          <Layout title="Test Title">
            <div>Test Content</div>
          </Layout>
        </BrowserRouter>
      </TestI18nProvider>
    );

    expect(screen.getByText('Test Title')).toBeInTheDocument();
    expect(screen.getByText('Test Content')).toBeInTheDocument();
    expect(screen.getByTestId('cookie-consent-mock')).toBeInTheDocument();
  });

  it('renders Layout without title', () => {
    render(
      <TestI18nProvider>
        <BrowserRouter>
          <Layout>
            <div>Test Content</div>
          </Layout>
        </BrowserRouter>
      </TestI18nProvider>
    );

    expect(screen.getByText('Clubs')).toBeInTheDocument(); // Default title from Header
    expect(screen.getByText('Test Content')).toBeInTheDocument();
    expect(screen.getByTestId('cookie-consent-mock')).toBeInTheDocument();
  });
});