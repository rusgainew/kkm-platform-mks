'use client';

import { useState, useCallback } from 'react';
import { X, ChevronDown } from 'lucide-react';

export interface FilterConfig {
  id: string;
  label: string;
  type: 'text' | 'select' | 'date' | 'dateRange' | 'checkbox';
  options?: Array<{ value: string; label: string }>;
  value: any;
  onChange: (value: any) => void;
  placeholder?: string;
}

interface AdvancedFilterProps {
  filters: FilterConfig[];
  onApply?: () => void;
  onReset?: () => void;
  onSave?: (name: string) => void;
}

export function AdvancedFilter({
  filters,
  onApply,
  onReset,
  onSave,
}: AdvancedFilterProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [savedFilters, setSavedFilters] = useState<Array<{ name: string; filters: Record<string, any> }>>([]);
  const [filterName, setFilterName] = useState('');

  const handleSaveFilter = useCallback(() => {
    if (!filterName.trim()) return;

    const currentFilters = filters.reduce(
      (acc, filter) => {
        acc[filter.id] = filter.value;
        return acc;
      },
      {} as Record<string, any>
    );

    setSavedFilters([...savedFilters, { name: filterName, filters: currentFilters }]);
    setFilterName('');
    localStorage.setItem(
      'savedFilters',
      JSON.stringify([...savedFilters, { name: filterName, filters: currentFilters }])
    );
  }, [filterName, filters, savedFilters]);

  const handleLoadFilter = useCallback((saved: any) => {
    filters.forEach((filter) => {
      if (saved.filters[filter.id] !== undefined) {
        filter.onChange(saved.filters[filter.id]);
      }
    });
  }, [filters]);

  const activeFiltersCount = filters.filter((f) => f.value).length;

  return (
    <div className="space-y-4">
      {/* Кнопка фильтра */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
      >
        <ChevronDown className={`w-4 h-4 transition ${isOpen ? 'rotate-180' : ''}`} />
        Фильтры
        {activeFiltersCount > 0 && (
          <span className="ml-2 px-2 py-1 bg-blue-600 text-white text-xs rounded-full">
            {activeFiltersCount}
          </span>
        )}
      </button>

      {/* Панель фильтров */}
      {isOpen && (
        <div className="bg-white border border-gray-200 rounded-lg p-4 space-y-4">
          {/* Фильтры */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {filters.map((filter) => (
              <div key={filter.id} className="space-y-2">
                <label className="block text-sm font-medium text-gray-700">
                  {filter.label}
                </label>

                {filter.type === 'text' && (
                  <input
                    type="text"
                    placeholder={filter.placeholder}
                    value={filter.value || ''}
                    onChange={(e) => filter.onChange(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                )}

                {filter.type === 'select' && (
                  <select
                    value={filter.value || ''}
                    onChange={(e) => filter.onChange(e.target.value || null)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="">Все</option>
                    {filter.options?.map((opt) => (
                      <option key={opt.value} value={opt.value}>
                        {opt.label}
                      </option>
                    ))}
                  </select>
                )}

                {filter.type === 'date' && (
                  <input
                    type="date"
                    value={filter.value || ''}
                    onChange={(e) => filter.onChange(e.target.value || null)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                )}

                {filter.type === 'dateRange' && (
                  <div className="flex gap-2">
                    <input
                      type="date"
                      value={filter.value?.from || ''}
                      onChange={(e) =>
                        filter.onChange({
                          ...filter.value,
                          from: e.target.value || null,
                        })
                      }
                      placeholder="От"
                      className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <input
                      type="date"
                      value={filter.value?.to || ''}
                      onChange={(e) =>
                        filter.onChange({
                          ...filter.value,
                          to: e.target.value || null,
                        })
                      }
                      placeholder="До"
                      className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                )}

                {filter.type === 'checkbox' && filter.options && (
                  <div className="space-y-2">
                    {filter.options.map((opt) => (
                      <label key={opt.value} className="flex items-center gap-2">
                        <input
                          type="checkbox"
                          checked={filter.value?.includes(opt.value) || false}
                          onChange={(e) => {
                            const newValue = filter.value || [];
                            if (e.target.checked) {
                              filter.onChange([...newValue, opt.value]);
                            } else {
                              filter.onChange(newValue.filter((v: string) => v !== opt.value));
                            }
                          }}
                          className="rounded border-gray-300"
                        />
                        <span className="text-sm">{opt.label}</span>
                      </label>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>

          {/* Кнопки действий */}
          <div className="flex gap-2 pt-4 border-t">
            {onApply && (
              <button
                onClick={() => {
                  onApply();
                  setIsOpen(false);
                }}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                Применить
              </button>
            )}

            {onReset && (
              <button
                onClick={onReset}
                className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition"
              >
                Сбросить
              </button>
            )}

            {/* Сохранить фильтр */}
            <div className="ml-auto flex gap-2">
              <input
                type="text"
                placeholder="Название фильтра"
                value={filterName}
                onChange={(e) => setFilterName(e.target.value)}
                className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                onClick={handleSaveFilter}
                disabled={!filterName.trim()}
                className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 transition"
              >
                Сохранить
              </button>
            </div>
          </div>

          {/* Сохраненные фильтры */}
          {savedFilters.length > 0 && (
            <div className="pt-4 border-t space-y-2">
              <p className="text-sm font-medium text-gray-700">Сохраненные фильтры</p>
              <div className="flex flex-wrap gap-2">
                {savedFilters.map((saved, index) => (
                  <button
                    key={index}
                    onClick={() => handleLoadFilter(saved)}
                    className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm hover:bg-gray-200 transition flex items-center gap-2"
                  >
                    {saved.name}
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
