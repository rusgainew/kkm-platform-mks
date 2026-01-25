'use client';

import React from 'react';
import { Locale } from '@/i18n/config';
import { useTranslationsSync } from '@/lib/hooks/useTranslations';
import { ThemeSwitcher, ThemeSwitcherCompact, ThemeSwitcherDropdown } from '@/components/theme/ThemeSwitcher';
import { useThemeCustom } from '@/lib/hooks/useTheme';
import {
  ThemeExampleCard,
  ThemeExampleButton,
  ThemeExampleForm,
  ThemeExampleGradient,
  ThemeExampleTable,
  ThemeExampleAlert,
} from '@/components/theme/ThemeExamples';

export default function ThemeDemoPage({ params }: { params: { locale: Locale } }) {
  const locale = params.locale;
  const t = useTranslationsSync(locale);
  const { isDark, mounted, theme } = useThemeCustom();

  if (!mounted) {
    return (
      <div className="p-8 text-center">
        <p className="text-gray-600 dark:text-gray-300">{t('common.loading')}</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white dark:bg-gray-900 transition">
      {/* Фиксированный заголовок */}
      <header className="sticky top-0 z-40 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm">
        <div className="max-w-7xl mx-auto px-8 py-4 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
              🌙 Dark Mode Demo
            </h1>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Current theme: <span className="font-semibold">{theme}</span> ({isDark ? 'Dark' : 'Light'})
            </p>
          </div>
          <ThemeSwitcher />
        </div>
      </header>

      {/* Основной контент */}
      <main className="max-w-7xl mx-auto px-8 py-12 space-y-12">
        {/* Переключатели темы */}
        <section>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
            🎛️ Theme Switchers
          </h2>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-6">
              <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-4">
                Full Switcher
              </h3>
              <ThemeSwitcher />
            </div>

            <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-6">
              <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-4">
                Compact Switcher
              </h3>
              <div className="flex gap-4">
                <ThemeSwitcherCompact />
                <ThemeSwitcherCompact />
                <ThemeSwitcherCompact />
              </div>
            </div>

            <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-6 md:col-span-2">
              <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-4">
                Dropdown Switcher
              </h3>
              <ThemeSwitcherDropdown />
            </div>
          </div>
        </section>

        {/* Примеры компонентов */}
        <section>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
            💻 UI Components Examples
          </h2>

          <div className="space-y-6">
            {/* Карточки */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Cards
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <ThemeExampleCard />
                <ThemeExampleCard />
                <ThemeExampleCard />
              </div>
            </div>

            {/* Кнопки */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Buttons
              </h3>
              <ThemeExampleButton />
            </div>

            {/* Форма */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Form Inputs
              </h3>
              <ThemeExampleForm />
            </div>

            {/* Градиент */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Gradient
              </h3>
              <ThemeExampleGradient />
            </div>

            {/* Таблица */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Table
              </h3>
              <ThemeExampleTable />
            </div>

            {/* Алерты */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Alerts
              </h3>
              <ThemeExampleAlert />
            </div>
          </div>
        </section>

        {/* Информация */}
        <section className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6">
          <h3 className="text-lg font-bold text-blue-900 dark:text-blue-200 mb-4">
            ℹ️ Об этой демонстрации
          </h3>
          <ul className="text-sm text-blue-800 dark:text-blue-300 space-y-2">
            <li>✅ Темная и светлая темы полностью интегрированы</li>
            <li>✅ Выбор темы сохраняется в localStorage</li>
            <li>✅ Поддержка системной темы (автоматический выбор)</li>
            <li>✅ Плавные переходы между темами</li>
            <li>✅ Все компоненты оптимизированы для обеих тем</li>
            <li>✅ Используется next-themes + Tailwind CSS</li>
            <li>✅ Полная поддержка TypeScript</li>
          </ul>
        </section>

        {/* Рекомендации */}
        <section>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
            💡 Best Practices
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-6">
              <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-3">
                ✓ Рекомендуется
              </h3>
              <ul className="text-sm text-gray-700 dark:text-gray-300 space-y-2">
                <li>• Используйте Tailwind dark: prefix</li>
                <li>• Тестируйте в обеих темах</li>
                <li>• Проверяйте контраст текста</li>
                <li>• Используйте semantic colors</li>
              </ul>
            </div>
            <div className="bg-red-50 dark:bg-red-900/20 rounded-lg p-6">
              <h3 className="text-lg font-bold text-red-900 dark:text-red-200 mb-3">
                ✕ Избегайте
              </h3>
              <ul className="text-sm text-red-800 dark:text-red-300 space-y-2">
                <li>• Жестко закодированные цвета</li>
                <li>• Низкий контраст текста</li>
                <li>• Изображения без fallback</li>
                <li>• Игнорирование доступности</li>
              </ul>
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
