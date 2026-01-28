/**
 * API Client Exports
 * Централизованный экспорт всех API клиентов
 */

// ESF Invoice API (Priority 1) - новое
export { InvoiceAPIClient, invoiceAPI, InvoiceAPIError } from "./invoice";
export type {
  Invoice,
  InvoiceDetail,
  CreateInvoiceRequest,
  UpdateInvoiceRequest,
  InvoiceFilters,
  InvoiceOperationResponse,
} from "./invoice";

// Существующие API клиенты - экспортируем только основные функции
export * from "./invoices"; // Внутренний invoice API

// Импортируем типы из основных модулей, используя namespace для избежания конфликтов
export type * from "./bank-accounts";
export type * from "./catalog";
export type * from "./documents";
export type * from "./foreign-companies";
export type * from "./users";
export type * from "./dashboard";
