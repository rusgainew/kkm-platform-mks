import React, { createContext, useContext, useState } from 'react';

interface TabsContextValue {
  activeTab: string;
  setActiveTab: (value: string) => void;
  variant: 'line' | 'pills';
}

const TabsContext = createContext<TabsContextValue | undefined>(undefined);

const useTabsContext = () => {
  const context = useContext(TabsContext);
  if (!context) {
    throw new Error('Tabs components must be used within a Tabs provider');
  }
  return context;
};

export interface TabsProps {
  /** Активная вкладка по умолчанию */
  defaultValue: string;
  /** Контролируемое значение активной вкладки */
  value?: string;
  /** Callback при изменении вкладки */
  onChange?: (value: string) => void;
  /** Вариант отображения */
  variant?: 'line' | 'pills';
  /** Дочерние элементы */
  children: React.ReactNode;
  /** Дополнительные CSS классы */
  className?: string;
}

export const Tabs: React.FC<TabsProps> = ({
  defaultValue,
  value,
  onChange,
  variant = 'line',
  children,
  className = '',
}) => {
  const [internalValue, setInternalValue] = useState(defaultValue);
  const activeTab = value !== undefined ? value : internalValue;

  const setActiveTab = (newValue: string) => {
    if (value === undefined) {
      setInternalValue(newValue);
    }
    onChange?.(newValue);
  };

  return (
    <TabsContext.Provider value={{ activeTab, setActiveTab, variant }}>
      <div className={`w-full ${className}`}>{children}</div>
    </TabsContext.Provider>
  );
};

export interface TabsListProps {
  /** Дочерние элементы (TabsTrigger) */
  children: React.ReactNode;
  /** Растянуть на всю ширину */
  fullWidth?: boolean;
  /** Дополнительные CSS классы */
  className?: string;
}

export const TabsList: React.FC<TabsListProps> = ({
  children,
  fullWidth = false,
  className = '',
}) => {
  const { variant } = useTabsContext();

  return (
    <div
      role="tablist"
      className={`
        flex gap-1
        ${variant === 'line' ? 'border-b border-gray-200 dark:border-gray-700' : 'bg-gray-100 dark:bg-gray-800 p-1 rounded-lg'}
        ${fullWidth ? 'w-full' : ''}
        ${className}
      `}
    >
      {children}
    </div>
  );
};

export interface TabsTriggerProps {
  /** Значение вкладки */
  value: string;
  /** Контент кнопки */
  children: React.ReactNode;
  /** Отключенное состояние */
  disabled?: boolean;
  /** Дополнительные CSS классы */
  className?: string;
}

export const TabsTrigger: React.FC<TabsTriggerProps> = ({
  value,
  children,
  disabled = false,
  className = '',
}) => {
  const { activeTab, setActiveTab, variant } = useTabsContext();
  const isActive = activeTab === value;

  const handleClick = () => {
    if (!disabled) {
      setActiveTab(value);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      if (!disabled) {
        setActiveTab(value);
      }
    }
  };

  if (variant === 'line') {
    return (
      <button
        role="tab"
        aria-selected={isActive}
        aria-disabled={disabled}
        disabled={disabled}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
        className={`
          px-4 py-2 text-sm font-medium transition-colors
          border-b-2 -mb-[1px]
          ${
            isActive
              ? 'text-blue-600 dark:text-blue-400 border-blue-600 dark:border-blue-400'
              : 'text-gray-600 dark:text-gray-400 border-transparent hover:text-gray-900 dark:hover:text-gray-200'
          }
          ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
          focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
          ${className}
        `}
      >
        {children}
      </button>
    );
  }

  return (
    <button
      role="tab"
      aria-selected={isActive}
      aria-disabled={disabled}
      disabled={disabled}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
      className={`
        px-4 py-2 text-sm font-medium rounded-md transition-colors
        ${
          isActive
            ? 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white shadow-sm'
            : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
        }
        ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
        focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
        ${className}
      `}
    >
      {children}
    </button>
  );
};

export interface TabsContentProps {
  /** Значение вкладки */
  value: string;
  /** Контент вкладки */
  children: React.ReactNode;
  /** Дополнительные CSS классы */
  className?: string;
}

export const TabsContent: React.FC<TabsContentProps> = ({
  value,
  children,
  className = '',
}) => {
  const { activeTab } = useTabsContext();

  if (activeTab !== value) {
    return null;
  }

  return (
    <div role="tabpanel" className={`py-4 ${className}`}>
      {children}
    </div>
  );
};
