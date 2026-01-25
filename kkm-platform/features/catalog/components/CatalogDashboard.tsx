'use client';

import React, { useEffect, useState } from 'react';
import { useCatalog } from '@/hooks/useCatalog';
import { Package, TrendingUp, AlertCircle, ShoppingCart } from 'lucide-react';
import CatalogList from './CatalogList';

export const CatalogDashboard: React.FC = () => {
  const {
    items,
    filteredItems,
    isLoading,
    fetchCatalog,
    totalInventoryValue,
    lowStockItems,
  } = useCatalog();

  useEffect(() => {
    fetchCatalog();
  }, [fetchCatalog]);

  const activeItems = items.filter((item: any) => item.active).length;
  const inactiveItems = items.filter((item: any) => !item.active).length;

  const stats = [
    {
      title: 'Всего товаров',
      value: items.length,
      icon: Package,
      color: 'text-blue-500',
      bgColor: 'bg-blue-50',
    },
    {
      title: 'Активные',
      value: activeItems,
      icon: TrendingUp,
      color: 'text-green-500',
      bgColor: 'bg-green-50',
    },
    {
      title: 'Стоимость склада',
      value: `$${totalInventoryValue?.toFixed(2) || '0.00'}`,
      icon: ShoppingCart,
      color: 'text-purple-500',
      bgColor: 'bg-purple-50',
    },
    {
      title: 'Низкий запас',
      value: lowStockItems?.length || 0,
      icon: AlertCircle,
      color: 'text-red-500',
      bgColor: 'bg-red-50',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, index) => {
          const Icon = stat.icon;
          return (
            <div key={index} className="bg-white border border-gray-200 rounded-lg shadow-sm">
              <div className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600">{stat.title}</p>
                    <p className="text-2xl font-bold mt-2">{stat.value}</p>
                  </div>
                  <div className={`${stat.bgColor} p-3 rounded-lg`}>
                    <Icon className={`w-6 h-6 ${stat.color}`} />
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Low Stock Alert */}
      {lowStockItems && lowStockItems.length > 0 && (
        <div className="border-2 border-orange-200 bg-orange-50 rounded-lg">
          <div className="p-6">
            <h3 className="text-lg font-semibold text-orange-900 flex items-center gap-2">
              <AlertCircle size={20} />
              Товары с низким запасом
            </h3>
            <p className="text-sm text-orange-700 mt-1">
              {lowStockItems.length} товаров требует переордера
            </p>
            <div className="space-y-2 mt-4">
              {lowStockItems.slice(0, 5).map((item: any) => (
                <div key={item.id} className="flex items-center justify-between p-2 bg-white rounded">
                  <span className="font-medium">{item.name}</span>
                  <span className="text-sm text-gray-600">
                    {item.stock} / {item.reorderLevel}
                  </span>
                </div>
              ))}
              {lowStockItems.length > 5 && (
                <p className="text-sm text-gray-600 p-2">
                  и еще {lowStockItems.length - 5} товаров...
                </p>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Catalog List */}
      <CatalogList />
    </div>
  );
};
