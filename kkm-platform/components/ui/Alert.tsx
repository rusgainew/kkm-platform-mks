/**
 * Alert Component
 * For displaying important messages and notifications
 * @version 1.0
 * @date 2026-01-28
 */

'use client';

import React from 'react';
import { AlertCircle, CheckCircle2, Info, XCircle, X } from 'lucide-react';

export interface AlertProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'info' | 'success' | 'warning' | 'danger';
  title?: string;
  dismissible?: boolean;
  onDismiss?: () => void;
  icon?: React.ReactNode;
  hideIcon?: boolean;
}

const variantClasses = {
  info: 'bg-blue-500/10 border-blue-500/30 text-blue-400',
  success: 'bg-emerald-500/10 border-emerald-500/30 text-emerald-400',
  warning: 'bg-yellow-500/10 border-yellow-500/30 text-yellow-400',
  danger: 'bg-red-500/10 border-red-500/30 text-red-400',
};

const iconMap = {
  info: Info,
  success: CheckCircle2,
  warning: AlertCircle,
  danger: XCircle,
};

export function Alert({
  variant = 'info',
  title,
  dismissible = false,
  onDismiss,
  icon,
  hideIcon = false,
  className = '',
  children,
  ...props
}: AlertProps) {
  const IconComponent = iconMap[variant];

  return (
    <div
      role="alert"
      className={`flex gap-3 p-4 rounded-lg border ${variantClasses[variant]} ${className}`}
      {...props}
    >
      {!hideIcon && (
        <div className="flex-shrink-0">
          {icon || <IconComponent className="w-5 h-5" />}
        </div>
      )}
      <div className="flex-1">
        {title && <h4 className="font-semibold mb-1">{title}</h4>}
        <div className="text-sm opacity-90">{children}</div>
      </div>
      {dismissible && onDismiss && (
        <button
          type="button"
          onClick={onDismiss}
          className="flex-shrink-0 opacity-70 hover:opacity-100 transition-opacity focus:outline-none"
          aria-label="Dismiss alert"
        >
          <X className="w-4 h-4" />
        </button>
      )}
    </div>
  );
}
