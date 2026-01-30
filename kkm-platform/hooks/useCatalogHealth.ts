"use client";

import { useEffect, useState } from "react";

export interface CatalogHealthStatus {
  status: "ok" | "error";
  backend?: string;
  message?: string;
  error?: string;
  suggestion?: string;
}

export function useCatalogHealth() {
  const [health, setHealth] = useState<CatalogHealthStatus | null>(null);
  const [checking, setChecking] = useState(false);

  const check = async () => {
    setChecking(true);
    try {
      const response = await fetch("/api/v1/health");
      const data = await response.json();
      setHealth(data);
    } catch (error) {
      setHealth({
        status: "error",
        error: error instanceof Error ? error.message : "Unknown error",
        message: "Failed to check catalog health",
      });
    } finally {
      setChecking(false);
    }
  };

  useEffect(() => {
    check();
  }, []);

  return { health, checking, check };
}
