/**
 * Generic API Response Types & Utilities
 * Provides type-safe wrapper for all API responses
 */

import { ApiErrorCode, isValidEnumValue } from "@/types/enums";

/**
 * Generic API Response Wrapper
 * Standardizes all API responses across the application
 */
export interface APIResponse<T = unknown> {
  /** Whether the request was successful */
  success: boolean;

  /** Response data (null if error) */
  data: T | null;

  /** Error information (only present if success=false) */
  error?: {
    code: ApiErrorCode | string;
    message: string;
    details?: Record<string, unknown>;
  };

  /** Response metadata */
  meta?: {
    request_id: string;
    timestamp: number;
    version: string;
    path?: string;
  };
}

/**
 * Paginated List Response
 * Used for endpoints that return lists with pagination
 */
export interface ListResponse<T> extends APIResponse<T[]> {
  data: T[];
  meta: {
    request_id: string;
    timestamp: number;
    version: string;
    page: number;
    page_size: number;
    total_count: number;
    total_pages: number;
  };
}

/**
 * Single Item Response
 */
export interface ItemResponse<T> extends APIResponse<T> {
  data: T;
}

/**
 * Boolean Success Response
 * For endpoints that just return success/failure
 */
export interface BooleanResponse extends APIResponse<boolean> {
  data: boolean;
}

/**
 * Error Response
 * Typed error response
 */
export interface ErrorResponse extends APIResponse<null> {
  success: false;
  data: null;
  error: {
    code: ApiErrorCode | string;
    message: string;
    details?: Record<string, unknown>;
  };
}

/**
 * Type guard: Check if response is successful
 */
export function isSuccessResponse<T>(
  response: APIResponse<T>,
): response is APIResponse<T> & { data: T; success: true } {
  return response.success === true && response.data !== null;
}

/**
 * Type guard: Check if response is error
 */
export function isErrorResponse<T>(
  response: APIResponse<T>,
): response is ErrorResponse {
  return response.success === false || response.data === null;
}

/**
 * Extract data from response or throw error
 */
export function getResponseData<T>(response: APIResponse<T>): T {
  if (isErrorResponse(response)) {
    throw new Error(response.error?.message || "Unknown API error");
  }
  if (response.data === null) {
    throw new Error("No data in response");
  }
  return response.data;
}

/**
 * Extract data safely, returning null on error
 */
export function getResponseDataSafe<T>(response: APIResponse<T>): T | null {
  try {
    return getResponseData(response);
  } catch {
    return null;
  }
}

/**
 * Create successful response
 */
export function createSuccessResponse<T>(
  data: T,
  meta?: Omit<APIResponse["meta"], "timestamp">,
): APIResponse<T> {
  return {
    success: true,
    data,
    meta: {
      timestamp: Math.floor(Date.now() / 1000),
      request_id: generateRequestId(),
      version: "1.0.0",
      ...meta,
    },
  };
}

/**
 * Create error response
 */
export function createErrorResponse(
  message: string,
  code: ApiErrorCode | string = "INTERNAL_SERVER_ERROR",
  details?: Record<string, unknown>,
): ErrorResponse {
  return {
    success: false,
    data: null,
    error: {
      code,
      message,
      details,
    },
    meta: {
      timestamp: Math.floor(Date.now() / 1000),
      request_id: generateRequestId(),
      version: "1.0.0",
    },
  };
}

/**
 * Create paginated list response
 */
export function createListResponse<T>(
  data: T[],
  pagination: {
    page: number;
    page_size: number;
    total_count: number;
  },
): ListResponse<T> {
  const total_pages = Math.ceil(pagination.total_count / pagination.page_size);

  return {
    success: true,
    data,
    meta: {
      timestamp: Math.floor(Date.now() / 1000),
      request_id: generateRequestId(),
      version: "1.0.0",
      page: pagination.page,
      page_size: pagination.page_size,
      total_count: pagination.total_count,
      total_pages,
    },
  } as ListResponse<T>;
}

