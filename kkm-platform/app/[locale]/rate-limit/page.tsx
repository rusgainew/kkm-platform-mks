"use client";

import { useRateLimitStatus } from "@/hooks/useRateLimit";
import { RateLimitCard } from "@/components/rate-limit/RateLimitCard";
import { RateLimitHistory } from "@/components/rate-limit/RateLimitHistory";
import { Activity, TrendingUp, AlertTriangle, Check } from "lucide-react";
import { Card } from "@/components/ui/Card";

export default function RateLimitPage() {
  const { statuses, history, clearHistory } = useRateLimitStatus();

  const totalRequests = statuses.reduce((sum, s) => sum + s.requestsUsed, 0);
  const limitedEndpoints = statuses.filter((s) => s.isLimited).length;
  const warningEndpoints = statuses.filter((s) => s.isNearLimit && !s.isLimited).length;
  const healthyEndpoints = statuses.filter((s) => !s.isNearLimit).length;

  return (
    <div className="container mx-auto py-6 space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">Rate Limiting Dashboard</h1>
        <p className="text-muted-foreground mt-2">
          Мониторинг лимитов API запросов в реальном времени
        </p>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-muted-foreground">Всего запросов</p>
              <p className="text-2xl font-bold mt-1">{totalRequests}</p>
            </div>
            <Activity className="h-8 w-8 text-blue-500" />
          </div>
        </Card>

        <Card className="p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-muted-foreground">Здоровые</p>
              <p className="text-2xl font-bold mt-1 text-green-600">{healthyEndpoints}</p>
            </div>
            <Check className="h-8 w-8 text-green-500" />
          </div>
        </Card>

        <Card className="p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-muted-foreground">Предупреждения</p>
              <p className="text-2xl font-bold mt-1 text-yellow-600">{warningEndpoints}</p>
            </div>
            <AlertTriangle className="h-8 w-8 text-yellow-500" />
          </div>
        </Card>

        <Card className="p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-muted-foreground">Заблокированы</p>
              <p className="text-2xl font-bold mt-1 text-red-600">{limitedEndpoints}</p>
            </div>
            <TrendingUp className="h-8 w-8 text-red-500" />
          </div>
        </Card>
      </div>

      {/* Active Endpoints */}
      <div>
        <h2 className="text-xl font-semibold mb-4">Активные Endpoints</h2>
        {statuses.length === 0 ? (
          <Card className="p-12">
            <div className="text-center text-muted-foreground">
              <Activity className="h-16 w-16 mx-auto mb-4 opacity-50" />
              <p className="text-lg">Нет активных endpoint&#x27;ов</p>
              <p className="text-sm mt-2">
                Информация появится после первого API запроса
              </p>
            </div>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {statuses
              .sort((a, b) => b.percentageUsed - a.percentageUsed)
              .map((status) => (
                <RateLimitCard key={status.endpoint} status={status} />
              ))}
          </div>
        )}
      </div>

      {/* History */}
      <div>
        <h2 className="text-xl font-semibold mb-4">История запросов</h2>
        <RateLimitHistory history={history} onClear={clearHistory} />
      </div>
    </div>
  );
}
