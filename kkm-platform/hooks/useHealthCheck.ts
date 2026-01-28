"use client";

import { useState, useEffect, useCallback, useMemo } from "react";
import type {
  ServiceHealth,
  HealthCheckHistory,
  ServiceStatus,
} from "@/types/health";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// Список микросервисов для мониторинга
const SERVICES = [
  { name: "API Gateway", endpoint: "/health" },
  { name: "User Service", endpoint: "/health" },
  { name: "Invoice Service", endpoint: "/health" },
  { name: "Catalog Service", endpoint: "/health" },
  { name: "Document Service", endpoint: "/health" },
  { name: "Company Service", endpoint: "/health" },
  { name: "Analytics Service", endpoint: "/health" },
];

/**
 * Hook для мониторинга здоровья сервисов
 */
export function useHealthCheck(autoRefresh = true, intervalMs = 30000) {
  const [services, setServices] = useState<ServiceHealth[]>([]);
  const [history, setHistory] = useState<HealthCheckHistory[]>([]);
  const [isChecking, setIsChecking] = useState(false);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);

  const checkServiceHealth = async (
    name: string,
    endpoint: string,
  ): Promise<ServiceHealth> => {
    const startTime = Date.now();
    try {
      const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
        cache: "no-store",
        signal: AbortSignal.timeout(5000), // 5s timeout
      });

      const responseTime = Date.now() - startTime;
      const data = await response.json();

      let status: ServiceStatus = "unknown";
      if (response.ok) {
        status = data.status === "healthy" ? "healthy" : "degraded";
      } else {
        status = "unhealthy";
      }

      return {
        name,
        status,
        endpoint,
        lastCheck: new Date(),
        responseTime,
        version: data.version,
        uptime: data.uptime,
        dependencies: data.dependencies,
        metrics: data.metrics,
        message: data.message,
      };
    } catch (error) {
      const responseTime = Date.now() - startTime;
      return {
        name,
        status: "unhealthy",
        endpoint,
        lastCheck: new Date(),
        responseTime,
        message: error instanceof Error ? error.message : "Connection failed",
      };
    }
  };

  const checkAllServices = useCallback(async () => {
    setIsChecking(true);
    try {
      const results = await Promise.all(
        SERVICES.map((service) =>
          checkServiceHealth(service.name, service.endpoint),
        ),
      );

      setServices(results);
      setLastUpdate(new Date());

      // Добавляем в историю
      const historyEntries: HealthCheckHistory[] = results.map((result) => ({
        timestamp: new Date(),
        serviceName: result.name,
        status: result.status,
        responseTime: result.responseTime || 0,
      }));

      setHistory((prev) => {
        const newHistory = [...historyEntries, ...prev].slice(0, 500); // Последние 500
        // Сохраняем в localStorage
        if (typeof window !== "undefined") {
          localStorage.setItem(
            "healthCheckHistory",
            JSON.stringify(newHistory),
          );
        }
        return newHistory;
      });
    } finally {
      setIsChecking(false);
    }
  }, []);

  // Загружаем историю из localStorage
  useEffect(() => {
    if (typeof window !== "undefined") {
      const stored = localStorage.getItem("healthCheckHistory");
      if (stored) {
        try {
          const parsed = JSON.parse(stored) as Array<{
            timestamp: string;
            serviceName: string;
            status: ServiceStatus;
            responseTime: number;
          }>;
          setHistory(
            parsed.map((entry) => ({
              ...entry,
              timestamp: new Date(entry.timestamp),
            })),
          );
        } catch (e) {
          console.error("Failed to parse health check history:", e);
        }
      }
    }
  }, []);

  // Первоначальная проверка
  useEffect(() => {
    checkAllServices();
  }, [checkAllServices]);

  // Автоматическое обновление
  useEffect(() => {
    if (!autoRefresh) return;

    const interval = setInterval(() => {
      checkAllServices();
    }, intervalMs);

    return () => clearInterval(interval);
  }, [autoRefresh, intervalMs, checkAllServices]);

  const clearHistory = useCallback(() => {
    setHistory([]);
    if (typeof window !== "undefined") {
      localStorage.removeItem("healthCheckHistory");
    }
  }, []);

  return {
    services,
    history,
    isChecking,
    lastUpdate,
    checkAllServices,
    clearHistory,
  };
}

/**
 * Hook для получения сводки по здоровью системы
 */
export function useHealthSummary(services: ServiceHealth[]) {
  const summary = useMemo(() => {
    if (services.length === 0) {
      return {
        total: 0,
        healthy: 0,
        degraded: 0,
        unhealthy: 0,
        unknown: 0,
        averageResponseTime: 0,
        uptime: 0,
      };
    }

    const total = services.length;
    const healthy = services.filter((s) => s.status === "healthy").length;
    const degraded = services.filter((s) => s.status === "degraded").length;
    const unhealthy = services.filter((s) => s.status === "unhealthy").length;
    const unknown = services.filter((s) => s.status === "unknown").length;

    const averageResponseTime =
      services.reduce((sum, s) => sum + (s.responseTime || 0), 0) / total;

    const uptime = (healthy / total) * 100;

    return {
      total,
      healthy,
      degraded,
      unhealthy,
      unknown,
      averageResponseTime,
      uptime,
    };
  }, [services]);

  return summary;
}
