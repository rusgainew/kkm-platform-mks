'use client';

import React, { Component, ReactNode, ErrorInfo } from 'react';
import Link from 'next/link';
import { AlertTriangle, RefreshCw, Home } from 'lucide-react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
}

/**
 * Error Boundary компонент для перехвата ошибок в React приложении
 * Отображает fallback UI при возникновении ошибки в дочерних компонентах
 * Логирует ошибки для отладки
 */
export default class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
    };
  }

  static getDerivedStateFromError(error: Error): State {
    // Обновляем state для отображения fallback UI
    return {
      hasError: true,
      error,
      errorInfo: null,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Логируем ошибку в консоль для разработки
    console.error('[ErrorBoundary] Перехвачена ошибка:', error);
    console.error('[ErrorBoundary] Информация об ошибке:', errorInfo);

    // Обновляем state с полной информацией об ошибке
    this.setState({
      error,
      errorInfo,
    });

    // Вызываем callback если он был передан
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }

    // TODO: Отправить ошибку на сервер логирования (Sentry, LogRocket и т.д.)
  }

  handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
    });
  };

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen bg-linear-to-br from-gray-950 via-gray-900 to-gray-950 flex items-center justify-center p-4">
          <div className="max-w-md w-full">
            {/* Error Icon */}
            <div className="flex justify-center mb-6">
              <div className="w-16 h-16 bg-red-900/20 rounded-full flex items-center justify-center">
                <AlertTriangle className="w-8 h-8 text-red-400" />
              </div>
            </div>

            {/* Error Title */}
            <h1 className="text-2xl font-bold text-white text-center mb-2">
              Что-то пошло не так
            </h1>

            {/* Error Description */}
            <p className="text-gray-400 text-center mb-6">
              Произошла непредвиденная ошибка. Пожалуйста, попробуйте снова или вернитесь на главную страницу.
            </p>

            {/* Error Details (только в разработке) */}
            {process.env.NODE_ENV === 'development' && this.state.error && (
              <div className="mb-6 p-4 bg-red-900/10 border border-red-800 rounded-lg">
                <p className="text-red-300 text-sm font-mono wrap-break-word">
                  <span className="font-semibold">Ошибка:</span> {this.state.error.toString()}
                </p>
                {this.state.errorInfo && (
                  <details className="mt-3 cursor-pointer">
                    <summary className="text-red-300 text-xs font-semibold hover:text-red-200">
                      Stack trace
                    </summary>
                    <pre className="mt-2 text-red-300/70 text-xs overflow-auto max-h-40 whitespace-pre-wrap">
                      {this.state.errorInfo.componentStack}
                    </pre>
                  </details>
                )}
              </div>
            )}

            {/* Actions */}
            <div className="flex gap-3">
              <button
                onClick={this.handleReset}
                className="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition"
              >
                <RefreshCw className="w-4 h-4" />
                Повторить
              </button>
              <Link
                href="/"
                className="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-gray-800 hover:bg-gray-700 text-white rounded-lg font-medium transition"
              >
                <Home className="w-4 h-4" />
                На главную
              </Link>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
