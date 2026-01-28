'use client';

import React, { memo, useMemo, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { 
  DollarSign, 
  ShoppingCart, 
  TrendingUp, 
  Zap,
  Activity,
  AlertTriangle,
  LogOut
} from 'lucide-react';
import { useLogoutMutation } from '@/lib/hooks/useAuthApi';
import StatsCard from './StatsCard';
import SalesChart from './SalesChart';
import TerminalsTable from './TerminalsTable';
import InventoryAlerts from './InventoryAlerts';
import StoreHeatmap from './StoreHeatmap';
import RevenueTrendChart from './RevenueTrendChart';
import RealtimeIndicator from './RealtimeIndicator';
import CatalogDashboardWidget from './CatalogDashboardWidget';
import { DashboardStats, SalesData, Terminal, InventoryItem } from '@/types';
import type { User } from '@/types/entities';

interface StoreData {
  id: string;
  name: string;
  revenue: number;
  coordinates: { x: number; y: number };
}

interface RevenueTrendData {
  date: string;
  revenue: number;
  transactions: number;
}

interface DashboardProps {
  user: User | null;
  stats: DashboardStats;
  salesData: SalesData[];
  terminals: Terminal[];
  inventoryItems: InventoryItem[];
  storeLocations: StoreData[];
  revenueTrend: RevenueTrendData[];
}

/**
 * Modern Dashboard Component with comprehensive analytics
 * Displays KPIs, revenue trends, sales analytics, and real-time terminal status
 */
const Dashboard = memo(function Dashboard({
  user,
  stats,
  salesData,
  terminals,
  inventoryItems,
  storeLocations,
  revenueTrend,
}: DashboardProps) {
  const router = useRouter();
  const logoutMutation = useLogoutMutation();

  const handleLogout = useCallback(async () => {
    await logoutMutation.mutateAsync();
    router.push('/auth');
  }, [logoutMutation, router]);

  // Вычисляем дополнительные метрики
  const metrics = useMemo(() => ({
    conversionRate: ((stats.totalTransactions / (stats.totalTransactions * 1.2)) * 100).toFixed(1),
    terminalUsageRate: ((stats.activeTerminals / stats.totalTerminals) * 100).toFixed(0),
    criticalAlerts: inventoryItems.filter(item => item.currentStock < item.minThreshold).length,
  }), [stats, inventoryItems]);

  return (
    <div className="min-h-screen bg-linear-to-br from-gray-950 via-gray-900 to-gray-950 p-6 md:p-8">
      <div className="max-w-7xl mx-auto">
        
        {/* Header Section */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4 flex-1">
              <div className="w-12 h-12 bg-linear-to-br from-blue-500 to-emerald-500 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/30">
                <Activity className="w-6 h-6 text-white" />
              </div>
              <div className="flex-1">
                <h1 className="text-3xl md:text-4xl font-bold text-white">
                  {user ? `Добро пожаловать, ${user.name}!` : 'Дашборд'}
                </h1>
                <div className="flex items-center gap-4 mt-2">
                  <p className="text-gray-400 text-sm">Мониторинг сети кассовых систем в реальном времени</p>
                  {user && (
                    <div className="px-3 py-1 bg-blue-900/30 border border-blue-700 rounded-lg">
                      <p className="text-xs text-blue-300">
                        <span className="font-semibold">Роль:</span> {user.role}
                      </p>
                    </div>
                  )}
                </div>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <div className="hidden md:flex items-center gap-2 px-4 py-2 bg-gray-800/50 rounded-lg border border-gray-700">
                <div className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse"></div>
                <span className="text-gray-300 text-sm">Всё работает нормально</span>
              </div>
              <button
                onClick={handleLogout}
                disabled={logoutMutation.isPending}
                className="flex items-center gap-2 px-4 py-2 rounded-lg transition text-sm font-semibold text-white bg-red-600/80 hover:bg-red-600 active:bg-red-700 border border-red-500 shadow-lg shadow-red-600/30 disabled:opacity-70 disabled:cursor-not-allowed"
                title="Выход из аккаунта"
              >
                <LogOut size={18} />
                <span className="hidden md:inline">{logoutMutation.isPending ? 'Выход...' : 'Выход'}</span>
              </button>
            </div>
          </div>
        </div>

        {/* KPI Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <StatsCard
            title="Общая выручка"
            value={`${stats.totalRevenue.toLocaleString('ru-RU')} ₽`}
            icon={<DollarSign className="w-8 h-8 text-blue-500" />}
            trend={{ value: 12.5, isPositive: true }}
            bgGradient="from-blue-600/10 to-blue-600/5"
          />
          <StatsCard
            title="Активные кассы"
            value={`${stats.activeTerminals}/${stats.totalTerminals}`}
            subtitle={`${metrics.terminalUsageRate}% загруженности`}
            icon={<Zap className="w-8 h-8 text-emerald-500" />}
            trend={{ value: 5, isPositive: true }}
            bgGradient="from-emerald-600/10 to-emerald-600/5"
          />
          <StatsCard
            title="Средний чек"
            value={`${stats.avgCheck.toFixed(0)} ₽`}
            icon={<ShoppingCart className="w-8 h-8 text-purple-500" />}
            trend={{ value: 8.3, isPositive: true }}
            bgGradient="from-purple-600/10 to-purple-600/5"
          />
          <StatsCard
            title="Критические алерты"
            value={metrics.criticalAlerts}
            icon={<AlertTriangle className="w-8 h-8 text-red-500" />}
            trend={{ value: 0, isPositive: false }}
            bgGradient="from-red-600/10 to-red-600/5"
          />
        </div>

        {/* Charts Section */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
          {/* Left Column - Main Charts */}
          <div className="lg:col-span-2 space-y-8">
            {/* Revenue Trend Chart */}
            <div className="rounded-lg border border-gray-800 bg-gray-900/50 backdrop-blur-sm p-6 hover:border-gray-700 transition">
              <h2 className="text-xl font-bold text-white mb-6">Тренд выручки</h2>
              <RevenueTrendChart data={revenueTrend} />
            </div>

            {/* Sales Chart */}
            <div className="rounded-lg border border-gray-800 bg-gray-900/50 backdrop-blur-sm p-6 hover:border-gray-700 transition">
              <h2 className="text-xl font-bold text-white mb-6">Аналитика продаж</h2>
              <SalesChart data={salesData} />
            </div>

            {/* Store Heatmap */}
            <div className="rounded-lg border border-gray-800 bg-gray-900/50 backdrop-blur-sm p-6 hover:border-gray-700 transition">
              <h2 className="text-xl font-bold text-white mb-6">Карта выручки по точкам</h2>
              <StoreHeatmap stores={storeLocations} />
            </div>
          </div>

          {/* Right Column - Sidebar */}
          <div className="space-y-8">
            {/* Alerts Section */}
            <div>
              <InventoryAlerts items={inventoryItems} />
            </div>

            {/* Quick Stats */}
            <div className="rounded-lg border border-gray-800 bg-gray-900/50 backdrop-blur-sm p-6 hover:border-gray-700 transition">
              <h3 className="text-lg font-bold text-white mb-4">Сводка</h3>
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Времени отслеживания</span>
                  <span className="text-white font-semibold">24/7</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Последнее обновление</span>
                  <span className="text-white font-semibold">2 сек назад</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Доступность системы</span>
                  <span className="text-emerald-400 font-semibold">99.95%</span>
                </div>
                <div className="h-px bg-gray-800"></div>
                <div className="text-xs text-gray-500">
                  Данные обновляются в реальном времени с интервалом 2 секунды
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Catalog Section */}
        <div className="rounded-lg border border-gray-800 bg-gray-900/50 backdrop-blur-sm p-6 hover:border-gray-700 transition">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-bold text-white">Каталог товаров</h2>
            <a href="/catalog" className="text-sm text-blue-400 hover:text-blue-300">
              Перейти в каталог →
            </a>
          </div>
          <CatalogDashboardWidget />
        </div>

        {/* Terminals Table Section */}
        <div className="rounded-lg border border-gray-800 bg-gray-900/50 backdrop-blur-sm p-6 hover:border-gray-700 transition">
          <h2 className="text-xl font-bold text-white mb-6">Статус кассовых аппаратов</h2>
          <TerminalsTable terminals={terminals} />
        </div>
      </div>
    </div>
  );
});

Dashboard.displayName = 'Dashboard';

export default Dashboard;
