'use client';

import { useEffect } from 'react';

interface MetricsData {
  route: string;
  timestamp: number;
  // Core Web Vitals
  LCP?: number; // Largest Contentful Paint
  FID?: number; // First Input Delay
  CLS?: number; // Cumulative Layout Shift
  INP?: number; // Interaction to Next Paint
  TTFB?: number; // Time to First Byte
  // Custom metrics
  apiCallDuration?: number;
  renderTime?: number;
}

/**
 * Собирает метрики производительности
 */
export class PerformanceMonitor {
  private metrics: MetricsData[] = [];
  private readonly MAX_METRICS = 100;

  /**
   * Запускает мониторинг Web Vitals
   */
  startWebVitalsMonitoring() {
    if (typeof window === 'undefined') return;

    // LCP - Largest Contentful Paint
    try {
      const lcpObserver = new PerformanceObserver((entryList) => {
        const entries = entryList.getEntries();
        const lastEntry = entries[entries.length - 1];
        this.recordMetric({ LCP: lastEntry.startTime });
      });
      lcpObserver.observe({ entryTypes: ['largest-contentful-paint'] });
    } catch (e) {
      console.warn('LCP observer not supported');
    }

    // CLS - Cumulative Layout Shift
    try {
      let clsValue = 0;
      const clsObserver = new PerformanceObserver((entryList) => {
        for (const entry of entryList.getEntries()) {
          if ((entry as any).hadRecentInput) continue;
          clsValue += (entry as any).value;
          this.recordMetric({ CLS: clsValue });
        }
      });
      clsObserver.observe({ entryTypes: ['layout-shift'] });
    } catch (e) {
      console.warn('CLS observer not supported');
    }

    // INP - Interaction to Next Paint (новый метрик для FID)
    try {
      const inpObserver = new PerformanceObserver((entryList) => {
        const entries = entryList.getEntries();
        const lastEntry = entries[entries.length - 1];
        this.recordMetric({ INP: (lastEntry as any).processingDuration });
      });
      inpObserver.observe({ entryTypes: ['event'] });
    } catch (e) {
      console.warn('INP observer not supported');
    }

    // TTFB - Time to First Byte
    if (typeof window !== 'undefined') {
      const perfData = window.performance.timing;
      const ttfb = perfData.responseStart - perfData.fetchStart;
      if (ttfb > 0) {
        this.recordMetric({ TTFB: ttfb });
      }
    }
  }

  /**
   * Записывает метрику
   */
  recordMetric(data: Partial<MetricsData>) {
    const metric: MetricsData = {
      route: typeof window !== 'undefined' ? window.location.pathname : '',
      timestamp: Date.now(),
      ...data,
    };

    this.metrics.push(metric);
    if (this.metrics.length > this.MAX_METRICS) {
      this.metrics.shift();
    }

    // Отправить на аналитику если нужно
    this.sendToAnalytics(metric);
  }

  /**
   * Измеряет время выполнения функции
   */
  async measureAsync<T>(name: string, fn: () => Promise<T>): Promise<T> {
    const start = performance.now();
    try {
      const result = await fn();
      const duration = performance.now() - start;
      this.recordMetric({ apiCallDuration: duration });
      console.log(`[${name}] completed in ${duration.toFixed(2)}ms`);
      return result;
    } catch (error) {
      const duration = performance.now() - start;
      console.error(`[${name}] failed after ${duration.toFixed(2)}ms`, error);
      throw error;
    }
  }

  /**
   * Получить все метрики
   */
  getMetrics(): MetricsData[] {
    return [...this.metrics];
  }

  /**
   * Очистить метрики
   */
  clearMetrics() {
    this.metrics = [];
  }

  /**
   * Отправить метрику на сервер/аналитику
   */
  private sendToAnalytics(metric: MetricsData) {
    // Реализуй отправку на Sentry, DataDog или другой сервис
    // fetch('/api/metrics', { method: 'POST', body: JSON.stringify(metric) })
  }
}

export const performanceMonitor = new PerformanceMonitor();

/**
 * Hook для мониторинга производительности
 */
export function usePerformanceMonitoring() {
  useEffect(() => {
    performanceMonitor.startWebVitalsMonitoring();
  }, []);

  return performanceMonitor;
}

/**
 * Hook для измерения render time
 */
export function useRenderTime(componentName: string) {
  useEffect(() => {
    const start = performance.now();
    return () => {
      const duration = performance.now() - start;
      performanceMonitor.recordMetric({
        renderTime: duration,
      });
      if (duration > 100) {
        console.warn(
          `[${componentName}] slow render: ${duration.toFixed(2)}ms`
        );
      }
    };
  }, [componentName]);
}

/**
 * Утилита для профилирования
 */
export function profile<T>(name: string, fn: () => T): T {
  if (typeof window === 'undefined') {
    return fn();
  }

  performance.mark(`${name}-start`);
  const result = fn();
  performance.mark(`${name}-end`);
  performance.measure(name, `${name}-start`, `${name}-end`);

  const measure = performance.getEntriesByName(name)[0];
  console.log(`[${name}] ${(measure as PerformanceMeasure).duration.toFixed(2)}ms`);

  return result;
}
