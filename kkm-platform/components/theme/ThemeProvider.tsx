'use client';

import React, { ReactNode } from 'react';
import { ThemeProvider as NextThemeProvider } from 'next-themes';

interface ThemeProviderProps {
  children: ReactNode;
}

/**
 * Theme Provider для приложения
 * Обеспечивает управление темой через next-themes
 * - Автоматическое определение системной темы
 * - Сохранение выбора в localStorage
 * - Синхронизация между вкладками браузера
 */
export function ThemeProvider({ children }: ThemeProviderProps) {
  return (
    <NextThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      enableColorScheme
      disableTransitionOnChange
      storageKey="kkm-theme"
    >
      {children}
    </NextThemeProvider>
  );
}
