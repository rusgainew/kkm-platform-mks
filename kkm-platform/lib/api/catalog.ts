"use client";

/**
 * API клиент для работы с каталогом товаров/услуг
 * Использует типы из @/types/entities для полного соответствия backend API
 *
 * @version 2.0
 * @date 2026-01-27
 */

import { bearerAuth } from "@/lib/auth/bearer";
import type {
  CatalogItem,
  CreateCatalogItemRequest,
  UpdateCatalogItemRequest,
  ListResponse,
  APIResponse,
  PaginationRequest,
} from "@/types/entities";

const getApiBase = () => {
  if (process.env.NEXT_PUBLIC_API_URL) {
    return process.env.NEXT_PUBLIC_API_URL;
  }
  if (typeof window !== "undefined") {
    const hostname = window.location.hostname;
    const isLocalhost = hostname === "localhost" || hostname === "127.0.0.1";
    const baseUrl = isLocalhost
      ? `http://${hostname}/api/v1`
      : `${window.location.origin}/api/v1`;
    return baseUrl;
  }
  return "http://localhost/api/v1";
};

/**
 * Helper function for API requests with auth
 */
async function apiRequest<T>(path: string, options?: RequestInit): Promise<T> {
  const API_BASE = getApiBase();
  const url = `${API_BASE}${path}`;
  const authHeaders = bearerAuth();

  console.log("[apiRequest:catalog] URL:", url);

  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...authHeaders,
      ...(options?.headers || {}),
    },
    credentials: "include",
  });

  console.log("[apiRequest:catalog] Статус ответа:", response.status);

  if (response.status === 401) {
    throw new Error("Неавторизовано - пожалуйста, повторите вход");
  }

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Ошибка API: ${response.status} ${errorText}`);
  }

  const data = await response.json();
  return data;
}

// ============================================================================
// API METHODS
// ============================================================================

/**
 * List catalog items with pagination
 */
export async function listCatalogItems(
  params?: PaginationRequest,
): Promise<ListResponse<CatalogItem>> {
  const searchParams = new URLSearchParams({
    page: (params?.page || 1).toString(),
    page_size: (params?.page_size || 20).toString(),
  });

  return apiRequest<ListResponse<CatalogItem>>(
    `/catalog?${searchParams.toString()}`,
    {
      method: "GET",
    },
  );
}

/**
 * Get a single catalog item by ID
 */
export async function getCatalogItem(
  itemId: string,
): Promise<APIResponse<CatalogItem>> {
  return apiRequest<APIResponse<CatalogItem>>(`/catalog/${itemId}`, {
    method: "GET",
  });
}

/**
 * Create a new catalog item
 */
export async function createCatalogItem(
  data: CreateCatalogItemRequest,
): Promise<APIResponse<CatalogItem>> {
  return apiRequest<APIResponse<CatalogItem>>("/catalog", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

/**
 * Update a catalog item
 */
export async function updateCatalogItem(
  itemId: string,
  data: UpdateCatalogItemRequest,
): Promise<APIResponse<CatalogItem>> {
  return apiRequest<APIResponse<CatalogItem>>(`/catalog/${itemId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

/**
 * Delete a catalog item
 */
export async function deleteCatalogItem(
  itemId: string,
): Promise<APIResponse<void>> {
  return apiRequest<APIResponse<void>>(`/catalog/${itemId}`, {
    method: "DELETE",
  });
}

/**
 * Search catalog items by name or number
 */
export async function searchCatalogItems(
  query: string,
  params?: PaginationRequest,
): Promise<ListResponse<CatalogItem>> {
  const searchParams = new URLSearchParams({
    q: query,
    page: (params?.page || 1).toString(),
    page_size: (params?.page_size || 20).toString(),
  });

  return apiRequest<ListResponse<CatalogItem>>(
    `/catalog/search?${searchParams.toString()}`,
    {
      method: "GET",
    },
  );
}

/**
 * Get catalog items by TNVED code
 */
export async function getCatalogItemsByTnved(
  tnvedCode: string,
): Promise<ListResponse<CatalogItem>> {
  return apiRequest<ListResponse<CatalogItem>>(`/catalog/tnved/${tnvedCode}`, {
    method: "GET",
  });
}
