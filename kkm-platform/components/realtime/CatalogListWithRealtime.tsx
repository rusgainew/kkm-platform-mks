'use client';

import React, { useEffect, useState } from 'react';
import { useRealtimeCatalog, RealtimeProduct } from '@/store/realtime';

const PRODUCT_CATEGORIES = [
  { id: 'all', name: 'Все товары' },
  { id: 'electronics', name: '🖥️ Электроника' },
  { id: 'supplies', name: '📎 Расходники' },
  { id: 'furniture', name: '🪑 Мебель' },
];

export function CatalogListWithRealtime() {
  const products = useRealtimeCatalog((state) => state.getAllProducts());
  const [displayProducts, setDisplayProducts] = useState<RealtimeProduct[]>([]);
  const [category, setCategory] = useState<string>('all');
  const [search, setSearch] = useState<string>('');
  const [sortBy, setSortBy] = useState<'name' | 'price' | 'stock'>('name');

  useEffect(() => {
    let filtered = products;

    // Фильтр по категории (берется из SKU префикса для примера)
    if (category !== 'all') {
      filtered = filtered.filter((p) => p.sku.toLowerCase().startsWith(category[0]));
    }

    // Поиск по названию
    if (search) {
      filtered = filtered.filter((p) =>
        p.name.toLowerCase().includes(search.toLowerCase()) ||
        p.sku.toLowerCase().includes(search.toLowerCase())
      );
    }

    // Сортировка
    if (sortBy === 'price') {
      filtered = filtered.sort((a, b) => a.price - b.price);
    } else if (sortBy === 'stock') {
      filtered = filtered.sort((a, b) => a.stock - b.stock);
    } else {
      filtered = filtered.sort((a, b) => a.name.localeCompare(b.name));
    }

    setDisplayProducts(filtered);
  }, [products, category, search, sortBy]);

  const getStockStatus = (stock: number) => {
    if (stock === 0) return { color: 'text-red-600', icon: '❌', label: 'Нет в наличии' };
    if (stock < 5) return { color: 'text-orange-600', icon: '⚠️', label: 'Критически мало' };
    if (stock < 20) return { color: 'text-yellow-600', icon: '⏱️', label: 'Мало' };
    return { color: 'text-green-600', icon: '✅', label: 'В наличии' };
  };

  const totalProducts = products.length;
  const totalValue = products.reduce((sum, p) => sum + p.price * p.stock, 0);
  const lowStockCount = products.filter((p) => p.stock < 20).length;

  return (
    <div className="space-y-4">
      {/* Заголовок и статистика */}
      <div>
        <h2 className="text-2xl font-bold text-gray-900">
          📦 Каталог товаров
          {displayProducts.length > 0 && (
            <span className="text-sm font-normal text-gray-500 ml-2">
              ({displayProducts.length})
            </span>
          )}
        </h2>
        <p className="text-sm text-gray-600">
          <span className="inline-flex items-center gap-1">
            <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
            Реал-тайм синхронизация товаров и цен
          </span>
        </p>
      </div>

      {/* Карточки статистики */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div className="bg-gradient-to-br from-blue-50 to-blue-100 p-4 rounded-lg">
          <div className="text-2xl font-bold text-blue-900">{totalProducts}</div>
          <div className="text-xs text-blue-700">Всего товаров</div>
        </div>
        <div className="bg-gradient-to-br from-green-50 to-green-100 p-4 rounded-lg">
          <div className="text-2xl font-bold text-green-900">
            ₽{(totalValue / 1000).toFixed(0)}K
          </div>
          <div className="text-xs text-green-700">Общая стоимость</div>
        </div>
        <div className="bg-gradient-to-br from-yellow-50 to-yellow-100 p-4 rounded-lg">
          <div className="text-2xl font-bold text-yellow-900">{lowStockCount}</div>
          <div className="text-xs text-yellow-700">Низкий остаток</div>
        </div>
        <div className="bg-gradient-to-br from-purple-50 to-purple-100 p-4 rounded-lg">
          <div className="text-2xl font-bold text-purple-900">Live</div>
          <div className="text-xs text-purple-700">Синхронизация</div>
        </div>
      </div>

      {/* Поиск и фильтры */}
      <div className="flex gap-2 flex-col md:flex-row">
        <input
          type="text"
          placeholder="Поиск по названию или SKU..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />

        <select
          value={sortBy}
          onChange={(e) => setSortBy(e.target.value as any)}
          className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="name">По названию</option>
          <option value="price">По цене</option>
          <option value="stock">По остатку</option>
        </select>

        <div className="flex gap-1 flex-wrap">
          {PRODUCT_CATEGORIES.map((cat) => (
            <button
              key={cat.id}
              onClick={() => setCategory(cat.id)}
              className={`px-3 py-2 text-sm rounded-lg transition ${
                category === cat.id
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              {cat.name}
            </button>
          ))}
        </div>
      </div>

      {/* Сетка товаров */}
      {displayProducts.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {displayProducts.map((product) => {
            const stockStatus = getStockStatus(product.stock);
            return (
              <div
                key={product.id}
                className="border border-gray-200 rounded-lg p-4 hover:shadow-lg transition"
              >
                {/* Заголовок товара */}
                <div className="flex justify-between items-start mb-2">
                  <div className="flex-1">
                    <h3 className="font-semibold text-gray-900">{product.name}</h3>
                    <p className="text-xs text-gray-500">SKU: {product.sku}</p>
                  </div>
                </div>

                {/* Цена и остаток */}
                <div className="flex justify-between items-end mb-3">
                  <div>
                    <p className="text-sm text-gray-600">Цена</p>
                    <p className="text-xl font-bold text-gray-900">
                      ₽{product.price.toLocaleString('ru-RU')}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-gray-600">Остаток</p>
                    <p className={`text-xl font-bold ${stockStatus.color}`}>
                      {product.stock} шт
                    </p>
                  </div>
                </div>

                {/* Статус остатка */}
                <div className={`flex items-center gap-2 p-2 rounded-lg ${stockStatus.color} bg-opacity-10 border border-current border-opacity-30 mb-3`}>
                  <span className="text-lg">{stockStatus.icon}</span>
                  <span className="text-xs font-medium">{stockStatus.label}</span>
                </div>

                {/* Дата обновления */}
                <p className="text-xs text-gray-500 mb-3">
                  Обновлено:{' '}
                  {new Date(product.updatedAt).toLocaleString('ru-RU', {
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </p>

                {/* Действия */}
                <div className="flex gap-2">
                  <button className="flex-1 px-3 py-2 bg-blue-500 text-white text-sm rounded-lg hover:bg-blue-600 transition">
                    Подробно
                  </button>
                  <button className="flex-1 px-3 py-2 bg-gray-200 text-gray-700 text-sm rounded-lg hover:bg-gray-300 transition">
                    Редактировать
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        <div className="text-center py-12 bg-gray-50 rounded-lg">
          <p className="text-gray-600">Товары не найдены</p>
        </div>
      )}
    </div>
  );
}
