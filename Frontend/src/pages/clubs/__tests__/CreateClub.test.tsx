import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom';
import CreateClub from '../CreateClub';
import { TestI18nProvider } from '../../../test/i18n-test-utils';

// Mock useNavigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate,
}));

// Mock useAuth
const mockPost = vi.fn();
vi.mock('../../../hooks/useAuth', () => ({
  useAuth: () => ({
    api: {
      post: mockPost,
    },
  }),
}));

// Mock Layout component
vi.mock('../../../components/layout/Layout', () => ({
  default: ({ children, title }: { children: React.ReactNode; title: string }) => (
    <div data-testid="layout">
      <h1>{title}</h1>
      {children}
    </div>
  ),
}));

const renderWithProviders = (component: React.ReactElement) => {
  return render(
    <TestI18nProvider>
      {component}
    </TestI18nProvider>
  );
};

describe('CreateClub', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders correctly', () => {
    renderWithProviders(<CreateClub />);
    
    expect(screen.getByTestId('layout')).toBeInTheDocument();
    expect(screen.getAllByText('Create New Club')[0]).toBeInTheDocument(); // Get first occurrence
    expect(screen.getByLabelText(/club name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /create club/i })).toBeInTheDocument();
  });

  it('updates input values when user types', () => {
    renderWithProviders(<CreateClub />);
    
    const clubNameInput = screen.getByLabelText(/club name/i) as HTMLInputElement;
    const descriptionInput = screen.getByLabelText(/description/i) as HTMLTextAreaElement;
    
    fireEvent.change(clubNameInput, { target: { value: 'Test Club' } });
    fireEvent.change(descriptionInput, { target: { value: 'Test Description' } });
    
    expect(clubNameInput.value).toBe('Test Club');
    expect(descriptionInput.value).toBe('Test Description');
  });

  it('successfully creates a club and redirects', async () => {
    const mockResponse = { data: { ID: '123', Name: 'Test Club' } };
    mockPost.mockResolvedValue(mockResponse);
    
    renderWithProviders(<CreateClub />);
    
    const clubNameInput = screen.getByLabelText(/club name/i);
    const descriptionInput = screen.getByLabelText(/description/i);
    const submitButton = screen.getByRole('button', { name: /create club/i });
    
    fireEvent.change(clubNameInput, { target: { value: 'Test Club' } });
    fireEvent.change(descriptionInput, { target: { value: 'Test Description' } });
    fireEvent.click(submitButton);
    
    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith('/api/v2/Clubs', { 
        Name: 'Test Club', 
        Description: 'Test Description' 
      });
    });
    
    const successMessages = screen.getAllByText('Club created successfully!');
    expect(successMessages.length).toBeGreaterThan(0);
    const messageDiv = successMessages.find((el) => el.classList.contains('message'));
    expect(messageDiv).toBeInTheDocument();
    
    // Wait for navigation timeout
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/clubs/123');
    }, { timeout: 1500 });
  });
});
