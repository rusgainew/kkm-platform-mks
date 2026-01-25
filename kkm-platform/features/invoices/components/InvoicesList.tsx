'use client';

import React, { useState } from 'react';
import { FileText, Plus, Edit2, Trash2, Search, Loader2 } from 'lucide-react';
import Link from 'next/link';

interface Invoice {
  id: string;
  number: string;
  date: string;
  amount: number;
  status: 'draft' | 'sent' | 'paid' | 'overdue';
  company_name: string;
}

export default function InvoicesList() {
  const [invoices] = useState<Invoice[]>([]);
  const [searchQuery, setSearchQuery] = useState('');

  const getStatusColor = (status: Invoice['status']) => {
    switch (status) {
      case 'paid':
        return 'bg-green-900/20 text-green-400';
      case 'sent':
        return 'bg-blue-900/20 text-blue-400';
      case 'overdue':
        return 'bg-red-900/20 text-red-400';
      default:
        return 'bg-gray-800/20 text-gray-400';
    }
  };

  const getStatusLabel = (status: Invoice['status']) => {
    const labels = {
      draft: 'Черновик',
      sent: 'Отправлена',
      paid: 'Оплачена',
      overdue: 'Просрочена',
    };
    return labels[status];
  };

  const filteredInvoices = invoices.filter(i =>
    i.number.includes(searchQuery) ||
    i.company_name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
      <div className="p-6 border-b border-gray-800 flex items-center justify-between">
        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
          <FileText className="w-6 h-6 text-orange-400" />
          Счета
        </h2>
        <Link
          href="/invoices/create"
          className="flex items-center gap-2 px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700"
        >
          <Plus size={18} /> Создать счет
        </Link>
      </div>

      <div className="p-6 border-b border-gray-800">
        <input
          type="text"
          placeholder="Поиск по номеру или компании..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:border-orange-500 outline-none"
        />
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-800/50">
            <tr>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Номер</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Компания</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Дата</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Сумма</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Статус</th>
              <th className="px-6 py-3 text-center text-sm font-semibold text-gray-300">Действия</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800">
            {filteredInvoices.map((invoice) => (
              <tr key={invoice.id} className="hover:bg-gray-800/50 transition-colors">
                <td className="px-6 py-4 text-white font-medium">{invoice.number}</td>
                <td className="px-6 py-4 text-gray-300">{invoice.company_name}</td>
                <td className="px-6 py-4 text-gray-300">{new Date(invoice.date).toLocaleDateString('ru-RU')}</td>
                <td className="px-6 py-4 text-white font-medium">{invoice.amount.toLocaleString('ru-RU')} ₽</td>
                <td className="px-6 py-4">
                  <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(invoice.status)}`}>
                    {getStatusLabel(invoice.status)}
                  </span>
                </td>
                <td className="px-6 py-4 flex items-center justify-center gap-2">
                  <Link
                    href={`/invoices/${invoice.id}/edit`}
                    className="p-2 hover:bg-orange-600/20 text-orange-400 rounded transition-colors"
                  >
                    <Edit2 size={18} />
                  </Link>
                  <Link
                    href={`/invoices/${invoice.id}/delete`}
                    className="p-2 hover:bg-red-600/20 text-red-400 rounded transition-colors"
                  >
                    <Trash2 size={18} />
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {filteredInvoices.length === 0 && (
          <div className="p-8 text-center text-gray-400">
            {invoices.length === 0 ? 'Нет счетов' : 'Счета не найдены'}
          </div>
        )}
      </div>
    </div>
  );
}
