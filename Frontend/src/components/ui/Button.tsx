import React from 'react';
import './Button.css';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'accept' | 'cancel' | 'maybe';
  size?: 'sm' | 'md' | 'lg';
  fullWidth?: boolean;
  counter?: number;
  counterColor?: string;
  children: React.ReactNode;
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  fullWidth = false,
  counter,
  counterColor = 'var(--color-cancel)',
  className = '',
  children,
  ...props
}) => {
  const baseClass = 'ui-button';
  const variantClass = `ui-button--${variant}`;
  const sizeClass = `ui-button--${size}`;
  const fullWidthClass = fullWidth ? 'ui-button--full-width' : '';
  const counterClass = counter && counter > 0 ? 'ui-button--with-counter' : '';

  const classes = [
    baseClass,
    variantClass,
    sizeClass,
    fullWidthClass,
    counterClass,
    className
  ].filter(Boolean).join(' ');

  return (
    <button className={classes} {...props}>
      <span className="ui-button__content">
        {children}
      </span>
      {counter && counter > 0 && (
        <span 
          className="ui-button__counter" 
          style={{ backgroundColor: counterColor }}
          {...(counter >= 10 ? { 'data-large': true } : {})}
        >
          {counter > 99 ? '99+' : counter}
        </span>
      )}
    </button>
  );
};

export default Button;
