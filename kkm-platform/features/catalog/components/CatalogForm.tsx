'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeft, Loader2, AlertCircle, CheckCircle, Package, DollarSign } from 'lucide-react';
import type { CatalogItem, CreateCatalogItemRequest, UpdateCatalogItemRequest } from '@/types/entities';
import { isValidTnvedCode, UNITS_OF_MEASURE, CURRENCIES } from '@/types/entities';

interface CatalogFormProps {
  initialData?: CatalogItem;
  onSubmit: (data: CreateCatalogItemRequest | UpdateCatalogItemRequest) => Promise<CatalogItem>;
  mode?: 'create' | 'edit';
}

interface FormData {
  name: string;
  number: string;
  description?: string;
  tnved_code: string;
  price: string; // Строка для контроля ввода
  currency: string;
  unit: string;
}

interface ValidationErrors {
  name?: string;
  number?: string;
  tnved_code?: string;
  price?: string;
  currency?: string;
  unit?: string;
}

export default function CatalogForm({ initialData, onSubmit, mode = 'create' }: CatalogFormProps) {
  const router = useRouter();
  const [formData, setFormData] = useState<FormData>({
    name: initialData?.name || '',
    number: initialData?.number || '',
    description: initialData?.description || '',
    tnved_code: initialData?.tnved_code || '',
    price: initialData?.price?.toString() || '',
    currency: initialData?.currency || 'KGS',
    unit: initialData?.unit || '796', // По умолчанию "шт"
  });

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState('');
  const [errors, setErrors] = useState<ValidationErrors>({});

  /**
   * Форматирование ТНВЭД кода (XXXX XX XXXX)
   */
  const formatTnvedCode = (value: string): string => {
    const cleaned = value.replace(/\s/g, '');
    if (cleaned.length <= 4) return cleaned;
    if (cleaned.length <= 6) return `${cleaned.slice(0, 4)} ${cleaned.slice(4)}`;
    return `${cleaned.slice(0, 4)} ${cleaned.slice(4, 6)} ${cleaned.slice(6, 10)}`;
  };

  /**
   * Валидация цены
   */
  const validatePrice = (price: string): boolean => {
    if (!price.trim()) return false;
    const numPrice = parseFloat(price);
    return !isNaN(numPrice) && numPrice >= 0;
  };

  /**
   * Полная валидация формы
   */
  const validateForm = (): boolean => {
    const newErrors: ValidationErrors = {};

    // Название (обязательно, минимум 2 символа)
    if (!formData.name?.trim()) {
      newErrors.name = 'Название товара/услуги обязательно';
    } else if (formData.name.trim().length < 2) {
      newErrors.name = 'Название должно быть не менее 2 символов';
    }

    // Номер по каталогу (обязательно)
    if (!formData.number?.trim()) {
      newErrors.number = 'Номер по каталогу обязателен';
    }

    // ТНВЭД код (обязательно, 10 цифр)
    if (!formData.tnved_code?.trim()) {
      newErrors.tnved_code = 'Код ТНВЭД обязателен для ЭСФ';
    } else if (!isValidTnvedCode(formData.tnved_code)) {
      newErrors.tnved_code = 'Код ТНВЭД должен содержать 10 цифр';
    }

    // Цена (обязательно, должна быть числом >= 0)
    if (!formData.price?.trim()) {
      newErrors.price = 'Цена обязательна';
    } else if (!validatePrice(formData.price)) {
      newErrors.price = 'Цена должна быть положительным числом';
    }

    // Валюта (обязательно)
    if (!formData.currency) {
      newErrors.currency = 'Выберите валюту';
    }

    // Единица измерения (обязательно)
    if (!formData.unit) {
      newErrors.unit = 'Выберите единицу измерения';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      setError('Пожалуйста, исправьте ошибки в форме');
      return;
    }

    setIsLoading(true);
    setError(null);
    setSuccessMessage('');

    try {
      // Подготовка данных для отправки
      const submitData: CreateCatalogItemRequest | UpdateCatalogItemRequest = {
        name: formData.name.trim(),
        number: formData.number.trim(),
        tnved_code: formData.tnved_code.replace(/\s/g, ''),
        price: parseFloat(formData.price),
        currency: formData.currency,
        unit: formData.unit,
        ...(formData.description && { description: formData.description.trim() }),
      };

      await onSubmit(submitData);
      
      setSuccessMessage(mode === 'create' ? 'Товар успешно добавлен в каталог' : 'Товар успешно обновлен');
      
      // Перенаправление через 1.5 секунды
      setTimeout(() => {
        router.push('/catalog');
      }, 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при сохранении товара');
      console.error('[CatalogForm] Error:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleInputChange = (field: keyof FormData, value: string) => {
    // Специальная обработка для ТНВЭД кода
    if (field === 'tnved_code') {
      const formatted = formatTnvedCode(value);
      setFormData(prev => ({ ...prev, [field]: formatted }));
    } else {
      setFormData(prev => ({ ...prev, [field]: value }));
    }
    
    // Очистка ошибки для этого поля при изменении (только для полей с валидацией)
    if (field in errors) {
      setErrors(prev => {
        const newErrors = { ...prev };
        delete newErrors[field as keyof ValidationErrors];
        return newErrors;
      });
    }
  };

  // Получить символ валюты
  const getCurrencySymbol = (code: string): string => {
    return CURRENCIES.find(c => c.code === code)?.symbol || code;
  };

  // Получить название единицы измерения
  const getUnitName = (code: string): string => {
    return UNITS_OF_MEASURE.find(u => u.code === code)?.name || code;
  };

  return (
    <div className="min-h-screen bg-linear-to-b from-gray-900 to-black">
      {/* Header */}
      <div className="bg-gray-800/50 border-b border-gray-700 sticky top-0 z-10 backdrop-blur-sm">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <button
            onClick={() => router.back()}
            disabled={isLoading}
            className="flex items-center gap-2 text-gray-400 hover:text-white mb-4 transition-colors disabled:opacity-50"
          >
            <ArrowLeft size={20} />
            Вернуться
          </button>
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <Package size={32} />
            {mode === 'create' ? 'Добавление товара в каталог' : 'Редактирование товара'}
          </h1>
        </div>
      </div>

      {/* Form Content */}
      <div className="max-w-4xl mx-auto px-4 py-8">
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Success Message */}
          {successMessage && (
            <div className="bg-emerald-900/20 border border-emerald-800 rounded-lg p-4 flex items-start gap-3 animate-in fade-in duration-300">
              <CheckCircle size={20} className="text-emerald-400 shrink-0 mt-0.5" />
              <p className="text-emerald-300">{successMessage}</p>
            </div>
          )}

          {/* Error Message */}
          {error && (
            <div className="bg-red-900/20 border border-red-800 rounded-lg p-4 flex items-start gap-3 animate-in fade-in duration-300">
              <AlertCircle size={20} className="text-red-400 shrink-0 mt-0.5" />
              <p className="text-red-300">{error}</p>
            </div>
          )}

          {/* Main Information Section */}
          <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6 space-y-5">
            <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              <Package size={20} />
              Основная информация
            </h3>

            {/* Item Name */}
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-300 mb-2">
                Наименование товара/услуги <span className="text-red-400">*</span>
              </label>
              <input
                id="name"
                type="text"
                value={formData.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                disabled={isLoading}
                placeholder="Например: Ноутбук Lenovo ThinkPad"
                className={`w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border outline-none transition-colors ${
                  errors.name
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white placeholder-gray-500 disabled:opacity-50`}
              />
              {errors.name && (
                <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.name}
                </p>
              )}
            </div>

            {/* Catalog Number */}
            <div>
              <label htmlFor="number" className="block text-sm font-medium text-gray-300 mb-2">
                Номер по каталогу <span className="text-red-400">*</span>
              </label>
              <input
                id="number"
                type="text"
                value={formData.number}
                onChange={(e) => handleInputChange('number', e.target.value)}
                disabled={isLoading}
                placeholder="TOV-001"
                className={`w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border outline-none transition-colors ${
                  errors.number
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white placeholder-gray-500 disabled:opacity-50`}
              />
              {errors.number && (
                <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.number}
                </p>
              )}
            </div>

            {/* Description */}
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-300 mb-2">
                Описание (опционально)
              </label>
              <textarea
                id="description"
                value={formData.description || ''}
                onChange={(e) => handleInputChange('description', e.target.value)}
                disabled={isLoading}
                placeholder="Подробное описание товара или услуги"
                rows={3}
                className="w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border border-gray-600 focus:border-blue-500 outline-none transition-colors resize-none text-white placeholder-gray-500 disabled:opacity-50"
              />
            </div>

            {/* TNVED Code */}
            <div>
              <label htmlFor="tnved_code" className="block text-sm font-medium text-gray-300 mb-2">
                Код ТНВЭД <span className="text-red-400">*</span>
              </label>
              <input
                id="tnved_code"
                type="text"
                value={formData.tnved_code}
                onChange={(e) => handleInputChange('tnved_code', e.target.value)}
                disabled={isLoading}
                placeholder="XXXX XX XXXX"
                maxLength={12} // 10 цифр + 2 пробела
                className={`w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border outline-none transition-colors ${
                  errors.tnved_code
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white placeholder-gray-500 disabled:opacity-50 font-mono`}
              />
              {errors.tnved_code && (
                <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.tnved_code}
                </p>
              )}
              <p className="text-gray-500 text-xs mt-1">
                Код ТН ВЭД обязателен для создания ЭСФ (10 цифр, автоформатирование)
              </p>
            </div>
          </div>

          {/* Pricing Section */}
          <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6 space-y-5">
            <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              <DollarSign size={20} />
              Ценообразование
            </h3>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              {/* Price */}
              <div>
                <label htmlFor="price" className="block text-sm font-medium text-gray-300 mb-2">
                  Цена <span className="text-red-400">*</span>
                </label>
                <div className="relative">
                  <input
                    id="price"
                    type="number"
                    step="0.01"
                    min="0"
                    value={formData.price}
                    onChange={(e) => handleInputChange('price', e.target.value)}
                    disabled={isLoading}
                    placeholder="0.00"
                    className={`w-full px-4 py-2.5 pr-16 rounded-lg bg-gray-900/50 border outline-none transition-colors ${
                      errors.price
                        ? 'border-red-600 focus:border-red-500'
                        : 'border-gray-600 focus:border-blue-500'
                    } text-white placeholder-gray-500 disabled:opacity-50`}
                  />
                  <span className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 font-medium">
                    {getCurrencySymbol(formData.currency)}
                  </span>
                </div>
                {errors.price && (
                  <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                    <AlertCircle size={14} />
                    {errors.price}
                  </p>
                )}
              </div>

              {/* Currency */}
              <div>
                <label htmlFor="currency" className="block text-sm font-medium text-gray-300 mb-2">
                  Валюта <span className="text-red-400">*</span>
                </label>
                <select
                  id="currency"
                  value={formData.currency}
                  onChange={(e) => handleInputChange('currency', e.target.value)}
                  disabled={isLoading}
                  className={`w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border outline-none transition-colors ${
                    errors.currency
                      ? 'border-red-600 focus:border-red-500'
                      : 'border-gray-600 focus:border-blue-500'
                  } text-white disabled:opacity-50`}
                >
                  {CURRENCIES.map((curr) => (
                    <option key={curr.code} value={curr.code}>
                      {curr.symbol} {curr.name} ({curr.code})
                    </option>
                  ))}
                </select>
                {errors.currency && (
                  <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                    <AlertCircle size={14} />
                    {errors.currency}
                  </p>
                )}
              </div>
            </div>

            {/* Unit of Measure */}
            <div>
              <label htmlFor="unit" className="block text-sm font-medium text-gray-300 mb-2">
                Единица измерения (код ОКЕИ) <span className="text-red-400">*</span>
              </label>
              <select
                id="unit"
                value={formData.unit}
                onChange={(e) => handleInputChange('unit', e.target.value)}
                disabled={isLoading}
                className={`w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border outline-none transition-colors ${
                  errors.unit
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white disabled:opacity-50`}
              >
                {UNITS_OF_MEASURE.map((unit) => (
                  <option key={unit.code} value={unit.code}>
                    {unit.name} (код {unit.code})
                  </option>
                ))}
              </select>
              {errors.unit && (
                <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.unit}
                </p>
              )}
              <p className="text-gray-500 text-xs mt-1">
                Выбранная единица: {getUnitName(formData.unit)} (код ОКЕИ: {formData.unit})
              </p>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-4 pt-4">
            <button
              type="submit"
              disabled={isLoading}
              className={`flex-1 px-6 py-3 rounded-lg font-medium flex items-center justify-center gap-2 transition-colors ${
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
                mode === 'create' ? 'Добавить в каталог' : 'Сохранить изменения'
              )}
            </button>
            <button
              type="button"
              onClick={() => router.push('/catalog')}
              disabled={isLoading}
              className="px-6 py-3 rounded-lg bg-gray-700 text-gray-200 hover:bg-gray-600 font-medium transition-colors disabled:opacity-50"
            >
              Отмена
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
