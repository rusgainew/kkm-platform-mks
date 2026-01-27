"use client";

import { useQuery } from "@tanstack/react-query";
import { analyticsAPI } from "@/lib/api/analytics";
import type { AnalyticsFilters } from "@/types/analytics";

/**
 * Хук для получения всех данных Dashboard
 */
export function useDashboardData(filters: AnalyticsFilters) {
  return useQuery({
    queryKey: ["dashboard-data", filters],
    queryFn: () => analyticsAPI.getAnalyticsData(filters),
    staleTime: 5 * 60 * 1000, // 5 минут
    refetchOnWindowFocus: true,
  });
}

/**
 * Хук для получения статистики Dashboard
 */
export function useDashboardStats(filters: AnalyticsFilters) {
  return useQuery({
    queryKey: ["dashboard-stats", filters],
    queryFn: () => analyticsAPI.getStats(filters),
    staleTime: 5 * 60 * 1000,
  });
}

/**
 * Хук для получения данных графика продаж
 */
export function useSalesChart(filters: AnalyticsFilters) {
  return useQuery({
    queryKey: ["sales-chart", filters],
    queryFn: () => analyticsAPI.getSalesChart(filters),
    staleTime: 5 * 60 * 1000,
  });
}

/**
 * Хук для получения статистики по статусам накладных
 */
export function useInvoiceStatusStats(filters: AnalyticsFilters) {
  return useQuery({
    queryKey: ["invoice-status-stats", filters],
    queryFn: () => analyticsAPI.getStatusStats(filters),
    staleTime: 5 * 60 * 1000,
  });
}

/**
 * Хук для получения статистики по типам операций
 */
export function useOperationTypeStats(filters: AnalyticsFilters) {
  return useQuery({
    queryKey: ["operation-type-stats", filters],
    queryFn: () => analyticsAPI.getOperationTypeStats(filters),
    staleTime: 5 * 60 * 1000,
  });
}

/**
 * Хук для получения топ контрагентов
 */
export function useTopContractors(
  filters: AnalyticsFilters,
  limit: number = 10,
) {
  return useQuery({
    queryKey: ["top-contractors", filters, limit],
    queryFn: () => analyticsAPI.getTopContractors(filters, limit),
    staleTime: 5 * 60 * 1000,
  });
}

/**
 * Хук для получения выручки по месяцам
 */
export function useRevenueByMonth(filters: AnalyticsFilters) {
  return useQuery({
    queryKey: ["revenue-by-month", filters],
    queryFn: () => analyticsAPI.getRevenueByMonth(filters),
    staleTime: 5 * 60 * 1000,
  });
}
