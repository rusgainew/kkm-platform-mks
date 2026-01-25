'use client';

import React, { useState } from 'react';
import { Loader2 } from 'lucide-react';

interface BankAccountFormProps {
  initialData?: any;
  onSubmit?: (data: any) => Promise<void>;
}

export default function BankAccountForm({ initialData, onSubmit }: BankAccountFormProps) {
  const [formData, setFormData] = useState({
    bank_name: initialData?.bank_name || '',
    bik: initialData?.bik || '',
    account_number: initialData?.account_number || '',
    correspondent_account: initialData?.correspondent_account || '',
    company_id: initialData?.company_id || '',
    is_primary: initialData?.is_primary || false,
  });

  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.bank_name.trim()) newErrors.bank_name = 'Название банка обязательно';
    if (!formData.bik.trim() || formData.bik.length !== 9) {
      newErrors.bik = 'БИК должен содержать 9 цифр';
    }
    if (!formData.account_number.trim() || formData.account_number.length !== 20) {
      newErrors.account_number = 'Номер счета должен содержать 20 цифр';
    }
    if (!formData.correspondent_account.trim() || formData.correspondent_account.length !== 20) {
      newErrors.correspondent_account = 'Коррсчет должен содержать 20 цифр';
    }

    return newErrors;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const newErrors = validateForm();

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setIsLoading(true);
    setErrors({});

    try {
      if (onSubmit) {
        await onSubmit(formData);
      }
    } catch (err) {
      setErrors({
        submit: err instanceof Error ? err.message : 'Ошибка отправки',
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-gray-900 rounded-lg border border-gray-800 p-6 space-y-6">
      <h2 className="text-2xl font-bold text-white">
        {initialData?.id ? 'Редактировать банковский счет' : 'Добавить банковский счет'}
      </h2>

      {errors.submit && (
        <div className="p-4 bg-red-900/20 border border-red-800 text-red-300 rounded-lg">
          {errors.submit}
        </div>
      )}

      {/* Банковские данные */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-300">Банковские реквизиты</h3>
        
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Название банка *
          </label>
          <input
            type="text"
            value={formData.bank_name}
            onChange={(e) => setFormData({ ...formData, bank_name: e.target.value })}
            className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-green-500 outline-none"
            placeholder="ПАО Сбербанк"
          />
          {errors.bank_name && <p className="text-red-400 text-sm mt-1">{errors.bank_name}</p>}
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              БИК *
            </label>
            <input
              type="text"
              value={formData.bik}
              onChange={(e) => setFormData({ ...formData, bik: e.target.value.replace(/\D/g, '').slice(0, 9) })}
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-green-500 outline-none font-mono"
              placeholder="044525225"
              maxLength={9}
            />
            {errors.bik && <p className="text-red-400 text-sm mt-1">{errors.bik}</p>}
            <p className="text-gray-400 text-xs mt-1">Требуется 9 цифр</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Коррсчет *
            </label>
            <input
              type="text"
              value={formData.correspondent_account}
              onChange={(e) => setFormData({ 
                ...formData, 
                correspondent_account: e.target.value.replace(/\D/g, '').slice(0, 20) 
              })}
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-green-500 outline-none font-mono"
              placeholder="30101810400000000225"
              maxLength={20}
            />
            {errors.correspondent_account && <p className="text-red-400 text-sm mt-1">{errors.correspondent_account}</p>}
            <p className="text-gray-400 text-xs mt-1">Требуется 20 цифр</p>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Расчетный счет *
          </label>
          <input
            type="text"
            value={formData.account_number}
            onChange={(e) => setFormData({ 
              ...formData, 
              account_number: e.target.value.replace(/\D/g, '').slice(0, 20) 
            })}
            className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-green-500 outline-none font-mono"
            placeholder="40702810800003000002"
            maxLength={20}
          />
          {errors.account_number && <p className="text-red-400 text-sm mt-1">{errors.account_number}</p>}
          <p className="text-gray-400 text-xs mt-1">Требуется 20 цифр</p>
        </div>
      </div>

      {/* Компания и статус */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-300">Настройки</h3>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Компания
          </label>
          <select
            value={formData.company_id}
            onChange={(e) => setFormData({ ...formData, company_id: e.target.value })}
            className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-green-500 outline-none"
          >
            <option value="">Выберите компанию</option>
          </select>
        </div>

        <label className="flex items-center gap-3 p-4 bg-gray-800 rounded-lg cursor-pointer hover:bg-gray-700 transition-colors">
          <input
            type="checkbox"
            checked={formData.is_primary}
            onChange={(e) => setFormData({ ...formData, is_primary: e.target.checked })}
            className="w-5 h-5"
          />
          <span className="text-gray-300 font-medium">
            Основной счет для платежей по умолчанию
          </span>
        </label>
      </div>

      {/* Информация */}
      <div className="p-4 bg-blue-900/20 border border-blue-800 rounded-lg">
        <p className="text-blue-300 text-sm">
          <strong>Важно:</strong> Все реквизиты должны быть заполнены корректно. Убедитесь в правильности БИК и номера счета перед сохранением.
        </p>
      </div>

      {/* Кнопки */}
      <div className="flex gap-3">
        <button
          type="submit"
          disabled={isLoading}
          className="flex-1 px-4 py-2 bg-green-600 text-white rounded-lg font-medium hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
        >
          {isLoading ? (
            <>
              <Loader2 className="animate-spin" size={18} />
              Сохранение...
            </>
          ) : (
            'Сохранить счет'
          )}
        </button>
      </div>
    </form>
  );
}
