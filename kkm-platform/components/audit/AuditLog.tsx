'use client';

import React, { useEffect, useState } from 'react';
import { useAuditStore } from '@/store/audit';
import { AuditLogEntry } from '@/lib/hooks/useAudit';

interface AuditLogProps {
  userId?: string;
  resourceType?: string;
  limit?: number;
}

export function AuditLog({ userId, resourceType, limit = 50 }: AuditLogProps) {
  const { logs, searchLogs } = useAuditStore();
  const [displayLogs, setDisplayLogs] = useState<AuditLogEntry[]>([]);
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState<'all' | 'success' | 'failure'>('all');
  const [sortBy, setSortBy] = useState<'newest' | 'oldest'>('newest');

  useEffect(() => {
    const updateDisplay = () => {
      let filtered = logs;

      // Фильтр по пользователю
      if (userId) {
        filtered = filtered.filter((log) => log.userId === userId);
      }

      // Фильтр по типу ресурса
      if (resourceType) {
        filtered = filtered.filter((log) => log.resource === resourceType);
      }

      // Поиск
      if (search) {
        filtered = searchLogs(search);
      }

      // Фильтр по статусу
      if (filter === 'success') {
        filtered = filtered.filter((log) => log.status === 'success');
      } else if (filter === 'failure') {
        filtered = filtered.filter((log) => log.status === 'failure');
      }

      // Сортировка
      if (sortBy === 'oldest') {
        filtered = filtered.sort((a, b) => a.timestamp - b.timestamp);
      } else {
        filtered = filtered.sort((a, b) => b.timestamp - a.timestamp);
      }

      // Ограничиваем количество
      setDisplayLogs(filtered.slice(0, limit));
    };

    updateDisplay();
  }, [logs, userId, resourceType, search, filter, sortBy, limit, searchLogs]);

  const getStatusIcon = (status: string) => {
    return status === 'success' ? '✅' : '❌';
  };

  const getActionIcon = (action: string) => {
    const icons: Record<string, string> = {
      create: '➕',
      read: '👁️',
      update: '✏️',
      delete: '🗑️',
      export: '📤',
      import: '📥',
      login: '🔓',
      logout: '🔐',
      permission_change: '🔑',
      role_change: '👤',
    };
    return icons[action] || '📋';
  };

  const getResourceColor = (resource: string) => {
    const colors: Record<string, string> = {
      invoice: 'bg-blue-100 text-blue-800',
      product: 'bg-green-100 text-green-800',
      company: 'bg-purple-100 text-purple-800',
      user: 'bg-orange-100 text-orange-800',
      settings: 'bg-gray-100 text-gray-800',
      report: 'bg-yellow-100 text-yellow-800',
      document: 'bg-red-100 text-red-800',
    };
    return colors[resource] || 'bg-gray-100 text-gray-800';
  };

  return (
    <div className="space-y-4">
      {/* Заголовок */}
      <div className="mb-4">
        <h3 className="text-xl font-bold text-gray-900">📋 Логирование действий</h3>
        <p className="text-sm text-gray-600 mt-1">
          Всего записей: <span className="font-semibold">{logs.length}</span>
        </p>
      </div>

      {/* Контролы фильтрации */}
      <div className="flex flex-col gap-3">
        <div className="flex gap-2 flex-wrap">
          <input
            type="text"
            placeholder="Поиск по имени, email, ресурсу..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-lg text-sm flex-1 min-w-64 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />

          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value as any)}
            className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">Все события</option>
            <option value="success">Успешные</option>
            <option value="failure">Ошибки</option>
          </select>

          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value as any)}
            className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="newest">Новые сверху</option>
            <option value="oldest">Старые сверху</option>
          </select>
        </div>
      </div>

      {/* Таблица логов */}
      {displayLogs.length > 0 ? (
        <div className="overflow-x-auto border border-gray-200 rounded-lg">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">
                  Действие
                </th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">
                  Пользователь
                </th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">
                  Ресурс
                </th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">
                  Результат
                </th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">
                  Время
                </th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">
                  Детали
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {displayLogs.map((log) => (
                <tr key={log.id} className="hover:bg-gray-50 transition">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <span>{getActionIcon(log.action)}</span>
                      <span className="font-medium capitalize">{log.action}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <div>
                      <p className="font-medium text-gray-900">{log.userName}</p>
                      <p className="text-xs text-gray-500">{log.userEmail}</p>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex flex-col gap-1">
                      <span
                        className={`inline-block px-2 py-1 text-xs font-medium rounded w-fit ${getResourceColor(
                          log.resource
                        )}`}
                      >
                        {log.resource}
                      </span>
                      <span className="text-xs text-gray-600">{log.resourceId}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <span>{getStatusIcon(log.status)}</span>
                      <span className="capitalize font-medium">
                        {log.status === 'success' ? 'OK' : 'Error'}
                      </span>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-gray-600 whitespace-nowrap">
                    <div className="text-xs">
                      <p>
                        {new Date(log.timestamp).toLocaleString('ru-RU', {
                          month: '2-digit',
                          day: '2-digit',
                          hour: '2-digit',
                          minute: '2-digit',
                          second: '2-digit',
                        })}
                      </p>
                      {log.duration && (
                        <p className="text-gray-500">{log.duration}ms</p>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    {log.error ? (
                      <details className="cursor-pointer">
                        <summary className="text-red-600 hover:underline font-medium">
                          Ошибка
                        </summary>
                        <p className="text-xs text-gray-600 mt-1 p-2 bg-red-50 rounded">
                          {log.error}
                        </p>
                      </details>
                    ) : log.newValue && Object.keys(log.newValue).length > 0 ? (
                      <details className="cursor-pointer">
                        <summary className="text-blue-600 hover:underline font-medium">
                          Изменения
                        </summary>
                        <div className="text-xs text-gray-600 mt-1 p-2 bg-blue-50 rounded max-h-40 overflow-y-auto">
                          <pre className="whitespace-pre-wrap break-words">
                            {JSON.stringify(log.newValue, null, 2)}
                          </pre>
                        </div>
                      </details>
                    ) : (
                      <span className="text-gray-400">-</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="text-center py-8 bg-gray-50 rounded-lg">
          <p className="text-gray-600">Записей не найдено</p>
        </div>
      )}
    </div>
  );
}

/**
 * Компонент для экспорта логов
 */
export function AuditLogExport() {
  const { exportLogs } = useAuditStore();

  const handleExport = (format: 'json' | 'csv') => {
    const data = exportLogs(format);
    const blob = new Blob([data], {
      type: format === 'json' ? 'application/json' : 'text/csv',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `audit-logs-${Date.now()}.${format}`;
    link.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="flex gap-2">
      <button
        onClick={() => handleExport('json')}
        className="px-4 py-2 bg-blue-500 text-white text-sm rounded-lg hover:bg-blue-600 transition"
      >
        📥 JSON
      </button>
      <button
        onClick={() => handleExport('csv')}
        className="px-4 py-2 bg-green-500 text-white text-sm rounded-lg hover:bg-green-600 transition"
      >
        📊 CSV
      </button>
    </div>
  );
}
