'use client';

import React, { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { getCompany, deleteCompany } from '@/lib/api/companies';
import { Loader2, AlertTriangle } from 'lucide-react';

interface Company {
  id: string;
  name: string;
  inn: string;
  email: string;
}

export default function DeleteCompanyPage() {
  const params = useParams();
  const router = useRouter();
  const companyId = params.id as string;
  const [company, setCompany] = useState<Company | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isConfirmed, setIsConfirmed] = useState(false);

  useEffect(() => {
    const loadCompany = async () => {
      if (!companyId) {
        setError('ID компании не найден');
        setIsLoading(false);
        return;
      }

      try {
        setIsLoading(true);
        const data = await getCompany(companyId);
        setCompany(data as any);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Ошибка загрузки');
      } finally {
        setIsLoading(false);
      }
    };

    loadCompany();
  }, [companyId]);

  const handleDelete = async () => {
    if (!isConfirmed) return;

    try {
      setIsDeleting(true);
      setError(null);
      await deleteCompany(companyId);
      setTimeout(() => router.push('/companies'), 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка удаления');
    } finally {
      setIsDeleting(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Loader2 className="animate-spin text-blue-400" size={32} />
      </div>
    );
  }

  if (error && !company) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-red-400">{error}</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <div className="bg-gray-900 rounded-lg border border-gray-800 p-6">
          <div className="flex items-center gap-3 mb-6 p-4 bg-red-900/20 border border-red-800 rounded-lg">
            <AlertTriangle className="text-red-400" size={24} />
            <div>
              <h2 className="text-xl font-bold text-red-400">Удалить компанию</h2>
              <p className="text-red-300 text-sm">Это действие необратимо</p>
            </div>
          </div>

          {company && (
            <div className="space-y-4 mb-6">
              <div className="p-4 bg-gray-800 rounded-lg">
                <p className="text-gray-400 text-sm mb-1">Название компании</p>
                <p className="text-white font-medium">{company.name}</p>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="p-4 bg-gray-800 rounded-lg">
                  <p className="text-gray-400 text-sm mb-1">ИНН</p>
                  <p className="text-white font-medium">{company.inn}</p>
                </div>
                <div className="p-4 bg-gray-800 rounded-lg">
                  <p className="text-gray-400 text-sm mb-1">Email</p>
                  <p className="text-white font-medium">{company.email}</p>
                </div>
              </div>
            </div>
          )}

          {error && (
            <div className="mb-4 p-3 bg-red-900/20 border border-red-800 text-red-300 rounded-lg text-sm">
              {error}
            </div>
          )}

          <label className="flex items-center gap-3 mb-6 p-4 bg-gray-800 rounded-lg">
            <input
              type="checkbox"
              checked={isConfirmed}
              onChange={(e) => setIsConfirmed(e.target.checked)}
              className="w-5 h-5"
            />
            <span className="text-gray-300">
              Я уверен, что хочу безвозвратно удалить эту компанию
            </span>
          </label>

          <div className="flex gap-3">
            <button
              onClick={handleDelete}
              disabled={!isConfirmed || isDeleting}
              className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg font-medium hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
            >
              {isDeleting ? (
                <>
                  <Loader2 className="animate-spin" size={18} />
                  Удаление...
                </>
              ) : (
                'Удалить компанию'
              )}
            </button>
            <button
              onClick={() => router.back()}
              className="flex-1 px-4 py-2 bg-gray-800 text-white rounded-lg font-medium hover:bg-gray-700 transition-colors"
            >
              Отмена
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
