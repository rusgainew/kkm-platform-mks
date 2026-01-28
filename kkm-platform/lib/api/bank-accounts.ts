"use client";

/**
 * API клиент для работы с банковскими счетами
 * Использует типы из @/types/entities для полного соответствия backend API
 *
 * @version 2.0
 * @date 2026-01-27
 */

import { bearerAuth } from "@/lib/auth/bearer";
import { getApiBase } from "./client";
import type {
  BankAccount,
  CreateBankAccountRequest,
  UpdateBankAccountRequest,
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

  console.log("[apiRequest:bank-accounts] URL:", url);

  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...authHeaders,
      ...(options?.headers || {}),
    },
    credentials: "include",
  });

  console.log("[apiRequest:bank-accounts] Статус ответа:", response.status);

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
 * List bank accounts with pagination
 */
export async function listBankAccounts(
  params?: PaginationRequest,
): Promise<ListResponse<BankAccount>> {
  const searchParams = new URLSearchParams({
    page: (params?.page || 1).toString(),
    page_size: (params?.page_size || 20).toString(),
  });

  return apiRequest<ListResponse<BankAccount>>(
    `/bank-accounts?${searchParams.toString()}`,
    {
      method: "GET",
    },
  );
}

/**
 * Get a single bank account by ID
 */
export async function getBankAccount(
  accountId: string,
): Promise<APIResponse<BankAccount>> {
  return apiRequest<APIResponse<BankAccount>>(`/bank-accounts/${accountId}`, {
    method: "GET",
  });
}

/**
 * Create a new bank account
 */
export async function createBankAccount(
  data: CreateBankAccountRequest,
): Promise<APIResponse<BankAccount>> {
  return apiRequest<APIResponse<BankAccount>>("/bank-accounts", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

/**
 * Update a bank account
 */
export async function updateBankAccount(
  accountId: string,
  data: UpdateBankAccountRequest,
): Promise<APIResponse<BankAccount>> {
  return apiRequest<APIResponse<BankAccount>>(`/bank-accounts/${accountId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

/**
 * Delete a bank account
 */
export async function deleteBankAccount(
  accountId: string,
): Promise<APIResponse<void>> {
  return apiRequest<APIResponse<void>>(`/bank-accounts/${accountId}`, {
    method: "DELETE",
  });
}

/**
 * Toggle bank account active status
 */
export async function toggleBankAccountStatus(
  accountId: string,
  isActive: boolean,
): Promise<APIResponse<BankAccount>> {
  return apiRequest<APIResponse<BankAccount>>(
    `/bank-accounts/${accountId}/status`,
    {
      method: "PATCH",
      body: JSON.stringify({ is_active: isActive }),
    },
  );
}

/**
 * Get bank accounts by organization
 */
export async function getBankAccountsByOrganization(
  organizationId: string,
): Promise<ListResponse<BankAccount>> {
  return apiRequest<ListResponse<BankAccount>>(
    `/organizations/${organizationId}/bank-accounts`,
    {
      method: "GET",
    },
  );
}
