'use client';

import React, { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useApiToken } from '@/lib/hooks/useApiToken';
import { getInvoiceById, deleteInvoice } from '@/lib/api/invoices';
import { Loader2, AlertTriangle } from 'lucide-react';

interface Invoice {
  id: string;
  number: string;
  date: string;
  total: number;
}

export default function DeleteInvoicePage() {
  const params = useParams();
  const router = useRouter();
  const invoiceId = params.id as string;
  const [invoice, setInvoice] = useState<Invoice | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isConfirmed, setIsConfirmed] = useState(false);
  const token = useApiToken();

  useEffect(() => {
    const loadInvoice = async () => {
      if (!token) {
        setError('Требуется аутентификация');
        setIsLoading(false);
        return;
      }

      try {
        setIsLoading(true);
        const data = await getInvoiceById(invoiceId, token);
        setInvoice(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Ошибка загрузки');
      } finally {
        setIsLoading(false);
      }
    };

    loadInvoice();
  }, [invoiceId, token]);

  const handleDelete = async () => {
    if (!token || !isConfirmed) return;

    try {
      setIsDeleting(true);
      setError(null);
      await deleteInvoice(invoiceId, token);
      setTimeout(() => router.push('/invoices'), 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка удаления');
    } finally {
      setIsDeleting(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Loader2 className="animate-spin text-orange-400" size={32} />
      </div>
    );
  }

  if (error && !invoice) {
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
              <h2 className="text-xl font-bold text-red-400">Удалить счет</h2>
              <p className="text-red-300 text-sm">Это действие необратимо</p>
            </div>
          </div>

          {invoice && (
            <div className="space-y-4 mb-6">
              <div className="p-4 bg-gray-800 rounded-lg">
                <p className="text-gray-400 text-sm mb-1">Номер счета</p>
                <p className="text-white font-medium">{invoice.number}</p>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="p-4 bg-gray-800 rounded-lg">
                  <p className="text-gray-400 text-sm mb-1">Дата</p>
                  <p className="text-white font-medium">{new Date(invoice.date).toLocaleDateString('ru-RU')}</p>
                </div>
                <div className="p-4 bg-gray-800 rounded-lg">
                  <p className="text-gray-400 text-sm mb-1">Сумма</p>
                  <p className="text-white font-medium">{invoice.total.toLocaleString('ru-RU')} ₽</p>
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
              Я уверен, что хочу безвозвратно удалить этот счет
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
                'Удалить счет'
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
