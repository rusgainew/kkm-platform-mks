"use client";

import dynamic from "next/dynamic";
import React, { ReactNode } from "react";

interface DynamicImportOptions {
  loading?: () => ReactNode;
  ssr?: boolean;
}

/**
 * Создать динамический компонент с loading fallback
 */
export function createDynamicComponent(
  importFn: () => Promise<any>,
  options: DynamicImportOptions = {}
) {
  return dynamic(importFn, {
    loading:
      options.loading ||
      (() =>
        React.createElement(
          "div",
          { className: "flex items-center justify-center h-64" },
          React.createElement("div", {
            className:
              "animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600",
          })
        )),
    ssr: options.ssr ?? true,
  });
}

/**
 * Динамические компоненты (большие, используются редко)
 */

// Графики и аналитика
export const DynamicCharts = dynamic(
  () =>
    import("@/components/dashboard/Charts").then((mod) => ({
      default: mod.RevenueChart,
    })),
  {
    ssr: false,
    loading: () =>
      React.createElement("div", {
        className: "h-80 bg-gray-200 rounded animate-pulse",
      }),
  }
);

export const DynamicPieChart = dynamic(
  () =>
    import("@/components/dashboard/Charts").then((mod) => ({
      default: mod.InvoiceStatusChart,
    })),
  { ssr: false }
);

export const DynamicInventoryChart = dynamic(
  () =>
    import("@/components/dashboard/Charts").then((mod) => ({
      default: mod.InventoryChart,
    })),
  { ssr: false }
);

// Фильтры
export const DynamicAdvancedFilter = dynamic(
  () =>
    import("@/components/filters/AdvancedFilter").then((mod) => ({
      default: mod.AdvancedFilter,
    })),
  { ssr: true }
);

// Экспорт
export const DynamicExportButtons = dynamic(
  () =>
    import("@/components/export/ExportButtons").then((mod) => ({
      default: mod.ExportButtons,
    })),
  { ssr: true }
);

// Bulk actions
export const DynamicBulkActions = dynamic(
  () =>
    import("@/components/bulk/BulkActions").then((mod) => ({
      default: mod.BulkActions,
    })),
  { ssr: true }
);

/**
 * Динамический импорт с retry для надежности
 */
export async function dynamicImportWithRetry<T>(
  importFn: () => Promise<T>,
  maxRetries = 3
): Promise<T> {
  let lastError: Error | null = null;

  for (let i = 0; i < maxRetries; i++) {
    try {
      return await importFn();
    } catch (error) {
      lastError = error as Error;
      // Экспоненциальная задержка перед повтором
      await new Promise((resolve) =>
        setTimeout(resolve, Math.pow(2, i) * 1000)
      );
    }
  }

  throw lastError || new Error("Failed to import after retries");
}
