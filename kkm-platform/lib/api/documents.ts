"use client";

/**
 * API клиент для работы с документами
 * Поддерживает загрузку файлов, метаданные, статусы и версионирование
 *
 * @version 2.0
 * @date 2026-01-27
 */

import { bearerAuth } from "@/lib/auth/bearer";
import { getApiBase } from "./client";
import type {
  Document,
  DocumentStatus,
  DocumentEntry as DocEntry,
  CreateDocumentRequest,
  UpdateDocumentRequest,
} from "@/types/entities";

/**
 * Helper function for API requests with auth
 */
async function apiRequest<T>(path: string, options?: RequestInit): Promise<T> {
  const API_BASE = getApiBase();
  const url = `${API_BASE}${path}`;
  const authHeaders = bearerAuth();

  console.log("[apiRequest] URL:", url);
  console.log(
    "[apiRequest] Заголовки аутентификации присутствуют:",
    !!authHeaders.Authorization,
  );
  if (authHeaders.Authorization) {
    console.log(
      "[apiRequest] Длина заголовка Authorization:",
      authHeaders.Authorization.length,
    );
  }

  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...authHeaders,
      ...(options?.headers || {}),
    },
    credentials: "include",
  });

  console.log("[apiRequest] Статус ответа:", response.status);
  if (response.status === 401) {
    console.warn("[apiRequest] Got 401 Unauthorized");
    throw new Error("Unauthorized");
  }

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(
      `API Error: ${response.status} ${errorText || response.statusText}`,
    );
  }

  return response.json();
}

/**
 * Helper for multipart file uploads
 */
async function uploadRequest<T>(path: string, formData: FormData): Promise<T> {
  const API_BASE = getApiBase();
  const url = `${API_BASE}${path}`;
  const authHeaders = bearerAuth();

  console.log("[uploadRequest:documents] URL:", url);

  const response = await fetch(url, {
    method: "POST",
    headers: {
      ...authHeaders,
      // Don't set Content-Type for FormData - browser will set it with boundary
    },
    body: formData,
    credentials: "include",
  });

  console.log("[uploadRequest:documents] Статус ответа:", response.status);

  if (response.status === 401) {
    throw new Error("Неавторизовано - пожалуйста, повторите вход");
  }

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Ошибка загрузки: ${response.status} ${errorText}`);
  }

  const data = await response.json();
  return data;
}

// Legacy types for backward compatibility
export interface DocumentEntry {
  id: string;
  document_id: string;
  key: string;
  value: string;
  created_at: number;
  updated_at: number;
}

export type {
  Document,
  DocumentStatus,
  CreateDocumentRequest,
  UpdateDocumentRequest,
};

export interface PageInfo {
  page: number;
  size: number;
  total_count: number;
}

export interface DocumentListResponse {
  documents: Document[];
  page_info: PageInfo;
}

export interface ListDocumentsResponse extends DocumentListResponse {}

// ============================================================================
// API METHODS
// ============================================================================

/**
 * List documents with pagination and filtering
 */
export async function listDocuments(
  page: number = 1,
  pageSize: number = 10,
  status?: string,
): Promise<ListDocumentsResponse | Document[]> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
    ...(status && { status }),
  });

  return apiRequest<ListDocumentsResponse | Document[]>(
    `/documents-query?${params.toString()}`,
    {
      method: "GET",
    },
  );
}

/**
 * Get a single document by ID
 */
export async function getDocument(documentId: string): Promise<Document> {
  return apiRequest<Document>(`/documents-query/${documentId}`, {
    method: "GET",
  });
}

/**
 * Search documents
 */
export async function searchDocuments(
  query: string,
  status?: string,
  documentType?: string,
): Promise<Document[]> {
  const params = new URLSearchParams({
    query,
    ...(status && { status }),
    ...(documentType && { document_type: documentType }),
  });

  return apiRequest<Document[]>(
    `/documents-query/search?${params.toString()}`,
    {
      method: "GET",
    },
  );
}

/**
 * Create a new document
 */
export async function createDocument(
  data: CreateDocumentRequest,
): Promise<Document> {
  return apiRequest<Document>("/documents", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

/**
 * Update a document
 */
export async function updateDocument(
  documentId: string,
  data: UpdateDocumentRequest,
): Promise<Document> {
  return apiRequest<Document>(`/documents/${documentId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

/**
 * Send a document for approval
 */
export async function sendDocument(documentId: string): Promise<Document> {
  return apiRequest<Document>(`/documents/${documentId}/send`, {
    method: "POST",
  });
}

/**
 * Approve a document
 */
export async function approveDocument(documentId: string): Promise<Document> {
  return apiRequest<Document>(`/documents/${documentId}/approve`, {
    method: "POST",
  });
}

/**
 * Reject a document
 */
export async function rejectDocument(
  documentId: string,
  reason?: string,
): Promise<Document> {
  return apiRequest<Document>(`/documents/${documentId}/reject`, {
    method: "POST",
    body: JSON.stringify({ reason }),
  });
}

/**
 * Archive a document
 */
export async function archiveDocument(documentId: string): Promise<Document> {
  return apiRequest<Document>(`/documents/${documentId}/archive`, {
    method: "POST",
  });
}

/**
 * Get pending approval documents
 */
export async function getPendingApprovalDocuments(): Promise<Document[]> {
  return apiRequest<Document[]>("/documents-query/pending-approval", {
    method: "GET",
  });
}

/**
 * Delete a document
 */
export async function deleteDocument(documentId: string): Promise<void> {
  return apiRequest<void>(`/documents/${documentId}`, {
    method: "DELETE",
  });
}

/**
 * Upload file attachment to document
 */
export async function uploadDocumentFile(
  documentId: string,
  file: File,
  metadata?: Record<string, string>,
): Promise<{ file_id: string; url: string }> {
  const formData = new FormData();
  formData.append("file", file);

  if (metadata) {
    formData.append("metadata", JSON.stringify(metadata));
  }

  return uploadRequest<{ file_id: string; url: string }>(
    `/documents/${documentId}/files`,
    formData,
  );
}

/**
 * Download document file
 */
export async function downloadDocumentFile(
  documentId: string,
  fileId: string,
): Promise<Blob> {
  const API_BASE = getApiBase();
  const url = `${API_BASE}/documents/${documentId}/files/${fileId}`;
  const authHeaders = bearerAuth();

  const response = await fetch(url, {
    method: "GET",
    headers: {
      ...authHeaders,
    },
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error(`Ошибка загрузки файла: ${response.status}`);
  }

  return response.blob();
}

/**
 * Get document entries (metadata key-value pairs)
 */
export async function getDocumentEntries(
  documentId: string,
): Promise<DocumentEntry[]> {
  return apiRequest<DocumentEntry[]>(`/documents/${documentId}/entries`, {
    method: "GET",
  });
}

/**
 * Add/Update document entry
 */
export async function upsertDocumentEntry(
  documentId: string,
  key: string,
  value: string,
): Promise<DocumentEntry> {
  return apiRequest<DocumentEntry>(`/documents/${documentId}/entries`, {
    method: "POST",
    body: JSON.stringify({ key, value }),
  });
}
