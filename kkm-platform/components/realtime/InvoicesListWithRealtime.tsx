'use client';

import React, { useEffect, useState } from 'react';
import { useRealtimeInvoices, RealtimeInvoice } from '@/store/realtime';
import toast from 'react-hot-toast';

export function InvoicesListWithRealtime() {
  const invoices = useRealtimeInvoices((state) => state.getAllInvoices());
  const [displayInvoices, setDisplayInvoices] = useState<RealtimeInvoice[]>([]);
  const [filter, setFilter] = useState<string>('all');
  const [search, setSearch] = useState<string>('');

  useEffect(() => {
    let filtered = invoices;

    // Фильтр по статусу
    if (filter !== 'all') {
      filtered = filtered.filter((inv) => inv.status === filter);
    }

    // Поиск по номеру
    if (search) {
      filtered = filtered.filter((inv) =>
        inv.number.toLowerCase().includes(search.toLowerCase())
      );
    }

    // Сортировка по дате (новые сверху)
    filtered = filtered.sort(
      (a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
    );

    setDisplayInvoices(filtered);
  }, [invoices, filter, search]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'paid':
        return 'bg-green-100 text-green-800';
      case 'issued':
        return 'bg-blue-100 text-blue-800';
      case 'draft':
        return 'bg-gray-100 text-gray-800';
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'paid':
        return '✅ Оплачен';
      case 'issued':
        return '📤 Отправлен';
      case 'draft':
        return '📝 Черновик';
      case 'cancelled':
        return '❌ Отменен';
      default:
        return status;
    }
  };

  return (
    <div className="space-y-4">
      {/* Заголовок */}
      <div>
        <h2 className="text-2xl font-bold text-gray-900">
          📄 Счета-фактуры
          {displayInvoices.length > 0 && (
            <span className="text-sm font-normal text-gray-500 ml-2">
              ({displayInvoices.length})
            </span>
          )}
        </h2>
        <p className="text-sm text-gray-600">
          {invoices.length > 0 ? (
            <>
              <span className="inline-flex items-center gap-1">
                <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                Реал-тайм синхронизация включена
              </span>
            </>
          ) : (
            'Загрузка счетов...'
          )}
        </p>
      </div>

      {/* Поиск и фильтры */}
      <div className="flex gap-2 flex-wrap">
        <input
          type="text"
          placeholder="Поиск по номеру счета..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />

        <div className="flex gap-2">
          {['all', 'draft', 'issued', 'paid', 'cancelled'].map((status) => (
            <button
              key={status}
              onClick={() => setFilter(status)}
              className={`px-3 py-2 text-sm rounded-lg transition ${
                filter === status
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              {status === 'all'
                ? 'Все'
                : status === 'draft'
                  ? 'Черновики'
                  : status === 'issued'
                    ? 'Отправлены'
                    : status === 'paid'
                      ? 'Оплачены'
                      : 'Отменены'}
            </button>
          ))}
        </div>
      </div>

      {/* Таблица счетов */}
      {displayInvoices.length > 0 ? (
        <div className="overflow-x-auto border border-gray-200 rounded-lg">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                  Номер
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                  Компания
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                  Сумма
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                  Статус
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                  Обновлено
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                  Действия
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {displayInvoices.map((invoice) => (
                <tr key={invoice.id} className="hover:bg-gray-50 transition">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">
                    #{invoice.number}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {invoice.company || '-'}
                  </td>
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">
                    ₽{invoice.amount.toLocaleString('ru-RU')}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(invoice.status)}`}>
                      {getStatusLabel(invoice.status)}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {new Date(invoice.updatedAt).toLocaleString('ru-RU', {
                      year: 'numeric',
                      month: '2-digit',
                      day: '2-digit',
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <button className="text-blue-600 hover:text-blue-700 font-medium">
                      Просмотр
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="text-center py-12 bg-gray-50 rounded-lg">
          <p className="text-gray-600">Счетов не найдено</p>
        </div>
      )}
    </div>
  );
}
