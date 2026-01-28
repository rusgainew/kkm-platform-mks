"use client";

import { Card } from "@/components/ui/Card";
import { Progress } from "@/components/ui/progress";
import { AlertTriangle, Clock, Check, XCircle } from "lucide-react";
import type { RateLimitStatus } from "@/types/rate-limit";

interface RateLimitCardProps {
  status: RateLimitStatus;
}

export function RateLimitCard({ status }: RateLimitCardProps) {
  const getStatusColor = () => {
    if (status.isLimited) return "text-red-600 dark:text-red-400";
    if (status.isNearLimit) return "text-yellow-600 dark:text-yellow-400";
    return "text-green-600 dark:text-green-400";
  };

  const getStatusIcon = () => {
    if (status.isLimited)
      return <XCircle className="h-5 w-5 text-red-600 dark:text-red-400" />;
    if (status.isNearLimit)
      return <AlertTriangle className="h-5 w-5 text-yellow-600 dark:text-yellow-400" />;
    return <Check className="h-5 w-5 text-green-600 dark:text-green-400" />;
  };

  const getProgressColor = () => {
    if (status.isLimited) return "bg-red-500";
    if (status.isNearLimit) return "bg-yellow-500";
    return "bg-green-500";
  };

  const formatResetTime = () => {
    const now = new Date();
    const diff = status.resetAt.getTime() - now.getTime();

    if (diff < 0) return "сброшен";

    const minutes = Math.floor(diff / 60000);
    const seconds = Math.floor((diff % 60000) / 1000);

    if (minutes > 0) {
      return `через ${minutes}м ${seconds}с`;
    }
    return `через ${seconds}с`;
  };

  return (
    <Card className="p-4">
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            {getStatusIcon()}
            <h3 className="font-medium text-sm truncate">{status.endpoint}</h3>
          </div>
          <p className="text-xs text-muted-foreground">
            {status.requestsUsed} / {status.requestsLimit} запросов
          </p>
        </div>
        <div className="text-right">
          <p className={`text-xl font-bold ${getStatusColor()}`}>
            {status.percentageUsed.toFixed(0)}%
          </p>
        </div>
      </div>

      <Progress value={status.percentageUsed} className={`h-2 mb-2 ${getProgressColor()}`} />

      <div className="flex items-center justify-between text-xs text-muted-foreground">
        <div className="flex items-center gap-1">
          <Clock className="h-3 w-3" />
          <span>Сброс: {formatResetTime()}</span>
        </div>
        <span>{status.requestsLimit - status.requestsUsed} доступно</span>
      </div>

      {status.isLimited && (
        <div className="mt-2 p-2 bg-red-50 dark:bg-red-950/20 rounded text-xs text-red-700 dark:text-red-300">
          Лимит превышен. Дождитесь сброса.
        </div>
      )}

      {status.isNearLimit && !status.isLimited && (
        <div className="mt-2 p-2 bg-yellow-50 dark:bg-yellow-950/20 rounded text-xs text-yellow-700 dark:text-yellow-300">
          Приближение к лимиту. Снизьте частоту запросов.
        </div>
      )}
    </Card>
  );
}
