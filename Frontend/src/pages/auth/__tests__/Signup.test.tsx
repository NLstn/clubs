import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import Signup from '../Signup';
import { AuthContext } from '../../../context/AuthContext';

const mockNavigate = vi.fn();

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

vi.mock('../../../hooks/useCurrentUser', () => ({
  useCurrentUser: () => ({
    user: { ID: 'user-123', Email: 'test@example.com', FirstName: 'Test', LastName: 'User' },
    loading: false,
    error: null,
  }),
}));

const mockApi = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  patch: vi.fn(),
  delete: vi.fn(),
};

const mockAuthContext = {
  isAuthenticated: true,
  accessToken: 'test-token',
  refreshToken: 'test-refresh-token',
  login: vi.fn(),
  logout: vi.fn(),
  api: mockApi,
};

const SignupWrapper = ({ children }: { children: React.ReactNode }) => (
  <BrowserRouter>
    <AuthContext.Provider value={mockAuthContext}>
      {children}
    </AuthContext.Provider>
  </BrowserRouter>
);

describe('Signup', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders signup form with first name and last name fields', () => {
    render(
      <SignupWrapper>
        <Signup />
      </SignupWrapper>
    );

    expect(screen.getByText('Complete Your Profile')).toBeInTheDocument();
    expect(screen.getByLabelText('First Name *')).toBeInTheDocument();
    expect(screen.getByLabelText('Last Name *')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Complete Profile' })).toBeInTheDocument();
  });

  it('disables submit button when fields are empty', () => {
    render(
      <SignupWrapper>
        <Signup />
      </SignupWrapper>
    );

    const submitButton = screen.getByRole('button', { name: 'Complete Profile' });
    expect(submitButton).toBeDisabled();
  });

  it('enables submit button when both fields are filled', () => {
    render(
      <SignupWrapper>
        <Signup />
      </SignupWrapper>
    );

    const firstNameInput = screen.getByLabelText('First Name *');
    const lastNameInput = screen.getByLabelText('Last Name *');
    const submitButton = screen.getByRole('button', { name: 'Complete Profile' });

    fireEvent.change(firstNameInput, { target: { value: 'John' } });
    fireEvent.change(lastNameInput, { target: { value: 'Doe' } });

    expect(submitButton).not.toBeDisabled();
  });

  it('submits form with correct data and navigates to dashboard', async () => {
    mockApi.patch = vi.fn().mockResolvedValue({ data: {} });

    render(
      <SignupWrapper>
        <Signup />
      </SignupWrapper>
    );

    const firstNameInput = screen.getByLabelText('First Name *');
    const lastNameInput = screen.getByLabelText('Last Name *');
    const submitButton = screen.getByRole('button', { name: 'Complete Profile' });

    fireEvent.change(firstNameInput, { target: { value: 'John' } });
    fireEvent.change(lastNameInput, { target: { value: 'Doe' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockApi.patch).toHaveBeenCalledWith('/api/v2/Users(\'user-123\')', {
        FirstName: 'John',
        LastName: 'Doe',
      });
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });

  it('shows error message when API call fails', async () => {
    mockApi.patch = vi.fn().mockRejectedValue(new Error('API Error'));

    render(
      <SignupWrapper>
        <Signup />
      </SignupWrapper>
    );

    const firstNameInput = screen.getByLabelText('First Name *');
    const lastNameInput = screen.getByLabelText('Last Name *');
    const submitButton = screen.getByRole('button', { name: 'Complete Profile' });

    fireEvent.change(firstNameInput, { target: { value: 'John' } });
    fireEvent.change(lastNameInput, { target: { value: 'Doe' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText('Failed to update profile. Please try again.')).toBeInTheDocument();
    });
  });
});