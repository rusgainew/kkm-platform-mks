/**
 * API Validators
 * Runtime type validation for API responses and requests
 */

import {
  UserStatus,
  UserRole,
  DocumentStatus,
  DocumentType,
  InvoiceStatus,
  CurrencyCode,
  CompanyStatus,
  isValidEnumValue,
} from "@/types/enums";

/**
 * Validate email format
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * Validate UUID format (v4)
 */
export function isValidUUID(uuid: string): boolean {
  const uuidRegex =
    /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
  return uuidRegex.test(uuid);
}

/**
 * Validate phone number (basic format)
 */
export function isValidPhone(phone: string): boolean {
  // Accepts formats like +7-123-456-7890 or 1234567890
  const phoneRegex = /^[\d\-\+\(\)\s]{6,20}$/;
  return phoneRegex.test(phone);
}

/**
 * Validate Unix timestamp
 */
export function isValidTimestamp(timestamp: number): boolean {
  return Number.isInteger(timestamp) && timestamp > 0;
}

/**
 * Validate pagination parameters
 */
export function isValidPageParams(page: number, pageSize: number): boolean {
  return (
    Number.isInteger(page) &&
    Number.isInteger(pageSize) &&
    page >= 1 &&
    pageSize >= 1 &&
    pageSize <= 100 // Max page size
  );
}

/**
 * Validate monetary amount
 */
export function isValidAmount(amount: number): boolean {
  return typeof amount === "number" && amount >= 0 && isFinite(amount);
}

/**
 * Validate User object
 */
export function isValidUser(user: any): boolean {
  return (
    typeof user === "object" &&
    typeof user.user_id === "string" &&
    isValidUUID(user.user_id) &&
    typeof user.email === "string" &&
    isValidEmail(user.email) &&
    typeof user.is_active === "boolean" &&
    isValidTimestamp(user.created_at) &&
    isValidTimestamp(user.updated_at) &&
    isValidEnumValue(user.role, UserRole)
  );
}

/**
 * Validate Company object
 */
export function isValidCompany(company: any): boolean {
  return (
    typeof company === "object" &&
    typeof company.company_id === "string" &&
    isValidUUID(company.company_id) &&
    typeof company.name === "string" &&
    company.name.length > 0 &&
    typeof company.inn === "string" &&
    company.inn.length > 0 &&
    typeof company.address === "string" &&
    isValidTimestamp(company.created_at) &&
    isValidTimestamp(company.updated_at)
  );
}

/**
 * Validate BankAccount object
 */
export function isValidBankAccount(account: any): boolean {
  return (
    typeof account === "object" &&
    typeof account.account_id === "string" &&
    isValidUUID(account.account_id) &&
    typeof account.company_id === "string" &&
    isValidUUID(account.company_id) &&
    typeof account.account_number === "string" &&
    account.account_number.length > 0 &&
    typeof account.bank_code === "string" &&
    account.bank_code.length > 0 &&
    typeof account.is_default === "boolean" &&
    isValidTimestamp(account.created_at) &&
    isValidTimestamp(account.updated_at)
  );
}

/**
 * Validate Catalog/Product object
 */
export function isValidProduct(product: any): boolean {
  return (
    typeof product === "object" &&
    typeof product.id === "string" &&
    isValidUUID(product.id) &&
    typeof product.name === "string" &&
    product.name.length > 0 &&
    typeof product.price === "number" &&
    isValidAmount(product.price) &&
    typeof product.quantity === "number" &&
    product.quantity >= 0 &&
    typeof product.tnved_code === "string" &&
    isValidTimestamp(product.created_at) &&
    isValidTimestamp(product.updated_at)
  );
}

/**
 * Validate Invoice object
 */
