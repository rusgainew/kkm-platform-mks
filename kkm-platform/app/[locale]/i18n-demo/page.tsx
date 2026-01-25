'use client';

import React from 'react';
import { Locale } from '@/i18n/config';
import { useTranslationsSync } from '@/lib/hooks/useTranslations';
import { LanguageSwitcher, LanguageSwitcherDropdown } from '@/components/i18n/LanguageSwitcher';

/**
 * Demo страница для тестирования i18n
 */
export default function I18nDemoPage({ params }: { params: { locale: Locale } }) {
  const locale = params.locale;
  const t = useTranslationsSync(locale);

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        {/* Заголовок */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            🌍 {t('settings.language')}
          </h1>
          <p className="text-gray-600">
            Демонстрация интернационализации приложения
          </p>
        </div>

        {/* Переключатели языка */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-xl font-bold text-gray-900 mb-4">
            {t('settings.currentLanguage')}: <span className="text-blue-600">{locale.toUpperCase()}</span>
          </h2>

          <div className="space-y-4">
            <div>
              <p className="text-gray-700 font-medium mb-2">Horizontal Switcher:</p>
              <LanguageSwitcher currentLocale={locale} showFlags showText />
            </div>

            <div>
              <p className="text-gray-700 font-medium mb-2">Dropdown Switcher:</p>
              <LanguageSwitcherDropdown currentLocale={locale} />
            </div>
          </div>
        </div>

        {/* Примеры переводов */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          {/* Навигация */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-4">
              📱 {t('settings.language')} - Navigation
            </h3>
            <ul className="space-y-2 text-sm text-gray-700">
              <li>• {t('nav.dashboard')}</li>
              <li>• {t('nav.invoices')}</li>
              <li>• {t('nav.catalog')}</li>
              <li>• {t('nav.companies')}</li>
              <li>• {t('nav.users')}</li>
              <li>• {t('nav.admin')}</li>
            </ul>
          </div>

          {/* Общие переводы */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-4">
              ✅ Common Translations
            </h3>
            <ul className="space-y-2 text-sm text-gray-700">
              <li>• {t('common.loading')}</li>
              <li>• {t('common.save')}</li>
              <li>• {t('common.delete')}</li>
              <li>• {t('common.cancel')}</li>
              <li>• {t('common.confirm')}</li>
              <li>• {t('common.search')}</li>
            </ul>
          </div>

          {/* Счета */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-4">
              📄 {t('invoices.title')}
            </h3>
            <ul className="space-y-2 text-sm text-gray-700">
              <li>• {t('invoices.draft')}</li>
              <li>• {t('invoices.issued')}</li>
              <li>• {t('invoices.paid')}</li>
              <li>• {t('invoices.cancelled')}</li>
              <li>• {t('invoices.markAsPaid')}</li>
              <li>• {t('invoices.downloadPDF')}</li>
            </ul>
          </div>

          {/* Каталог */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-4">
              📦 {t('catalog.title')}
            </h3>
            <ul className="space-y-2 text-sm text-gray-700">
              <li>• {t('catalog.inStock')}</li>
              <li>• {t('catalog.lowStock')}</li>
              <li>• {t('catalog.outOfStock')}</li>
              <li>• {t('catalog.filterByCategory')}</li>
              <li>• {t('catalog.sortByPrice')}</li>
              <li>• {t('catalog.searchByName')}</li>
            </ul>
          </div>
        </div>

        {/* Информация о текущем языке */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
          <h3 className="text-lg font-bold text-blue-900 mb-4">ℹ️ Информация о локализации</h3>
          <ul className="text-sm text-blue-800 space-y-2">
            <li>✅ Текущий язык: <span className="font-semibold">{locale.toUpperCase()}</span></li>
            <li>✅ Поддерживаемые языки: <span className="font-semibold">Русский (RU), English (EN)</span></li>
            <li>✅ Язык определяется из: URL, cookies, Accept-Language header</li>
            <li>✅ Переводы кэшируются для быстродействия</li>
            <li>✅ Выбор языка сохраняется в cookies на 1 год</li>
            <li>✅ Управление переводами: <a href={`/${locale}/admin/translations`} className="text-blue-600 font-semibold hover:underline">Admin Panel</a></li>
          </ul>
        </div>

        {/* Примеры использования */}
        <div className="mt-8 bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-bold text-gray-900 mb-4">💻 Примеры кода</h3>
          
          <div className="space-y-4 text-sm">
            <div>
              <p className="font-mono bg-gray-100 p-3 rounded text-gray-800">
                &lt;LanguageSwitcher currentLocale=&quot;ru&quot; /&gt;
              </p>
              <p className="text-gray-600 mt-1">Переключатель языков</p>
            </div>

            <div>
              <p className="font-mono bg-gray-100 p-3 rounded text-gray-800">
                const t = useTranslationsSync(locale);<br />
                t(&apos;nav.dashboard&apos;)
              </p>
              <p className="text-gray-600 mt-1">Получение перевода</p>
            </div>

            <div>
              <p className="font-mono bg-gray-100 p-3 rounded text-gray-800">
                t(&apos;validation.minLength&apos;, {'{'} min: 8 {'}'})
              </p>
              <p className="text-gray-600 mt-1">Интерполяция переменных</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
