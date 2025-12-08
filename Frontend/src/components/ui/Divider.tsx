import React from 'react';
import './Divider.css';

export interface DividerProps {
    /** Optional text to display in the center of the divider */
    text?: string;
    /** Spacing variant for the divider */
    spacing?: 'sm' | 'md' | 'lg';
    /** Additional CSS class name */
    className?: string;
}

/**
 * Divider component for visual separation of content sections
 * 
 * @example
 * // Simple divider
 * <Divider />
 * 
 * @example
 * // Divider with text
 * <Divider text="OR" />
 * 
 * @example
 * // Custom spacing
 * <Divider spacing="lg" />
 */
const Divider: React.FC<DividerProps> = ({ 
    text, 
    spacing = 'md',
    className = '' 
}) => {
    const dividerClass = text 
        ? `divider divider--text divider--${spacing}` 
        : `divider divider--simple divider--${spacing}`;

    return (
        <div className={`${dividerClass} ${className}`.trim()}>
            {text && <span className="divider__text">{text}</span>}
        </div>
    );
};

export default Divider;
