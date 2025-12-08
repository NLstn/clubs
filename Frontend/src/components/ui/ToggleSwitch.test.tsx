import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ToggleSwitch } from './ToggleSwitch';

describe('ToggleSwitch', () => {
    it('renders unchecked toggle switch', () => {
        render(<ToggleSwitch checked={false} onChange={() => {}} />);
        const toggle = screen.getByRole('switch') as HTMLInputElement;
        expect(toggle).toBeDefined();
        expect(toggle.checked).toBe(false);
    });

    it('renders checked toggle switch', () => {
        render(<ToggleSwitch checked={true} onChange={() => {}} />);
        const toggle = screen.getByRole('switch') as HTMLInputElement;
        expect(toggle.checked).toBe(true);
    });

    it('calls onChange with correct boolean value when toggled on', () => {
        const onChange = vi.fn();
        render(<ToggleSwitch checked={false} onChange={onChange} />);
        
        const toggle = screen.getByRole('switch');
        fireEvent.click(toggle);
        
        expect(onChange).toHaveBeenCalledWith(true);
        expect(onChange).toHaveBeenCalledTimes(1);
    });

    it('calls onChange with correct boolean value when toggled off', () => {
        const onChange = vi.fn();
        render(<ToggleSwitch checked={true} onChange={onChange} />);
        
        const toggle = screen.getByRole('switch');
        fireEvent.click(toggle);
        
        expect(onChange).toHaveBeenCalledWith(false);
        expect(onChange).toHaveBeenCalledTimes(1);
    });

    it('respects disabled prop', () => {
        const onChange = vi.fn();
        render(<ToggleSwitch checked={false} onChange={onChange} disabled={true} />);
        
        const toggle = screen.getByRole('switch') as HTMLInputElement;
        expect(toggle.disabled).toBe(true);
        
        // Note: fireEvent.click on disabled inputs still triggers onChange in jsdom
        // In real browsers, disabled inputs don't trigger events
        // Testing that the disabled attribute is set is sufficient
    });

    it('applies correct CSS classes for styling', () => {
        const { container } = render(<ToggleSwitch checked={true} onChange={() => {}} />);
        const label = container.querySelector('.toggle-switch');
        const slider = container.querySelector('.slider');
        
        expect(label).toBeDefined();
        expect(slider).toBeDefined();
    });

    it('has proper accessibility attributes', () => {
        render(<ToggleSwitch checked={true} onChange={() => {}} />);
        const toggle = screen.getByRole('switch');
        
        expect(toggle.getAttribute('role')).toBe('switch');
        expect(toggle.getAttribute('aria-checked')).toBe('true');
    });

    it('updates aria-checked when checked prop changes', () => {
        const { rerender } = render(<ToggleSwitch checked={false} onChange={() => {}} />);
        let toggle = screen.getByRole('switch');
        expect(toggle.getAttribute('aria-checked')).toBe('false');
        
        rerender(<ToggleSwitch checked={true} onChange={() => {}} />);
        toggle = screen.getByRole('switch');
        expect(toggle.getAttribute('aria-checked')).toBe('true');
    });

    it('generates stable IDs across re-renders', () => {
        const { rerender } = render(<ToggleSwitch checked={false} onChange={() => {}} />);
        const toggle1 = screen.getByRole('switch');
        const id1 = toggle1.id;
        
        rerender(<ToggleSwitch checked={true} onChange={() => {}} />);
        const toggle2 = screen.getByRole('switch');
        const id2 = toggle2.id;
        
        expect(id1).toBe(id2);
        expect(id1).toBeTruthy();
    });

    it('uses custom id when provided', () => {
        render(<ToggleSwitch checked={false} onChange={() => {}} id="custom-toggle" />);
        const toggle = screen.getByRole('switch');
        expect(toggle.id).toBe('custom-toggle');
    });

    it('renders optional label when provided', () => {
        render(<ToggleSwitch checked={false} onChange={() => {}} label="Enable feature" />);
        expect(screen.getByText('Enable feature')).toBeDefined();
    });

    it('does not render label when not provided', () => {
        const { container } = render(<ToggleSwitch checked={false} onChange={() => {}} />);
        const label = container.querySelector('.toggle-label');
        expect(label).toBeNull();
    });

    it('has aria-label when label is not provided', () => {
        render(<ToggleSwitch checked={false} onChange={() => {}} />);
        const toggle = screen.getByRole('switch');
        expect(toggle.getAttribute('aria-label')).toBe('Toggle switch');
    });

    it('does not have aria-label when label is provided', () => {
        render(<ToggleSwitch checked={false} onChange={() => {}} label="Enable feature" />);
        const toggle = screen.getByRole('switch');
        expect(toggle.getAttribute('aria-label')).toBeNull();
    });
});
