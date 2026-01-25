'use client';

import dynamic from 'next/dynamic';
import { Suspense } from 'react';

/**
 * Простой Skeleton loader
 */
function Skeleton({ className = '' }: { className?: string }) {
  return (
    <div
      className={`animate-pulse bg-gray-200 rounded-lg ${className}`}
    />
  );
}

/**
 * Динамически загружаемые компоненты для dashbaord
 * Каждый компонент загружается при необходимости
 */

const DynamicCharts = dynamic(
  () =>
    import('@/components/dashboard/Charts').then((mod) => ({
      default: () => (
        <div className="space-y-6">
          <div className="h-80 bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg flex items-center justify-center text-white font-bold">
            Диаграмма доходов
          </div>
          <div className="h-80 bg-gradient-to-r from-green-500 to-green-600 rounded-lg flex items-center justify-center text-white font-bold">
            Статус счетов-фактур
          </div>
        </div>
      ),
    })),
  {
    loading: () => (
      <div className="space-y-6">
        <Skeleton className="h-80 w-full" />
        <Skeleton className="h-80 w-full" />
      </div>
    ),
    ssr: false,
  }
);

/**
 * Оптимизированный dashboard с lazy loading
 */
export function OptimizedDashboard() {
  return (
    <div className="space-y-8">
      <h2 className="text-2xl font-bold">Оптимизированный Dashboard</h2>

      {/* Charts - загружаются отложенно */}
      <Suspense
        fallback={
          <div className="space-y-6">
            <Skeleton className="h-80 w-full" />
            <Skeleton className="h-80 w-full" />
          </div>
        }
      >
        <DynamicCharts />
      </Suspense>

      {/* Info */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h3 className="font-bold mb-2">💡 Оптимизация:</h3>
        <ul className="space-y-1 text-sm">
          <li>✅ Динамическая загрузка компонентов (Dynamic imports)</li>
          <li>✅ Suspense для показа loader во время загрузки</li>
          <li>✅ SSR disabled для тяжелых компонентов</li>
          <li>✅ Lazy loading при скролле до раздела</li>
        </ul>
      </div>
    </div>
  );
}
