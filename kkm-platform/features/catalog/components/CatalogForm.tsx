'use client';

import React, { useState, useEffect } from 'react';
import { Loader2, AlertCircle, CheckCircle } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useCatalog } from '@/hooks';

interface CatalogFormProps {
  initialData?: any;
  itemId?: string;
}

export default function CatalogForm({ initialData, itemId }: CatalogFormProps) {
  const router = useRouter();
  const { createItem, updateItem, error: storeError, isLoading: storeLoading } = useCatalog();
  
  const [formData, setFormData] = useState({
    name: initialData?.name || '',
    sku: initialData?.sku || initialData?.number || '',
    category: initialData?.category || initialData?.tnved_code || '',
    description: initialData?.description || '',
    price: initialData?.price || 0,
    currency: initialData?.currency || 'RUB',
    unit: initialData?.unit || 'шт',
    costPrice: initialData?.costPrice || initialData?.cost || 0,
    stock: initialData?.stock || initialData?.quantity || 0,
    reorderLevel: initialData?.reorderLevel || initialData?.min_quantity || 0,
    active: initialData?.active !== undefined ? initialData.active : true,
  });

  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [successMessage, setSuccessMessage] = useState('');

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim() || formData.name.length < 2) {
      newErrors.name = 'Название должно содержать минимум 2 символа';
    }
    if (!formData.sku.trim()) newErrors.sku = 'Артикул обязателен';
    if (!formData.category.trim()) newErrors.category = 'Категория обязательна';
    if (formData.price <= 0) newErrors.price = 'Цена должна быть больше 0';
    if (formData.costPrice < 0) newErrors.costPrice = 'Себестоимость не может быть отрицательной';

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
    setSuccessMessage('');

    try {
      if (itemId) {
        await updateItem(itemId, formData);
        setSuccessMessage('Товар успешно обновлен');
        setTimeout(() => router.push('/catalog'), 1500);
      } else {
        await createItem(formData);
        setSuccessMessage('Товар успешно создан');
        setFormData({
          name: '',
          sku: '',
          category: '',
          description: '',
          price: 0,
          currency: 'RUB',
          unit: 'шт',
          costPrice: 0,
          stock: 0,
          reorderLevel: 0,
          active: true,
        });
        setTimeout(() => router.push('/catalog'), 1500);
      }
    } catch (err) {
      setErrors({
        submit: err instanceof Error ? err.message : 'Ошибка отправки',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const margin = formData.price - formData.costPrice;
  const marginPercent = formData.price > 0 ? ((margin / formData.price) * 100).toFixed(1) : 0;

  return (
    <form onSubmit={handleSubmit} className="bg-gray-900 rounded-lg border border-gray-800 p-6 space-y-6">
      <h2 className="text-2xl font-bold text-white">
        {itemId ? 'Редактировать товар' : 'Добавить товар'}
      </h2>

      {successMessage && (
        <div className="p-4 bg-green-900/20 border border-green-800 text-green-300 rounded-lg flex items-center gap-2">
          <CheckCircle size={18} />
          {successMessage}
        </div>
      )}

      {(errors.submit || storeError) && (
        <div className="p-4 bg-red-900/20 border border-red-800 text-red-300 rounded-lg flex items-center gap-2">
          <AlertCircle size={18} />
          {errors.submit || storeError}
        </div>
      )}

      {/* Основная информация */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-300">Основная информация</h3>
        
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Название товара *
          </label>
          <input
            type="text"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
            placeholder="Название товара"
          />
          {errors.name && <p className="text-red-400 text-sm mt-1">{errors.name}</p>}
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Артикул (SKU) *
            </label>
            <input
              type="text"
              value={formData.sku}
              onChange={(e) => setFormData({ ...formData, sku: e.target.value })}
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
              placeholder="SKU-001"
            />
            {errors.sku && <p className="text-red-400 text-sm mt-1">{errors.sku}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Категория *
            </label>
            <input
              type="text"
              value={formData.category}
              onChange={(e) => setFormData({ ...formData, category: e.target.value })}
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
              placeholder="Категория"
            />
            {errors.category && <p className="text-red-400 text-sm mt-1">{errors.category}</p>}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Описание
          </label>
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
            placeholder="Описание товара"
            rows={3}
          />
        </div>
      </div>

      {/* Цены и остатки */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-300">Цены и остатки</h3>
        
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Себестоимость
            </label>
            <input
              type="number"
              value={formData.costPrice}
              onChange={(e) => setFormData({ ...formData, costPrice: parseFloat(e.target.value) || 0 })}
              min="0"
              step="0.01"
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
              placeholder="0"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Цена продажи *
            </label>
            <input
              type="number"
              value={formData.price}
              onChange={(e) => setFormData({ ...formData, price: parseFloat(e.target.value) || 0 })}
              min="0"
              step="0.01"
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
              placeholder="0"
            />
            {errors.price && <p className="text-red-400 text-sm mt-1">{errors.price}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Валюта
            </label>
            <select
              value={formData.currency}
              onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
            >
              <option value="RUB">RUB (₽)</option>
              <option value="USD">USD ($)</option>
              <option value="EUR">EUR (€)</option>
              <option value="KZT">KZT (₸)</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Единица измерения
            </label>
            <input
              type="text"
              value={formData.unit}
              onChange={(e) => setFormData({ ...formData, unit: e.target.value })}
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
              placeholder="шт, м, л и т.д."
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Остаток
            </label>
            <input
              type="number"
              value={formData.stock}
              onChange={(e) => setFormData({ ...formData, stock: parseInt(e.target.value) || 0 })}
              min="0"
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
              placeholder="0"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Минимум для переказа
            </label>
            <input
              type="number"
              value={formData.reorderLevel}
              onChange={(e) => setFormData({ ...formData, reorderLevel: parseInt(e.target.value) || 0 })}
              min="0"
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-purple-500 outline-none"
              placeholder="0"
            />
          </div>
        </div>
      </div>

      {/* Показатели */}
      {formData.price > 0 && (
        <div className="grid grid-cols-3 gap-4 p-4 bg-purple-900/20 border border-purple-800 rounded-lg">
          <div>
            <p className="text-gray-400 text-sm mb-1">Маржа</p>
            <p className="text-2xl font-bold text-purple-400">
              {margin.toLocaleString('ru-RU')} ₽
            </p>
          </div>
          <div>
            <p className="text-gray-400 text-sm mb-1">Маржа %</p>
            <p className="text-2xl font-bold text-purple-400">
              {marginPercent}%
            </p>
          </div>
          <div>
            <p className="text-gray-400 text-sm mb-1">Стоимость остатка</p>
            <p className="text-2xl font-bold text-purple-400">
              {(formData.stock * formData.price).toLocaleString('ru-RU')} ₽
            </p>
          </div>
        </div>
      )}

      {/* Кнопки */}
      <div className="flex gap-3">
        <button
          type="submit"
          disabled={isLoading || storeLoading}
          className="flex-1 px-4 py-2 bg-purple-600 text-white rounded-lg font-medium hover:bg-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
        >
          {isLoading || storeLoading ? (
            <>
              <Loader2 className="animate-spin" size={18} />
              Сохранение...
            </>
          ) : (
            'Сохранить товар'
          )}
        </button>
      </div>
    </form>
  );
}
