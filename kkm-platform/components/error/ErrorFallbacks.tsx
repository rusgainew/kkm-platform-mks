'use client';

import React from 'react';
import { AlertCircle, RotateCcw } from 'lucide-react';

interface ErrorFallbackProps {
  error: Error;
  reset: () => void;
  title?: string;
  description?: string;
}

/**
 * Fallback компонент для ошибок в таблице компаний
 * Отображается когда произойдет ошибка при загрузке/отображении компаний
 */
export const CompaniesErrorFallback: React.FC<ErrorFallbackProps> = ({
  error,
  reset,
  title = 'Ошибка при загрузке компаний',
  description = 'Не удалось загрузить список компаний. Пожалуйста, попробуйте снова.',
}) => (
  <div className="rounded-lg border border-red-800 bg-red-900/10 p-6 text-center">
    <div className="flex justify-center mb-4">
      <AlertCircle className="w-12 h-12 text-red-400" />
    </div>
    <h3 className="text-lg font-semibold text-white mb-2">{title}</h3>
    <p className="text-red-300 text-sm mb-4">{description}</p>
    {process.env.NODE_ENV === 'development' && (
      <p className="text-red-400/70 text-xs mb-4 font-mono">{error.message}</p>
    )}
    <button
      onClick={reset}
      className="inline-flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg font-medium transition"
    >
      <RotateCcw className="w-4 h-4" />
      Повторить
    </button>
  </div>
);

/**
 * Fallback компонент для ошибок в таблице пользователей
 */
export const UsersErrorFallback: React.FC<ErrorFallbackProps> = ({
  error,
  reset,
  title = 'Ошибка при загрузке пользователей',
  description = 'Не удалось загрузить список пользователей. Пожалуйста, попробуйте снова.',
}) => (
  <div className="rounded-lg border border-red-800 bg-red-900/10 p-6 text-center">
    <div className="flex justify-center mb-4">
      <AlertCircle className="w-12 h-12 text-red-400" />
    </div>
    <h3 className="text-lg font-semibold text-white mb-2">{title}</h3>
    <p className="text-red-300 text-sm mb-4">{description}</p>
    {process.env.NODE_ENV === 'development' && (
      <p className="text-red-400/70 text-xs mb-4 font-mono">{error.message}</p>
    )}
    <button
      onClick={reset}
      className="inline-flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg font-medium transition"
    >
      <RotateCcw className="w-4 h-4" />
      Повторить
    </button>
  </div>
);

/**
 * Fallback компонент для ошибок в дашборде
 */
export const DashboardErrorFallback: React.FC<ErrorFallbackProps> = ({
  error,
  reset,
  title = 'Ошибка при загрузке дашборда',
  description = 'Не удалось загрузить данные дашборда. Пожалуйста, попробуйте снова.',
}) => (
  <div className="rounded-lg border border-red-800 bg-red-900/10 p-6 text-center">
    <div className="flex justify-center mb-4">
      <AlertCircle className="w-12 h-12 text-red-400" />
    </div>
    <h3 className="text-lg font-semibold text-white mb-2">{title}</h3>
    <p className="text-red-300 text-sm mb-4">{description}</p>
    {process.env.NODE_ENV === 'development' && (
      <p className="text-red-400/70 text-xs mb-4 font-mono">{error.message}</p>
    )}
    <button
      onClick={reset}
      className="inline-flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg font-medium transition"
    >
      <RotateCcw className="w-4 h-4" />
      Повторить
    </button>
  </div>
);

/**
 * Fallback компонент для ошибок в POS системе
 */
export const POSErrorFallback: React.FC<ErrorFallbackProps> = ({
  error,
  reset,
  title = 'Ошибка в POS системе',
  description = 'Произошла ошибка при обработке операции. Пожалуйста, попробуйте снова.',
}) => (
  <div className="rounded-lg border border-red-800 bg-red-900/10 p-6 text-center">
    <div className="flex justify-center mb-4">
      <AlertCircle className="w-12 h-12 text-red-400" />
    </div>
    <h3 className="text-lg font-semibold text-white mb-2">{title}</h3>
    <p className="text-red-300 text-sm mb-4">{description}</p>
    {process.env.NODE_ENV === 'development' && (
      <p className="text-red-400/70 text-xs mb-4 font-mono">{error.message}</p>
    )}
    <button
      onClick={reset}
      className="inline-flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg font-medium transition"
    >
      <RotateCcw className="w-4 h-4" />
      Повторить
    </button>
  </div>
);
