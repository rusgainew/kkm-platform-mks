/**
 * Document status constants and utilities
 */

export const DOCUMENT_STATUS = {
  DRAFT: "draft",
  SENT: "sent",
  APPROVED: "approved",
  REJECTED: "rejected",
  ARCHIVED: "archived",
} as const;

export type DocumentStatus =
  (typeof DOCUMENT_STATUS)[keyof typeof DOCUMENT_STATUS];

export const STATUS_LABELS: Record<DocumentStatus, string> = {
  draft: "Черновик",
  sent: "Отправлено",
  approved: "Одобрено",
  rejected: "Отклонено",
  archived: "В архиве",
};

export const STATUS_COLORS: Record<DocumentStatus, string> = {
  draft: "bg-gray-500/20 text-gray-400",
  sent: "bg-blue-500/20 text-blue-400",
  approved: "bg-green-500/20 text-green-400",
  rejected: "bg-red-500/20 text-red-400",
  archived: "bg-gray-600/20 text-gray-500",
};

export const ORGANIZATIONS = [
  { id: "company-1", name: "ООО Компания 1" },
  { id: "company-2", name: "ООО Компания 2" },
  { id: "company-3", name: "ООО Компания 3" },
];
