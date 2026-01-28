'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeft, AlertCircle, CheckCircle, Building2 } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import type { Company, CreateCompanyRequest, UpdateCompanyRequest, CompanyStatus } from '@/types/entities';
import { isValidEmail } from '@/types/entities';

interface CompanyFormProps {
  initialData?: Company;
  onSubmit: (data: CreateCompanyRequest | UpdateCompanyRequest) => Promise<Company>;
  mode?: 'create' | 'edit';
}

interface FormData {
  name: string;
  tin: string;
  kpp?: string;
  ogrn?: string;
  address: string;
  phone: string;
  email: string;
  website?: string;
  status?: CompanyStatus;
}

interface ValidationErrors {
  name?: string;
  tin?: string;
  address?: string;
  phone?: string;
  email?: string;
}

export default function CompanyForm({ initialData, onSubmit, mode = 'create' }: CompanyFormProps) {
  const router = useRouter();
  const [formData, setFormData] = useState<FormData>({
    name: initialData?.name || '',
    tin: initialData?.tin || '',
    kpp: initialData?.kpp || '',
    ogrn: initialData?.ogrn || '',
    address: initialData?.address || '',
    phone: initialData?.phone || '',
    email: initialData?.email || '',
    website: initialData?.website || '',
    status: initialData?.status || 'active',
  });

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState('');
  const [errors, setErrors] = useState<ValidationErrors>({});

  /**
   * Валидация ИНН
   * Для юрлиц - 10 цифр, для ИП - 12 цифр
   */
  const validateTIN = (tin: string): boolean => {
    const cleaned = tin.replace(/\s/g, '');
    return /^\d{10}$/.test(cleaned) || /^\d{12}$/.test(cleaned);
  };

  /**
   * Валидация телефона (базовая)
   * Поддерживает форматы: +7XXXXXXXXXX, 8XXXXXXXXXX, 7XXXXXXXXXX
   */
  const validatePhone = (phone: string): boolean => {
    const cleaned = phone.replace(/[\s\-()]/g, '');
    return /^[\+]?[78]\d{10}$/.test(cleaned);
  };

  /**
   * Полная валидация формы
   */
  const validateForm = (): boolean => {
    const newErrors: ValidationErrors = {};

    // Название компании (обязательно, минимум 2 символа)
    if (!formData.name?.trim()) {
      newErrors.name = 'Название компании обязательно';
    } else if (formData.name.trim().length < 2) {
      newErrors.name = 'Название должно быть не менее 2 символов';
    }

    // ИНН (обязательно, должен быть валидным)
    if (!formData.tin?.trim()) {
      newErrors.tin = 'ИНН обязателен';
    } else if (!validateTIN(formData.tin)) {
      newErrors.tin = 'ИНН должен содержать 10 (юр. лицо) или 12 (ИП) цифр';
    }

    // Адрес (обязательно, минимум 10 символов)
    if (!formData.address?.trim()) {
      newErrors.address = 'Адрес обязателен';
    } else if (formData.address.trim().length < 10) {
      newErrors.address = 'Адрес должен быть не менее 10 символов';
    }

    // Телефон (обязательно, должен быть валидным)
    if (!formData.phone?.trim()) {
      newErrors.phone = 'Телефон обязателен';
    } else if (!validatePhone(formData.phone)) {
      newErrors.phone = 'Введите корректный номер телефона (+7XXXXXXXXXX)';
    }

    // Email (обязательно, должен быть валидным)
    if (!formData.email?.trim()) {
      newErrors.email = 'Email обязателен';
    } else if (!isValidEmail(formData.email)) {
      newErrors.email = 'Введите корректный email';
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
      const submitData: CreateCompanyRequest | UpdateCompanyRequest = {
        name: formData.name.trim(),
        tin: formData.tin.trim(),
        address: formData.address.trim(),
        phone: formData.phone.trim(),
        email: formData.email.trim(),
        ...(formData.kpp && { kpp: formData.kpp.trim() }),
        ...(formData.ogrn && { ogrn: formData.ogrn.trim() }),
        ...(formData.website && { website: formData.website.trim() }),
        ...(mode === 'edit' && formData.status && { status: formData.status }),
      };

      await onSubmit(submitData);
      
      setSuccessMessage(mode === 'create' ? 'Компания успешно создана' : 'Компания успешно обновлена');
      
      // Перенаправление через 1.5 секунды
      setTimeout(() => {
        router.push('/companies');
      }, 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при сохранении компании');
      console.error('[CompanyForm] Error:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleInputChange = (field: keyof FormData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    // Очистка ошибки для этого поля при изменении
    if (errors[field as keyof ValidationErrors]) {
      setErrors(prev => {
        const newErrors = { ...prev };
        delete newErrors[field as keyof ValidationErrors];
        return newErrors;
      });
    }
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
            <Building2 size={32} />
            {mode === 'create' ? 'Создание компании' : 'Редактирование компании'}
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
              <Building2 size={20} />
              Основная информация
            </h3>

            {/* Company Name */}
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-300 mb-2">
                Название компании <span className="text-red-400">*</span>
              </label>
              <input
                id="name"
                type="text"
                value={formData.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                disabled={isLoading}
                placeholder="ООО &quot;Название компании&quot;"
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

            {/* TIN (ИНН) */}
            <div>
              <label htmlFor="tin" className="block text-sm font-medium text-gray-300 mb-2">
                ИНН <span className="text-red-400">*</span>
              </label>
              <input
                id="tin"
                type="text"
                value={formData.tin}
                onChange={(e) => handleInputChange('tin', e.target.value)}
                disabled={isLoading}
                placeholder="1234567890 или 123456789012"
                maxLength={12}
                className={`w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border outline-none transition-colors ${
                  errors.tin
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white placeholder-gray-500 disabled:opacity-50`}
              />
              {errors.tin && (
                <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.tin}
                </p>
              )}
              <p className="text-gray-500 text-xs mt-1">10 цифр для юр. лица, 12 цифр для ИП</p>
            </div>

            {/* KPP (optional) */}
            <div>
              <label htmlFor="kpp" className="block text-sm font-medium text-gray-300 mb-2">
                КПП (для юридических лиц)
              </label>
              <input
                id="kpp"
                type="text"
                value={formData.kpp || ''}
                onChange={(e) => handleInputChange('kpp', e.target.value)}
                disabled={isLoading}
                placeholder="123456789"
                maxLength={9}
                className="w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border border-gray-600 focus:border-blue-500 outline-none transition-colors text-white placeholder-gray-500 disabled:opacity-50"
              />
            </div>

            {/* OGRN (optional) */}
            <div>
              <label htmlFor="ogrn" className="block text-sm font-medium text-gray-300 mb-2">
                ОГРН
              </label>
              <input
                id="ogrn"
                type="text"
                value={formData.ogrn || ''}
                onChange={(e) => handleInputChange('ogrn', e.target.value)}
                disabled={isLoading}
                placeholder="1234567890123 или 123456789012345"
                maxLength={15}
                className="w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border border-gray-600 focus:border-blue-500 outline-none transition-colors text-white placeholder-gray-500 disabled:opacity-50"
              />
              <p className="text-gray-500 text-xs mt-1">13 цифр для юр. лица, 15 цифр для ИП (ОГРНИП)</p>
            </div>
          </div>

          {/* Contact Information Section */}
          <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6 space-y-5">
            <h3 className="text-lg font-semibold text-white mb-4">Контактная информация</h3>

            {/* Address */}
            <div>
              <label htmlFor="address" className="block text-sm font-medium text-gray-300 mb-2">
                Юридический адрес <span className="text-red-400">*</span>
              </label>
              <textarea
                id="address"
                value={formData.address}
                onChange={(e) => handleInputChange('address', e.target.value)}
                disabled={isLoading}
                placeholder="г. Москва, ул. Примерная, д. 1, офис 100"
                rows={3}
                className={`w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border outline-none transition-colors resize-none ${
                  errors.address
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white placeholder-gray-500 disabled:opacity-50`}
              />
              {errors.address && (
                <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.address}
                </p>
              )}
            </div>

            {/* Phone */}
            <div>
              <label htmlFor="phone" className="block text-sm font-medium text-gray-300 mb-2">
                Телефон <span className="text-red-400">*</span>
              </label>
              <input
                id="phone"
                type="tel"
                value={formData.phone}
                onChange={(e) => handleInputChange('phone', e.target.value)}
                disabled={isLoading}
                placeholder="+7 (999) 123-45-67"
                className={`w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border outline-none transition-colors ${
                  errors.phone
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white placeholder-gray-500 disabled:opacity-50`}
              />
              {errors.phone && (
                <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.phone}
                </p>
              )}
            </div>

            {/* Email */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-300 mb-2">
                Email <span className="text-red-400">*</span>
              </label>
              <input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) => handleInputChange('email', e.target.value)}
                disabled={isLoading}
                placeholder="info@company.ru"
                className={`w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border outline-none transition-colors ${
                  errors.email
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white placeholder-gray-500 disabled:opacity-50`}
              />
              {errors.email && (
                <p className="text-red-400 text-sm mt-1.5 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.email}
                </p>
              )}
            </div>

            {/* Website (optional) */}
            <div>
              <label htmlFor="website" className="block text-sm font-medium text-gray-300 mb-2">
                Веб-сайт
              </label>
              <input
                id="website"
                type="url"
                value={formData.website || ''}
                onChange={(e) => handleInputChange('website', e.target.value)}
                disabled={isLoading}
                placeholder="https://company.ru"
                className="w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border border-gray-600 focus:border-blue-500 outline-none transition-colors text-white placeholder-gray-500 disabled:opacity-50"
              />
            </div>
          </div>

          {/* Status (only in edit mode) */}
          {mode === 'edit' && (
            <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6">
              <h3 className="text-lg font-semibold text-white mb-4">Статус компании</h3>
              <div>
                <label htmlFor="status" className="block text-sm font-medium text-gray-300 mb-2">
                  Статус
                </label>
                <select
                  id="status"
                  value={formData.status}
                  onChange={(e) => handleInputChange('status', e.target.value as CompanyStatus)}
                  disabled={isLoading}
                  className="w-full px-4 py-2.5 rounded-lg bg-gray-900/50 border border-gray-600 focus:border-blue-500 outline-none transition-colors text-white disabled:opacity-50"
                >
                  <option value="active">Активна</option>
                  <option value="inactive">Неактивна</option>
                  <option value="suspended">Приостановлена</option>
                </select>
              </div>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex gap-4 pt-4">
            <Button
              type="submit"
              disabled={isLoading}
              loading={isLoading}
              variant="primary"
              className="flex-1"
            >
              {mode === 'create' ? 'Создать компанию' : 'Сохранить изменения'}
            </Button>
            <button
              type="button"
              onClick={() => router.push('/companies')}
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
