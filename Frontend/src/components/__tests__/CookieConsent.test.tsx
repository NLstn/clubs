import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { beforeEach, afterEach, describe, it, expect, vi } from 'vitest';
import CookieConsent from '../CookieConsent';

describe('CookieConsent', () => {
  // Mock localStorage
  const localStorageMock = {
    getItem: vi.fn(),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
  };

  beforeEach(() => {
    Object.defineProperty(window, 'localStorage', { value: localStorageMock });
    // Mock window.alert
    vi.stubGlobal('alert', vi.fn());
    localStorageMock.getItem.mockClear();
    localStorageMock.setItem.mockClear();
  });

  afterEach(() => {
    localStorageMock.getItem.mockClear();
    localStorageMock.setItem.mockClear();
    vi.restoreAllMocks();
  });

  it('renders the banner when no consent is stored', () => {
    localStorageMock.getItem.mockReturnValue(null);
    
    render(<CookieConsent />);
    
    expect(screen.getByTestId('cookie-consent-banner')).toBeInTheDocument();
    expect(screen.getByText(/We use cookies to improve your experience/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Accept' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Learn More' })).toBeInTheDocument();
  });

  it('does not render the banner when consent is already given', () => {
    localStorageMock.getItem.mockReturnValue('accepted');
    
    render(<CookieConsent />);
    
    expect(screen.queryByTestId('cookie-consent-banner')).not.toBeInTheDocument();
  });

  it('hides the banner and stores consent when Accept is clicked', async () => {
    localStorageMock.getItem.mockReturnValue(null);
    
    render(<CookieConsent />);
    
    const acceptButton = screen.getByRole('button', { name: 'Accept' });
    fireEvent.click(acceptButton);

    await waitFor(() => {
      expect(localStorageMock.setItem).toHaveBeenCalledWith('cookie-consent', 'accepted');
      expect(screen.queryByTestId('cookie-consent-banner')).not.toBeInTheDocument();
    });
  });

  it('shows alert when Learn More is clicked', () => {
    localStorageMock.getItem.mockReturnValue(null);
    
    render(<CookieConsent />);
    
    const learnMoreButton = screen.getByRole('button', { name: 'Learn More' });
    fireEvent.click(learnMoreButton);

    expect(window.alert).toHaveBeenCalledWith(
      expect.stringContaining('This website uses cookies to enhance your browsing experience')
    );
  });

  it('checks localStorage on mount', () => {
    localStorageMock.getItem.mockReturnValue(null);
    
    render(<CookieConsent />);
    
    expect(localStorageMock.getItem).toHaveBeenCalledWith('cookie-consent');
  });

  it('renders with correct text content', () => {
    localStorageMock.getItem.mockReturnValue(null);
    
    render(<CookieConsent />);
    
    expect(screen.getByText('We use cookies to improve your experience on our site. By using our site, you accept our use of cookies.')).toBeInTheDocument();
  });

  it('banner is initially hidden and shows only when no consent exists', async () => {
    localStorageMock.getItem.mockReturnValue(null);
    
    const { rerender } = render(<CookieConsent />);
    
    // Should be visible when no consent
    expect(screen.getByTestId('cookie-consent-banner')).toBeInTheDocument();
    
    // Simulate consent being given
    localStorageMock.getItem.mockReturnValue('accepted');
    rerender(<CookieConsent />);
    
    // Should still be visible since it was already rendered and consent check happens only on mount
    expect(screen.getByTestId('cookie-consent-banner')).toBeInTheDocument();
  });
});