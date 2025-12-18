import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Button } from '../Button';

describe('Button Component', () => {
  it('renders button with children', () => {
    render(<Button>Click Me</Button>);
    expect(screen.getByRole('button', { name: /click me/i })).toBeInTheDocument();
  });

  it('renders with primary variant by default', () => {
    render(<Button>Test</Button>);
    const button = screen.getByRole('button');
    expect(button).toHaveClass('ui-button--primary');
  });

  it('renders with different variants', () => {
    const { rerender } = render(<Button variant="secondary">Test</Button>);
    expect(screen.getByRole('button')).toHaveClass('ui-button--secondary');

    rerender(<Button variant="accept">Test</Button>);
    expect(screen.getByRole('button')).toHaveClass('ui-button--accept');

    rerender(<Button variant="cancel">Test</Button>);
    expect(screen.getByRole('button')).toHaveClass('ui-button--cancel');

    rerender(<Button variant="maybe">Test</Button>);
    expect(screen.getByRole('button')).toHaveClass('ui-button--maybe');
  });

  it('renders with different sizes', () => {
    const { rerender } = render(<Button size="sm">Test</Button>);
    expect(screen.getByRole('button')).toHaveClass('ui-button--sm');

    rerender(<Button size="md">Test</Button>);
    expect(screen.getByRole('button')).toHaveClass('ui-button--md');

    rerender(<Button size="lg">Test</Button>);
    expect(screen.getByRole('button')).toHaveClass('ui-button--lg');
  });

  it('renders with counter badge when counter is provided and greater than 0', () => {
    render(<Button counter={5}>View Requests</Button>);
    const button = screen.getByRole('button');
    expect(button).toHaveClass('ui-button--with-counter');
    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('does not render counter badge when counter is 0', () => {
    render(<Button counter={0}>View Requests</Button>);
    const button = screen.getByRole('button');
    expect(button).not.toHaveClass('ui-button--with-counter');
    expect(screen.queryByText('0')).not.toBeInTheDocument();
  });

  it('does not render counter badge when counter is undefined', () => {
    render(<Button>View Requests</Button>);
    const button = screen.getByRole('button');
    expect(button).not.toHaveClass('ui-button--with-counter');
  });

  it('displays "99+" when counter is greater than 99', () => {
    render(<Button counter={150}>View Requests</Button>);
    expect(screen.getByText('99+')).toBeInTheDocument();
    expect(screen.queryByText('150')).not.toBeInTheDocument();
  });

  it('applies data-large attribute for counters >= 10', () => {
    const { container, rerender } = render(<Button counter={10}>Test</Button>);
    let counterBadge = container.querySelector('.ui-button__counter');
    expect(counterBadge).toHaveAttribute('data-large');

    rerender(<Button counter={5}>Test</Button>);
    counterBadge = container.querySelector('.ui-button__counter');
    expect(counterBadge).not.toHaveAttribute('data-large');
  });

  it('uses custom counter color when provided', () => {
    const customColor = '#ff5500';
    const { container } = render(
      <Button counter={3} counterColor={customColor}>Test</Button>
    );
    const counterBadge = container.querySelector('.ui-button__counter');
    expect(counterBadge).toHaveStyle({ backgroundColor: customColor });
  });

  it('uses default red color for counter when counterColor is not provided', () => {
    const { container } = render(<Button counter={3}>Test</Button>);
    const counterBadge = container.querySelector('.ui-button__counter');
    expect(counterBadge).toHaveStyle({ backgroundColor: 'var(--color-cancel)' });
  });

  it('renders as full width when fullWidth prop is true', () => {
    render(<Button fullWidth>Test</Button>);
    expect(screen.getByRole('button')).toHaveClass('ui-button--full-width');
  });

  it('passes through standard button props', () => {
    const handleClick = vi.fn();
    render(
      <Button onClick={handleClick} disabled aria-label="Custom Label">
        Test
      </Button>
    );
    const button = screen.getByRole('button');
    expect(button).toBeDisabled();
    expect(button).toHaveAttribute('aria-label', 'Custom Label');
  });

  describe('State-based functionality', () => {
    it('renders in idle state by default', () => {
      render(<Button>Test</Button>);
      const button = screen.getByRole('button');
      expect(button).not.toHaveClass('ui-button--loading');
      expect(button).not.toHaveClass('ui-button--success');
      expect(button).not.toHaveClass('ui-button--error');
      expect(button).not.toBeDisabled();
    });

    it('renders loading state with spinner', () => {
      const { container } = render(<Button state="loading">Save</Button>);
      const button = screen.getByRole('button');
      
      expect(button).toHaveClass('ui-button--loading');
      expect(button).toBeDisabled();
      expect(button).toHaveAttribute('aria-busy', 'true');
      
      const spinner = container.querySelector('.ui-button__spinner');
      expect(spinner).toBeInTheDocument();
      expect(screen.getByText('Save')).toBeInTheDocument();
    });

    it('renders success state with checkmark', () => {
      render(<Button state="success" successMessage="Saved!">Save</Button>);
      const button = screen.getByRole('button');
      
      expect(button).toHaveClass('ui-button--success');
      expect(button).toBeDisabled();
      expect(screen.getByText('Saved!')).toBeInTheDocument();
      expect(screen.getByText('✓')).toBeInTheDocument();
    });

    it('uses children as fallback when no success message provided', () => {
      render(<Button state="success">Save</Button>);
      expect(screen.getByText('Save')).toBeInTheDocument();
      expect(screen.getByText('✓')).toBeInTheDocument();
    });

    it('renders error state with X icon', () => {
      render(<Button state="error" errorMessage="Failed!">Save</Button>);
      const button = screen.getByRole('button');
      
      expect(button).toHaveClass('ui-button--error');
      expect(button).not.toBeDisabled(); // Error state should allow retry
      expect(screen.getByText('Failed!')).toBeInTheDocument();
      expect(screen.getByText('✕')).toBeInTheDocument();
    });

    it('uses children as fallback when no error message provided', () => {
      render(<Button state="error">Save</Button>);
      expect(screen.getByText('Save')).toBeInTheDocument();
      expect(screen.getByText('✕')).toBeInTheDocument();
    });

    it('renders cancel button during loading state when onCancel is provided', () => {
      const onCancel = vi.fn();
      render(
        <Button state="loading" onCancel={onCancel}>
          Saving...
        </Button>
      );

      const cancelButton = screen.getByLabelText('Cancel');
      expect(cancelButton).toBeInTheDocument();
      expect(cancelButton).toHaveClass('ui-button__cancel');
    });

    it('calls onCancel when cancel button is clicked', () => {
      const onCancel = vi.fn();
      render(
        <Button state="loading" onCancel={onCancel}>
          Saving...
        </Button>
      );

      const cancelButton = screen.getByLabelText('Cancel');
      cancelButton.click();
      
      expect(onCancel).toHaveBeenCalledTimes(1);
    });

    it('does not render cancel button when onCancel is not provided', () => {
      render(<Button state="loading">Saving...</Button>);
      
      const cancelButton = screen.queryByLabelText('Cancel');
      expect(cancelButton).not.toBeInTheDocument();
    });

    it('hides counter badge during non-idle states', () => {
      const { rerender } = render(
        <Button counter={5} state="loading">
          Save
        </Button>
      );
      expect(screen.queryByText('5')).not.toBeInTheDocument();

      rerender(
        <Button counter={5} state="success">
          Save
        </Button>
      );
      expect(screen.queryByText('5')).not.toBeInTheDocument();

      rerender(
        <Button counter={5} state="error">
          Save
        </Button>
      );
      expect(screen.queryByText('5')).not.toBeInTheDocument();

      rerender(
        <Button counter={5} state="idle">
          Save
        </Button>
      );
      expect(screen.getByText('5')).toBeInTheDocument();
    });

    it('respects explicit disabled prop even in idle state', () => {
      render(<Button disabled>Test</Button>);
      expect(screen.getByRole('button')).toBeDisabled();
    });
  });
});
