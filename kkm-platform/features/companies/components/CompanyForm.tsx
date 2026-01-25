'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeft, Loader2, AlertCircle, CheckCircle, Building2 } from 'lucide-react';

interface CompanyCreateFormProps {
  initialData?: {
    name: string;
    inn: string;
    phone: string;
    email: string;
    legal_address: string;
  };
  onSubmit: (data: any) => Promise<void>;
}

export default function CompanyCreateForm({ initialData, onSubmit }: CompanyCreateFormProps) {
  const router = useRouter();
  const [formData, setFormData] = useState(initialData || {
    name: '',
    inn: '',
    kpp: '',
    ogrn: '',
    legal_address: '',
    actual_address: '',
    phone: '',
    email: '',
    director_name: '',
    accountant_name: '',
    website: '',
  });

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name?.trim()) newErrors.name = 'Название компании обязательно';
    if (!formData.inn?.trim()) newErrors.inn = 'ИНН обязателен';
    if (!formData.phone?.trim()) newErrors.phone = 'Телефон обязателен';
    if (!formData.email?.trim()) newErrors.email = 'Email обязателен';
    if (!formData.legal_address?.trim()) newErrors.legal_address = 'Юридический адрес обязателен';

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    setIsLoading(true);
    setError(null);
    setSuccessMessage('');

    try {
      await onSubmit(formData);
      setSuccessMessage('Компания успешно сохранена');
      setTimeout(() => router.push('/companies'), 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при сохранении');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-linear-to-b from-gray-900 to-black">
      <div className="bg-gray-800/50 border-b border-gray-700 sticky top-0 z-10">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <button
            onClick={() => router.back()}
            className="flex items-center gap-2 text-gray-400 hover:text-white mb-4"
          >
            <ArrowLeft size={20} /> Вернуться
          </button>
          <h1 className="text-3xl font-bold text-white flex items-center gap-2">
            <Building2 size={32} /> Управление компаниями
          </h1>
        </div>
      </div>

      <div className="max-w-4xl mx-auto px-4 py-8">
        <form onSubmit={handleSubmit} className="space-y-8">
          {successMessage && (
            <div className="bg-emerald-900/20 border border-emerald-800 rounded-lg p-4 flex items-start gap-3">
              <CheckCircle size={20} className="text-emerald-400 shrink-0 mt-0.5" />
              <p className="text-emerald-300">{successMessage}</p>
            </div>
          )}

          {error && (
            <div className="bg-red-900/20 border border-red-800 rounded-lg p-4 flex items-start gap-3">
              <AlertCircle size={20} className="text-red-400 shrink-0 mt-0.5" />
              <p className="text-red-300">{error}</p>
            </div>
          )}

          <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6 space-y-4">
            <h3 className="text-lg font-semibold text-white mb-4">Основные данные</h3>

            {['name', 'inn', 'kpp', 'ogrn', 'legal_address', 'actual_address', 'phone', 'email', 'director_name', 'accountant_name', 'website'].map((field) => (
              <div key={field}>
                <label className="block text-sm font-medium text-gray-300 mb-2 capitalize">
                  {field.replace(/_/g, ' ')}
                </label>
                <input
                  type="text"
                  value={(formData as any)[field] || ''}
                  onChange={(e) => setFormData({ ...formData, [field]: e.target.value })}
                  disabled={isLoading}
                  className={`w-full px-4 py-2 rounded-lg bg-gray-900/50 border outline-none ${
                    errors[field]
                      ? 'border-red-600 focus:border-red-500'
                      : 'border-gray-600 focus:border-blue-500'
                  } text-white disabled:opacity-50`}
                />
                {errors[field] && <p className="text-red-400 text-sm mt-1">{errors[field]}</p>}
              </div>
            ))}
          </div>

          <div className="flex gap-4">
            <button
              type="submit"
              disabled={isLoading}
              className={`flex-1 px-6 py-3 rounded-lg font-medium flex items-center justify-center gap-2 ${
                isLoading
                  ? 'bg-gray-700 text-gray-400 cursor-not-allowed'
                  : 'bg-blue-600 text-white hover:bg-blue-700'
              }`}
            >
              {isLoading ? (
                <>
                  <Loader2 size={18} className="animate-spin" />
                  Сохранение...
                </>
              ) : (
                'Сохранить'
              )}
            </button>
            <button
              type="button"
              onClick={() => router.push('/companies')}
              disabled={isLoading}
              className="px-6 py-3 rounded-lg bg-gray-700 text-gray-200 hover:bg-gray-600"
            >
              Отмена
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
