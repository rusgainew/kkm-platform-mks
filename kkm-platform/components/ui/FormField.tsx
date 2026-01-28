/**
 * Form Field Components
 * Reusable form inputs with consistent styling
 * @version 1.0
 * @date 2026-01-28
 */

'use client';

import React from 'react';

export interface FormFieldProps {
  label: string;
  error?: string;
  required?: boolean;
  children: React.ReactNode;
}

export function FormField({ label, error, required, children }: FormFieldProps) {
  return (
    <div>
      <label className="block text-sm font-medium text-gray-300 mb-2">
        {label}
        {required && <span className="text-red-400 ml-1">*</span>}
      </label>
      {children}
      {error && <p className="mt-1 text-sm text-red-400">{error}</p>}
    </div>
  );
}

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  error?: boolean;
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ error, className = '', disabled, ...props }, ref) => {
    const baseClasses = 'w-full px-4 py-2 bg-gray-800 border rounded-lg transition-colors';
    const stateClasses = disabled
      ? 'border-gray-700 text-gray-400 cursor-not-allowed'
      : error
      ? 'border-red-600 text-white focus:ring-2 focus:ring-red-500 focus:border-red-500'
      : 'border-gray-700 text-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500';

    return (
      <input
        ref={ref}
        className={`${baseClasses} ${stateClasses} ${className}`}
        disabled={disabled}
        {...props}
      />
    );
  }
);

Input.displayName = 'Input';

export interface TextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  error?: boolean;
}

export const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ error, className = '', disabled, ...props }, ref) => {
    const baseClasses = 'w-full px-4 py-2 bg-gray-800 border rounded-lg transition-colors';
    const stateClasses = disabled
      ? 'border-gray-700 text-gray-400 cursor-not-allowed'
      : error
      ? 'border-red-600 text-white focus:ring-2 focus:ring-red-500 focus:border-red-500'
      : 'border-gray-700 text-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500';

    return (
      <textarea
        ref={ref}
        className={`${baseClasses} ${stateClasses} ${className}`}
        disabled={disabled}
        {...props}
      />
    );
  }
);

Textarea.displayName = 'Textarea';

export interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  error?: boolean;
}

export const Select = React.forwardRef<HTMLSelectElement, SelectProps>(
  ({ error, className = '', disabled, children, ...props }, ref) => {
    const baseClasses = 'w-full px-4 py-2 bg-gray-800 border rounded-lg transition-colors';
    const stateClasses = disabled
      ? 'border-gray-700 text-gray-400 cursor-not-allowed'
      : error
      ? 'border-red-600 text-white focus:ring-2 focus:ring-red-500 focus:border-red-500'
      : 'border-gray-700 text-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500';

    return (
      <select
        ref={ref}
        className={`${baseClasses} ${stateClasses} ${className}`}
        disabled={disabled}
        {...props}
      >
        {children}
      </select>
    );
  }
);

Select.displayName = 'Select';
