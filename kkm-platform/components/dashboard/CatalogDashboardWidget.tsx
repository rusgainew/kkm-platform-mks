'use client';

import React, { useEffect } from 'react';
import { useCatalog } from '@/hooks/useCatalog';
import { useCatalogHealth } from '@/hooks/useCatalogHealth';
import { Package, AlertCircle, Loader, CheckCircle, XCircle } from 'lucide-react';

export default function CatalogDashboardWidget() {
  const { items, fetchCatalog, lowStockItems, isLoading, error } = useCatalog();
  const { health } = useCatalogHealth();

  useEffect(() => {
    fetchCatalog().catch(err => {
      console.error('Ошибка загрузки каталога в виджете:', err);
    });
  }, [fetchCatalog]);

  const activeItems = items.filter((item: any) => item.active).length;

  if (error) {
    return (
      <div className="space-y-3">
        <div className="bg-red-900/20 border border-red-700/50 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <AlertCircle size={18} className="text-red-400 flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <p className="text-sm font-medium text-red-300">Ошибка при загрузке каталога</p>
              <p className="text-xs text-red-300/80 mt-1 font-mono break-words">{error}</p>
            </div>
          </div>
        </div>
        
        {health?.status === 'error' && (
          <div className="bg-orange-900/20 border border-orange-700/50 rounded-lg p-4">
            <div className="flex items-start gap-3">
              <XCircle size={18} className="text-orange-400 flex-shrink-0 mt-0.5" />
              <div className="flex-1">
                <p className="text-sm font-medium text-orange-300">Бэкэнд API недоступен</p>
                <p className="text-xs text-orange-300/80 mt-1">{health.message}</p>
                {health.suggestion && (
                  <p className="text-xs text-orange-300/70 mt-2 italic">{health.suggestion}</p>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    );
  }

  if (isLoading && items.length === 0) {
    return (
      <div className="flex items-center justify-center p-8">
        <Loader size={20} className="text-blue-400 animate-spin mr-2" />
        <p className="text-gray-400">Загрузка каталога...</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Stats */}
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-gray-800 rounded-lg p-4">
          <p className="text-sm text-gray-400">Всего товаров</p>
          <p className="text-2xl font-bold text-white">{items.length}</p>
        </div>
        <div className="bg-gray-800 rounded-lg p-4">
          <p className="text-sm text-gray-400">Активные</p>
          <p className="text-2xl font-bold text-green-400">{activeItems}</p>
        </div>
        <div className="bg-gray-800 rounded-lg p-4">
          <p className="text-sm text-gray-400">Низкий запас</p>
          <p className="text-2xl font-bold text-red-400">{lowStockItems?.length || 0}</p>
        </div>
      </div>

      {/* Low Stock Items */}
      {lowStockItems && lowStockItems.length > 0 && (
        <div className="bg-red-900/20 border border-red-700/50 rounded-lg p-4">
          <div className="flex items-center gap-2 mb-3">
            <AlertCircle size={16} className="text-red-400" />
            <p className="text-sm font-semibold text-red-400">Требуется переордер</p>
          </div>
          <div className="space-y-2">
            {lowStockItems.slice(0, 3).map((item: any) => (
              <div key={item.id} className="flex items-center justify-between text-sm">
                <span className="text-gray-300">{item.name}</span>
                <span className="text-red-400">{item.stock} / {item.reorderLevel}</span>
              </div>
            ))}
            {lowStockItems.length > 3 && (
              <p className="text-xs text-gray-400 pt-2">
                и еще {lowStockItems.length - 3} товаров...
              </p>
            )}
          </div>
        </div>
      )}

      {/* Recent Items */}
      {items.length > 0 && (
        <div className="bg-gray-800/50 rounded-lg p-4">
          <h3 className="text-sm font-semibold text-white mb-3 flex items-center gap-2">
            <Package size={14} />
            Последние товары
          </h3>
          <div className="space-y-2">
            {items.slice(-3).map((item: any) => (
              <div key={item.id} className="flex items-center justify-between text-sm">
                <span className="text-gray-300">{item.name}</span>
                <span className="text-gray-500 text-xs">{item.sku}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Empty State */}
      {items.length === 0 && !isLoading && (
        <div className="bg-gray-800/50 rounded-lg p-4 text-center">
          <Package size={32} className="text-gray-500 mx-auto mb-2 opacity-50" />
          <p className="text-gray-400 text-sm">Каталог товаров пуст</p>
        </div>
      )}
    </div>
  );
}
