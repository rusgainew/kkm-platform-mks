'use client';

import React, { useTransition } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { Locale, locales, localeInfo } from '@/i18n/config';
import { useTranslationsSync } from '@/lib/hooks/useTranslations';

interface LanguageSwitcherProps {
  currentLocale: Locale;
  className?: string;
  showFlags?: boolean;
  showText?: boolean;
}

export function LanguageSwitcher({
  currentLocale,
  className = '',
  showFlags = true,
  showText = true,
}: LanguageSwitcherProps) {
  const router = useRouter();
  const pathname = usePathname();
  const [isPending, startTransition] = useTransition();
  const t = useTranslationsSync(currentLocale);

  const handleLanguageChange = (newLocale: Locale) => {
    if (newLocale === currentLocale) return;

    startTransition(() => {
      // Заменяем текущий locale на новый в пути
      const newPathname = pathname.replace(`/${currentLocale}`, `/${newLocale}`);
      
      // Сохраняем выбор языка в cookies
      document.cookie = `NEXT_LOCALE=${newLocale}; path=/; max-age=31536000`;
      
      router.push(newPathname);
    });
  };

  return (
    <div className={`flex gap-2 ${className}`}>
      {locales.map((locale) => {
        const info = localeInfo[locale];
        const isActive = locale === currentLocale;

        return (
          <button
            key={locale}
            onClick={() => handleLanguageChange(locale)}
            disabled={isPending}
            title={info.nativeName}
            className={`flex items-center gap-1 px-3 py-2 rounded-lg transition ${
              isActive
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300 disabled:opacity-50'
            }`}
          >
            {showFlags && <span>{info.flag}</span>}
            {showText && <span className="text-sm font-medium">{locale.toUpperCase()}</span>}
          </button>
        );
      })}
    </div>
  );
}

/**
 * Компонент выпадающего меню языков
 */
export function LanguageSwitcherDropdown({
  currentLocale,
  className = '',
}: {
  currentLocale: Locale;
  className?: string;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const [isPending, startTransition] = useTransition();
  const [isOpen, setIsOpen] = React.useState(false);
  const t = useTranslationsSync(currentLocale);

  const handleLanguageChange = (newLocale: Locale) => {
    if (newLocale === currentLocale) return;

    startTransition(() => {
      const newPathname = pathname.replace(`/${currentLocale}`, `/${newLocale}`);
      document.cookie = `NEXT_LOCALE=${newLocale}; path=/; max-age=31536000`;
      router.push(newPathname);
      setIsOpen(false);
    });
  };

  const currentInfo = localeInfo[currentLocale];

  return (
    <div className={`relative ${className}`}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        disabled={isPending}
        className="flex items-center gap-2 px-3 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg transition disabled:opacity-50"
      >
        <span>{currentInfo.flag}</span>
        <span className="text-sm font-medium">{currentInfo.nativeName}</span>
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
        <div className="absolute top-full right-0 mt-2 bg-white border border-gray-200 rounded-lg shadow-lg z-50">
          {locales.map((locale) => {
            const info = localeInfo[locale];
            const isActive = locale === currentLocale;

            return (
              <button
                key={locale}
                onClick={() => handleLanguageChange(locale)}
                disabled={isPending || isActive}
                className={`w-full flex items-center gap-3 px-4 py-2 text-left transition ${
                  isActive
                    ? 'bg-blue-50 text-blue-600 font-medium'
                    : 'hover:bg-gray-50 disabled:opacity-50'
                } ${locale !== locales[locales.length - 1] ? 'border-b border-gray-100' : ''}`}
              >
                <span>{info.flag}</span>
                <span className="flex-1">{info.nativeName}</span>
                {isActive && <span>✓</span>}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
