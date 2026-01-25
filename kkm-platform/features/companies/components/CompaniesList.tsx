'use client';

import React, { useState, useEffect } from 'react';
import { Building2, Plus, Edit2, Trash2, Search, Loader2 } from 'lucide-react';
import Link from 'next/link';
import { listCompanies } from '@/lib/api/companies';

interface Company {
  id: string;
  name: string;
  inn: string;
  email: string;
  phone: string;
  legal_address: string;
}

export default function CompaniesList() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadCompanies = async () => {
      try {
        setIsLoading(true);
        setError(null);
        const response = await listCompanies();
        const data = response?.data || [];
        setCompanies(data as any);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Ошибка загрузки');
      } finally {
        setIsLoading(false);
      }
    };

    loadCompanies();
  }, []);

  const filteredCompanies = companies.filter(c =>
    c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    c.inn.includes(searchQuery)
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="animate-spin text-blue-400" size={32} />
      </div>
    );
  }

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
      <div className="p-6 border-b border-gray-800 flex items-center justify-between">
        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
          <Building2 className="w-6 h-6 text-blue-400" />
          Управление компаниями
        </h2>
        <Link
          href="/companies/create"
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          <Plus size={18} /> Добавить компанию
        </Link>
      </div>

      {error && (
        <div className="p-4 bg-red-900/20 border-b border-red-800 text-red-300">
          {error}
        </div>
      )}

      <div className="p-6 border-b border-gray-800">
        <input
          type="text"
          placeholder="Поиск по названию или ИНН..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:border-blue-500 outline-none"
        />
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-800/50">
            <tr>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Название</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">ИНН</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Email</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Телефон</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Адрес</th>
              <th className="px-6 py-3 text-center text-sm font-semibold text-gray-300">Действия</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800">
            {filteredCompanies.map((company) => (
              <tr key={company.id} className="hover:bg-gray-800/50 transition-colors">
                <td className="px-6 py-4 text-white font-medium">{company.name}</td>
                <td className="px-6 py-4 text-gray-300">{company.inn}</td>
                <td className="px-6 py-4 text-gray-300">{company.email}</td>
                <td className="px-6 py-4 text-gray-300">{company.phone}</td>
                <td className="px-6 py-4 text-gray-300 text-sm">{company.legal_address}</td>
                <td className="px-6 py-4 flex items-center justify-center gap-2">
                  <Link
                    href={`/companies/${company.id}/edit`}
                    className="p-2 hover:bg-blue-600/20 text-blue-400 rounded transition-colors"
                  >
                    <Edit2 size={18} />
                  </Link>
                  <Link
                    href={`/companies/${company.id}/delete`}
                    className="p-2 hover:bg-red-600/20 text-red-400 rounded transition-colors"
                  >
                    <Trash2 size={18} />
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {filteredCompanies.length === 0 && (
          <div className="p-8 text-center text-gray-400">
            {companies.length === 0 ? 'Нет компаний' : 'Компании не найдены'}
          </div>
        )}
      </div>
    </div>
  );
}
