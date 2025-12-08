import React, { CSSProperties, ReactNode } from 'react';
import './Card.css';

export interface CardProps {
  children: ReactNode;
  clickable?: boolean;
  onClick?: () => void;
  variant?: 'default' | 'light' | 'dark' | 'white';
  padding?: 'sm' | 'md' | 'lg';
  className?: string;
  style?: CSSProperties;
  hover?: boolean;
}

export const Card: React.FC<CardProps> = ({
  children,
  clickable = false,
  onClick,
  variant = 'default',
  padding = 'md',
  className = '',
  style,
  hover = false,
}) => {
  const classes = [
    'card',
    `card--${variant}`,
    `card--padding-${padding}`,
    clickable || onClick ? 'card--clickable' : '',
    hover ? 'card--hover' : '',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classes} onClick={onClick} style={style}>
      {children}
    </div>
  );
};
