'use client';

import { TrendingUp, TrendingDown, DollarSign, Users, ShoppingCart, AlertCircle } from 'lucide-react';

export interface KPIMetric {
  label: string;
  value: string | number;
  icon: React.ReactNode;
  trend?: number; // процент
  color: 'blue' | 'green' | 'red' | 'yellow';
}

interface KPIDashboardProps {
  metrics: KPIMetric[];
}

const colorClasses = {
  blue: 'bg-blue-50 text-blue-700',
  green: 'bg-green-50 text-green-700',
  red: 'bg-red-50 text-red-700',
  yellow: 'bg-yellow-50 text-yellow-700',
};

const iconColorClasses = {
  blue: 'bg-blue-100 text-blue-600',
  green: 'bg-green-100 text-green-600',
  red: 'bg-red-100 text-red-600',
  yellow: 'bg-yellow-100 text-yellow-600',
};

export function KPIDashboard({ metrics }: KPIDashboardProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      {metrics.map((metric, index) => (
        <div
          key={index}
          className={`rounded-lg shadow p-6 ${colorClasses[metric.color]}`}
        >
          <div className="flex items-start justify-between">
            <div>
              <p className="text-sm font-medium opacity-75">{metric.label}</p>
              <p className="text-2xl font-bold mt-2">{metric.value}</p>
              {metric.trend !== undefined && (
                <div className="flex items-center gap-1 mt-2 text-sm">
                  {metric.trend > 0 ? (
                    <TrendingUp className="w-4 h-4" />
                  ) : (
                    <TrendingDown className="w-4 h-4" />
                  )}
                  <span>{Math.abs(metric.trend)}% {metric.trend > 0 ? 'вверх' : 'вниз'}</span>
                </div>
              )}
            </div>
            <div className={`rounded-full p-3 ${iconColorClasses[metric.color]}`}>
              {metric.icon}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

/**
 * Готовые KPI метрики для финансового дашборда
 */
export function createFinanceKPIs(data: {
  totalRevenue: number;
  paidAmount: number;
  pendingAmount: number;
  overallUsers: number;
  activeUsers: number;
  totalOrders: number;
  lowStockItems: number;
}): KPIMetric[] {
  return [
    {
      label: 'Общий доход',
      value: `$${data.totalRevenue.toFixed(2)}`,
      icon: <DollarSign className="w-6 h-6" />,
      trend: 12,
      color: 'blue',
    },
    {
      label: 'Оплачено',
      value: `$${data.paidAmount.toFixed(2)}`,
      icon: <TrendingUp className="w-6 h-6" />,
      trend: 8,
      color: 'green',
    },
    {
      label: 'В ожидании',
      value: `$${data.pendingAmount.toFixed(2)}`,
      icon: <AlertCircle className="w-6 h-6" />,
      trend: -5,
      color: 'yellow',
    },
    {
      label: 'Активные пользователи',
      value: data.activeUsers,
      icon: <Users className="w-6 h-6" />,
      trend: 15,
      color: 'green',
    },
  ];
}
