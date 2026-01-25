'use client';

import { useAuthStore } from '@/store/authStore';
import { BarChart3, Users, ShoppingCart, TrendingUp, AlertCircle } from 'lucide-react';
import Link from 'next/link';
import ProtectedRoute from '@/components/auth/ProtectedRoute';
import SPADashboardLayout from '@/components/layout/SPADashboardLayout';
import { StatCard, ManagerAlert, PerformanceMetric } from '@/components/manager';

export default function ManagerDashboardPage() {
  const user = useAuthStore((state) => state.user);

  // Проверка, что это менеджер
  if (user && user.role !== 'store_manager') {
    return (
      <ProtectedRoute requiredPermission="view_reports">
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <AlertCircle className="w-12 h-12 text-red-500 mx-auto mb-4" />
            <h1 className="text-2xl font-bold text-white mb-2">Доступ запрещен</h1>
            <p className="text-gray-400 mb-6">Эта страница доступна только для менеджеров</p>
            <Link href="/dashboard" className="text-blue-500 hover:text-blue-400">
              Вернуться в панель управления
            </Link>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  return (
    <ProtectedRoute requiredPermission="view_own_store">
      <SPADashboardLayout currentPage="manager-dashboard">
        <div className="space-y-6 p-6">
          {/* Header */}
          <div className="flex flex-col md:flex-row md:items-center md:justify-between">
            <div>
              <h1 className="text-3xl font-bold text-white mb-2">Панель менеджера</h1>
              <p className="text-gray-400">
                {user?.storeId ? `Магазин ID: ${user.storeId}` : 'Управляйте вашей точкой'}
              </p>
            </div>
            <div className="mt-4 md:mt-0 bg-linear-to-r from-blue-500/10 to-emerald-500/10 border border-blue-500/20 rounded-lg px-4 py-3">
              <p className="text-sm text-gray-400">Привет,</p>
              <p className="text-lg font-semibold text-white">{user?.name}</p>
            </div>
          </div>

          {/* Stats Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {/* Today's Revenue */}
            <StatCard
              title="Выручка сегодня"
              value="₽125,400"
              change="12.5%"
              icon={<TrendingUp className="w-5 h-5" />}
              trend="up"
            />

            {/* Active Terminals */}
            <StatCard
              title="Активные кассы"
              value="8 / 10"
              icon={<ShoppingCart className="w-5 h-5" />}
            />

            {/* Daily Transactions */}
            <StatCard
              title="Транзакций сегодня"
              value="1,245"
              change="8.3%"
              icon={<BarChart3 className="w-5 h-5" />}
              trend="up"
            />

            {/* Team Members */}
            <StatCard
              title="Сотрудники на смене"
              value="12"
              icon={<Users className="w-5 h-5" />}
            />
          </div>

          {/* Quick Actions */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
            {/* Today's Performance */}
            <div className="bg-gray-900/50 border border-gray-800 rounded-xl p-6">
              <h2 className="text-xl font-bold text-white mb-6 flex items-center gap-2">
                <TrendingUp className="w-5 h-5 text-blue-500" />
                Производительность сегодня
              </h2>
              <div className="space-y-4">
                <PerformanceMetric label="Пиковое время" value="14:00 - 15:30" />
                <PerformanceMetric label="Отклонено чеков" value="2" unit="(0.2%)" />
                <PerformanceMetric label="Среднее время чека" value="3м 42с" />
              </div>
            </div>

            {/* Quick Actions Menu */}
            <div className="bg-gray-900/50 border border-gray-800 rounded-xl p-6">
              <h2 className="text-xl font-bold text-white mb-6 flex items-center gap-2">
                <TrendingUp className="w-5 h-5 text-emerald-500" />
                Быстрые действия
              </h2>
              <div className="space-y-3">
                <Link
                  href="/inventory"
                  className="block p-3 bg-linear-to-r from-blue-500/10 to-emerald-500/10 border border-blue-500/20 rounded-lg hover:border-blue-500/40 transition-all text-white font-medium"
                >
                  📦 Проверить остатки
                </Link>
                <Link
                  href="/users"
                  className="block p-3 bg-linear-to-r from-purple-500/10 to-pink-500/10 border border-purple-500/20 rounded-lg hover:border-purple-500/40 transition-all text-white font-medium"
                >
                  👥 Управлять сотрудниками
                </Link>
                <Link
                  href="/analytics"
                  className="block p-3 bg-linear-to-r from-orange-500/10 to-red-500/10 border border-orange-500/20 rounded-lg hover:border-orange-500/40 transition-all text-white font-medium"
                >
                  📊 Просмотреть отчеты
                </Link>
              </div>
            </div>
          </div>

          {/* Recent Alerts */}
          <div className="bg-gray-900/50 border border-gray-800 rounded-xl p-6">
            <h2 className="text-xl font-bold text-white mb-4">Уведомления</h2>
            <div className="space-y-3">
              <ManagerAlert
                type="warning"
                title="Низкие остатки товара"
                message="Молоко и хлеб требуют пополнения"
              />
              <ManagerAlert
                type="info"
                title="Техническое обслуживание"
                message="Касса №5 требует обслуживания завтра"
              />
            </div>
          </div>
        </div>
      </SPADashboardLayout>
    </ProtectedRoute>
  );
}
