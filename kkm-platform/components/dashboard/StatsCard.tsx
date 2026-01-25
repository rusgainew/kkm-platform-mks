'use client';

import React, { memo } from 'react';
import { TrendingUp } from 'lucide-react';

interface StatsCardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  icon: React.ReactNode;
  trend?: {
    value: number;
    isPositive: boolean;
  };
  bgGradient?: string;
}

const StatsCard = memo(function StatsCard({ 
  title, 
  value, 
  subtitle, 
  icon, 
  trend,
  bgGradient = 'from-blue-600/10 to-blue-600/5'
}: StatsCardProps) {
  const gradientClass = `bg-linear-to-br ${bgGradient}`;
  return (
    <div className={`${gradientClass} rounded-lg border border-gray-800 p-6 shadow-lg hover:shadow-xl hover:border-gray-700 transition-all backdrop-blur-sm`}>
      <div className="flex justify-between items-start mb-4">
        <div className="flex-1">
          <p className="text-sm font-medium text-gray-400">{title}</p>
          <p className="text-3xl font-bold text-white mt-2 tracking-tight">{value}</p>
          {subtitle && <p className="text-xs text-gray-500 mt-1">{subtitle}</p>}
        </div>
        <div className="text-blue-400 opacity-80">{icon}</div>
      </div>

      {trend && (
        <div className="flex items-center gap-1">
          <TrendingUp
            className={`w-4 h-4 ${trend.isPositive ? 'text-emerald-400' : 'text-red-400'}`}
          />
          <span
            className={`text-sm font-semibold ${
              trend.isPositive ? 'text-emerald-400' : 'text-red-400'
            }`}
          >
            {trend.isPositive ? '+' : '-'}
            {Math.abs(trend.value)}%
          </span>
          <span className="text-xs text-gray-500">за день</span>
        </div>
      )}
    </div>
  );
});

export default StatsCard;
