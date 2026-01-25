'use client';

import React from 'react';

/**
 * Компонент для тестирования Error Boundary
 * ИСПОЛЬЗУЕТСЯ ТОЛЬКО ДЛЯ РАЗРАБОТКИ И ТЕСТИРОВАНИЯ
 */

interface ErrorTestComponentProps {
  shouldError?: boolean;
  errorMessage?: string;
}

/**
 * Компонент который выбрасывает ошибку для тестирования
 * Используется для проверки работы Error Boundary
 */
export const ErrorTestComponent: React.FC<ErrorTestComponentProps> = ({
  shouldError = false,
  errorMessage = 'Тестовая ошибка из ErrorTestComponent',
}) => {
  if (shouldError) {
    throw new Error(errorMessage);
  }

  return (
    <div className="p-4 bg-blue-900/20 border border-blue-800 rounded-lg">
      <p className="text-blue-300 text-sm">
        ErrorTestComponent готов к тестированию. Установите shouldError={true} для выброса ошибки.
      </p>
    </div>
  );
};

/**
 * Компонент который выбрасывает ошибку при клике на кнопку
 * для тестирования обработки ошибок в обработчиках событий
 */
export const EventErrorTestComponent: React.FC = () => {
  const handleClick = () => {
    // Это не будет перехвачено Error Boundary!
    // Error Boundary перехватывает ошибки при рендере, а не в обработчиках событий
    throw new Error('Ошибка в обработчике события');
  };

  return (
    <div className="p-4">
      <button
        onClick={handleClick}
        className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg font-medium"
      >
        Нажми для ошибки в обработчике
      </button>
      <p className="text-gray-400 text-sm mt-2">
        Эта ошибка НЕ будет перехвачена Error Boundary!
      </p>
    </div>
  );
};

/**
 * Компонент для тестирования различных типов ошибок
 */
export const MultiErrorTestComponent: React.FC = () => {
  const [errorType, setErrorType] = React.useState<string | null>(null);

  if (errorType === 'render') {
    throw new Error('Ошибка при рендере компонента');
  }

  if (errorType === 'undefined') {
    const obj: Record<string, unknown> | undefined = undefined;
    // Это вызовет ошибку: Cannot read property 'property' of undefined
    return <div>{String((obj as unknown as Record<string, unknown>)?.property)}</div>;
  }

  if (errorType === 'type') {
    const num: string | unknown = 'string';
    // Это вызовет ошибку: toFixed is not a function
    return <div>{(num as unknown as number).toFixed(2)}</div>;
  }

  return (
    <div className="p-4 border border-gray-700 rounded-lg space-y-2">
      <p className="text-gray-400 text-sm mb-3">Выберите тип ошибки для тестирования:</p>
      <button
        onClick={() => setErrorType('render')}
        className="block w-full px-3 py-2 bg-red-600 hover:bg-red-700 text-white rounded text-sm"
      >
        Ошибка при рендере
      </button>
      <button
        onClick={() => setErrorType('undefined')}
        className="block w-full px-3 py-2 bg-red-600 hover:bg-red-700 text-white rounded text-sm"
      >
        Undefined property access
      </button>
      <button
        onClick={() => setErrorType('type')}
        className="block w-full px-3 py-2 bg-red-600 hover:bg-red-700 text-white rounded text-sm"
      >
        Ошибка типа
      </button>
      <button
        onClick={() => setErrorType(null)}
        className="block w-full px-3 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded text-sm"
      >
        Сбросить
      </button>
    </div>
  );
};
