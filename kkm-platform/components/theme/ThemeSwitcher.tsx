'use client';

import React, { useEffect, useState } from 'react';
import { useTheme } from 'next-themes';

type Theme = 'light' | 'dark' | 'system';

/**
 * Компонент для переключения темы
 * Поддерживает: Light, Dark, System (автоматический)
 */
export function ThemeSwitcher() {
  const { theme, setTheme, systemTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  // Избегаем hydration mismatch
  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return null;

  const currentTheme = theme === 'system' ? systemTheme : theme;
  const isDark = currentTheme === 'dark';

  return (
    <div className="flex items-center gap-1 bg-gray-100 dark:bg-gray-800 rounded-lg p-1">
      {/* Light mode */}
      <button
        onClick={() => setTheme('light')}
        className={`flex items-center gap-1 px-3 py-2 rounded transition ${
          theme === 'light'
            ? 'bg-white dark:bg-gray-700 text-yellow-500 shadow'
            : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
        }`}
        title="Светлая тема"
      >
        <span className="text-lg">☀️</span>
        <span className="text-sm font-medium hidden sm:inline">Light</span>
      </button>

      {/* Dark mode */}
      <button
        onClick={() => setTheme('dark')}
        className={`flex items-center gap-1 px-3 py-2 rounded transition ${
          theme === 'dark'
            ? 'bg-white dark:bg-gray-700 text-blue-400 shadow'
            : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
        }`}
        title="Темная тема"
      >
        <span className="text-lg">🌙</span>
        <span className="text-sm font-medium hidden sm:inline">Dark</span>
      </button>

      {/* System (auto) */}
      <button
        onClick={() => setTheme('system')}
        className={`flex items-center gap-1 px-3 py-2 rounded transition ${
          theme === 'system'
            ? 'bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 shadow'
            : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
        }`}
        title="Система (автоматический выбор)"
      >
        <span className="text-lg">⚙️</span>
        <span className="text-sm font-medium hidden sm:inline">Auto</span>
      </button>
    </div>
  );
}

/**
 * Компактный переключатель (только иконка)
 */
export function ThemeSwitcherCompact() {
  const { theme, setTheme, systemTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return null;

  const currentTheme = theme === 'system' ? systemTheme : theme;
  const isDark = currentTheme === 'dark';

  return (
    <button
      onClick={() => setTheme(isDark ? 'light' : 'dark')}
      className="p-2 rounded-lg bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 transition"
      title={isDark ? 'Светлая тема' : 'Темная тема'}
    >
      {isDark ? (
        <span className="text-xl">☀️</span>
      ) : (
        <span className="text-xl">🌙</span>
      )}
    </button>
  );
}

/**
 * Dropdown переключатель
 */
export function ThemeSwitcherDropdown() {
  const { theme, setTheme, systemTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return null;

  const currentTheme = theme === 'system' ? systemTheme : theme;
  const themeLabels = {
    light: '☀️ Светлая',
    dark: '🌙 Темная',
    system: '⚙️ Система',
  };

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg transition"
      >
        <span className="text-lg">{currentTheme === 'dark' ? '🌙' : '☀️'}</span>
        <span className="text-sm font-medium">{themeLabels[theme as Theme]}</span>
        <svg
          className={`w-4 h-4 transition ${isOpen ? 'rotate-180' : ''}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute top-full right-0 mt-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg z-50">
          {(['light', 'dark', 'system'] as const).map((t) => (
            <button
              key={t}
              onClick={() => {
                setTheme(t);
                setIsOpen(false);
              }}
              className={`w-full flex items-center gap-3 px-4 py-2 text-left transition ${
                theme === t
                  ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400'
                  : 'hover:bg-gray-50 dark:hover:bg-gray-700/50'
              } ${t !== 'system' ? 'border-b border-gray-200 dark:border-gray-700' : ''}`}
            >
              <span className="text-lg">
                {t === 'light' ? '☀️' : t === 'dark' ? '🌙' : '⚙️'}
              </span>
              <span className="flex-1">{themeLabels[t]}</span>
              {theme === t && <span>✓</span>}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