export function isValidInvoice(invoice: any): boolean {
  return (
    typeof invoice === "object" &&
    typeof invoice.id === "string" &&
    isValidUUID(invoice.id) &&
    typeof invoice.invoice_number === "string" &&
    invoice.invoice_number.length > 0 &&
    typeof invoice.company_id === "string" &&
    isValidUUID(invoice.company_id) &&
    typeof invoice.total_amount === "number" &&
    isValidAmount(invoice.total_amount) &&
    typeof invoice.status === "string" &&
    isValidEnumValue(invoice.status, InvoiceStatus) &&
    isValidTimestamp(invoice.created_at) &&
    isValidTimestamp(invoice.updated_at)
  );
}

/**
 * Validate Document object
 */
export function isValidDocument(document: any): boolean {
  return (
    typeof document === "object" &&
    typeof document.id === "string" &&
    isValidUUID(document.id) &&
    typeof document.document_number === "string" &&
    document.document_number.length > 0 &&
    typeof document.company_id === "string" &&
    isValidUUID(document.company_id) &&
    typeof document.status === "string" &&
    isValidEnumValue(document.status, DocumentStatus) &&
    typeof document.created_by === "string" &&
    isValidUUID(document.created_by) &&
    Array.isArray(document.entries) &&
    isValidTimestamp(document.created_at) &&
    isValidTimestamp(document.updated_at) &&
    typeof document.version === "number" &&
    document.version >= 1
  );
}

/**
 * Validate DocumentEntry object
 */
export function isValidDocumentEntry(entry: any): boolean {
  return (
    typeof entry === "object" &&
    typeof entry.id === "string" &&
    isValidUUID(entry.id) &&
    typeof entry.document_id === "string" &&
    isValidUUID(entry.document_id) &&
    typeof entry.line_number === "number" &&
    entry.line_number >= 1 &&
    typeof entry.item_name === "string" &&
    entry.item_name.length > 0 &&
    typeof entry.quantity === "number" &&
    entry.quantity > 0 &&
    typeof entry.amount === "number" &&
    isValidAmount(entry.amount)
  );
}

/**
 * Validate ForeignCompany object
 */
export function isValidForeignCompany(company: any): boolean {
  return (
    typeof company === "object" &&
    typeof company.id === "string" &&
    isValidUUID(company.id) &&
    typeof company.name === "string" &&
    company.name.length > 0 &&
    typeof company.tax_id === "string" &&
    company.tax_id.length > 0 &&
    typeof company.country === "string" &&
    company.country.length === 2 && // ISO 3166-1 alpha-2
    typeof company.address === "string" &&
    company.address.length > 0 &&
    typeof company.contact_email === "string" &&
    isValidEmail(company.contact_email) &&
    typeof company.contact_phone === "string" &&
    isValidPhone(company.contact_phone) &&
    typeof company.is_active === "boolean" &&
    isValidTimestamp(company.created_at) &&
    isValidTimestamp(company.updated_at)
  );
}

/**
 * Validate API Response structure
 */
export function isValidApiResponse(response: any): boolean {
  return (
    typeof response === "object" &&
    typeof response.success === "boolean" &&
    (response.data === null || response.data !== undefined)
  );
}

/**
 * Validate paginated list response
 */
export function isValidPaginatedResponse(response: any): boolean {
  return (
    isValidApiResponse(response) &&
    Array.isArray(response.data) &&
    typeof response.meta === "object" &&
    typeof response.meta.total === "number" &&
    typeof response.meta.page === "number" &&
    typeof response.meta.page_size === "number"
  );
}

/**
 * Validate login credentials
 */
export function isValidLoginCredentials(
  email: string,
  password: string,
): boolean {
  return (
    typeof email === "string" &&
    isValidEmail(email) &&
    typeof password === "string" &&
    password.length >= 6
  );
}

/**
 * Validate password strength
 */
export function isStrongPassword(password: string): boolean {
  // At least 8 chars, 1 uppercase, 1 lowercase, 1 number, 1 special char
  const strongRegex =
    /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
  return strongRegex.test(password);
}

/**
 * Validate INN (Russian Tax ID)
 */
