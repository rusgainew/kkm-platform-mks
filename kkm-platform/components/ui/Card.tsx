/**
 * Card Component
 * Universal container for content sections
 * @version 1.0
 * @date 2026-01-28
 */

'use client';

import React from 'react';

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'bordered' | 'elevated';
  padding?: 'none' | 'sm' | 'md' | 'lg';
  hover?: boolean;
}

export interface CardHeaderProps extends Omit<React.HTMLAttributes<HTMLDivElement>, 'title'> {
  title?: string | React.ReactNode;
  subtitle?: string | React.ReactNode;
  action?: React.ReactNode;
}

export interface CardFooterProps extends React.HTMLAttributes<HTMLDivElement> {
  align?: 'left' | 'center' | 'right' | 'between';
}

const variantClasses = {
  default: 'bg-gray-900 border-gray-800',
  bordered: 'bg-gray-900 border-gray-700',
  elevated: 'bg-gray-900 border-gray-800 shadow-lg',
};

const paddingClasses = {
  none: '',
  sm: 'p-4',
  md: 'p-6',
  lg: 'p-8',
};

const footerAlignClasses = {
  left: 'justify-start',
  center: 'justify-center',
  right: 'justify-end',
  between: 'justify-between',
};

export function Card({
  variant = 'default',
  padding = 'md',
  hover = false,
  className = '',
  children,
  ...props
}: CardProps) {
  const baseClasses = 'rounded-lg border';
  const hoverClass = hover ? 'hover:shadow-lg hover:shadow-blue-500/20 transition-shadow' : '';

  return (
    <div
      className={`${baseClasses} ${variantClasses[variant]} ${paddingClasses[padding]} ${hoverClass} ${className}`}
      {...props}
    >
      {children}
    </div>
  );
}

export function CardHeader({
  title,
  subtitle,
  action,
  className = '',
  children,
  ...props
}: CardHeaderProps) {
  return (
    <div className={`flex items-start justify-between gap-4 mb-4 ${className}`} {...props}>
      <div className="flex-1">
        {title && (
          <h3 className="text-lg font-semibold text-white">
            {typeof title === 'string' ? title : title}
          </h3>
        )}
        {subtitle && (
          <p className="text-sm text-gray-400 mt-1">
            {typeof subtitle === 'string' ? subtitle : subtitle}
          </p>
        )}
        {children}
      </div>
      {action && <div className="flex-shrink-0">{action}</div>}
    </div>
  );
}

export function CardContent({
  className = '',
  children,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={className} {...props}>
      {children}
    </div>
  );
}

export function CardFooter({
  align = 'right',
  className = '',
  children,
  ...props
}: CardFooterProps) {
  return (
    <div
      className={`flex items-center gap-3 mt-4 pt-4 border-t border-gray-800 ${footerAlignClasses[align]} ${className}`}
      {...props}
    >
      {children}
    </div>
  );
}
