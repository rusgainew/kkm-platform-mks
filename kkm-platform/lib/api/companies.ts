"use client";

/**
 * API клиент для работы с компаниями
 * Использует типы из @/types/entities для полного соответствия backend API
 *
 * @version 2.0
 * @date 2026-01-27
 */

import { bearerAuth } from "@/lib/auth/bearer";
import { getApiBase } from "./client";
import type {
  Company,
  CreateCompanyRequest,
  UpdateCompanyRequest,
  AddMemberRequest,
  Employee,
  ListResponse,
  APIResponse,
  PaginationRequest,
} from "@/types/entities";

/**
 * Helper function for API requests with auth
 */
async function apiRequest<T>(path: string, options?: RequestInit): Promise<T> {
  const API_BASE = getApiBase();
  const url = `${API_BASE}${path}`;
  const authHeaders = bearerAuth();

  console.log("[apiRequest:companies] URL:", url);

  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...authHeaders,
      ...(options?.headers || {}),
    },
    credentials: "include",
  });

  console.log("[apiRequest:companies] Статус ответа:", response.status);

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
 * List companies with pagination
 */
export async function listCompanies(
  params?: PaginationRequest,
): Promise<ListResponse<Company>> {
  const searchParams = new URLSearchParams({
    page: (params?.page || 1).toString(),
    page_size: (params?.page_size || 20).toString(),
  });

  return apiRequest<ListResponse<Company>>(
    `/companies?${searchParams.toString()}`,
    {
      method: "GET",
    },
  );
}

/**
 * Get a single company by ID
 */
export async function getCompany(
  companyId: string,
): Promise<APIResponse<Company>> {
  return apiRequest<APIResponse<Company>>(`/companies/${companyId}`, {
    method: "GET",
  });
}

/**
 * Create a new company
 */
export async function createCompany(
  data: CreateCompanyRequest,
): Promise<APIResponse<Company>> {
  return apiRequest<APIResponse<Company>>("/companies", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

/**
 * Update a company
 */
export async function updateCompany(
  companyId: string,
  data: UpdateCompanyRequest,
): Promise<APIResponse<Company>> {
  return apiRequest<APIResponse<Company>>(`/companies/${companyId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

/**
 * Delete a company
 */
export async function deleteCompany(
  companyId: string,
): Promise<APIResponse<void>> {
  return apiRequest<APIResponse<void>>(`/companies/${companyId}`, {
    method: "DELETE",
  });
}

/**
 * Get company members list
 */
export async function getCompanyMembers(
  companyId: string,
): Promise<APIResponse<Employee[]>> {
  return apiRequest<APIResponse<Employee[]>>(
    `/companies/${companyId}/members`,
    {
      method: "GET",
    },
  );
}

/**
 * Add a member to company
 */
export async function addCompanyMember(
  companyId: string,
  data: AddMemberRequest,
): Promise<APIResponse<Employee>> {
  return apiRequest<APIResponse<Employee>>(`/companies/${companyId}/members`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}

/**
 * Remove a member from company
 */
export async function removeCompanyMember(
  companyId: string,
  userId: string,
): Promise<APIResponse<void>> {
  return apiRequest<APIResponse<void>>(
    `/companies/${companyId}/members/${userId}`,
    {
      method: "DELETE",
    },
  );
}

/**
 * Update member role in company
 */
export async function updateMemberRole(
  companyId: string,
  userId: string,
  role: string,
): Promise<APIResponse<Employee>> {
  return apiRequest<APIResponse<Employee>>(
    `/companies/${companyId}/members/${userId}`,
    {
      method: "PUT",
      body: JSON.stringify({ role }),
    },
  );
}
