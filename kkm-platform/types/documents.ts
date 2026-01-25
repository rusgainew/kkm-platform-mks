/**
 * Document-related type definitions
 */

import { Document } from "@/lib/api/documents";

export type DocumentStatus =
  | "draft"
  | "sent"
  | "approved"
  | "rejected"
  | "archived";

export interface DocumentContextData {
  documents: Document[];
  isLoading: boolean;
  error: string | null;
  selectedDocument: Document | null;
}

export interface DocumentFilters {
  status?: DocumentStatus;
  searchQuery?: string;
  organizationId?: string;
}
