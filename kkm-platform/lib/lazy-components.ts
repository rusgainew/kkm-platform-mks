import dynamic from "next/dynamic";
import React from "react";

/**
 * Skeleton loader для плейсхолдера
 */
export function Skeleton({ className = "" }: { className?: string }) {
  return React.createElement("div", {
    className: `animate-pulse bg-gray-200 rounded-lg ${className}`,
  });
}

// ===== ДИНАМИЧЕСКИЕ КОМПОНЕНТЫ =====

// Charts - загружаются при необходимости (SSR disabled)
export const DynamicRevenueChart = dynamic(
  () =>
    import("@/components/dashboard/Charts").then((mod) => ({
      default: mod.RevenueChart,
    })),
  {
    ssr: false,
    loading: () => React.createElement(Skeleton, { className: "h-96 w-full" }),
  }
);

export const DynamicInvoiceStatusChart = dynamic(
  () =>
    import("@/components/dashboard/Charts").then((mod) => ({
      default: mod.InvoiceStatusChart,
    })),
  {
    ssr: false,
    loading: () => React.createElement(Skeleton, { className: "h-96 w-full" }),
  }
);

export const DynamicInventoryChart = dynamic(
  () =>
    import("@/components/dashboard/Charts").then((mod) => ({
      default: mod.InventoryChart,
    })),
  {
    ssr: false,
    loading: () => React.createElement(Skeleton, { className: "h-96 w-full" }),
  }
);

// Filters - загружаются при надобности (с SSR)
export const DynamicAdvancedFilter = dynamic(
  () =>
    import("@/components/filters/AdvancedFilter").then((mod) => ({
      default: mod.AdvancedFilter,
    })),
  {
    ssr: true,
    loading: () => React.createElement(Skeleton, { className: "h-16 w-full" }),
  }
);

// Export - загружаются при надобности
export const DynamicExportButtons = dynamic(
  () =>
    import("@/components/export/ExportButtons").then((mod) => ({
      default: mod.ExportButtons,
    })),
  {
    ssr: true,
    loading: () => React.createElement(Skeleton, { className: "h-12 w-32" }),
  }
);

// Bulk Actions - загружаются при надобности
export const DynamicBulkActions = dynamic(
  () =>
    import("@/components/bulk/BulkActions").then((mod) => ({
      default: mod.BulkActions,
    })),
  {
    ssr: true,
    loading: () => React.createElement(Skeleton, { className: "h-12 w-40" }),
  }
);

// Advanced Dashboard - полная загрузка по требованию
export const DynamicOptimizedDashboard = dynamic(
  () =>
    import("@/components/dashboard/OptimizedDashboard").then((mod) => ({
      default: mod.OptimizedDashboard,
    })),
  {
    ssr: false,
    loading: () =>
      React.createElement(
        "div",
        { className: "space-y-6" },
        React.createElement(Skeleton, { className: "h-32 w-full" }),
        React.createElement(Skeleton, { className: "h-96 w-full" })
      ),
  }
);

/**
 * Обертка для ленивой загрузки компонента
 */
export function withLazyLoad(
  importFn: () => Promise<any>,
  options: { fallback?: React.ReactNode; ssr?: boolean } = {}
) {
  return dynamic(importFn, {
    loading: () =>
      options.fallback ||
      React.createElement(Skeleton, { className: "h-40 w-full" }),
    ssr: options.ssr ?? true,
  });
}

/**
 * Обертка для Suspense + динамического импорта
 */
export function LazyComponent({
  component: Component,
  fallback,
  ...props
}: {
  component: React.ComponentType<any>;
  fallback?: React.ReactNode;
  [key: string]: any;
}) {
  const defaultFallback = React.createElement(Skeleton, {
    className: "h-40 w-full",
  });

  return React.createElement(
    React.Suspense,
    { fallback: fallback || defaultFallback },
    React.createElement(Component, props)
  );
}
