import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom';
import CreateClub from '../CreateClub';

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

describe('CreateClub', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders correctly', () => {
    render(<CreateClub />);
    
    expect(screen.getByTestId('layout')).toBeInTheDocument();
    expect(screen.getAllByText('Create New Club')[0]).toBeInTheDocument(); // Get first occurrence
    expect(screen.getByLabelText(/club name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /create club/i })).toBeInTheDocument();
  });

  it('updates input values when user types', () => {
    render(<CreateClub />);
    
    const clubNameInput = screen.getByLabelText(/club name/i) as HTMLInputElement;
    const descriptionInput = screen.getByLabelText(/description/i) as HTMLTextAreaElement;
    
    fireEvent.change(clubNameInput, { target: { value: 'Test Club' } });
    fireEvent.change(descriptionInput, { target: { value: 'Test Description' } });
    
    expect(clubNameInput.value).toBe('Test Club');
    expect(descriptionInput.value).toBe('Test Description');
  });

  it('successfully creates a club and redirects', async () => {
    const mockResponse = { data: { id: '123', name: 'Test Club' } };
    mockPost.mockResolvedValue(mockResponse);
    
    render(<CreateClub />);
    
    const clubNameInput = screen.getByLabelText(/club name/i);
    const descriptionInput = screen.getByLabelText(/description/i);
    const submitButton = screen.getByRole('button', { name: /create club/i });
    
    fireEvent.change(clubNameInput, { target: { value: 'Test Club' } });
    fireEvent.change(descriptionInput, { target: { value: 'Test Description' } });
    fireEvent.click(submitButton);
    
    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith('/api/v1/clubs', { 
        name: 'Test Club', 
        description: 'Test Description' 
      });
    });
    
    expect(screen.getByText('Club created successfully!')).toBeInTheDocument();
    
    // Wait for navigation timeout
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/clubs/123');
    }, { timeout: 1500 });
  });
});
