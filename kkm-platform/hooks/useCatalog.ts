import { useCatalogStore } from '@/store';
import { useCallback } from 'react';

interface CatalogFilter {
  category?: string;
  search?: string;
  active?: boolean;
}

/**
 * useCatalog Hook - управление каталогом товаров
 * 
 * Примеры использования:
 * const { items, fetchCatalog, createItem, updateItem } = useCatalog();
 * 
 * // Fetch with filters
 * await fetchCatalog({ category: 'Electronics', active: true });
 * 
 * // Create
 * await createItem({ name, sku, price, costPrice, category });
 * 
 * // Update
 * await updateItem(itemId, { stock: 100 });
 * 
 * // Get statistics
 * const { totalInventoryValue, lowStockItems } = useCatalog();
 */
export const useCatalog = () => {
  const {
    items,
    filteredItems,
    categories,
    totalInventoryValue,
    lowStockItems,
    loading,
    error,
    filters,
    pagination,
    fetchCatalog,
    fetchItemById,
    createItem,
    updateItem,
    deleteItem,
    setFilters,
    setPagination,
    clearError,
    getItemById,
    getItemsByCategory,
    calculateProfit,
  } = useCatalogStore();

  const handleFetchCatalog = useCallback(
    async (newFilters?: CatalogFilter) => {
      try {
        await fetchCatalog(newFilters);
        return true;
      } catch (error) {
        return false;
      }
    },
    [fetchCatalog]
  );

  const handleCreateItem = useCallback(
    async (data: any) => {
      try {
        await createItem(data);
        return true;
      } catch (error) {
        return false;
      }
    },
    [createItem]
  );

  const handleUpdateItem = useCallback(
    async (id: string, data: any) => {
      try {
        await updateItem(id, data);
        return true;
      } catch (error) {
        return false;
      }
    },
    [updateItem]
  );

  const handleDeleteItem = useCallback(
    async (id: string) => {
      try {
        await deleteItem(id);
        return true;
      } catch (error) {
        return false;
      }
    },
    [deleteItem]
  );

  return {
    items,
    filteredItems,
    categories,
    totalInventoryValue,
    lowStockItems,
    isLoading: loading,
    error,
    filters,
    pagination,
    fetchCatalog: handleFetchCatalog,
    fetchItemById,
    createItem: handleCreateItem,
    updateItem: handleUpdateItem,
    deleteItem: handleDeleteItem,
    setFilters,
    setPagination,
    clearError,
    getItemById,
    getItemsByCategory,
    calculateProfit,
  };
};

export default useCatalog;
