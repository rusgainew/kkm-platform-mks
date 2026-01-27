"use client";

import { TrendingUp, TrendingDown, Minus } from "lucide-react";

export interface StatCardProps {
  label: string;
  value: string | number;
  change?: number;
  changeType?: "increase" | "decrease" | "neutral";
  format?: "currency" | "number" | "percentage";
  icon?: React.ReactNode;
  currency?: string;
  className?: string;
}

/**
 * Компонент карточки метрики
 * Отображает ключевую метрику с изменением относительно предыдущего периода
 */
export function StatCard({
  label,
  value,
  change,
  changeType = "neutral",
  format = "number",
  icon,
  currency = "KGS",
  className = "",
}: StatCardProps) {
  const formatValue = (val: string | number): string => {
    const numValue = typeof val === "string" ? parseFloat(val) : val;

    if (isNaN(numValue)) return String(val);

    switch (format) {
      case "currency":
        return new Intl.NumberFormat("ru-RU", {
          style: "currency",
          currency: currency,
          minimumFractionDigits: 0,
          maximumFractionDigits: 2,
        }).format(numValue);

      case "percentage":
        return `${numValue.toFixed(1)}%`;

      case "number":
      default:
        return new Intl.NumberFormat("ru-RU").format(numValue);
    }
  };

  const getTrendIcon = () => {
    switch (changeType) {
      case "increase":
        return <TrendingUp className="w-4 h-4" />;
      case "decrease":
        return <TrendingDown className="w-4 h-4" />;
      case "neutral":
      default:
        return <Minus className="w-4 h-4" />;
    }
  };

  const getTrendColor = () => {
    switch (changeType) {
      case "increase":
        return "text-green-600 dark:text-green-400";
      case "decrease":
        return "text-red-600 dark:text-red-400";
      case "neutral":
      default:
        return "text-gray-600 dark:text-gray-400";
    }
  };

  return (
    <div
      className={`bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6 ${className}`}
    >
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
          {label}
        </span>
        {icon && (
          <div className="text-gray-400 dark:text-gray-500">{icon}</div>
        )}
      </div>

      {/* Value */}
      <div className="mb-2">
        <span className="text-3xl font-bold text-gray-900 dark:text-gray-100">
          {formatValue(value)}
        </span>
      </div>

      {/* Change */}
      {change !== undefined && (
        <div className="flex items-center gap-1">
          <div className={`flex items-center gap-1 ${getTrendColor()}`}>
            {getTrendIcon()}
            <span className="text-sm font-medium">
              {Math.abs(change).toFixed(1)}%
            </span>
          </div>
          <span className="text-sm text-gray-500 dark:text-gray-400">
            за период
          </span>
        </div>
      )}
    </div>
  );
}
