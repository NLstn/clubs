import React from 'react';
import './Button.css';

export type ButtonState = 'idle' | 'loading' | 'success' | 'error';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'accept' | 'cancel' | 'maybe';
  size?: 'sm' | 'md' | 'lg';
  fullWidth?: boolean;
  counter?: number;
  counterColor?: string;
  children: React.ReactNode;
  
  // Feedback state props
  state?: ButtonState;
  successMessage?: string;
  errorMessage?: string;
  onCancel?: () => void; // Callback when cancel is clicked during loading state
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  fullWidth = false,
  counter,
  counterColor = 'var(--color-cancel)',
  className = '',
  children,
  state = 'idle',
  successMessage,
  errorMessage,
  onCancel,
  disabled,
  ...props
}) => {

  const baseClass = 'ui-button';
  const variantClass = `ui-button--${variant}`;
  const sizeClass = `ui-button--${size}`;
  const fullWidthClass = fullWidth ? 'ui-button--full-width' : '';
  const counterClass = counter != null && counter > 0 ? 'ui-button--with-counter' : '';
  const stateClass = state !== 'idle' ? `ui-button--${state}` : '';

  const classes = [
    baseClass,
    variantClass,
    sizeClass,
    fullWidthClass,
    counterClass,
    stateClass,
    className
  ].filter(Boolean).join(' ');

  // Determine if button should be disabled
  const isDisabled = disabled || state === 'loading' || state === 'success';

  // Render content based on state
  const renderContent = () => {
    if (state === 'loading') {
      return (
        <span className="ui-button__content">
          <span className="ui-button__spinner" aria-hidden="true"></span>
          <span>{children}</span>
        </span>
      );
    }

    if (state === 'success') {
      return (
        <span className="ui-button__content">
          <span className="ui-button__icon ui-button__icon--success" aria-hidden="true">✓</span>
          <span>{successMessage || children}</span>
        </span>
      );
    }

    if (state === 'error') {
      return (
        <span className="ui-button__content">
          <span className="ui-button__icon ui-button__icon--error" aria-hidden="true">✕</span>
          <span>{errorMessage || children}</span>
        </span>
      );
    }

    return (
      <span className="ui-button__content">
        {children}
      </span>
    );
  };

  return (
    <div className="ui-button-wrapper">
      <button 
        className={classes} 
        disabled={isDisabled}
        aria-busy={state === 'loading'}
        {...props}
      >
        {renderContent()}
        
        {/* Counter badge */}
        {counter != null && counter > 0 && state === 'idle' && (
          <span 
            className="ui-button__counter" 
            style={{ backgroundColor: counterColor }}
            {...(counter >= 10 ? { 'data-large': true } : {})}
          >
            {counter > 99 ? '99+' : counter}
          </span>
        )}
      </button>
      
      {/* Cancel button during loading state */}
      {state === 'loading' && onCancel && (
        <button
          type="button"
          className="ui-button__cancel"
          onClick={(e) => {
            e.stopPropagation();
            onCancel();
          }}
          aria-label="Cancel"
        >
          ✕
        </button>
      )}
    </div>
  );
};

export default Button;
