import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { beforeEach, describe, it, expect, vi } from 'vitest';
import Layout from '../Layout';
import { TestI18nProvider } from '../../../test/i18n-test-utils';

// Mock the useAuth hook
const mockUseAuth = vi.fn();
const mockLogout = vi.fn();
const mockNavigate = vi.fn();

vi.mock('../../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth()
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

describe('Layout', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseAuth.mockReturnValue({
      logout: mockLogout
    });
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