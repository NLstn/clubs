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
});
