/**
 * Companies Stat Card - Same style as Dashboard StatsCard
 */

'use client';

import React from 'react';

interface CompaniesStatCardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  icon?: React.ReactNode;
  trend?: {
    value: number;
    isPositive: boolean;
  };
  color?: 'blue' | 'green' | 'purple' | 'orange';
}

export default function CompaniesStatCard({
  title,
  value,
  subtitle,
  icon,
  trend,
  color = 'blue',
}: CompaniesStatCardProps) {
  const bgColor = {
    blue: 'bg-blue-900/20 border-blue-800',
    green: 'bg-green-900/20 border-green-800',
    purple: 'bg-purple-900/20 border-purple-800',
    orange: 'bg-orange-900/20 border-orange-800',
  }[color];

  return (
    <div className={`${bgColor} border rounded-xl p-6 backdrop-blur-sm transition-all hover:shadow-lg`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <p className="text-gray-400 text-sm font-medium mb-2">{title}</p>
          <p className="text-3xl font-bold text-white">{value}</p>
          {subtitle && <p className="text-gray-500 text-xs mt-1">{subtitle}</p>}
        </div>
        {icon && <div className="ml-4">{icon}</div>}
      </div>
      {trend && (
        <div className="mt-4 pt-4 border-t border-gray-700">
          <span
            className={`text-sm font-semibold ${
              trend.isPositive ? 'text-emerald-400' : 'text-red-400'
            }`}
          >
            {trend.isPositive ? '↑' : '↓'} {Math.abs(trend.value)}%
          </span>
        </div>
      )}
    </div>
  );
}
