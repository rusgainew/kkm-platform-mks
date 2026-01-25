'use client';

import React, { ReactNode, useState } from 'react';
import ErrorBoundary from './ErrorBoundary';

interface SafeClientComponentProps {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error) => void;
}

/**
 * Обертка для безопасного отображения клиентского компонента с обработкой ошибок
 * Используется для компонентов которые часто выбрасывают ошибки
 */
export const SafeClientComponent: React.FC<SafeClientComponentProps> = ({
  children,
  fallback,
  onError,
}) => {
  return (
    <ErrorBoundary
      onError={(error, errorInfo) => {
        console.error('[SafeClientComponent] Error:', error);
        onError?.(error);
      }}
    >
      {children}
    </ErrorBoundary>
  );
};

/**
 * HOC для обертки компонента с Error Boundary
 */
export const withErrorBoundary = <P extends object>(
  Component: React.ComponentType<P>,
  errorBoundaryProps?: {
    fallback?: ReactNode;
    onError?: (error: Error) => void;
  }
) => {
  const WrappedComponent = (props: P) => (
    <ErrorBoundary {...errorBoundaryProps}>
      <Component {...props} />
    </ErrorBoundary>
  );

  WrappedComponent.displayName = `withErrorBoundary(${Component.displayName || Component.name || 'Component'})`;

  return WrappedComponent;
};