export function isValidINN(inn: string): boolean {
  // Russian INN is either 10 or 12 digits
  const innRegex = /^\d{10}$|^\d{12}$/;
  return innRegex.test(inn);
}

/**
 * Validate TNVED code (Russian tariff code)
 * Format: XXXX.XX.XX where X is digit
 */
export function isValidTNVEDCode(code: string): boolean {
  const tnvedRegex = /^\d{4}\.\d{2}\.\d{2}$/;
  return tnvedRegex.test(code);
}

/**
 * Type guard for User
 */
export function isUser(
  value: unknown,
): value is { user_id: string; email: string } {
  return isValidUser(value);
}

/**
 * Type guard for Company
 */
export function isCompany(
  value: unknown,
): value is { company_id: string; name: string } {
  return isValidCompany(value);
}

/**
 * Type guard for Invoice
 */
export function isInvoice(
  value: unknown,
): value is { id: string; invoice_number: string } {
  return isValidInvoice(value);
}

/**
 * Type guard for Document
 */
export function isDocument(
  value: unknown,
): value is { id: string; document_number: string } {
  return isValidDocument(value);
}

/**
 * Type guard for ForeignCompany
 */
export function isForeignCompany(
  value: unknown,
): value is { id: string; name: string; tax_id: string } {
  return isValidForeignCompany(value);
}

/**
 * Batch validation results
 */
export interface ValidationResult {
  isValid: boolean;
  errors: Record<string, string[]>;
}

/**
 * Create validation error result
 */
export function createValidationError(
  errors: Record<string, string[]>,
): ValidationResult {
  return {
    isValid: false,
    errors,
  };
}

/**
 * Create successful validation result
 */
export function createValidationSuccess(): ValidationResult {
  return {
    isValid: true,
    errors: {},
  };
}

/**
 * Validate user creation payload
 */
export function validateUserCreation(data: any): ValidationResult {
  const errors: Record<string, string[]> = {};

  if (!data.email || !isValidEmail(data.email)) {
    errors.email = ["Invalid email format"];
  }

  if (!data.password || !isStrongPassword(data.password)) {
    errors.password = [
      "Password must be at least 8 characters with uppercase, lowercase, number, and special character",
    ];
  }

  if (!data.role || !isValidEnumValue(data.role, UserRole)) {
    errors.role = [
      `Role must be one of: ${Object.values(UserRole).join(", ")}`,
    ];
  }

  return Object.keys(errors).length === 0
    ? createValidationSuccess()
    : createValidationError(errors);
}

/**
 * Validate company creation payload
 */
export function validateCompanyCreation(data: any): ValidationResult {
  const errors: Record<string, string[]> = {};

  if (!data.name || data.name.length === 0) {
    errors.name = ["Company name is required"];
  }

  if (!data.inn || !isValidINN(data.inn)) {
    errors.inn = ["Invalid INN format (10 or 12 digits)"];
  }

  if (!data.address || data.address.length === 0) {
    errors.address = ["Address is required"];
  }

  return Object.keys(errors).length === 0
    ? createValidationSuccess()
    : createValidationError(errors);
}

/**
 * Validate invoice creation payload
 */
export function validateInvoiceCreation(data: any): ValidationResult {
  const errors: Record<string, string[]> = {};

  if (!data.invoice_number || data.invoice_number.length === 0) {
    errors.invoice_number = ["Invoice number is required"];
  }

  if (!data.company_id || !isValidUUID(data.company_id)) {
    errors.company_id = ["Invalid company ID"];
  }

  if (!data.total_amount || !isValidAmount(data.total_amount)) {
    errors.total_amount = ["Total amount must be a non-negative number"];
  }

  if (!Array.isArray(data.items) || data.items.length === 0) {
    errors.items = ["At least one invoice item is required"];
  }

  return Object.keys(errors).length === 0
    ? createValidationSuccess()
    : createValidationError(errors);
}