/**
 * Map HTTP status code to error code
 */
export function httpStatusToErrorCode(status: number): ApiErrorCode | string {
  const statusMap: Record<number, ApiErrorCode | string> = {
    400: "BAD_REQUEST",
    401: "UNAUTHORIZED",
    403: "FORBIDDEN",
    404: "NOT_FOUND",
    409: "CONFLICT",
    422: "UNPROCESSABLE_ENTITY",
    500: "INTERNAL_SERVER_ERROR",
    503: "SERVICE_UNAVAILABLE",
  };
  return statusMap[status] || "INTERNAL_SERVER_ERROR";
}

/**
 * Parse API error response
 */
export function parseApiError(error: unknown): {
  code: string;
  message: string;
  details?: Record<string, unknown>;
} {
  if (typeof error === "object" && error !== null) {
    const err = error as Record<string, unknown>;

    if ("response" in err) {
      const response = err.response as Record<string, unknown>;
      if ("data" in response) {
        const data = response.data as Record<string, unknown>;
        if ("error" in data && typeof data.error === "object") {
          const apiError = data.error as Record<string, unknown>;
          return {
            code: (apiError.code as string) || "INTERNAL_SERVER_ERROR",
            message: (apiError.message as string) || "Unknown error",
            details: (apiError.details as Record<string, unknown>) || undefined,
          };
        }
      }
    }

    if ("message" in err) {
      return {
        code: "INTERNAL_SERVER_ERROR",
        message: err.message as string,
      };
    }
  }

  return {
    code: "INTERNAL_SERVER_ERROR",
    message: "An unexpected error occurred",
  };
}

/**
 * Generate unique request ID
 */
function generateRequestId(): string {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * Transform API response with data transformer
 */
export function transformResponse<T, R>(
  response: APIResponse<T>,
  transformer: (data: T) => R,
): APIResponse<R> {
  if (isErrorResponse(response)) {
    return response as unknown as APIResponse<R>;
  }

  if (response.data === null) {
    return response as unknown as APIResponse<R>;
  }

  return {
    ...response,
    data: transformer(response.data),
  };
}

/**
 * Chain multiple response transformers
 */
export function chainTransformers<T, R1, R2>(
  response: APIResponse<T>,
  transformer1: (data: T) => R1,
  transformer2: (data: R1) => R2,
): APIResponse<R2> {
  const result = transformResponse(response, transformer1) as APIResponse<R1>;
  return transformResponse(result, transformer2);
}

/**
 * Retry response with exponential backoff
 */
export async function retryResponse<T>(
  fn: () => Promise<APIResponse<T>>,
  maxRetries: number = 3,
  initialDelayMs: number = 1000,
): Promise<APIResponse<T>> {
  let lastError: unknown;

  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      const response = await fn();
      if (isSuccessResponse(response)) {
        return response;
      }
      lastError = response.error;
    } catch (error) {
      lastError = error;
    }

    if (attempt < maxRetries - 1) {
      const delayMs = initialDelayMs * Math.pow(2, attempt);
      await new Promise((resolve) => setTimeout(resolve, delayMs));
    }
  }

  return createErrorResponse(
    lastError instanceof Error ? lastError.message : "Max retries exceeded",
  );
}

/**
 * Cache API response
 */
export class ResponseCache<T> {
  private cache = new Map<
    string,
    { response: APIResponse<T>; timestamp: number }
  >();
  private ttl: number; // milliseconds

  constructor(ttlSeconds: number = 300) {
    this.ttl = ttlSeconds * 1000;
  }

  get(key: string): APIResponse<T> | null {
    const cached = this.cache.get(key);
    if (!cached) return null;

    if (Date.now() - cached.timestamp > this.ttl) {
      this.cache.delete(key);
      return null;
    }

    return cached.response;
  }

