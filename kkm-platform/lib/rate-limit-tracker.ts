/**
 * Middleware для отслеживания Rate Limit из HTTP заголовков
 */

export class RateLimitTracker {
  private static instance: RateLimitTracker;
  private listeners: Set<(endpoint: string, headers: Headers) => void> =
    new Set();

  private constructor() {}

  static getInstance(): RateLimitTracker {
    if (!RateLimitTracker.instance) {
      RateLimitTracker.instance = new RateLimitTracker();
    }
    return RateLimitTracker.instance;
  }

  /**
   * Регистрирует слушателя для обновлений rate limit
   */
  subscribe(listener: (endpoint: string, headers: Headers) => void) {
    this.listeners.add(listener);
    return () => {
      this.listeners.delete(listener);
    };
  }

  /**
   * Обрабатывает ответ и извлекает информацию о rate limit
   */
  trackResponse(endpoint: string, response: Response) {
    const headers = response.headers;

    // Проверяем наличие заголовков rate limit
    if (
      headers.has("X-RateLimit-Limit") ||
      headers.has("X-RateLimit-Remaining") ||
      headers.has("X-RateLimit-Reset")
    ) {
      this.notifyListeners(endpoint, headers);
    }
  }

  private notifyListeners(endpoint: string, headers: Headers) {
    this.listeners.forEach((listener) => {
      try {
        listener(endpoint, headers);
      } catch (error) {
        console.error("Error in rate limit listener:", error);
      }
    });
  }
}

/**
 * Fetch wrapper с автоматическим отслеживанием rate limits
 */
export async function fetchWithRateLimitTracking(
  url: string,
  options?: RequestInit,
): Promise<Response> {
  const response = await fetch(url, options);

  // Извлекаем endpoint из URL
  const urlObj = new URL(url, window.location.origin);
  const endpoint = urlObj.pathname;

  // Отслеживаем rate limit
  const tracker = RateLimitTracker.getInstance();
  tracker.trackResponse(endpoint, response);

  return response;
}

/**
 * Декоратор для API функций с автоматическим отслеживанием rate limits
 */
export function withRateLimitTracking<
  T extends (...args: unknown[]) => Promise<unknown>,
>(fn: T, endpoint: string): T {
  return (async (...args: Parameters<T>) => {
    const result = await fn(...args);

    // Если result - это Response, обрабатываем его
    if (result instanceof Response) {
      const tracker = RateLimitTracker.getInstance();
      tracker.trackResponse(endpoint, result);
    }

    return result;
  }) as T;
}
