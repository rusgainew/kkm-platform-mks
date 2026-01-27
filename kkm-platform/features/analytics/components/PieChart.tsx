/* eslint-disable @typescript-eslint/no-explicit-any */
"use client";

import {
  PieChart as RechartsPieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

export interface PieChartProps {
  data: Array<Record<string, any>>;
  title?: string;
  height?: number;
  colors?: string[];
  className?: string;
}

const DEFAULT_COLORS = [
  "#3b82f6", // blue
  "#10b981", // green
  "#f59e0b", // amber
  "#ef4444", // red
  "#8b5cf6", // purple
  "#ec4899", // pink
];

/**
 * Компонент круговой диаграммы
 * Используется для отображения процентного соотношения
 */
export function PieChart({
  data,
  title,
  height = 300,
  colors = DEFAULT_COLORS,
  className = "",
}: PieChartProps) {
  const formatTooltip = (value: number | undefined) => {
    if (value === undefined) return "";
    return value.toLocaleString("ru-KZ");
  };

  const renderLabel = (props: { name: string; percentage: number }) => {
    return `${props.name}: ${props.percentage.toFixed(1)}%`;
  };

  return (
    <div className={`bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6 ${className}`}>
      {title && (
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
          {title}
        </h3>
      )}
      
      <ResponsiveContainer width="100%" height={height}>
        <RechartsPieChart>
          <Pie
            data={data}
            cx="50%"
            cy="50%"
            labelLine={false}
            label={renderLabel as any}
            outerRadius={80}
            fill="#8884d8"
            dataKey="value"
          >
            {data.map((entry: any, index: number) => (
              <Cell
                key={`cell-${index}`}
                fill={entry.color || colors[index % colors.length]}
              />
            ))}
          </Pie>
          <Tooltip
            formatter={formatTooltip}
            contentStyle={{
              backgroundColor: "#1f2937",
              border: "1px solid #374151",
              borderRadius: "0.5rem",
              color: "#f3f4f6",
            }}
          />
          <Legend
            wrapperStyle={{ color: "#6b7280" }}
          />
        </RechartsPieChart>
      </ResponsiveContainer>
    </div>
  );
}
