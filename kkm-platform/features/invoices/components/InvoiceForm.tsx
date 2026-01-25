'use client';

import React, { useState, useMemo } from 'react';
import { Plus, X, Loader2 } from 'lucide-react';

interface InvoiceItem {
  id: string;
  name: string;
  quantity: number;
  price: number;
}

interface InvoiceFormProps {
  initialData?: any;
  onSubmit?: (data: any) => Promise<void>;
}

export default function InvoiceForm({ initialData, onSubmit }: InvoiceFormProps) {
  const [formData, setFormData] = useState({
    number: initialData?.number || '',
    date: initialData?.date || new Date().toISOString().split('T')[0],
    company_id: initialData?.company_id || '',
    description: initialData?.description || '',
    items: initialData?.items || [] as InvoiceItem[],
    notes: initialData?.notes || '',
  });

  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [newItem, setNewItem] = useState({ name: '', quantity: 1, price: 0 });

  const total = useMemo(() => {
    return formData.items.reduce((sum: number, item: InvoiceItem) => sum + item.quantity * item.price, 0);
  }, [formData.items]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.number.trim()) newErrors.number = 'Номер счета обязателен';
    if (!formData.date) newErrors.date = 'Дата обязательна';
    if (!formData.company_id) newErrors.company_id = 'Компания обязательна';
    if (formData.items.length === 0) newErrors.items = 'Добавьте хотя бы один товар';

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

  const addItem = () => {
    if (!newItem.name.trim() || newItem.quantity <= 0 || newItem.price <= 0) {
      return;
    }

    const item: InvoiceItem = {
      id: Date.now().toString(),
      ...newItem,
    };

    setFormData({
      ...formData,
      items: [...formData.items, item],
    });

    setNewItem({ name: '', quantity: 1, price: 0 });
  };

  const removeItem = (id: string) => {
    setFormData({
      ...formData,
      items: formData.items.filter((item: any) => item.id !== id),
    });
  };

  return (
    <form onSubmit={handleSubmit} className="bg-gray-900 rounded-lg border border-gray-800 p-6 space-y-6">
      <h2 className="text-2xl font-bold text-white">
        {initialData?.id ? 'Редактировать счет' : 'Создать новый счет'}
      </h2>

      {errors.submit && (
        <div className="p-4 bg-red-900/20 border border-red-800 text-red-300 rounded-lg">
          {errors.submit}
        </div>
      )}

      {/* Основные данные */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-300">Основные данные</h3>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Номер счета *
            </label>
            <input
              type="text"
              value={formData.number}
              onChange={(e) => setFormData({ ...formData, number: e.target.value })}
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-blue-500 outline-none"
              placeholder="СЧ-001"
            />
            {errors.number && <p className="text-red-400 text-sm mt-1">{errors.number}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Дата *
            </label>
            <input
              type="date"
              value={formData.date}
              onChange={(e) => setFormData({ ...formData, date: e.target.value })}
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-blue-500 outline-none"
            />
            {errors.date && <p className="text-red-400 text-sm mt-1">{errors.date}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Компания *
            </label>
            <select
              value={formData.company_id}
              onChange={(e) => setFormData({ ...formData, company_id: e.target.value })}
              className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-blue-500 outline-none"
            >
              <option value="">Выберите компанию</option>
            </select>
            {errors.company_id && <p className="text-red-400 text-sm mt-1">{errors.company_id}</p>}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Описание
          </label>
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-blue-500 outline-none"
            placeholder="Описание счета"
            rows={3}
          />
        </div>
      </div>

      {/* Товары */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-300">Товары и услуги</h3>
        
        {errors.items && <p className="text-red-400 text-sm">{errors.items}</p>}

        <div className="bg-gray-800/50 p-4 rounded-lg space-y-3">
          <div className="grid grid-cols-5 gap-3">
            <input
              type="text"
              value={newItem.name}
              onChange={(e) => setNewItem({ ...newItem, name: e.target.value })}
              placeholder="Название товара"
              className="col-span-2 px-3 py-2 rounded bg-gray-700 border border-gray-600 text-white text-sm focus:border-blue-500 outline-none"
            />
            <input
              type="number"
              value={newItem.quantity}
              onChange={(e) => setNewItem({ ...newItem, quantity: parseFloat(e.target.value) || 1 })}
              min="1"
              className="px-3 py-2 rounded bg-gray-700 border border-gray-600 text-white text-sm focus:border-blue-500 outline-none"
              placeholder="Кол-во"
            />
            <input
              type="number"
              value={newItem.price}
              onChange={(e) => setNewItem({ ...newItem, price: parseFloat(e.target.value) || 0 })}
              min="0"
              className="px-3 py-2 rounded bg-gray-700 border border-gray-600 text-white text-sm focus:border-blue-500 outline-none"
              placeholder="Цена"
            />
            <button
              type="button"
              onClick={addItem}
              className="px-3 py-2 rounded bg-blue-600 hover:bg-blue-700 text-white font-medium text-sm transition-colors flex items-center justify-center gap-1"
            >
              <Plus size={16} /> Добавить
            </button>
          </div>
        </div>

        <div className="space-y-2">
          {formData.items.map((item: any) => (
            <div key={item.id} className="flex items-center justify-between p-3 bg-gray-800 rounded-lg">
              <div className="flex-1">
                <p className="text-white font-medium">{item.name}</p>
                <p className="text-gray-400 text-sm">
                  {item.quantity} × {item.price.toLocaleString('ru-RU')} ₽ = {(item.quantity * item.price).toLocaleString('ru-RU')} ₽
                </p>
              </div>
              <button
                type="button"
                onClick={() => removeItem(item.id)}
                className="p-2 hover:bg-red-600/20 text-red-400 rounded transition-colors"
              >
                <X size={18} />
              </button>
            </div>
          ))}
        </div>

        <div className="p-4 bg-blue-900/20 border border-blue-800 rounded-lg text-right">
          <p className="text-gray-400 text-sm mb-2">Итого:</p>
          <p className="text-3xl font-bold text-blue-400">
            {total.toLocaleString('ru-RU')} ₽
          </p>
        </div>
      </div>

      {/* Примечания */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-2">
          Примечания
        </label>
        <textarea
          value={formData.notes}
          onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
          className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-blue-500 outline-none"
          placeholder="Дополнительная информация"
          rows={3}
        />
      </div>

      {/* Кнопки */}
      <div className="flex gap-3">
        <button
          type="submit"
          disabled={isLoading}
          className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
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
