/**
 * API клиент для работы с иностранными компаниями
 * Соответствует новой API структуре с поддержкой пагинации и единообразного формата ответа
 */

"use client";

import { bearerAuth } from "@/lib/auth/bearer";
import { getApiBase } from "./client";

/**
 * Helper function for API requests with auth
 */
async function apiRequest<T>(path: string, options?: RequestInit): Promise<T> {
  const API_BASE = getApiBase();
  const url = `${API_BASE}${path}`;
  const authHeaders = bearerAuth();

  console.log("[apiRequest:foreign-companies] URL:", url);

  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...authHeaders,
      ...(options?.headers || {}),
    },
    credentials: "include",
  });

  console.log("[apiRequest:foreign-companies] Статус ответа:", response.status);

  if (response.status === 401) {
    throw new Error("Неавторизовано - пожалуйста, повторите вход");
  }

  if (!response.ok) {
    const errorText = await response.text();
    console.error("[apiRequest:foreign-companies] Ошибка:", errorText);
    throw new Error(`API Error: ${response.status} - ${errorText}`);
  }

  return response.json();
}

/**
 * ForeignCompany Interface
 * Represents an international company/vendor in the system
 */
export interface ForeignCompany {
  id: string; // Primary Key (UUID)
  name: string; // Company name
  tax_id: string; // Tax identification number (PIN/TIN)
  country: string; // ISO 3166-1 country code (e.g., "US", "CN")
  address: string; // Physical address
  contact_email: string; // Contact email address
  contact_phone: string; // Contact phone number
  currency: string; // ISO 4217 currency code (e.g., "USD", "EUR")
  is_active: boolean; // Whether company is active
  created_at: number; // Unix timestamp (seconds)
  updated_at: number; // Unix timestamp (seconds)
}

/**
 * Error information structure
 */
export interface ErrorInfo {
  code: string;
  message: string;
  details: string;
}

/**
 * Pagination metadata structure
 */
export interface PaginationMeta {
  page: number;
  page_size: number;
  total_count: number;
  total_pages: number;
}

/**
 * Response for list of foreign companies
 */
export interface ForeignCompaniesListResponse {
  data: ForeignCompany[];
  error?: ErrorInfo;
  meta: PaginationMeta;
  success: boolean;
}

/**
 * Response for single foreign company
 */
export interface ForeignCompanyResponse {
  data: ForeignCompany;
  error?: ErrorInfo;
  success: boolean;
}

/**
 * Request to create a foreign company
 */
export interface CreateForeignCompanyRequest {
  name: string;
  tax_id: string;
  country: string;
  address: string;
  contact_email: string;
  contact_phone: string;
  currency?: string;
}

/**
 * Request to update a foreign company
 */
export interface UpdateForeignCompanyRequest {
  name?: string;
  tax_id?: string;
  country?: string;
  address?: string;
  contact_email?: string;
  contact_phone?: string;
  currency?: string;
  is_active?: boolean;
}

/**
 * List all foreign companies with pagination
 * @param token - Authentication token
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @returns List of foreign companies
 */
export async function listForeignCompanies(
  token?: string,
  page: number = 1,
  pageSize: number = 20,
): Promise<ForeignCompaniesListResponse> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  });

  const response = await apiRequest<ForeignCompaniesListResponse>(
    `/foreign-companies?${params}`,
    {
      method: "GET",
    },
  );

  return response;
}

/**
 * Get a specific foreign company by ID
 * @param id - Foreign company ID (UUID)
 * @param token - Authentication token
 * @returns The foreign company
 */
export async function getForeignCompany(
  id: string,
  token?: string,
): Promise<ForeignCompanyResponse> {
  const response = await apiRequest<ForeignCompanyResponse>(
    `/foreign-companies/${id}`,
    {
      method: "GET",
    },
  );

  return response;
}

/**
 * Create a new foreign company
 * @param data - Foreign company data
 * @param token - Authentication token
 * @returns The created foreign company
 */
export async function createForeignCompany(
  data: CreateForeignCompanyRequest,
  token?: string,
): Promise<ForeignCompanyResponse> {
  const response = await apiRequest<ForeignCompanyResponse>(
    "/foreign-companies",
    {
      method: "POST",
      body: JSON.stringify(data),
    },
  );

  return response;
}

/**
 * Update an existing foreign company
 * @param id - Foreign company ID (UUID)
 * @param data - Foreign company update data
 * @param token - Authentication token
 * @returns The updated foreign company
 */
export async function updateForeignCompany(
  id: string,
  data: UpdateForeignCompanyRequest,
  token?: string,
): Promise<ForeignCompanyResponse> {
  const response = await apiRequest<ForeignCompanyResponse>(
    `/foreign-companies/${id}`,
    {
      method: "PUT",
      body: JSON.stringify(data),
    },
  );

  return response;
}

/**
 * Delete a foreign company
 * @param id - Foreign company ID (UUID)
 * @param token - Authentication token
 * @returns Success status
 */
export async function deleteForeignCompany(
  id: string,
  token?: string,
): Promise<{ success: boolean; error?: ErrorInfo }> {
  const response = await apiRequest<{ success: boolean; error?: ErrorInfo }>(
    `/foreign-companies/${id}`,
    {
      method: "DELETE",
    },
  );

  return response;
}

/**
 * Search foreign companies by name or tax_id
 * @param query - Search query
 * @param token - Authentication token
 * @returns List of matching foreign companies
 */
export async function searchForeignCompanies(
  query: string,
  token?: string,
): Promise<ForeignCompaniesListResponse> {
  const params = new URLSearchParams({
    q: query,
  });

  const response = await apiRequest<ForeignCompaniesListResponse>(
    `/foreign-companies/search?${params}`,
    {
      method: "GET",
    },
  );

  return response;
}

/**
 * Activate a foreign company
 * @param id - Foreign company ID (UUID)
 * @param token - Authentication token
 * @returns The updated foreign company
 */
export async function activateForeignCompany(
  id: string,
  token?: string,
): Promise<ForeignCompanyResponse> {
  const response = await apiRequest<ForeignCompanyResponse>(
    `/foreign-companies/${id}/activate`,
    {
      method: "PATCH",
    },
  );

  return response;
}

/**
 * Deactivate a foreign company
 * @param id - Foreign company ID (UUID)
 * @param token - Authentication token
 * @returns The updated foreign company
 */
export async function deactivateForeignCompany(
  id: string,
  token?: string,
): Promise<ForeignCompanyResponse> {
  const response = await apiRequest<ForeignCompanyResponse>(
    `/foreign-companies/${id}/deactivate`,
    {
      method: "PATCH",
    },
  );

  return response;
}
