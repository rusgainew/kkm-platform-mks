"use client";

import { useState, useMemo } from "react";
import type {
  RateLimitStatus,
  RateLimitHistoryEntry,
} from "@/types/rate-limit";

/**
 * Hook для отслеживания rate limit статуса
 */
export function useRateLimitStatus() {
  // Ленивая инициализация из localStorage
  const [statuses, setStatuses] = useState<Map<string, RateLimitStatus>>(() => {
    if (typeof window === "undefined") return new Map();

    const stored = localStorage.getItem("rateLimitStatuses");
    if (!stored) return new Map();

    try {
      const parsed = JSON.parse(stored) as Record<
        string,
        {
          endpoint: string;
          requestsUsed: number;
          requestsLimit: number;
          percentageUsed: number;
          resetAt: string;
          isNearLimit: boolean;
          isLimited: boolean;
        }
      >;
      return new Map(
        Object.entries(parsed).map(([key, value]) => [
          key,
          {
            ...value,
            resetAt: new Date(value.resetAt),
          },
        ]),
      );
    } catch (e) {
      console.error("Failed to parse rate limit statuses:", e);
      return new Map();
    }
  });

  const [history, setHistory] = useState<RateLimitHistoryEntry[]>(() => {
    if (typeof window === "undefined") return [];

    const storedHistory = localStorage.getItem("rateLimitHistory");
    if (!storedHistory) return [];

    try {
      const parsed = JSON.parse(storedHistory) as Array<{
        timestamp: string;
        endpoint: string;
        status: string;
        remainingRequests: number;
      }>;
      return parsed.map((entry) => ({
        ...entry,
        timestamp: new Date(entry.timestamp),
        status: (entry.status === "limited" ? "limited" : "success") as
          | "success"
          | "limited",
      }));
    } catch (e) {
      console.error("Failed to parse rate limit history:", e);
      return [];
    }
  });

  const updateStatus = (endpoint: string, headers: Headers) => {
    const limit = parseInt(headers.get("X-RateLimit-Limit") || "0", 10);
    const remaining = parseInt(headers.get("X-RateLimit-Remaining") || "0", 10);
    const reset = parseInt(headers.get("X-RateLimit-Reset") || "0", 10);

    if (limit === 0) return; // Нет информации о rate limit

    const requestsUsed = limit - remaining;
    const percentageUsed = (requestsUsed / limit) * 100;
    const resetAt = reset
      ? new Date(reset * 1000)
      : new Date(Date.now() + 60000);

    const status: RateLimitStatus = {
      endpoint,
      requestsUsed,
      requestsLimit: limit,
      percentageUsed,
      resetAt,
      isNearLimit: percentageUsed >= 80,
      isLimited: remaining === 0,
    };

    setStatuses((prev) => {
      const next = new Map(prev);
      next.set(endpoint, status);

      // Сохраняем в localStorage
      const obj = Object.fromEntries(next.entries());
      localStorage.setItem("rateLimitStatuses", JSON.stringify(obj));

      return next;
    });

    // Добавляем в историю
    const historyEntry: RateLimitHistoryEntry = {
      timestamp: new Date(),
      endpoint,
      status: remaining === 0 ? "limited" : "success",
      remainingRequests: remaining,
    };

    setHistory((prev) => {
      const next = [historyEntry, ...prev].slice(0, 100); // Храним последние 100
      localStorage.setItem("rateLimitHistory", JSON.stringify(next));
      return next;
    });
  };

  const getStatus = (endpoint: string): RateLimitStatus | null => {
    return statuses.get(endpoint) || null;
  };

  const getAllStatuses = (): RateLimitStatus[] => {
    return Array.from(statuses.values());
  };

  const clearHistory = () => {
    setHistory([]);
    localStorage.removeItem("rateLimitHistory");
  };

  return {
    statuses: getAllStatuses(),
    history,
    updateStatus,
    getStatus,
    clearHistory,
  };
}

/**
 * Hook для получения предупреждений о приближении к лимиту
 */
export function useRateLimitWarnings() {
  const { statuses } = useRateLimitStatus();

  const warnings = useMemo(() => {
    return statuses.filter((s) => s.isNearLimit && !s.isLimited);
  }, [statuses]);

  return warnings;
}

/**
 * Hook для проверки, не заблокирован ли endpoint
 */
export function useIsEndpointLimited(endpoint: string): boolean {
  const { getStatus } = useRateLimitStatus();
  const status = getStatus(endpoint);
  return status?.isLimited ?? false;
}
