"use client";

import { useState } from "react";
import { AlertTriangle, X } from "lucide-react";
import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { useRateLimitWarnings } from "@/hooks/useRateLimit";

export function RateLimitWarningBanner() {
  const warnings = useRateLimitWarnings();
  const [dismissed, setDismissed] = useState<Set<string>>(new Set());
  const [isMounted, setIsMounted] = useState(false);

  // Используем useLayoutEffect для синхронной инициализации перед рендером
  if (typeof window !== "undefined" && !isMounted) {
    setIsMounted(true);
  }

  if (!isMounted || warnings.length === 0) return null;

  const activeWarnings = warnings.filter((w) => !dismissed.has(w.endpoint));

  if (activeWarnings.length === 0) return null;

  return (
    <div className="fixed top-4 right-4 z-50 max-w-md space-y-2">
      {activeWarnings.map((warning) => (
        <Alert
          key={warning.endpoint}
          variant="warning"
          title="Приближение к лимиту API"
          icon={<AlertTriangle className="h-4 w-4" />}
          dismissible
          onDismiss={() => {
            setDismissed((prev) => new Set(prev).add(warning.endpoint));
          }}
        >
          <p className="text-sm mb-2">
            Endpoint: <code className="font-mono text-xs">{warning.endpoint}</code>
          </p>
          <p className="text-sm">
            Использовано: {warning.requestsUsed} / {warning.requestsLimit} (
            {warning.percentageUsed.toFixed(0)}%)
          </p>
        </Alert>
      ))}
    </div>
  );
}
