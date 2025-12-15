import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { act } from 'react';
import { ToastProvider } from '../../../context/ToastProvider';
import ToastContainer from '../ToastContainer';
import { useToast } from '../../../hooks/useToast';

// Component to test the useToast hook
function TestComponent() {
  const toast = useToast();

  return (
    <div>
      <button onClick={() => toast.success('Success message')}>
        Show Success
      </button>
      <button onClick={() => toast.error('Error message')}>
        Show Error
      </button>
      <button onClick={() => toast.info('Info message')}>
        Show Info
      </button>
      <button onClick={() => toast.warning('Warning message')}>
        Show Warning
      </button>
      <button onClick={() => toast.addToast('Custom message', 'success', 1000)}>
        Show Custom
      </button>
    </div>
  );
}

describe('Toast System', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.useRealTimers();
  });

  it('renders ToastContainer without crashing', () => {
    render(
      <ToastProvider>
        <ToastContainer />
      </ToastProvider>
    );
    // Container should not be visible when there are no toasts
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('displays a success toast', () => {
    render(
      <ToastProvider>
        <ToastContainer />
        <TestComponent />
      </ToastProvider>
    );

    const button = screen.getByText('Show Success');
    fireEvent.click(button);

    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText('Success message')).toBeInTheDocument();
    expect(screen.getByRole('alert')).toHaveClass('success');
  });

  it('displays an error toast', () => {
    render(
      <ToastProvider>
        <ToastContainer />
        <TestComponent />
      </ToastProvider>
    );

    const button = screen.getByText('Show Error');
    fireEvent.click(button);

    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText('Error message')).toBeInTheDocument();
    expect(screen.getByRole('alert')).toHaveClass('error');
  });

  it('displays an info toast', () => {
    render(
      <ToastProvider>
        <ToastContainer />
        <TestComponent />
      </ToastProvider>
    );

    const button = screen.getByText('Show Info');
    fireEvent.click(button);

    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText('Info message')).toBeInTheDocument();
    expect(screen.getByRole('alert')).toHaveClass('info');
  });

  it('displays a warning toast', () => {
    render(
      <ToastProvider>
        <ToastContainer />
        <TestComponent />
      </ToastProvider>
    );

    const button = screen.getByText('Show Warning');
    fireEvent.click(button);

    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText('Warning message')).toBeInTheDocument();
    expect(screen.getByRole('alert')).toHaveClass('warning');
  });

  it('displays multiple toasts simultaneously', () => {
    render(
      <ToastProvider>
        <ToastContainer />
        <TestComponent />
      </ToastProvider>
    );

    fireEvent.click(screen.getByText('Show Success'));
    fireEvent.click(screen.getByText('Show Error'));
    fireEvent.click(screen.getByText('Show Info'));

    const alerts = screen.getAllByRole('alert');
    expect(alerts).toHaveLength(3);
    expect(screen.getByText('Success message')).toBeInTheDocument();
    expect(screen.getByText('Error message')).toBeInTheDocument();
    expect(screen.getByText('Info message')).toBeInTheDocument();
  });

  it('removes toast when close button is clicked', () => {
    render(
      <ToastProvider>
        <ToastContainer />
        <TestComponent />
      </ToastProvider>
    );

    fireEvent.click(screen.getByText('Show Success'));
    
    expect(screen.getByText('Success message')).toBeInTheDocument();

    const closeButton = screen.getByRole('button', { name: /close notification/i });
    fireEvent.click(closeButton);

    // Wait for the animation to complete (300ms)
    act(() => {
      vi.advanceTimersByTime(300);
    });

    expect(screen.queryByText('Success message')).not.toBeInTheDocument();
  });

  it('auto-removes toast after duration', () => {
    render(
      <ToastProvider>
        <ToastContainer />
        <TestComponent />
      </ToastProvider>
    );

    fireEvent.click(screen.getByText('Show Custom'));
    
    expect(screen.getByText('Custom message')).toBeInTheDocument();

    // Fast-forward time to just before the toast should disappear
    act(() => {
      vi.advanceTimersByTime(700);
    });
    
    // Toast should start exit animation
    expect(screen.getByText('Custom message')).toBeInTheDocument();

    // Complete the exit animation and removal
    act(() => {
      vi.advanceTimersByTime(300);
    });

    expect(screen.queryByText('Custom message')).not.toBeInTheDocument();
  });

  it('throws error when useToast is used outside ToastProvider', () => {
    // Suppress console.error for this test
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});
    
    function InvalidComponent() {
      useToast();
      return null;
    }

    expect(() => {
      render(<InvalidComponent />);
    }).toThrow('useToast must be used within a ToastProvider');

    consoleError.mockRestore();
  });

  it('handles long messages properly', () => {
    const longMessage = 'This is a very long message that should be properly displayed in the toast notification without breaking the layout or causing any visual issues';
    
    function LongMessageComponent() {
      const toast = useToast();
      return (
        <button onClick={() => toast.info(longMessage)}>
          Show Long Message
        </button>
      );
    }

    render(
      <ToastProvider>
        <ToastContainer />
        <LongMessageComponent />
      </ToastProvider>
    );

    fireEvent.click(screen.getByText('Show Long Message'));
    
    expect(screen.getByText(longMessage)).toBeInTheDocument();
  });

  it('renders correct icons for each toast type', () => {
    render(
      <ToastProvider>
        <ToastContainer />
        <TestComponent />
      </ToastProvider>
    );

    fireEvent.click(screen.getByText('Show Success'));
    expect(screen.getByText('✓')).toBeInTheDocument();

    // Clear the success toast first
    const closeButton = screen.getByRole('button', { name: /close notification/i });
    fireEvent.click(closeButton);
    act(() => {
      vi.advanceTimersByTime(300);
    });

    fireEvent.click(screen.getByText('Show Error'));
    expect(screen.getByText('✕')).toBeInTheDocument();
  });
});
