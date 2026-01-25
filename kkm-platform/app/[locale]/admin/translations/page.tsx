'use client';

import React, { useEffect, useState } from 'react';
import { Locale } from '@/i18n/config';
import { useTranslationsSync, initializeTranslations } from '@/lib/hooks/useTranslations';
import { LanguageSwitcher } from '@/components/i18n/LanguageSwitcher';

interface TranslationEntry {
  key: string;
  value: string;
  locale: Locale;
}

/**
 * Admin страница для управления переводами
 */
export default function TranslationsPage({ params }: { params: { locale: Locale } }) {
  const locale = params.locale;
  const [translations, setTranslations] = useState<TranslationEntry[]>([]);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(true);
  const [editKey, setEditKey] = useState<string | null>(null);
  const [editValue, setEditValue] = useState('');
  const tSync = useTranslationsSync(locale);

  useEffect(() => {
    const loadTranslations = async () => {
      try {
        await initializeTranslations(locale);
        const messages = await import(`@/messages/${locale}.json`);

        // Раскрываем вложенный объект в плоский список
        const flattenMessages = (obj: Record<string, unknown>, prefix = ''): TranslationEntry[] => {
          const result: TranslationEntry[] = [];

          for (const [key, value] of Object.entries(obj)) {
            const fullKey = prefix ? `${prefix}.${key}` : key;

            if (typeof value === 'string') {
              result.push({
                key: fullKey,
                value,
                locale,
              });
            } else if (typeof value === 'object' && value !== null) {
              result.push(...flattenMessages(value as Record<string, unknown>, fullKey));
            }
          }

          return result;
        };

        const flattened = flattenMessages(messages.default);
        setTranslations(flattened);
      } catch (error) {
        console.error('Failed to load translations:', error);
      } finally {
        setLoading(false);
      }
    };

    loadTranslations();
  }, [locale]);

  const filtered = translations.filter(
    (t) =>
      t.key.toLowerCase().includes(search.toLowerCase()) ||
      t.value.toLowerCase().includes(search.toLowerCase())
  );

  const handleEdit = (key: string, value: string) => {
    setEditKey(key);
    setEditValue(value);
  };

  const handleSave = () => {
    // В реальном приложении здесь была бы интеграция с backend
    const updated = translations.map((t) =>
      t.key === editKey ? { ...t, value: editValue } : t
    );
    setTranslations(updated);
    setEditKey(null);
  };

  const handleExport = () => {
    const obj = {};
    for (const t of translations) {
      const keys = t.key.split('.');
      let current = obj as Record<string, unknown>;

      for (let i = 0; i < keys.length - 1; i++) {
        if (!current[keys[i]]) {
          current[keys[i]] = {};
        }
        current = current[keys[i]] as Record<string, unknown>;
      }

      current[keys[keys.length - 1]] = t.value;
    }

    const jsonStr = JSON.stringify(obj, null, 2);
    const blob = new Blob([jsonStr], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `translations-${locale}.json`;
    link.click();
    URL.revokeObjectURL(url);
  };

  if (loading) {
    return (
      <div className="p-8">
        <p className="text-gray-600">{tSync('common.loading')}</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-6xl mx-auto">
        {/* Заголовок */}
        <div className="mb-8">
          <div className="flex justify-between items-center mb-4">
            <h1 className="text-4xl font-bold text-gray-900">
              🌍 {tSync('admin.translations')}
            </h1>
            <LanguageSwitcher currentLocale={locale} />
          </div>
          <p className="text-gray-600">
            {tSync('admin.manageTranslations')} ({filtered.length} {tSync('common.total')})
          </p>
        </div>

        {/* Контролы */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <div className="flex flex-col sm:flex-row gap-4">
            <input
              type="text"
              placeholder={tSync('common.search')}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={handleExport}
              className="px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition"
            >
              📥 {tSync('common.export')}
            </button>
          </div>
        </div>

        {/* Таблица переводов */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                  Key
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                  {tSync('common.name')}
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                  {tSync('common.actions')}
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {filtered.map((entry) => (
                <tr key={entry.key} className="hover:bg-gray-50 transition">
                  <td className="px-6 py-4 text-sm font-mono text-gray-600">
                    {entry.key}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-900">
                    {editKey === entry.key ? (
                      <input
                        type="text"
                        value={editValue}
                        onChange={(e) => setEditValue(e.target.value)}
                        className="w-full px-3 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                    ) : (
                      entry.value
                    )}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    {editKey === entry.key ? (
                      <div className="flex gap-2">
                        <button
                          onClick={handleSave}
                          className="px-3 py-1 bg-green-500 text-white rounded hover:bg-green-600"
                        >
                          {tSync('common.save')}
                        </button>
                        <button
                          onClick={() => setEditKey(null)}
                          className="px-3 py-1 bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
                        >
                          {tSync('common.cancel')}
                        </button>
                      </div>
                    ) : (
                      <button
                        onClick={() => handleEdit(entry.key, entry.value)}
                        className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
                      >
                        ✏️ {tSync('common.edit')}
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Статистика */}
        <div className="grid grid-cols-3 gap-4 mt-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="text-3xl font-bold text-blue-600">{translations.length}</div>
            <p className="text-gray-600 mt-2">Всего строк</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="text-3xl font-bold text-green-600">{filtered.length}</div>
            <p className="text-gray-600 mt-2">Найдено</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="text-3xl font-bold text-purple-600">
              {translations.reduce((sum, t) => sum + t.value.length, 0)}
            </div>
            <p className="text-gray-600 mt-2">Всего символов</p>
          </div>
        </div>
      </div>
    </div>
  );
}
