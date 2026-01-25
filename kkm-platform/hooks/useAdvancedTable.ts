import { useState, useCallback, useMemo } from 'react';

export interface UseAdvancedTableOptions<T> {
  data: T[];
  filterFn?: (item: T, filters: Record<string, any>) => boolean;
  searchFields?: (keyof T)[];
}

export function useAdvancedTable<T extends { id: string }>({
  data,
  filterFn,
  searchFields = [],
}: UseAdvancedTableOptions<T>) {
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [filters, setFilters] = useState<Record<string, any>>({});
  const [searchQuery, setSearchQuery] = useState('');

  // Фильтрация данных
  const filteredData = useMemo(() => {
    let result = data;

    // Поиск
    if (searchQuery && searchFields.length > 0) {
      const query = searchQuery.toLowerCase();
      result = result.filter((item) =>
        searchFields.some((field) => {
          const value = item[field];
          return String(value).toLowerCase().includes(query);
        })
      );
    }

    // Пользовательская фильтрация
    if (filterFn) {
      result = result.filter((item) => filterFn(item, filters));
    }

    return result;
  }, [data, filterFn, filters, searchQuery, searchFields]);

  // Выбор элемента
  const toggleSelection = useCallback(
    (id: string) => {
      setSelectedIds((prev) => {
        const newSet = new Set(prev);
        if (newSet.has(id)) {
          newSet.delete(id);
        } else {
          newSet.add(id);
        }
        return newSet;
      });
    },
    []
  );

  // Выбрать все видимые
  const selectAll = useCallback(() => {
    setSelectedIds(new Set(filteredData.map((item) => item.id)));
  }, [filteredData]);

  // Очистить выбор
  const clearSelection = useCallback(() => {
    setSelectedIds(new Set());
  }, []);

  // Выбрать страницу
  const selectPage = useCallback((pageItems: T[]) => {
    setSelectedIds(new Set(pageItems.map((item) => item.id)));
  }, []);

  // Инвертировать выбор
  const invertSelection = useCallback(() => {
    const newSet = new Set<string>();
    filteredData.forEach((item) => {
      if (!selectedIds.has(item.id)) {
        newSet.add(item.id);
      }
    });
    setSelectedIds(newSet);
  }, [filteredData, selectedIds]);

  // Получить выбранные элементы
  const getSelectedItems = useCallback(() => {
    return filteredData.filter((item) => selectedIds.has(item.id));
  }, [filteredData, selectedIds]);

  return {
    // Данные
    filteredData,
    selectedIds: Array.from(selectedIds),
    selectedCount: selectedIds.size,
    isAllSelected: selectedIds.size === filteredData.length && filteredData.length > 0,

    // Фильтры
    filters,
    setFilters,
    searchQuery,
    setSearchQuery,

    // Функции выбора
    toggleSelection,
    selectAll,
    clearSelection,
    selectPage,
    invertSelection,
    getSelectedItems,
  };
}
