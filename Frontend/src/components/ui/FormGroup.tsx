import { ReactNode } from 'react';
import './FormGroup.css';

interface FormGroupProps {
    /**
     * The content to render inside the form group (typically input, textarea, select, or custom components)
     */
    children: ReactNode;
    /**
     * Optional className for additional styling
     */
    className?: string;
}

/**
 * FormGroup component - A reusable wrapper for form fields that provides consistent spacing and structure
 * 
 * @example
 * ```tsx
 * <FormGroup>
 *   <label htmlFor="email">Email</label>
 *   <input id="email" type="email" />
 * </FormGroup>
 * ```
 * 
 * @example
 * // With custom Input component
 * ```tsx
 * <FormGroup>
 *   <Input label="Username" value={username} onChange={handleChange} />
 * </FormGroup>
 * ```
 */
export const FormGroup = ({ children, className = '' }: FormGroupProps) => {
    const classes = ['form-group', className].filter(Boolean).join(' ');
    
    return (
        <div className={classes}>
            {children}
        </div>
    );
};

export default FormGroup;
