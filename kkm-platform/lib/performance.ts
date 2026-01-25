import { useMemo, useCallback, memo, useRef, useEffect, useState } from "react";

/**
 * Оптимизировать функцию с мемоизацией
 */
export function useMemoizedCallback<T extends (...args: any[]) => any>(
  callback: T,
  deps: any[],
): T {
  return useCallback(callback, deps) as T;
}

/**
 * Оптимизировать вычисления
 */
export function useMemoizedValue<T>(factory: () => T, deps: any[]): T {
  return useMemo(factory, deps);
}

/**
 * Обернуть компонент для оптимизации (memo)
 */
export function withMemo<P extends object>(
  Component: React.ComponentType<P>,
  areEqual?: (prevProps: P, nextProps: P) => boolean,
) {
  return memo(Component, areEqual);
}

/**
 * Debounce hook для оптимизации
 */
export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => clearTimeout(handler);
  }, [value, delay]);

  return debouncedValue;
}

/**
 * Throttle hook для частых событий
 */
export function useThrottle<T>(value: T, interval: number): T {
  const [throttledValue, setThrottledValue] = useState<T>(value);
  const lastRanRef = useRef<number>(0);

  useEffect(() => {
    if (lastRanRef.current === 0) {
      lastRanRef.current = Date.now();
    }
    const now = Date.now();
    if (now >= lastRanRef.current + interval) {
      lastRanRef.current = now;
      setThrottledValue(value);
    } else {
      const handler = setTimeout(() => {
        lastRanRef.current = Date.now();
        setThrottledValue(value);
      }, interval);

      return () => clearTimeout(handler);
    }
  }, [value, interval]);

  return throttledValue;
}

/**
 * Lazy load данные при скролле
 */
export function useIntersectionObserver<T extends HTMLElement>(
  ref: React.RefObject<T>,
  options: IntersectionObserverInit = {},
) {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting) {
        setIsVisible(true);
        observer.unobserve(entry.target);
      }
    }, options);

    if (ref.current) {
      observer.observe(ref.current);
    }

    return () => observer.disconnect();
  }, [ref, options]);

  return isVisible;
}

/**
 * Оптимизация таблиц с виртуализацией (простая версия)
 */
export function useVirtualization<T>(
  items: T[],
  itemHeight: number,
  containerHeight: number,
) {
  const [scrollOffset, setScrollOffset] = useState(0);

  const visibleRange = useMemo(() => {
    const startIndex = Math.floor(scrollOffset / itemHeight);
    const endIndex = Math.ceil((scrollOffset + containerHeight) / itemHeight);
    return {
      startIndex: Math.max(0, startIndex - 1),
      endIndex: Math.min(items.length, endIndex + 1),
    };
  }, [scrollOffset, itemHeight, containerHeight, items.length]);

  const visibleItems = useMemo(() => {
    return items.slice(visibleRange.startIndex, visibleRange.endIndex);
  }, [items, visibleRange]);

  const offsetY = visibleRange.startIndex * itemHeight;

  return {
    visibleItems,
    offsetY,
    totalHeight: items.length * itemHeight,
    onScroll: (offset: number) => setScrollOffset(offset),
  };
}

/**
 * Кэширование с TTL (time to live)
 */
export function createCache<K, V>(ttl: number = 60000) {
  const cache = new Map<K, { value: V; timestamp: number }>();

  return {
    get: (key: K): V | null => {
      const item = cache.get(key);
      if (!item) return null;

      if (Date.now() - item.timestamp > ttl) {
        cache.delete(key);
        return null;
      }

      return item.value;
    },

    set: (key: K, value: V) => {
      cache.set(key, { value, timestamp: Date.now() });
    },

    clear: () => cache.clear(),
  };
}

// Импорт React для useRef, useState, useEffect
import React from "react";
