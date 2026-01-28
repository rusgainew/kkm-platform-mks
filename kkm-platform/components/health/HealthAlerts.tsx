"use client";

import { Alert } from "@/components/ui/Alert";
import { AlertTriangle, XCircle } from "lucide-react";
import type { ServiceHealth } from "@/types/health";

interface HealthAlertsProps {
  services: ServiceHealth[];
}

export function HealthAlerts({ services }: HealthAlertsProps) {
  const unhealthyServices = services.filter(
    (s) => s.status === "unhealthy" || s.status === "degraded"
  );

  if (unhealthyServices.length === 0) {
    return null;
  }

  return (
    <div className="space-y-2">
      {unhealthyServices.map((service) => {
        const isUnhealthy = service.status === "unhealthy";
        return (
          <Alert
            key={service.name}
            variant={isUnhealthy ? "danger" : "warning"}
            title={`${service.name} - ${isUnhealthy ? "Недоступен" : "Деградация"}`}
            icon={isUnhealthy ? <XCircle className="h-4 w-4" /> : <AlertTriangle className="h-4 w-4" />}
          >
            {service.message || "Сервис не отвечает на запросы"}
            {service.responseTime && service.responseTime > 0 && ` (${service.responseTime}ms)`}
          </Alert>
        );
      })}
    </div>
  );
}
