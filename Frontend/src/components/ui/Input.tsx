import { forwardRef, InputHTMLAttributes, TextareaHTMLAttributes } from 'react';
import './Input.css';

interface BaseInputProps {
    label?: string;
    error?: string;
    helperText?: string;
    variant?: 'default' | 'outline' | 'filled';
    size?: 'sm' | 'md' | 'lg';
}

interface SingleLineInputProps extends BaseInputProps, Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
    multiline?: false;
}

interface MultiLineInputProps extends BaseInputProps, Omit<TextareaHTMLAttributes<HTMLTextAreaElement>, 'size'> {
    multiline: true;
}

type InputProps = SingleLineInputProps | MultiLineInputProps;

export const Input = forwardRef<HTMLInputElement | HTMLTextAreaElement, InputProps>(
    ({ 
        label, 
        error, 
        helperText, 
        variant = 'default', 
        size = 'md', 
        className = '', 
        id,
        multiline,
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
            multiline ? 'input-textarea' : '',
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
                {multiline ? (
                    <textarea
                        ref={ref as React.Ref<HTMLTextAreaElement>}
                        id={inputId}
                        className={inputClasses}
                        {...(props as TextareaHTMLAttributes<HTMLTextAreaElement>)}
                    />
                ) : (
                    <input
                        ref={ref as React.Ref<HTMLInputElement>}
                        id={inputId}
                        className={inputClasses}
                        {...(props as InputHTMLAttributes<HTMLInputElement>)}
                    />
                )}
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
