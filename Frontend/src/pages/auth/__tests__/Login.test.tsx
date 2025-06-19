import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import Login from '../Login';

// Mock environment variable
vi.stubEnv('VITE_API_HOST', 'http://localhost:3000');

// Mock fetch
global.fetch = vi.fn();

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  );
};

describe('Login Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders with proper styling classes', () => {
    renderWithRouter(<Login />);
    
    // Check if form elements are present and properly structured
    expect(screen.getByRole('heading', { name: /login/i })).toBeInTheDocument();
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /send magic link/i })).toBeInTheDocument();
  });
  
  it('displays instruction text', () => {
    renderWithRouter(<Login />);
    expect(screen.getByText(/enter your email to receive a magic link/i)).toBeInTheDocument();
  });
  
  it('has the correct container structure for styling', () => {
    const { container } = renderWithRouter(<Login />);
    
    // Check if the login container div exists
    const loginContainer = container.querySelector('.login-container');
    expect(loginContainer).toBeInTheDocument();
    
    // Check if the login box div exists
    const loginBox = container.querySelector('.login-box');
    expect(loginBox).toBeInTheDocument();
  });

  it('displays success message with correct styling when request succeeds', async () => {
    // Mock successful response
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: true,
    } as Response);

    renderWithRouter(<Login />);
    
    const emailInput = screen.getByLabelText(/email/i);
    const submitButton = screen.getByRole('button', { name: /send magic link/i });
    
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.click(submitButton);
    
    await waitFor(() => {
      const message = screen.getByText(/check your email for a login link/i);
      expect(message).toBeInTheDocument();
      expect(message).toHaveClass('message', 'success');
    });
  });

  it('displays error message with correct styling when request fails', async () => {
    // Mock failed response
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: false,
      text: () => Promise.resolve('Invalid email'),
    } as Response);

    renderWithRouter(<Login />);
    
    const emailInput = screen.getByLabelText(/email/i);
    const submitButton = screen.getByRole('button', { name: /send magic link/i });
    
    fireEvent.change(emailInput, { target: { value: 'invalid@example.com' } });
    fireEvent.click(submitButton);
    
    await waitFor(() => {
      const message = screen.getByText(/error: invalid email/i);
      expect(message).toBeInTheDocument();
      expect(message).toHaveClass('message', 'error');
    });
  });
});