"use client";

import {
  BarChart as RechartsBarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import type { ChartDataPoint } from "@/types/analytics";

export interface BarChartProps {
  data: ChartDataPoint[];
  xAxisKey?: string;
  yAxisKey?: string;
  barColor?: string;
  title?: string;
  height?: number;
  currency?: string;
  className?: string;
}

/**
 * Компонент столбчатого графика
 * Используется для сравнения значений по категориям
 */
export function BarChart({
  data,
  xAxisKey = "date",
  yAxisKey = "value",
  barColor = "#10b981",
  title,
  height = 300,
  currency = "KGS",
  className = "",
}: BarChartProps) {
  const formatYAxis = (value: number) => {
    if (value >= 1000000) {
      return `${(value / 1000000).toFixed(1)}М`;
    }
    if (value >= 1000) {
      return `${(value / 1000).toFixed(1)}К`;
    }
    return value.toString();
  };

  const formatTooltip = (value: number | undefined) => {
    if (value === undefined) return "";
    return new Intl.NumberFormat("ru-RU", {
      style: "currency",
      currency: currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 2,
    }).format(value);
  };

  return (
    <div className={`bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6 ${className}`}>
      {title && (
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
          {title}
        </h3>
      )}
      
      <ResponsiveContainer width="100%" height={height}>
        <RechartsBarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.1} />
          <XAxis
            dataKey={xAxisKey}
            stroke="#6b7280"
            tick={{ fill: "#6b7280" }}
            tickLine={{ stroke: "#6b7280" }}
          />
          <YAxis
            tickFormatter={formatYAxis}
            stroke="#6b7280"
            tick={{ fill: "#6b7280" }}
            tickLine={{ stroke: "#6b7280" }}
          />
          <Tooltip
            formatter={formatTooltip}
            contentStyle={{
              backgroundColor: "#1f2937",
              border: "1px solid #374151",
              borderRadius: "0.5rem",
              color: "#f3f4f6",
            }}
            labelStyle={{ color: "#9ca3af" }}
          />
          <Legend
            wrapperStyle={{ color: "#6b7280" }}
          />
          <Bar
            dataKey={yAxisKey}
            fill={barColor}
            radius={[4, 4, 0, 0]}
          />
        </RechartsBarChart>
      </ResponsiveContainer>
    </div>
  );
}
