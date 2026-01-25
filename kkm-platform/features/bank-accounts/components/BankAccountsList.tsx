'use client';

import React, { useState } from 'react';
import { CreditCard, Plus, Edit2, Trash2, Search } from 'lucide-react';
import Link from 'next/link';

interface BankAccount {
  id: string;
  bank_name: string;
  account_number: string;
  bik: string;
  correspondent_account: string;
  company_name: string;
  is_primary: boolean;
}

export default function BankAccountsList() {
  const [accounts] = useState<BankAccount[]>([]);
  const [searchQuery, setSearchQuery] = useState('');

  const filteredAccounts = accounts.filter(a =>
    a.bank_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    a.account_number.includes(searchQuery) ||
    a.bik.includes(searchQuery)
  );

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
      <div className="p-6 border-b border-gray-800 flex items-center justify-between">
        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
          <CreditCard className="w-6 h-6 text-green-400" />
          Банковские счета
        </h2>
        <Link
          href="/bank-accounts/create"
          className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
        >
          <Plus size={18} /> Добавить счет
        </Link>
      </div>

      <div className="p-6 border-b border-gray-800">
        <input
          type="text"
          placeholder="Поиск по банку, номеру или БИК..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:border-green-500 outline-none"
        />
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-800/50">
            <tr>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Банк</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Счет</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">БИК</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Корр. счет</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Статус</th>
              <th className="px-6 py-3 text-center text-sm font-semibold text-gray-300">Действия</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800">
            {filteredAccounts.map((account) => (
              <tr key={account.id} className="hover:bg-gray-800/50 transition-colors">
                <td className="px-6 py-4 text-white font-medium">{account.bank_name}</td>
                <td className="px-6 py-4 text-gray-300 font-mono">{account.account_number}</td>
                <td className="px-6 py-4 text-gray-300 font-mono">{account.bik}</td>
                <td className="px-6 py-4 text-gray-300 font-mono text-sm">{account.correspondent_account}</td>
                <td className="px-6 py-4">
                  {account.is_primary && (
                    <span className="px-3 py-1 rounded-full text-sm font-medium bg-green-900/20 text-green-400">
                      Основной
                    </span>
                  )}
                </td>
                <td className="px-6 py-4 flex items-center justify-center gap-2">
                  <Link
                    href={`/bank-accounts/${account.id}/edit`}
                    className="p-2 hover:bg-green-600/20 text-green-400 rounded transition-colors"
                  >
                    <Edit2 size={18} />
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {filteredAccounts.length === 0 && (
          <div className="p-8 text-center text-gray-400">
            {accounts.length === 0 ? 'Нет банковских счетов' : 'Счета не найдены'}
          </div>
        )}
      </div>
    </div>
  );
}
