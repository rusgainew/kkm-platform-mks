import { useQuery } from '@tanstack/react-query';
import { useEffect } from 'react';
import { useAuthStore } from '@/store/authStore';
import { useDashboardStore } from '@/store/dashboardStore';
import {
  fetchDashboardStats,
  fetchTerminals,
  fetchSalesData,
  fetchInventory,
} from '@/lib/api/dashboard';

export function useDashboardData() {
  const user = useAuthStore((state) => state.user);
  const hasPermission = useAuthStore((state) => state.hasPermission);
  const selectedStoreId = useDashboardStore((state) => state.selectedStoreId);
  const setStats = useDashboardStore((state) => state.setStats);
  const setTerminals = useDashboardStore((state) => state.setTerminals);
  const setSalesData = useDashboardStore((state) => state.setSalesData);
  const setInventoryItems = useDashboardStore((state) => state.setInventoryItems);
  const markLastUpdate = useDashboardStore((state) => state.markLastUpdate);

  // Определяем storeId на основе роли пользователя
  const storeId = hasPermission('view_all_stores') ? selectedStoreId || undefined : user?.storeId;

  // Real-time polling для статистики (каждые 10 секунд)
  const statsQuery = useQuery({
    queryKey: ['dashboard-stats', storeId],
    queryFn: () => fetchDashboardStats(storeId),
    refetchInterval: 10000,
    enabled: !!user,
  });

  // Real-time polling для терминалов (каждые 15 секунд)
  const terminalsQuery = useQuery({
    queryKey: ['terminals', storeId],
    queryFn: () => fetchTerminals(storeId),
    refetchInterval: 15000,
    enabled: !!user,
  });

  // Real-time polling для графика продаж (каждые 30 секунд)
  const salesQuery = useQuery({
    queryKey: ['sales-data', storeId],
    queryFn: () => fetchSalesData(storeId),
    refetchInterval: 30000,
    enabled: !!user,
  });

  // Инвентарь обновляется реже (каждые 30 секунд)
  const inventoryQuery = useQuery({
    queryKey: ['inventory', storeId],
    queryFn: () => fetchInventory(storeId),
    refetchInterval: 30000,
    enabled: !!user && hasPermission('manage_inventory'),
  });

  // Обновление store при получении новых данных
  useEffect(() => {
    if (statsQuery.data) {
      setStats(statsQuery.data);
      markLastUpdate();
    }
  }, [statsQuery.data]);

  useEffect(() => {
    if (terminalsQuery.data) setTerminals(terminalsQuery.data);
  }, [terminalsQuery.data]);

  useEffect(() => {
    if (salesQuery.data) setSalesData(salesQuery.data);
  }, [salesQuery.data]);

  useEffect(() => {
    if (inventoryQuery.data) setInventoryItems(inventoryQuery.data);
  }, [inventoryQuery.data]);

  return {
    stats: statsQuery.data,
    terminals: terminalsQuery.data,
    salesData: salesQuery.data,
    inventoryItems: inventoryQuery.data,
    isLoading:
      statsQuery.isLoading ||
      terminalsQuery.isLoading ||
      salesQuery.isLoading ||
      inventoryQuery.isLoading,
    isError:
      statsQuery.isError || terminalsQuery.isError || salesQuery.isError || inventoryQuery.isError,
  };
}
