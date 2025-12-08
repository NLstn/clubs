import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import { FormGroup } from './FormGroup';

describe('FormGroup', () => {
    it('renders children correctly', () => {
        const { container } = render(
            <FormGroup>
                <label htmlFor="test-input">Test Label</label>
                <input id="test-input" type="text" />
            </FormGroup>
        );
        
        const formGroup = container.querySelector('.form-group');
        expect(formGroup).toBeDefined();
        expect(formGroup?.querySelector('label')).toBeDefined();
        expect(formGroup?.querySelector('input')).toBeDefined();
    });

    it('applies custom className when provided', () => {
        const { container } = render(
            <FormGroup className="custom-class">
                <input type="text" />
            </FormGroup>
        );
        
        const formGroup = container.querySelector('.form-group');
        expect(formGroup?.classList.contains('custom-class')).toBe(true);
    });

    it('combines default and custom classNames properly', () => {
        const { container } = render(
            <FormGroup className="extra-class another-class">
                <input type="text" />
            </FormGroup>
        );
        
        const formGroup = container.querySelector('.form-group');
        expect(formGroup?.classList.contains('form-group')).toBe(true);
        expect(formGroup?.classList.contains('extra-class')).toBe(true);
        expect(formGroup?.classList.contains('another-class')).toBe(true);
    });

    it('renders with different types of children - label and input', () => {
        const { container } = render(
            <FormGroup>
                <label>Email</label>
                <input type="email" placeholder="Enter email" />
            </FormGroup>
        );
        
        const formGroup = container.querySelector('.form-group');
        const label = formGroup?.querySelector('label');
        const input = formGroup?.querySelector('input');
        
        expect(label?.textContent).toBe('Email');
        expect(input?.getAttribute('type')).toBe('email');
        expect(input?.getAttribute('placeholder')).toBe('Enter email');
    });

    it('renders with textarea', () => {
        const { container } = render(
            <FormGroup>
                <label>Description</label>
                <textarea placeholder="Enter description" />
            </FormGroup>
        );
        
        const formGroup = container.querySelector('.form-group');
        const textarea = formGroup?.querySelector('textarea');
        
        expect(textarea).toBeDefined();
        expect(textarea?.getAttribute('placeholder')).toBe('Enter description');
    });

    it('renders with custom components', () => {
        const CustomInput = () => <div className="custom-input">Custom Input</div>;
        
        const { container } = render(
            <FormGroup>
                <label>Custom Field</label>
                <CustomInput />
            </FormGroup>
        );
        
        const formGroup = container.querySelector('.form-group');
        const customInput = formGroup?.querySelector('.custom-input');
        
        expect(customInput).toBeDefined();
        expect(customInput?.textContent).toBe('Custom Input');
    });

    it('renders with multiple children', () => {
        const { container } = render(
            <FormGroup>
                <label>Username</label>
                <input type="text" />
                <span className="helper-text">Enter your username</span>
            </FormGroup>
        );
        
        const formGroup = container.querySelector('.form-group');
        expect(formGroup?.querySelector('label')).toBeDefined();
        expect(formGroup?.querySelector('input')).toBeDefined();
        expect(formGroup?.querySelector('.helper-text')).toBeDefined();
    });

    it('does not break with empty className', () => {
        const { container } = render(
            <FormGroup className="">
                <input type="text" />
            </FormGroup>
        );
        
        const formGroup = container.querySelector('.form-group');
        expect(formGroup).toBeDefined();
        expect(formGroup?.classList.contains('form-group')).toBe(true);
    });
});
