"use client";

import { useHealthCheck, useHealthSummary } from "@/hooks/useHealthCheck";
import {
  ServiceCard,
  HealthSummary,
  HealthChart,
  HealthAlerts,
} from "@/components/health";
import { Button } from "@/components/ui/Button";
import { RefreshCw, Power, PowerOff } from "lucide-react";
import { useState } from "react";

export default function HealthPage() {
  const [autoRefresh, setAutoRefresh] = useState(true);
  const { services, history, isChecking, checkAllServices } = useHealthCheck(autoRefresh, 30000);
  const summary = useHealthSummary(services);

  const handleRefresh = async () => {
    await checkAllServices();
  };

  const toggleAutoRefresh = () => {
    setAutoRefresh(!autoRefresh);
  };

  return (
    <div className="container mx-auto py-8 space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold">Health Checks Dashboard</h1>
          <p className="text-muted-foreground mt-1">
            Мониторинг состояния микросервисов
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={toggleAutoRefresh}
            className="gap-2"
          >
            {autoRefresh ? (
              <>
                <Power className="h-4 w-4 text-green-600" />
                <span>Авто: ВКЛ</span>
              </>
            ) : (
              <>
                <PowerOff className="h-4 w-4 text-gray-600" />
                <span>Авто: ВЫКЛ</span>
              </>
            )}
          </Button>
          <Button
            variant="secondary"
            size="sm"
            onClick={handleRefresh}
            disabled={isChecking}
            className="gap-2"
          >
            <RefreshCw className={`h-4 w-4 ${isChecking ? "animate-spin" : ""}`} />
            Обновить
          </Button>
        </div>
      </div>

      {/* Alerts */}
      <HealthAlerts services={services} />

      {/* Summary Stats */}
      <HealthSummary summary={summary} />

      {/* History Chart */}
      {history.length > 0 && (
        <div>
          <h2 className="text-xl font-semibold mb-4">График состояния</h2>
          <HealthChart history={history} />
        </div>
      )}

      {/* Service Cards */}
      <div>
        <h2 className="text-xl font-semibold mb-4">Сервисы</h2>
        {isChecking && services.length === 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {[...Array(6)].map((_, i) => (
              <div
                key={i}
                className="h-64 bg-muted animate-pulse rounded-lg"
              />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {services.map((service) => (
              <ServiceCard key={service.name} service={service} />
            ))}
          </div>
        )}
      </div>

      {/* Individual Service Charts */}
      {history.length > 0 && (
        <div>
          <h2 className="text-xl font-semibold mb-4">Детализация по сервисам</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {services.slice(0, 4).map((service) => (
              <HealthChart
                key={service.name}
                history={history}
                serviceName={service.name}
              />
            ))}
          </div>
        </div>
      )}

      {/* Last Update Time */}
      {services.length > 0 && (
        <div className="text-center text-sm text-muted-foreground">
          Последнее обновление:{" "}
          {new Intl.DateTimeFormat("ru-RU", {
            dateStyle: "short",
            timeStyle: "medium",
          }).format(new Date())}
        </div>
      )}
    </div>
  );
}
