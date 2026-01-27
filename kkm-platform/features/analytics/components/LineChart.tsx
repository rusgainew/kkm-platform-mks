"use client";

import {
  LineChart as RechartsLineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import type { ChartDataPoint } from "@/types/analytics";

export interface LineChartProps {
  data: ChartDataPoint[];
  xAxisKey?: string;
  yAxisKey?: string;
  lineColor?: string;
  title?: string;
  height?: number;
  currency?: string;
  className?: string;
}

/**
 * Компонент линейного графика
 * Используется для отображения трендов во времени
 */
export function LineChart({
  data,
  xAxisKey = "date",
  yAxisKey = "value",
  lineColor = "#3b82f6",
  title,
  height = 300,
  currency = "KGS",
  className = "",
}: LineChartProps) {
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
        <RechartsLineChart data={data}>
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
          <Line
            type="monotone"
            dataKey={yAxisKey}
            stroke={lineColor}
            strokeWidth={2}
            dot={{ fill: lineColor, r: 4 }}
            activeDot={{ r: 6 }}
          />
        </RechartsLineChart>
      </ResponsiveContainer>
    </div>
  );
}
