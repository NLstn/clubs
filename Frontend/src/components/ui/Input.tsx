import { forwardRef, InputHTMLAttributes } from 'react';
import './Input.css';

interface InputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
    label?: string;
    error?: string;
    helperText?: string;
    variant?: 'default' | 'outline' | 'filled';
    size?: 'sm' | 'md' | 'lg';
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
    ({ 
        label, 
        error, 
        helperText, 
        variant = 'default', 
        size = 'md', 
        className = '', 
        id,
        ...props 
    }, ref) => {
        const baseClasses = 'input-base';
        const variantClasses = {
            default: 'input-default',
            outline: 'input-outline',
            filled: 'input-filled'
        };
        const sizeClasses = {
            sm: 'input-sm',
            md: 'input-md',
            lg: 'input-lg'
        };
        
        const inputClasses = [
            baseClasses,
            variantClasses[variant],
            sizeClasses[size],
            error ? 'input-error' : '',
            className
        ].filter(Boolean).join(' ');

        const inputId = id || `input-${Math.random().toString(36).substr(2, 9)}`;

        return (
            <div className="input-container">
                {label && (
                    <label htmlFor={inputId} className="input-label">
                        {label}
                    </label>
                )}
                <input
                    ref={ref}
                    id={inputId}
                    className={inputClasses}
                    {...props}
                />
                {error && (
                    <span className="input-error-text">{error}</span>
                )}
                {helperText && !error && (
                    <span className="input-helper-text">{helperText}</span>
                )}
            </div>
        );
    }
);

Input.displayName = 'Input';

export default Input;
