/**
 * useAsync.ts - Async operations hooks
 * Reusable hooks for managing async state
 */

"use client";

import { useState, useCallback, useRef, useEffect } from "react";
import {
  APIResponse,
  isSuccessResponse,
  isErrorResponse,
  parseApiError,
} from "@/types/api-response";

/**
 * State for async operation
 */
export interface AsyncState<T, E = Record<string, string[]>> {
  data: T | null;
  loading: boolean;
  error: E | null;
}

/**
 * Options for useAsync hook
 */
export interface UseAsyncOptions<T> {
  immediate?: boolean;
  retries?: number;
  retryDelay?: number;
  onError?: (error: unknown) => void;
  onSuccess?: (data: T) => void;
  validator?: (data: unknown) => boolean;
  transformer?: (data: unknown) => T;
}

/**
 * useAsync Hook - Generic async operations
 */
export function useAsync<T>(
  asyncFunction: () => Promise<T>,
  options: UseAsyncOptions<T> = {},
): AsyncState<T> & {
  execute: () => Promise<T | null>;
  reset: () => void;
} {
  const {
    immediate = false,
    retries = 0,
    retryDelay = 1000,
    onError,
    onSuccess,
    validator,
    transformer,
  } = options;

  const [state, setState] = useState<AsyncState<T>>({
    data: null,
    loading: false,
    error: null,
  });

  const retriesRef = useRef(0);
  const mountedRef = useRef(true);

  const execute = useCallback(async (): Promise<T | null> => {
    setState({ data: null, loading: true, error: null });

    try {
      let result: T = await asyncFunction();

      if (validator && !validator(result)) {
        throw new Error("Validation failed");
      }

      if (transformer) {
        result = transformer(result);
      }

      if (mountedRef.current) {
        setState({ data: result, loading: false, error: null });
        onSuccess?.(result);
      }

      return result;
    } catch (error) {
      if (retriesRef.current < retries) {
        retriesRef.current++;
        await new Promise((resolve) => setTimeout(resolve, retryDelay));
        return execute();
      }

      const errorMessage =
        error instanceof Error ? error.message : "Unknown error";

      if (mountedRef.current) {
        setState({
          data: null,
          loading: false,
          error: { message: [errorMessage] },
        });
        onError?.(error);
      }

      return null;
    } finally {
      retriesRef.current = 0;
    }
  }, [
    asyncFunction,
    retries,
    retryDelay,
    onError,
    onSuccess,
    validator,
    transformer,
  ]);

  const reset = useCallback(() => {
    setState({ data: null, loading: false, error: null });
  }, []);

  useEffect(() => {
    if (immediate) {
      execute();
    }

    return () => {
      mountedRef.current = false;
    };
  }, [execute, immediate]);

  return { ...state, execute, reset };
}

/**
 * useAsyncAPI Hook - Specialized for API calls
 */
export function useAsyncAPI<T>(
  asyncFunction: () => Promise<APIResponse<T>>,
  options: Omit<UseAsyncOptions<T>, "validator" | "transformer"> = {},
): AsyncState<T> & {
  execute: () => Promise<T | null>;
  reset: () => void;
  refetch: () => Promise<T | null>;
} {
  const [state, setState] = useState<AsyncState<T>>({
    data: null,
    loading: false,
    error: null,
  });

  const retriesRef = useRef(0);
  const mountedRef = useRef(true);
  const {
    immediate = false,
    retries = 1,
    retryDelay = 1000,
    onError,
    onSuccess,
  } = options;

  const execute = useCallback(async (): Promise<T | null> => {
    setState({ data: null, loading: true, error: null });

    try {
      const response = await asyncFunction();

      if (isErrorResponse(response)) {
        throw new Error(response.error?.message || "API Error");
      }

      if (!isSuccessResponse(response) || response.data === null) {
        throw new Error("Invalid response");
      }

      if (mountedRef.current) {
        setState({ data: response.data, loading: false, error: null });
        onSuccess?.(response.data);
      }

      return response.data;
    } catch (error) {
      if (retriesRef.current < retries) {
        retriesRef.current++;
        await new Promise((resolve) => setTimeout(resolve, retryDelay));
        return execute();
      }

      const errorInfo = parseApiError(error);

      if (mountedRef.current) {
        setState({
          data: null,
          loading: false,
          error: { message: [errorInfo.message] },
        });
        onError?.(error);
      }

      return null;
    } finally {
      retriesRef.current = 0;
    }
  }, [asyncFunction, retries, retryDelay, onError, onSuccess]);

  const reset = useCallback(() => {
    setState({ data: null, loading: false, error: null });
  }, []);

  const refetch = useCallback(async (): Promise<T | null> => {
    return execute();
  }, [execute]);

  useEffect(() => {
    if (immediate) {
      execute();
    }

    return () => {
      mountedRef.current = false;
    };
  }, [execute, immediate]);

  return { ...state, execute, reset, refetch };
}

/**
 * useMutation Hook - For POST/PUT/DELETE operations
 */
export function useMutation<TData, TVariable>(
  mutationFn: (variables: TVariable) => Promise<TData>,
  options: {
    onSuccess?: (data: TData) => void;
    onError?: (error: unknown) => void;
  } = {},
): {
  mutate: (variables: TVariable) => Promise<TData | null>;
  loading: boolean;
  error: unknown;
  reset: () => void;
} {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<unknown>(null);

  const mutate = useCallback(
    async (variables: TVariable): Promise<TData | null> => {
      setLoading(true);
      setError(null);

      try {
        const result = await mutationFn(variables);
        options.onSuccess?.(result);
        return result;
      } catch (err) {
        setError(err);
        options.onError?.(err);
        return null;
      } finally {
        setLoading(false);
      }
    },
    [mutationFn, options],
  );

  const reset = useCallback(() => {
    setLoading(false);
    setError(null);
  }, []);

  return { mutate, loading, error, reset };
}

/**
 * useDebounce Hook - Debounce a value
 */
export function useDebounce<T>(value: T, delay: number = 500): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => clearTimeout(handler);
  }, [value, delay]);

  return debouncedValue;
}
