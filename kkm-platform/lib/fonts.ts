import { Inter, JetBrains_Mono } from 'next/font/google';

/**
 * Оптимизация шрифтов для быстрой загрузки
 * Используется variable fonts и preload стратегия
 */

// Основной шрифт - загружается с preload
export const inter = Inter({
  subsets: ['latin', 'cyrillic'],
  display: 'swap', // Быстрее показывает текст
  preload: true,
  variable: '--font-inter',
  fallback: ['system-ui', 'Arial'], // Fallback для быстрого рендеринга
  weight: ['400', '500', '600', '700'],
});

// Моноширинный шрифт для кода
export const jetBrainsMono = JetBrains_Mono({
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-jetbrains-mono',
  preload: false, // Загружается позже
  weight: ['400', '500', '600'],
});

/**
 * CSS переменные для шрифтов
 * Используй в globals.css:
 * 
 * :root {
 *   --font-inter: var(--font-inter);
 *   --font-jetbrains-mono: var(--font-jetbrains-mono);
 * }
 * 
 * body {
 *   font-family: var(--font-inter);
 * }
 * 
 * code, pre {
 *   font-family: var(--font-jetbrains-mono);
 * }
 */

export function getFontClassNames(): string {
  return `${inter.variable} ${jetBrainsMono.variable}`;
}
