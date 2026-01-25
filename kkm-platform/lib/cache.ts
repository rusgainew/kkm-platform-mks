/**
 * Клиентский кеш с TTL для оптимизации
 */

interface CacheEntry<T> {
  value: T;
  timestamp: number;
  ttl: number; // в миллисекундах
}

class Cache<T = any> {
  private cache: Map<string, CacheEntry<T>> = new Map();

  /**
   * Установить значение в кеш
   * @param key - ключ
   * @param value - значение
   * @param ttl - время жизни в миллисекундах (по умолчанию 5 минут)
   */
  set(key: string, value: T, ttl: number = 5 * 60 * 1000) {
    this.cache.set(key, {
      value,
      timestamp: Date.now(),
      ttl,
    });
  }

  /**
   * Получить значение из кеша
   * @returns значение или null если истекло время жизни
   */
  get(key: string): T | null {
    const entry = this.cache.get(key);
    if (!entry) return null;

    // Проверяем TTL
    if (Date.now() - entry.timestamp > entry.ttl) {
      this.cache.delete(key);
      return null;
    }

    return entry.value;
  }

  /**
   * Проверить наличие в кеше
   */
  has(key: string): boolean {
    return this.get(key) !== null;
  }

  /**
   * Удалить из кеша
   */
  delete(key: string): boolean {
    return this.cache.delete(key);
  }

  /**
   * Очистить весь кеш
   */
  clear() {
    this.cache.clear();
  }

  /**
   * Получить размер кеша
   */
  size(): number {
    return this.cache.size;
  }

  /**
   * Получить или установить значение
   */
  getOrSet(key: string, fn: () => T, ttl?: number): T {
    const cached = this.get(key);
    if (cached !== null) return cached;

    const value = fn();
    this.set(key, value, ttl);
    return value;
  }

  /**
   * Получить статистику
   */
  getStats() {
    let expired = 0;
    const now = Date.now();

    this.cache.forEach((entry) => {
      if (now - entry.timestamp > entry.ttl) {
        expired++;
      }
    });

    return {
      total: this.cache.size,
      expired,
      valid: this.cache.size - expired,
    };
  }
}

// Глобальный кеш для API запросов
export const apiCache = new Cache();

// Кеш для пользовательских данных (дольше живет)
export const userDataCache = new Cache();

// Кеш для каталога (статических данных)
export const catalogCache = new Cache();

/**
 * Хук для кеширования в локальное хранилище
 */
export function useLocalStorageCache(key: string, defaultValue: any = null) {
  if (typeof window === 'undefined') return defaultValue;

  try {
    const item = window.localStorage.getItem(key);
    return item ? JSON.parse(item) : defaultValue;
  } catch (error) {
    console.warn(`Failed to read localStorage key "${key}":`, error);
    return defaultValue;
  }
}

/**
 * Сохранить в локальное хранилище
 */
export function setLocalStorageCache(key: string, value: any, ttl?: number) {
  if (typeof window === 'undefined') return;

  try {
    const data = {
      value,
      timestamp: Date.now(),
      ttl: ttl || 7 * 24 * 60 * 60 * 1000, // 7 дней по умолчанию
    };
    window.localStorage.setItem(key, JSON.stringify(data));
  } catch (error) {
    console.warn(`Failed to set localStorage key "${key}":`, error);
  }
}

/**
 * Получить из локального хранилища с проверкой TTL
 */
export function getLocalStorageCache(key: string) {
  if (typeof window === 'undefined') return null;

  try {
    const item = window.localStorage.getItem(key);
    if (!item) return null;

    const data = JSON.parse(item);
    const now = Date.now();

    // Проверяем TTL
    if (now - data.timestamp > data.ttl) {
      window.localStorage.removeItem(key);
      return null;
    }

    return data.value;
  } catch (error) {
    console.warn(`Failed to get localStorage key "${key}":`, error);
    return null;
  }
}

/**
 * IndexedDB кеш для больших объемов данных
 */
export class IndexedDBCache {
  private dbName = 'kkm-cache';
  private storeName = 'data';
  private db: IDBDatabase | null = null;

  async init() {
    return new Promise<void>((resolve, reject) => {
      const request = indexedDB.open(this.dbName, 1);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => {
        this.db = request.result;
        resolve();
      };

      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result;
        if (!db.objectStoreNames.contains(this.storeName)) {
          db.createObjectStore(this.storeName, { keyPath: 'key' });
        }
      };
    });
  }

  async set(key: string, value: any, ttl: number = 24 * 60 * 60 * 1000) {
    if (!this.db) await this.init();

    return new Promise<void>((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], 'readwrite');
      const store = transaction.objectStore(this.storeName);

      const data = {
        key,
        value,
        timestamp: Date.now(),
        ttl,
      };

      store.put(data);
      transaction.oncomplete = () => resolve();
      transaction.onerror = () => reject(transaction.error);
    });
  }

  async get(key: string) {
    if (!this.db) await this.init();

    return new Promise<any>((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], 'readonly');
      const store = transaction.objectStore(this.storeName);
      const request = store.get(key);

      request.onsuccess = () => {
        const entry = request.result;
        if (!entry) {
          resolve(null);
          return;
        }

        // Проверяем TTL
        if (Date.now() - entry.timestamp > entry.ttl) {
          // Удалить если истекло
          this.delete(key);
          resolve(null);
        } else {
          resolve(entry.value);
        }
      };
      request.onerror = () => reject(request.error);
    });
  }

  async delete(key: string) {
    if (!this.db) await this.init();

    return new Promise<void>((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], 'readwrite');
      const store = transaction.objectStore(this.storeName);
      store.delete(key);

      transaction.oncomplete = () => resolve();
      transaction.onerror = () => reject(transaction.error);
    });
  }

  async clear() {
    if (!this.db) await this.init();

    return new Promise<void>((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], 'readwrite');
      const store = transaction.objectStore(this.storeName);
      store.clear();

      transaction.oncomplete = () => resolve();
      transaction.onerror = () => reject(transaction.error);
    });
  }
}

export const indexedDBCache = new IndexedDBCache();
