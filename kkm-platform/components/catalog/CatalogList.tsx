'use client';

import React, { useState, useEffect, useMemo } from 'react';
import {
  Package,
  Plus,
  Edit2,
  Trash2,
  Search,
  Loader2,
  AlertCircle,
  Download,
  Upload,
} from 'lucide-react';
import Link from 'next/link';
import { useCatalog } from '@/hooks';

interface CatalogItemAPI {
  id: string;
  name: string;
  sku: string;
  category?: string;
  price?: number;
  quantity?: number;
  active?: boolean;
  created_at?: string;
}

interface Product {
  id: string;
  name: string;
  sku: string;
  category: string;
  price: number;
  quantity: number;
  status: 'active' | 'inactive';
  createdAt: string;
}

const getStatusBadgeColor = (status: string) => {
  const colors: Record<string, string> = {
    active: 'bg-emerald-900/20 text-emerald-300 border border-emerald-800',
    inactive: 'bg-gray-800/50 text-gray-300 border border-gray-700',
  };
  return colors[status] || colors.inactive;
};

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    active: 'Активен',
    inactive: 'Неактивен',
  };
  return labels[status] || status;
};

const getStockStatus = (quantity: number) => {
  if (quantity === 0) return 'out-of-stock';
  if (quantity < 10) return 'low-stock';
  return 'in-stock';
};

const getStockBadgeColor = (quantity: number) => {
  const status = getStockStatus(quantity);
  const colors: Record<string, string> = {
    'in-stock': 'text-emerald-300',
    'low-stock': 'text-amber-300',
    'out-of-stock': 'text-red-300',
  };
  return colors[status] || 'text-gray-300';
};

