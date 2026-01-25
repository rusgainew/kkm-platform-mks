'use client';

import React from 'react';
import { Search, Filter, X } from 'lucide-react';

interface CatalogFilterProps {
  filters: any;
  onFilterChange: (filters: any) => void;
}

export const CatalogFilter: React.FC<CatalogFilterProps> = ({
  filters,
  onFilterChange,
}) => {
  const handleSearchChange = (value: string) => {
    onFilterChange({ ...filters, search: value });
  };

  const handleCategoryChange = (value: string) => {
    onFilterChange({ ...filters, category: value === 'all' ? undefined : value });
  };

  const handleActiveChange = (value: string) => {
    onFilterChange({
      ...filters,
      active: value === 'all' ? undefined : value === 'active',
    });
  };

  const handleClearFilters = () => {
    onFilterChange({});
  };

  const hasFilters = filters.search || filters.category || filters.active !== undefined;

  return (
    <div className="space-y-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
      <div className="flex items-center gap-2 mb-2">
        <Filter size={16} className="text-gray-600" />
        <h3 className="font-semibold text-sm">Фильтры</h3>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="relative">
          <Search size={16} className="absolute left-3 top-2.5 text-gray-400" />
          <input
            type="text"
            placeholder="Поиск по названию..."
            value={filters.search || ''}
            onChange={(e) => handleSearchChange(e.target.value)}
            className="pl-8 w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
          />
        </div>

        <select
          value={filters.category || 'all'}
          onChange={(e) => handleCategoryChange(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-md text-sm"
        >
          <option value="all">Все категории</option>
          <option value="electronics">Электроника</option>
          <option value="furniture">Мебель</option>
          <option value="office">Офисные товары</option>
        </select>

        <select
          value={filters.active === undefined ? 'all' : filters.active ? 'active' : 'inactive'}
          onChange={(e) => handleActiveChange(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-md text-sm"
        >
          <option value="all">Все</option>
          <option value="active">Активные</option>
          <option value="inactive">Неактивные</option>
        </select>

        {hasFilters && (
          <button
            onClick={handleClearFilters}
            className="px-3 py-2 border border-gray-300 rounded-md text-sm bg-white hover:bg-gray-50 flex items-center gap-2"
          >
            <X size={14} />
            Очистить
          </button>
        )}
      </div>
    </div>
  );
};
