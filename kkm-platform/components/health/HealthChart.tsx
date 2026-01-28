"use client";

import { Card } from "@/components/ui/Card";
import { Activity, TrendingUp, TrendingDown, Minus } from "lucide-react";
import type { HealthCheckHistory } from "@/types/health";
import { useMemo } from "react";

interface HealthChartProps {
  history: HealthCheckHistory[];
  serviceName?: string;
}

export function HealthChart({ history, serviceName }: HealthChartProps) {
  const filteredHistory = useMemo(() => {
    const filtered = serviceName
      ? history.filter((h) => h.serviceName === serviceName)
      : history;

    // Последние 50 записей
    return filtered.slice(0, 50).reverse();
  }, [history, serviceName]);

  const stats = useMemo(() => {
    if (filteredHistory.length === 0) return null;

    const avgResponseTime =
      filteredHistory.reduce((sum, h) => sum + h.responseTime, 0) /
      filteredHistory.length;

    const healthyCount = filteredHistory.filter((h) => h.status === "healthy").length;
    const uptimePercentage = (healthyCount / filteredHistory.length) * 100;

    const recentAvg =
      filteredHistory.slice(-10).reduce((sum, h) => sum + h.responseTime, 0) / 10;
    const olderAvg =
      filteredHistory.slice(0, 10).reduce((sum, h) => sum + h.responseTime, 0) / 10;
    const trend = recentAvg < olderAvg ? "up" : recentAvg > olderAvg ? "down" : "stable";

    return {
      avgResponseTime,
      uptimePercentage,
      trend,
      healthyCount,
      totalChecks: filteredHistory.length,
    };
  }, [filteredHistory]);

  if (!stats || filteredHistory.length === 0) {
    return (
      <Card className="p-6">
        <div className="text-center text-muted-foreground">
          <Activity className="h-12 w-12 mx-auto mb-2 opacity-50" />
          <p>Нет данных для отображения</p>
        </div>
      </Card>
    );
  }

  const maxResponseTime = Math.max(...filteredHistory.map((h) => h.responseTime));

  return (
    <Card className="p-4">
      <div className="mb-4">
        <h3 className="font-semibold text-lg mb-2">
          {serviceName ? `График: ${serviceName}` : "Общий график"}
        </h3>
        <div className="grid grid-cols-3 gap-4 text-sm">
          <div>
            <p className="text-muted-foreground">Средний ответ</p>
            <p className="text-2xl font-bold">{stats.avgResponseTime.toFixed(0)}ms</p>
          </div>
          <div>
            <p className="text-muted-foreground">Uptime</p>
            <p className="text-2xl font-bold text-green-600">
              {stats.uptimePercentage.toFixed(1)}%
            </p>
          </div>
          <div>
            <p className="text-muted-foreground">Тренд</p>
            <div className="flex items-center gap-1">
              {stats.trend === "up" && <TrendingUp className="h-5 w-5 text-green-600" />}
              {stats.trend === "down" && (
                <TrendingDown className="h-5 w-5 text-red-600" />
              )}
              {stats.trend === "stable" && <Minus className="h-5 w-5 text-gray-600" />}
              <span className="text-xl font-bold capitalize">{stats.trend}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Simple Bar Chart */}
      <div className="h-32 flex items-end gap-1">
        {filteredHistory.map((entry, idx) => {
          const height = (entry.responseTime / maxResponseTime) * 100;
          const isHealthy = entry.status === "healthy";

          return (
            <div
              key={`${entry.serviceName}-${entry.timestamp.getTime()}-${idx}`}
              className="flex-1 flex flex-col justify-end"
              title={`${entry.serviceName}: ${entry.responseTime}ms (${entry.status})`}
            >
              <div
                className={`w-full rounded-t transition-all ${
                  isHealthy ? "bg-green-500" : "bg-red-500"
                }`}
                style={{ height: `${height}%` }}
              />
            </div>
          );
        })}
      </div>

      <div className="mt-2 flex items-center justify-between text-xs text-muted-foreground">
        <span>
          {new Intl.DateTimeFormat("ru-RU", {
            hour: "2-digit",
            minute: "2-digit",
          }).format(filteredHistory[0].timestamp)}
        </span>
        <span>
          {stats.healthyCount} / {stats.totalChecks} проверок успешно
        </span>
        <span>
          {new Intl.DateTimeFormat("ru-RU", {
            hour: "2-digit",
            minute: "2-digit",
          }).format(filteredHistory[filteredHistory.length - 1].timestamp)}
        </span>
      </div>
    </Card>
  );
}
