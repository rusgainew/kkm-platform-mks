/**
 * Spinner Component
 * Loading indicator with different sizes and variants
 * @version 1.0
 * @date 2026-01-28
 */

'use client';

import React from 'react';
import { Loader2 } from 'lucide-react';

export interface SpinnerProps extends React.HTMLAttributes<HTMLDivElement> {
  size?: 'sm' | 'md' | 'lg' | 'xl';
  variant?: 'primary' | 'secondary' | 'white';
  centered?: boolean;
  label?: string;
}

const sizeClasses = {
  sm: 'w-4 h-4',
  md: 'w-6 h-6',
  lg: 'w-8 h-8',
  xl: 'w-12 h-12',
};

const colorClasses = {
  primary: 'text-blue-500',
  secondary: 'text-gray-400',
  white: 'text-white',
};

export function Spinner({
  size = 'md',
  variant = 'primary',
  centered = false,
  label,
  className = '',
  ...props
}: SpinnerProps) {
  const spinner = (
    <div className={`inline-flex flex-col items-center gap-2 ${className}`} {...props}>
      <Loader2 className={`${sizeClasses[size]} ${colorClasses[variant]} animate-spin`} />
      {label && <span className="text-sm text-gray-400">{label}</span>}
    </div>
  );

  if (centered) {
    return (
      <div className="flex items-center justify-center min-h-[200px]" {...props}>
        {spinner}
      </div>
    );
  }

  return spinner;
}

export function SpinnerOverlay({
  label,
  className = '',
  ...props
}: Omit<SpinnerProps, 'centered'>) {
  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
      <Spinner size="lg" variant="white" label={label} className={className} {...props} />
    </div>
  );
}
