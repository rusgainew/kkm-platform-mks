/* eslint-disable @typescript-eslint/no-explicit-any */
"use client";

import { useState } from "react";
import { TrendingUp, FileText, DollarSign, Users } from "lucide-react";
import { PeriodFilter } from "@/features/analytics/components/PeriodFilter";
import { StatCard } from "@/features/analytics/components/StatCard";
import { LineChart } from "@/features/analytics/components/LineChart";
import { BarChart } from "@/features/analytics/components/BarChart";
import { PieChart } from "@/features/analytics/components/PieChart";
import {
  useDashboardData,
  useSalesChart,
  useInvoiceStatusStats,
  useOperationTypeStats,
  useTopContractors,
  useRevenueByMonth,
} from "@/features/analytics/hooks/useDashboardData";
import type { TimePeriod, DateRange, AnalyticsFilters } from "@/types/analytics";

/**
 * Главная страница Dashboard с аналитикой
 */
export default function DashboardPage() {
  const [selectedPeriod, setSelectedPeriod] = useState<TimePeriod>("month");
  const [dateRange, setDateRange] = useState<DateRange | undefined>();

  // Формируем фильтры
  const filters: AnalyticsFilters = {
    period: selectedPeriod,
    ...(dateRange && { startDate: dateRange.startDate, endDate: dateRange.endDate }),
  };

  // Загружаем данные
  const { data: dashboardData, isLoading: isDashboardLoading } = useDashboardData(filters);
  const { data: salesChart, isLoading: isSalesLoading } = useSalesChart(filters);
  const { data: statusStats, isLoading: isStatusLoading } = useInvoiceStatusStats(filters);
  const { data: operationStats, isLoading: isOperationLoading } = useOperationTypeStats(filters);
  const { data: topContractors, isLoading: isContractorsLoading } = useTopContractors(filters, 5);
  const { data: revenueByMonth, isLoading: isRevenueLoading } = useRevenueByMonth(filters);

  const handlePeriodChange = (period: TimePeriod, range?: DateRange) => {
    setSelectedPeriod(period);
    setDateRange(range);
  };

  const isLoading =
    isDashboardLoading ||
    isSalesLoading ||
    isStatusLoading ||
    isOperationLoading ||
    isContractorsLoading ||
    isRevenueLoading;

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2">
            Dashboard
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Обзор ключевых метрик и аналитика продаж
          </p>
        </div>

        {/* Period filter */}
        <PeriodFilter
          selectedPeriod={selectedPeriod}
          onPeriodChange={handlePeriodChange}
          className="mb-6"
        />

        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        ) : (
          <>
            {/* Key metrics */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
              <StatCard
                label="Общая выручка"
                value={dashboardData?.stats.totalRevenue || 0}
                change={dashboardData?.stats.revenueChange}
                changeType={
                  (dashboardData?.stats.revenueChange || 0) > 0 ? "increase" : "decrease"
                }
                format="currency"
                icon={<DollarSign className="w-6 h-6" />}
              />
              <StatCard
                label="Количество накладных"
                value={dashboardData?.stats.totalInvoices || 0}
                change={dashboardData?.stats.invoiceCountChange}
                changeType={
                  (dashboardData?.stats.invoiceCountChange || 0) > 0 ? "increase" : "decrease"
                }
                format="number"
                icon={<FileText className="w-6 h-6" />}
              />
              <StatCard
                label="Средний чек"
                value={dashboardData?.stats.averageInvoiceAmount || 0}
                change={dashboardData?.stats.averageAmountChange}
                changeType={
                  (dashboardData?.stats.averageAmountChange || 0) > 0 ? "increase" : "decrease"
                }
                format="currency"
                icon={<TrendingUp className="w-6 h-6" />}
              />
              <StatCard
                label="Активные контрагенты"
                value={dashboardData?.stats.activeContractors || 0}
                change={dashboardData?.stats.contractorsChange}
                changeType={
                  (dashboardData?.stats.contractorsChange || 0) > 0 ? "increase" : "decrease"
                }
                format="number"
                icon={<Users className="w-6 h-6" />}
              />
            </div>

            {/* Charts row 1 */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
              <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                <LineChart
                  data={salesChart?.data || []}
                  xAxisKey="date"
                  yAxisKey="amount"
                  lineColor="#3b82f6"
                  title="Динамика продаж"
                  height={300}
                  currency="KZT"
                />
              </div>

              <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                <BarChart
                  data={revenueByMonth?.data || []}
                  xAxisKey="month"
                  yAxisKey="revenue"
                  barColor="#10b981"
                  title="Выручка по месяцам"
                  height={300}
                  currency="KZT"
                />
              </div>
            </div>

            {/* Charts row 2 */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
              <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                <PieChart
                  data={statusStats?.data || []}
                  title="Статусы накладных"
                  height={300}
                />
              </div>

              <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                <PieChart
                  data={operationStats?.data || []}
                  title="Типы операций"
                  height={300}
                />
              </div>

              <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Топ контрагенты
                </h3>
                <div className="space-y-3">
                  {(topContractors || []).map((contractor: any, index: number) => (
                    <div
                      key={contractor.contractorId}
                      className="flex items-center justify-between"
                    >
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center text-blue-600 dark:text-blue-300 font-semibold text-sm">
                          {index + 1}
                        </div>
                        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                          {contractor.contractorName}
                        </span>
                      </div>
                      <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                        {contractor.totalAmount.toLocaleString("ru-KZ")} ₸
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
