"use client";

import { Card } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { CheckCircle2, AlertTriangle, XCircle, HelpCircle, Clock, Zap } from "lucide-react";
import type { ServiceHealth, ServiceStatus } from "@/types/health";

interface ServiceCardProps {
  service: ServiceHealth;
}

export function ServiceCard({ service }: ServiceCardProps) {
  const getStatusConfig = (status: ServiceStatus) => {
    switch (status) {
      case "healthy":
        return {
          icon: CheckCircle2,
          color: "text-green-600 dark:text-green-400",
          bgColor: "bg-green-50 dark:bg-green-950/20",
          badgeVariant: "success" as const,
          label: "Здоров",
        };
      case "degraded":
        return {
          icon: AlertTriangle,
          color: "text-yellow-600 dark:text-yellow-400",
          bgColor: "bg-yellow-50 dark:bg-yellow-950/20",
          badgeVariant: "warning" as const,
          label: "Снижена",
        };
      case "unhealthy":
        return {
          icon: XCircle,
          color: "text-red-600 dark:text-red-400",
          bgColor: "bg-red-50 dark:bg-red-950/20",
          badgeVariant: "danger" as const,
          label: "Не работает",
        };
      default:
        return {
          icon: HelpCircle,
          color: "text-gray-600 dark:text-gray-400",
          bgColor: "bg-gray-50 dark:bg-gray-950/20",
          badgeVariant: "neutral" as const,
          label: "Неизвестно",
        };
    }
  };

  const config = getStatusConfig(service.status);
  const StatusIcon = config.icon;

  const formatUptime = (seconds?: number) => {
    if (!seconds) return "N/A";
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    if (days > 0) return `${days}д ${hours}ч`;
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}ч ${minutes}м`;
  };

  return (
    <Card className={`p-4 transition-all hover:shadow-md ${config.bgColor}`}>
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-2">
          <StatusIcon className={`h-5 w-5 ${config.color}`} />
          <h3 className="font-semibold text-lg">{service.name}</h3>
        </div>
        <Badge variant={config.badgeVariant}>{config.label}</Badge>
      </div>

      <div className="space-y-2 text-sm">
        {/* Response Time */}
        {service.responseTime !== undefined && (
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground flex items-center gap-1">
              <Zap className="h-3 w-3" />
              Ответ:
            </span>
            <span className="font-medium">{service.responseTime}ms</span>
          </div>
        )}

        {/* Uptime */}
        {service.uptime !== undefined && (
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground flex items-center gap-1">
              <Clock className="h-3 w-3" />
              Uptime:
            </span>
            <span className="font-medium">{formatUptime(service.uptime)}</span>
          </div>
        )}

        {/* Version */}
        {service.version && (
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Версия:</span>
            <span className="font-mono text-xs">{service.version}</span>
          </div>
        )}

        {/* Endpoint */}
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Endpoint:</span>
          <code className="text-xs bg-muted px-1 py-0.5 rounded">{service.endpoint}</code>
        </div>

        {/* Last Check */}
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>Проверено:</span>
          <span>
            {new Intl.DateTimeFormat("ru-RU", {
              hour: "2-digit",
              minute: "2-digit",
              second: "2-digit",
            }).format(service.lastCheck)}
          </span>
        </div>

        {/* Message */}
        {service.message && (
          <div className="mt-2 p-2 bg-muted rounded text-xs">
            <p className="text-muted-foreground">{service.message}</p>
          </div>
        )}

        {/* Metrics */}
        {service.metrics && (
          <div className="mt-3 pt-3 border-t">
            <p className="text-xs font-medium mb-2">Метрики:</p>
            <div className="grid grid-cols-2 gap-2 text-xs">
              {service.metrics.cpu !== undefined && (
                <div>
                  <span className="text-muted-foreground">CPU:</span>
                  <span className="ml-1 font-medium">{service.metrics.cpu.toFixed(1)}%</span>
                </div>
              )}
              {service.metrics.memory !== undefined && (
                <div>
                  <span className="text-muted-foreground">Memory:</span>
                  <span className="ml-1 font-medium">{service.metrics.memory}MB</span>
                </div>
              )}
              {service.metrics.requests !== undefined && (
                <div>
                  <span className="text-muted-foreground">Requests:</span>
                  <span className="ml-1 font-medium">{service.metrics.requests}</span>
                </div>
              )}
              {service.metrics.errors !== undefined && (
                <div>
                  <span className="text-muted-foreground">Errors:</span>
                  <span className="ml-1 font-medium text-red-600">
                    {service.metrics.errors}
                  </span>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Dependencies */}
        {service.dependencies && service.dependencies.length > 0 && (
          <div className="mt-3 pt-3 border-t">
            <p className="text-xs font-medium mb-2">Зависимости:</p>
            <div className="space-y-1">
              {service.dependencies.map((dep, idx) => (
                <div key={idx} className="flex items-center justify-between text-xs">
                  <span className="text-muted-foreground">{dep.name}</span>
                  <Badge
                    variant={dep.status === "healthy" ? "success" : "danger"}
                    className="h-5 text-xs"
                  >
                    {dep.status}
                  </Badge>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </Card>
  );
}