  set(key: string, response: APIResponse<T>): void {
    this.cache.set(key, {
      response,
      timestamp: Date.now(),
    });
  }

  clear(): void {
    this.cache.clear();
  }

  invalidate(pattern: string | RegExp): void {
    const regex = typeof pattern === "string" ? new RegExp(pattern) : pattern;
    for (const key of this.cache.keys()) {
      if (regex.test(key)) {
        this.cache.delete(key);
      }
    }
  }
}

// ============================================================================
// ESF API Специфичные типы ответов (Приоритет 1)
// ============================================================================

/**
 * ESF API Response - специфичный формат ответов от ESF API
 * Отличается от стандартного APIResponse структурой
 *
 * ESF API использует свой формат вместо нашего стандартного
 */
export interface ESFAPIResponse<T = unknown> {
  /** Уникальный ID ответа от ESF API */
  responseId: string;

  /** ID исходного запроса (опционально) */
  requestId?: string;

  /** ID пользователя который сделал запрос (опционально) */
  userId?: string;

  /** ID клиента (опционально) */
  clientId?: string;

  /** Данные ответа */
  data?: T;

  /** Сообщение об ошибке (если есть) */
  message?: string;

  /** Ошибки валидации по полям */
  errors?: Record<string, string | string[]>;

  /** Детали ошибки */
  details?: Record<string, unknown>;
}

/**
 * Адаптер для преобразования ESF ответа в наш стандартный формат
 */
export function adaptESFResponse<T>(
  esfResponse: ESFAPIResponse<T>,
): APIResponse<T> {
  const hasError = !!esfResponse.message || !!esfResponse.errors;

  return {
    success: !hasError,
    data: esfResponse.data ?? null,
    error: hasError
      ? {
          code: "ESF_API_ERROR",
          message: esfResponse.message || "ESF API returned an error",
          details: {
            errors: esfResponse.errors,
            details: esfResponse.details,
          },
        }
      : undefined,
    meta: {
      request_id: esfResponse.requestId || esfResponse.responseId,
      timestamp: Math.floor(Date.now() / 1000),
      version: "1.0.0",
    },
  };
}

/**
 * Типизированные ответы для ESF API операций
 */

/**
 * Ответ на создание счета-фактуры
 */
export interface ESFCreateInvoiceResponse extends ESFAPIResponse {
  data?: {
    responseId: string;
    documentUuid: string;
    invoiceNumber?: string;
    status: string;
  };
}

/**
 * Ответ на получение счета-фактуры
 */
export interface ESFInvoiceResponse extends ESFAPIResponse {
  data?: {
    documentUuid: string;
    invoiceNumber: string;
    status: string;
    totalAmount: number;
    [key: string]: unknown;
  };
}

/**
 * Ответ на список счетов-фактур
 */
export interface ESFInvoiceListResponse extends ESFAPIResponse {
  data?: Array<{
    documentUuid: string;
    invoiceNumber: string;
    status: string;
    totalAmount: number;
    createdDate?: string;
    [key: string]: unknown;
  }>;
}

/**
 * Ответ на операцию с подписью/акцептом/отклонением
 */
export interface ESFInvoiceActionResponse extends ESFAPIResponse {
  data?: {
    documentUuid: string;
    action: string;
    status: string;
    timestamp?: string;
  };
}

/**
 * Ответ на получение справочника
 */
export interface ESFDictionaryResponse extends ESFAPIResponse {
  data?: Array<{
    code: string;
    name: string;
    id?: number;
    [key: string]: unknown;
  }>;
}

/**
 * Ответ на получение каталога товаров
 */
export interface ESFCatalogResponse extends ESFAPIResponse {
  data?: Array<{
    id: number;
    name: string;
    tnvedCode?: string;
    gkedCode?: string;
    unitClassification: { code: string; name: string };
    [key: string]: unknown;
  }>;
}
