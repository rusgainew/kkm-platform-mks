'use client';

import React from 'react';

/**
 * Примеры стилизации для светлой и темной темы
 */

export function ThemeExampleCard() {
  return (
    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow dark:shadow-md">
      <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">
        📝 Пример карточки
      </h3>
      <p className="text-gray-600 dark:text-gray-300">
        Эта карточка поддерживает светлую и темную темы
      </p>
    </div>
  );
}

export function ThemeExampleButton() {
  return (
    <div className="flex gap-4 flex-wrap">
      <button className="px-4 py-2 bg-blue-500 hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700 text-white rounded-lg transition">
        Primary Button
      </button>
      <button className="px-4 py-2 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-900 dark:text-white rounded-lg transition">
        Secondary Button
      </button>
      <button className="px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white hover:bg-gray-50 dark:hover:bg-gray-800 rounded-lg transition">
        Outline Button
      </button>
    </div>
  );
}

export function ThemeExampleForm() {
  return (
    <form className="space-y-4 max-w-md">
      <div>
        <label className="block text-sm font-medium text-gray-900 dark:text-white mb-1">
          Email
        </label>
        <input
          type="email"
          className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="your@email.com"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-900 dark:text-white mb-1">
          Message
        </label>
        <textarea
          className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
          rows={4}
          placeholder="Your message..."
        />
      </div>
    </form>
  );
}

export function ThemeExampleGradient() {
  return (
    <div className="bg-gradient-to-r from-blue-500 to-purple-600 dark:from-blue-600 dark:to-purple-700 rounded-lg p-8 text-white">
      <h3 className="text-2xl font-bold mb-2">Gradient Example</h3>
      <p className="text-blue-100 dark:text-blue-200">
        Этот градиент адаптирован для обеих тем
      </p>
    </div>
  );
}

export function ThemeExampleTable() {
  return (
    <div className="overflow-x-auto border border-gray-200 dark:border-gray-700 rounded-lg">
      <table className="w-full">
        <thead className="bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
          <tr>
            <th className="px-4 py-2 text-left text-sm font-semibold text-gray-900 dark:text-white">
              Name
            </th>
            <th className="px-4 py-2 text-left text-sm font-semibold text-gray-900 dark:text-white">
              Status
            </th>
            <th className="px-4 py-2 text-left text-sm font-semibold text-gray-900 dark:text-white">
              Date
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
          {['Item 1', 'Item 2', 'Item 3'].map((item, i) => (
            <tr
              key={i}
              className="hover:bg-gray-50 dark:hover:bg-gray-800 transition"
            >
              <td className="px-4 py-2 text-sm text-gray-900 dark:text-white">
                {item}
              </td>
              <td className="px-4 py-2 text-sm">
                <span className="px-2 py-1 bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 rounded text-xs font-medium">
                  Active
                </span>
              </td>
              <td className="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
                Jan 16, 2026
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export function ThemeExampleAlert() {
  return (
    <div className="space-y-3">
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 text-blue-800 dark:text-blue-200 rounded-lg">
        ℹ️ Info alert with dark mode support
      </div>
      <div className="p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 text-green-800 dark:text-green-200 rounded-lg">
        ✓ Success alert with dark mode support
      </div>
      <div className="p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 text-yellow-800 dark:text-yellow-200 rounded-lg">
        ⚠ Warning alert with dark mode support
      </div>
      <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-800 dark:text-red-200 rounded-lg">
        ✕ Error alert with dark mode support
      </div>
    </div>
  );
}
