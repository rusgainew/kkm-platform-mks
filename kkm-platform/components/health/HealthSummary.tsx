"use client";

import { Card } from "@/components/ui/Card";
import { CheckCircle2, AlertTriangle, XCircle, Activity } from "lucide-react";
import type { HealthSummary as HealthSummaryType } from "@/types/health";

interface HealthSummaryProps {
  summary: HealthSummaryType;
}

export function HealthSummary({ summary }: HealthSummaryProps) {
  const stats = [
    {
      label: "Всего сервисов",
      value: summary.total,
      icon: Activity,
      color: "text-blue-600",
      bgColor: "bg-blue-50 dark:bg-blue-950",
    },
    {
      label: "Работают",
      value: summary.healthy,
      icon: CheckCircle2,
      color: "text-green-600",
      bgColor: "bg-green-50 dark:bg-green-950",
    },
    {
      label: "Деградация",
      value: summary.degraded,
      icon: AlertTriangle,
      color: "text-yellow-600",
      bgColor: "bg-yellow-50 dark:bg-yellow-950",
    },
    {
      label: "Недоступны",
      value: summary.unhealthy,
      icon: XCircle,
      color: "text-red-600",
      bgColor: "bg-red-50 dark:bg-red-950",
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {stats.map((stat) => {
        const Icon = stat.icon;
        return (
          <Card key={stat.label} className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground mb-1">{stat.label}</p>
                <p className="text-3xl font-bold">{stat.value}</p>
              </div>
              <div className={`p-3 rounded-full ${stat.bgColor}`}>
                <Icon className={`h-6 w-6 ${stat.color}`} />
              </div>
            </div>
          </Card>
        );
      })}
      
      {/* Additional metrics */}
      <Card className="p-6 md:col-span-2 lg:col-span-4">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <p className="text-sm text-muted-foreground mb-1">Средний ответ</p>
            <p className="text-2xl font-bold">{summary.averageResponseTime.toFixed(0)}ms</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground mb-1">Общий Uptime</p>
            <p className="text-2xl font-bold text-green-600">
              {summary.uptime.toFixed(2)}%
            </p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground mb-1">Неизвестные</p>
            <p className="text-2xl font-bold text-gray-600">{summary.unknown}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground mb-1">Статус системы</p>
            <p
              className={`text-2xl font-bold ${
                summary.unhealthy > 0
                  ? "text-red-600"
                  : summary.degraded > 0
                    ? "text-yellow-600"
                    : "text-green-600"
              }`}
            >
              {summary.unhealthy > 0
                ? "Критично"
                : summary.degraded > 0
                  ? "Деградация"
                  : "Отлично"}
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
}
