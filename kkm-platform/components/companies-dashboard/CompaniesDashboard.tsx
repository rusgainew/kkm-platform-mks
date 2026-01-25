/**
 * Дашборд компаний
 * Отображает аналитику и статистику по компаниям:
 * - Статистические карточки (всего, активные, неактивные, заблокированные, члены)
 * - График топ компании по активности
 * - Круговую диаграмму распределения по статусам
 * - Таблицу недавней активности
 * - Переключение между режимом дашборда и списком компаний
 */

'use client';

import React, { useState } from 'react';
import CompaniesStatCard from './CompaniesStatCard';
import CompaniesOverviewChart from './CompaniesOverviewChart';
import CompaniesStatusDistribution from './CompaniesStatusDistribution';
import RecentCompaniesActivity from './RecentCompaniesActivity';
import CompaniesList from '@/components/companies/CompaniesList';
import { useListCompaniesQuery } from '@/lib/hooks/useCompaniesApi';
import { DashboardSkeleton } from '@/components/loading';

export default function CompaniesDashboard() {
  const [viewMode, setViewMode] = useState<'dashboard' | 'list'>('dashboard');
  const { data: companies = [], isLoading } = useListCompaniesQuery();

  // Calculate stats
  const totalCompanies = companies.length;
  const activeCompanies = companies.filter((c) => c.status === 'active').length;
  const inactiveCompanies = companies.filter((c) => c.status === 'inactive').length;
  const suspendedCompanies = companies.filter((c) => c.status === 'suspended').length;
  const totalMembers = companies.reduce((sum, c) => sum + (c.member_count || 0), 0);

  if (isLoading) {
    return (
      <div className="p-8">
        <div className="max-w-7xl mx-auto">
          <DashboardSkeleton />
        </div>
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="max-w-7xl mx-auto">
        {/* Header with Tabs */}
        <div className="mb-8 flex items-start justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white flex items-center gap-3">
              <span className="w-2 h-8 bg-linear-to-b from-blue-400 to-emerald-400 rounded-full"></span>
              Компании
            </h1>
            <p className="text-gray-400 mt-2 ml-5">Управление организациями</p>
          </div>
          <div className="flex items-center gap-2 bg-gray-800/50 rounded-lg p-1">
            <button
              onClick={() => setViewMode('dashboard')}
              className={`px-4 py-2 rounded transition-colors ${
                viewMode === 'dashboard'
                  ? 'bg-blue-600 text-white'
                  : 'text-gray-400 hover:text-white'
              }`}
            >
              Дашборд
            </button>
            <button
              onClick={() => setViewMode('list')}
              className={`px-4 py-2 rounded transition-colors ${
                viewMode === 'list'
                  ? 'bg-blue-600 text-white'
                  : 'text-gray-400 hover:text-white'
              }`}
            >
              Список
            </button>
          </div>
        </div>

        {viewMode === 'list' ? (
          <CompaniesList />
        ) : (
          <>
            {/* Stats Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
              <CompaniesStatCard
                title="Всего компаний"
                value={totalCompanies}
                trend={{ value: 12.5, isPositive: true }}
                color="blue"
              />
              <CompaniesStatCard
                title="Активных компаний"
                value={activeCompanies}
                subtitle="компаний online"
                trend={{ value: 5, isPositive: true }}
                color="green"
              />
              <CompaniesStatCard
                title="Всего членов"
                value={totalMembers}
                subtitle="сотрудников"
                trend={{ value: 8.3, isPositive: true }}
                color="purple"
              />
              <CompaniesStatCard
                title="Заблокированных"
                value={suspendedCompanies}
                trend={{ value: 3.2, isPositive: false }}
                color="orange"
              />
            </div>

            {/* Charts Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
              {/* Left Column - Charts */}
              <div className="lg:col-span-2">
                <CompaniesOverviewChart
                  data={companies.map((c) => ({
                    name: c.name,
                    members: c.member_count || 0,
                  }))}
                />
              </div>

              {/* Right Column - Status Distribution */}
              <CompaniesStatusDistribution
                active={activeCompanies}
                inactive={inactiveCompanies}
                suspended={suspendedCompanies}
              />
            </div>

            {/* Recent Activity */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
              <div className="lg:col-span-2">
                <RecentCompaniesActivity companies={companies} />
              </div>

              {/* Summary Stats */}
              <div className="bg-gray-800/30 border border-gray-700 rounded-xl p-6 backdrop-blur-sm">
                <h3 className="text-lg font-bold text-white mb-6">Сводка</h3>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-gray-400">Активные</span>
                    <span className="text-lg font-semibold text-emerald-400">
                      {activeCompanies}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-400">Неактивные</span>
                    <span className="text-lg font-semibold text-gray-400">
                      {inactiveCompanies}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-400">Заблокированные</span>
                    <span className="text-lg font-semibold text-orange-400">
                      {suspendedCompanies}
                    </span>
                  </div>
                  <div className="pt-4 border-t border-gray-700 flex items-center justify-between">
                    <span className="text-gray-300 font-medium">Процент активных</span>
                    <span className="text-lg font-bold text-blue-400">
                      {totalCompanies > 0
                        ? ((activeCompanies / totalCompanies) * 100).toFixed(1)
                        : 0}
                      %
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
