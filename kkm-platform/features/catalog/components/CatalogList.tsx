'use client';

import React, { useState, useEffect } from 'react';
import { Package, Plus, Edit2, Trash2, Search, Filter, AlertCircle } from 'lucide-react';
import Link from 'next/link';
import { useCatalog } from '@/hooks';

export default function CatalogList() {
  const { items, categories, isLoading, error, fetchCatalog, setFilters } = useCatalog();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');

  useEffect(() => {
    fetchCatalog();
  }, []);

  const products = items;

  const filteredProducts = products.filter(p =>
    (p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.sku.toLowerCase().includes(searchQuery.toLowerCase())) &&
    (!selectedCategory || p.category === selectedCategory)
  );

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
      <div className="p-6 border-b border-gray-800 flex items-center justify-between">
        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
          <Package className="w-6 h-6 text-purple-400" />
          Каталог товаров
        </h2>
        <Link
          href="/catalog/create"
          className="flex items-center gap-2 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 disabled:opacity-50"
          aria-disabled={isLoading}
        >
          <Plus size={18} /> Добавить товар
        </Link>
      </div>

      {error && (
        <div className="p-4 bg-red-900/20 border-b border-red-800 text-red-300 flex items-center gap-2">
          <AlertCircle size={18} />
          {error}
        </div>
      )}

      <div className="p-6 border-b border-gray-800 space-y-4">
        <input
          type="text"
          placeholder="Поиск по названию или артикулу..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:border-purple-500 outline-none"
        />
        
        {categories.length > 0 && (
          <div className="flex items-center gap-2 flex-wrap">
            <Filter size={18} className="text-gray-400" />
            <button
              onClick={() => setSelectedCategory('')}
              className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                !selectedCategory
                  ? 'bg-purple-600 text-white'
                  : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
              }`}
            >
              Все
            </button>
            {categories.map(cat => (
              <button
                key={cat}
                onClick={() => setSelectedCategory(cat)}
                className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                  selectedCategory === cat
                    ? 'bg-purple-600 text-white'
                    : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
                }`}
              >
                {cat}
              </button>
            ))}
          </div>
        )}
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-800/50">
            <tr>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Название</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Артикул</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-300">Категория</th>
              <th className="px-6 py-3 text-right text-sm font-semibold text-gray-300">Цена</th>
              <th className="px-6 py-3 text-right text-sm font-semibold text-gray-300">Остаток</th>
              <th className="px-6 py-3 text-center text-sm font-semibold text-gray-300">Действия</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800">
            {isLoading ? (
              <tr>
                <td colSpan={6} className="px-6 py-8 text-center text-gray-400">
                  Загрузка товаров...
                </td>
              </tr>
            ) : filteredProducts.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-6 py-8 text-center text-gray-400">
                  {products.length === 0 ? 'Нет товаров' : 'Товары не найдены'}
                </td>
              </tr>
            ) : (
              filteredProducts.map((product) => (
                <tr key={product.id} className="hover:bg-gray-800/50 transition-colors">
                  <td className="px-6 py-4 text-white font-medium">{product.name}</td>
                  <td className="px-6 py-4 text-gray-300 font-mono">{product.sku}</td>
                  <td className="px-6 py-4 text-gray-300">{product.category}</td>
                  <td className="px-6 py-4 text-white font-medium text-right">{product.price.toLocaleString('ru-RU')} ₽</td>
                  <td className={`px-6 py-4 text-right font-medium ${
                    product.stock > 10 ? 'text-green-400' :
                    product.stock > 0 ? 'text-yellow-400' :
                    'text-red-400'
                  }`}>
                    {product.stock}
                  </td>
                  <td className="px-6 py-4 flex items-center justify-center gap-2">
                    <Link
                      href={`/catalog/${product.id}/edit`}
                      className="p-2 hover:bg-purple-600/20 text-purple-400 rounded transition-colors"
                    >
                      <Edit2 size={18} />
                    </Link>
                    <button
                      onClick={() => {
                        if (confirm('Вы уверены?')) {
                          // Delete will be handled by store
                        }
                      }}
                      className="p-2 hover:bg-red-600/20 text-red-400 rounded transition-colors"
                    >
                      <Trash2 size={18} />
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
