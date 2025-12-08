import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import Divider from '../Divider';
import '@testing-library/jest-dom';

describe('Divider', () => {
  it('renders simple divider by default', () => {
    const { container } = render(<Divider />);
    const divider = container.querySelector('.divider');
    
    expect(divider).toBeInTheDocument();
    expect(divider).toHaveClass('divider--simple');
    expect(divider).toHaveClass('divider--md'); // default spacing
  });

  it('renders text divider when text prop is provided', () => {
    render(<Divider text="OR" />);
    
    const text = screen.getByText('OR');
    expect(text).toBeInTheDocument();
    expect(text).toHaveClass('divider__text');
    
    const divider = text.parentElement;
    expect(divider).toHaveClass('divider--text');
    expect(divider).not.toHaveClass('divider--simple');
  });

  it('applies small spacing variant', () => {
    const { container } = render(<Divider spacing="sm" />);
    const divider = container.querySelector('.divider');
    
    expect(divider).toHaveClass('divider--sm');
  });

  it('applies medium spacing variant', () => {
    const { container } = render(<Divider spacing="md" />);
    const divider = container.querySelector('.divider');
    
    expect(divider).toHaveClass('divider--md');
  });

  it('applies large spacing variant', () => {
    const { container } = render(<Divider spacing="lg" />);
    const divider = container.querySelector('.divider');
    
    expect(divider).toHaveClass('divider--lg');
  });

  it('applies custom className', () => {
    const { container } = render(<Divider className="custom-divider" />);
    const divider = container.querySelector('.divider');
    
    expect(divider).toHaveClass('custom-divider');
  });

  it('applies custom className with text divider', () => {
    render(<Divider text="OR" className="login-divider" />);
    
    const text = screen.getByText('OR');
    const divider = text.parentElement;
    
    expect(divider).toHaveClass('login-divider');
    expect(divider).toHaveClass('divider--text');
  });

  it('renders without text span when no text provided', () => {
    const { container } = render(<Divider />);
    const span = container.querySelector('.divider__text');
    
    expect(span).not.toBeInTheDocument();
  });

  it('combines spacing and text props correctly', () => {
    render(<Divider text="AND" spacing="lg" />);
    
    const text = screen.getByText('AND');
    const divider = text.parentElement;
    
    expect(divider).toHaveClass('divider--text');
    expect(divider).toHaveClass('divider--lg');
  });
});