export default function CatalogList() {
  const { items, isLoading, error, fetchCatalog } = useCatalog();
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoadingLocal, setIsLoadingLocal] = useState(true);
  const [apiError, setApiError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [stockFilter, setStockFilter] = useState<string>('all');

  // Load products from API
  useEffect(() => {
    const loadProducts = async () => {
      try {
        setIsLoadingLocal(true);
        setApiError(null);

        // Fetch catalog from API
        await fetchCatalog();
      } catch (err) {
        let errorMessage = 'Ошибка загрузки товаров';

        if (err instanceof Error) {
          errorMessage = err.message;
        } else if (typeof err === 'string') {
          errorMessage = err;
        }

        console.error('[CatalogList] Error loading products:', errorMessage);
        setApiError(errorMessage);
        
        // Fallback: показываем пустой список вместо крашения
        setProducts([]);
      } finally {
        setIsLoadingLocal(false);
      }
    };

    loadProducts();
    // Зависимость от fetchCatalog может вызвать бесконечный loop
    // Добавляем проверку на наличие ошибки для предотвращения повторных запросов
  }, []);

  // Transform API items to Product format
  useEffect(() => {
    if (items && Array.isArray(items) && items.length > 0) {
      const transformedProducts: Product[] = items.map((item: CatalogItemAPI) => ({
        id: item.id,
        name: item.name || 'Без названия',
        sku: item.sku || 'N/A',
        category: item.category || 'Без категории',
        price: item.price || 0,
        quantity: item.quantity || 0,
        status: item.active ? 'active' : 'inactive',
        createdAt: item.created_at || new Date().toISOString(),
      }));
      setProducts(transformedProducts);
    } else if (!isLoading && items.length === 0 && !apiError) {
      // Если API успешно ответил но вернул пустой массив
      setProducts([]);
    }
  }, [items, isLoading, apiError]);

  // Фильтрация товаров
  const filteredProducts = useMemo(() => {
    return products.filter((product) => {
      const matchesSearch =
        product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        product.sku.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesStatus =
        statusFilter === 'all' || product.status === statusFilter;

      const matchesStock =
        stockFilter === 'all' || getStockStatus(product.quantity) === stockFilter;

      return matchesSearch && matchesStatus && matchesStock;
    });
  }, [products, searchQuery, statusFilter, stockFilter]);

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
      {/* Header */}
      <div className="p-6 border-b border-gray-800">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-2xl font-bold text-white flex items-center gap-2">
            <Package className="w-6 h-6 text-blue-400" />
            Каталог товаров
          </h2>
          <div className="flex items-center gap-2">
            <button className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-semibold transition-colors">
              <Download size={20} />
              Экспортировать CSV
            </button>
            <Link
              href="/catalog/import"
              className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-semibold transition-colors"
            >
              <Upload size={20} />
              Импорт
            </Link>
            <Link
              href="/catalog/create"
              className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold transition-colors"
            >
              <Plus size={20} />
              Добавить товар
            </Link>
          </div>
        </div>

        {/* Filters */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500 w-5 h-5" />
            <input
              placeholder="Поиск по названию или артикулу..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500 transition-colors"
              type="text"
            />
          </div>
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-blue-500 transition-colors"
          >
            <option value="all">Все статусы</option>
            <option value="active">Активен</option>
            <option value="inactive">Неактивен</option>
          </select>
          <select
            value={stockFilter}
            onChange={(e) => setStockFilter(e.target.value)}
            className="px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-blue-500 transition-colors"
          >
            <option value="all">Все товары</option>
            <option value="in-stock">Есть в наличии</option>
            <option value="low-stock">Заканчивается</option>
            <option value="out-of-stock">Нет в наличии</option>
          </select>
        </div>
      </div>

      {/* Error */}
      {apiError && (
        <div className="p-4 bg-red-900/20 border-b border-red-800 text-red-300 flex items-center gap-3">
          <AlertCircle size={18} className="flex-shrink-0" />
          <div className="flex-1">
            <p className="font-semibold">Ошибка загрузки товаров</p>
            <p className="text-sm">{apiError}</p>
            <p className="text-xs mt-1 text-red-400">Проверьте подключение к backend. Backend должен быть доступен на /api/catalogs-query</p>
          </div>
          <button
            onClick={() => {
              setApiError(null);
              setProducts([]);
              setTimeout(() => window.location.reload(), 300);
            }}
            className="flex-shrink-0 px-3 py-1 bg-red-800 hover:bg-red-700 rounded text-white text-xs font-semibold transition-colors"
          >
            Повтор
          </button>
        </div>
      )}

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-800/50 border-b border-gray-700">
            <tr>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Название
              </th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Артикул
              </th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Категория
              </th>
              <th className="px-6 py-4 text-right text-sm font-semibold text-gray-300">
                Цена
              </th>
              <th className="px-6 py-4 text-right text-sm font-semibold text-gray-300">
                Остаток
              </th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Статус
              </th>
              <th className="px-6 py-4 text-center text-sm font-semibold text-gray-300">
                Действия
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800">
            {(() => {
              if (isLoading) {
                return (
                  <tr key="loading">
                    <td colSpan={7} className="px-6 py-8 text-center">
                      <div className="flex items-center justify-center gap-2 text-gray-400">
                        <Loader2 size={18} className="animate-spin" />
                        Загрузка товаров...
                      </div>
                    </td>
                  </tr>
                );
              }

              if (filteredProducts.length === 0) {
                return (
                  <tr key="empty">
                    <td colSpan={7} className="px-6 py-8 text-center">
                      <div className="text-gray-400">
                        <Package className="w-12 h-12 mx-auto mb-2 opacity-50" />
                        <p className="font-medium">
                          {searchQuery || statusFilter !== 'all' || stockFilter !== 'all'
                            ? 'Товары не найдены по выбранным критериям'
                            : 'Нет доступных товаров'}
                        </p>
                        {(searchQuery || statusFilter !== 'all' || stockFilter !== 'all') && (
                          <p className="text-sm mt-1">
                            Попробуйте изменить фильтры
                          </p>
                        )}
                      </div>
                    </td>
                  </tr>
                );
              }

              return filteredProducts.map((product, index) => (
                <tr key={`${product.id}-${index}`} className="hover:bg-gray-800/50 transition-colors">
                  <td className="px-6 py-4 text-sm text-white font-medium">
                    {product.name}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-300 font-mono">
                    {product.sku}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-300">
                    {product.category}
                  </td>
                  <td className="px-6 py-4 text-sm text-right text-white font-medium">
                    ₽ {product.price.toLocaleString('ru-RU')}
                  </td>
                  <td className={`px-6 py-4 text-sm text-right font-medium ${getStockBadgeColor(product.quantity)}`}>
                    {product.quantity} шт
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusBadgeColor(product.status)}`}>
                      {getStatusLabel(product.status)}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-center">
                    <div className="flex items-center justify-center gap-2">
                      <Link
                        href={`/catalog/${product.id}/edit`}
                        className="p-2 text-gray-400 hover:text-blue-400 hover:bg-gray-800 rounded transition-colors"
                        title="Редактировать"
                      >
                        <Edit2 size={18} />
                      </Link>
                      <button
                        className="p-2 text-gray-400 hover:text-red-400 hover:bg-gray-800 rounded transition-colors"
                        title="Удалить"
                      >
                        <Trash2 size={18} />
                      </button>
                    </div>
                  </td>
                </tr>
              ));
            })()}
          </tbody>
        </table>
      </div>

      {/* Footer */}
      <div className="px-6 py-4 bg-gray-800/50 border-t border-gray-800">
        <p className="text-sm text-gray-400">
          Показано <span className="text-white font-semibold">{filteredProducts.length}</span> из{' '}
          <span className="text-white font-semibold">{products.length}</span> товаров
        </p>
      </div>
    </div>
  );
}
