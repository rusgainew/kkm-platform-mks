"use client";

import { useEffect } from "react";
import { RateLimitTracker } from "@/lib/rate-limit-tracker";
import { useRateLimitStatus } from "@/hooks/useRateLimit";

/**
 * Провайдер для автоматического отслеживания rate limits
 */
export function RateLimitProvider({ children }: { children: React.ReactNode }) {
  const { updateStatus } = useRateLimitStatus();

  useEffect(() => {
    const tracker = RateLimitTracker.getInstance();

    // Подписываемся на обновления
    const unsubscribe = tracker.subscribe((endpoint, headers) => {
      updateStatus(endpoint, headers);
    });

    return unsubscribe;
  }, [updateStatus]);

  return <>{children}</>;
}
