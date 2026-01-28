/**
 * Централизованные типы для всех сущностей системы
 * Соответствуют API Gateway models (services/api-gateway/internal/domain/models/models.go)
 *
 * @version 1.0
 * @date 2026-01-27
 */

// Import and re-export auth types for User interface
import type { Permission, AuthTokens } from "./auth";
import { ROLE_CONFIGS } from "./auth"; // Value export, not type
export type { Permission, AuthTokens };
export { ROLE_CONFIGS };

// ============================================================================
// USER TYPES
// ============================================================================

export interface User {
  // API fields (from backend)
  user_id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  is_active: boolean;
  created_at: number; // Unix timestamp
  updated_at: number; // Unix timestamp
  status: string;

  // Frontend convenience fields (for components)
  id?: string; // Alias for user_id or extracted from JWT
  name?: string; // Computed from first_name + last_name
  firstName?: string; // Alias for first_name
  lastName?: string; // Alias for last_name
  permissions?: Permission[]; // Computed from role
  storeId?: string; // Optional store association for managers
}

export type UserRole =
  | "admin"
  | "manager"
  | "cashier"
  | "employee"
  | "store_manager";

export interface UpdateUserRequest {
  first_name?: string;
  last_name?: string;
  role?: UserRole;
  is_active?: boolean;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in?: number;
  refresh_expires_in?: number;
  token_type?: string;
  user: User;
  timestamp: number;
}

export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in?: number;
  refresh_expires_in?: number;
  token_type?: string;
  timestamp?: number;
}

// ============================================================================
// COMPANY (ORGANIZATION) TYPES
// ============================================================================

export interface Company {
  company_id: string; // или id для обратной совместимости
  id?: string; // Алиас для company_id
  name: string;
  tin: string; // ИНН
  kpp?: string; // КПП (опционально)
  ogrn?: string; // ОГРН (опционально)
  address: string; // Юридический адрес
  phone: string;
  email: string;
  website?: string;
  description?: string;
  owner_id: string;
  member_count?: number;
  status: CompanyStatus;
  created_at: number;
  updated_at: number;
}

export type CompanyStatus = "active" | "inactive" | "suspended";

export interface CreateCompanyRequest {
  name: string;
  tin: string;
  kpp?: string;
  ogrn?: string;
  address: string;
  phone: string;
  email: string;
  website?: string;
  description?: string;
}

export interface UpdateCompanyRequest {
  name?: string;
  tin?: string;
  kpp?: string;
  ogrn?: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  description?: string;
  status?: CompanyStatus;
}

export interface Employee {
  id: string;
  user_id: string;
  organization_id: string;
  email: string; // Email сотрудника
  first_name: string; // Имя
  last_name: string; // Фамилия
  role: EmployeeRole;
  status: string;
  position?: string;
  department?: string;
  joined_at: number;
  last_active_at?: number;
}

export type EmployeeRole =
  | "owner"
  | "admin"
  | "manager"
  | "employee"
  | "cashier";

export interface AddMemberRequest {
  email: string; // Email нового участника
  role: string; // Роль (преобразуется в EmployeeRole на backend)
}

export interface RemoveMemberRequest {
  organization_id: string;
  member_id: string;
}

// ============================================================================
// CATALOG TYPES
// ============================================================================

export interface CatalogItem {
  id: string;
  name: string;
  number: string; // Номер по каталогу
  description?: string; // Описание
  tnved_code: string; // Код ТНВЭД (обязательно для ЭСФ!)
  category?: string; // @deprecated - использовать tnved_code
  price: number;
  currency: string; // "KGS", "USD", "RUB", "EUR"
  unit: string; // Единица измерения (код ОКЕИ)
  created_at: number;
  updated_at: number;
}

export interface CreateCatalogItemRequest {
  name: string;
  number: string;
  description?: string;
  tnved_code: string; // Обязательно!
  price: number;
  currency: string;
  unit: string;
}

export interface UpdateCatalogItemRequest {
  name?: string;
  number?: string;
  description?: string;
  tnved_code?: string;
  price?: number;
  currency?: string;
  unit?: string;
}

// Справочник единиц измерения (ОКЕИ)
export interface UnitOfMeasure {
  code: string; // Код ОКЕИ
  name: string; // Название
}

export const UNITS_OF_MEASURE: UnitOfMeasure[] = [
  { code: "796", name: "шт" }, // штука
  { code: "006", name: "м" }, // метр
  { code: "055", name: "м²" }, // квадратный метр
  { code: "113", name: "м³" }, // кубический метр
  { code: "166", name: "кг" }, // килограмм
  { code: "112", name: "л" }, // литр
  { code: "212", name: "т" }, // тонна
  { code: "778", name: "компл" }, // комплект
  { code: "715", name: "упак" }, // упаковка
];

// Справочник валют
export interface Currency {
  code: string;
  name: string;
  symbol: string;
}

export const CURRENCIES: Currency[] = [
  { code: "KGS", name: "Кыргызский сом", symbol: "с" },
  { code: "USD", name: "Доллар США", symbol: "$" },
  { code: "RUB", name: "Российский рубль", symbol: "₽" },
  { code: "EUR", name: "Евро", symbol: "€" },
  { code: "CNY", name: "Китайский юань", symbol: "¥" },
];

// ============================================================================
// BANK ACCOUNT TYPES
// ============================================================================

export interface BankAccount {
  id: string;
  account_number: string; // Номер счета (20 цифр для КР)
  bank_name: string;
  bank_code: string; // БИК банка
  currency: string;
  owner_id: string; // ID организации
  is_active: boolean;
  created_at: number;
  updated_at: number;
}

export interface CreateBankAccountRequest {
  account_number: string;
  bank_name: string;
  bank_code: string;
  currency: string;
  owner_id: string;
}

export interface UpdateBankAccountRequest {
  account_number?: string;
  bank_name?: string;
  bank_code?: string;
  currency?: string;
  is_active?: boolean;
}

// ============================================================================
// DOCUMENT TYPES
// ============================================================================

export interface Document {
  id: string;
  organization_id: string;
  title: string;
  content: string;
  status: DocumentStatus;
  created_by: string;
  assigned_to?: string;
  created_at: number;
  updated_at: number;
  status_changed_at: number;
  version: number;
  entries?: DocumentEntry[];
}

export type DocumentStatus =
  | "draft"
  | "pending"
  | "sent"
  | "approved"
  | "rejected"
  | "archived";

export interface DocumentEntry {
  id: string;
  document_id: string;
  key: string;
  value: string;
  created_at: number;
  updated_at: number;
}

export interface CreateDocumentRequest {
  organization_id: string;
  title: string;
  content: string;
  assigned_to?: string;
  entries?: Array<{
    key: string;
    value: string;
  }>;
}

export interface UpdateDocumentRequest {
  title?: string;
  content?: string;
  status?: DocumentStatus;
  assigned_to?: string;
  entries?: Array<{
    key: string;
    value: string;
  }>;
}

// ============================================================================
// FOREIGN COMPANY TYPES
// ============================================================================

export interface ForeignCompany {
  id: number;
  pin: string; // Иностранный ИНН/Tax ID
  full_name: string;
  country_code: string; // ISO 3166-1 alpha-2 (RU, CN, US...)
  address?: string;
  created_at: number;
  updated_at: number;
}

export interface CreateForeignCompanyRequest {
  pin: string;
  full_name: string;
  country_code: string; // Обязательно 2 символа (ISO)
  address?: string;
}

export interface UpdateForeignCompanyRequest {
  pin?: string;
  full_name?: string;
  country_code?: string;
  address?: string;
}

// Справочник стран (популярные)
export interface Country {
  code: string; // ISO 3166-1 alpha-2
  name: string;
  nameEn: string;
}

export const POPULAR_COUNTRIES: Country[] = [
  { code: "RU", name: "Россия", nameEn: "Russia" },
  { code: "CN", name: "Китай", nameEn: "China" },
  { code: "US", name: "США", nameEn: "United States" },
  { code: "KZ", name: "Казахстан", nameEn: "Kazakhstan" },
  { code: "UZ", name: "Узбекистан", nameEn: "Uzbekistan" },
  { code: "TJ", name: "Таджикистан", nameEn: "Tajikistan" },
  { code: "TR", name: "Турция", nameEn: "Turkey" },
  { code: "DE", name: "Германия", nameEn: "Germany" },
  { code: "GB", name: "Великобритания", nameEn: "United Kingdom" },
  { code: "FR", name: "Франция", nameEn: "France" },
];

// ============================================================================
// INVOICE TYPES (Basic - для совместимости с существующим кодом)
// ============================================================================

export interface Invoice {
  id: string;
  invoice_number: string;
  invoice_date: string;
  delivery_date: string;
  total_amount: number;
  is_resident: boolean;
  note?: string;
  status: string;
  created_at: number;
  updated_at: number;
}

export interface InvoiceDetail {
  id: string;
  invoice_uuid: string;
  catalog_code: string;
  quantity: number;
  price: number;
  amount: number;
}

// ============================================================================
// PAGINATION & RESPONSE TYPES
// ============================================================================

export interface PageInfo {
  page: number;
  size: number;
  total_count: number;
}

export interface PaginationRequest {
  page?: number;
  page_size?: number;
}

export interface ListResponse<T> {
  items?: T[];
  users?: T[]; // Для ListUsersResponse
  companies?: T[]; // Для ListCompaniesResponse
  documents?: T[]; // Для ListDocumentsResponse
  page_info: PageInfo;
}

export interface APIResponse<T = any> {
  success: boolean;
  data?: T;
  error?: APIError;
  meta?: MetaData;
}

export interface APIError {
  code: string;
  message: string;
  details?: string;
}

export interface MetaData {
  page?: number;
  page_size?: number;
  total_count?: number;
  total_pages?: number;
}

// ============================================================================
// VALIDATION HELPERS
// ============================================================================

/**
 * Проверка валидности UUID
 */
export function isValidUUID(uuid: string): boolean {
  const uuidRegex =
    /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
  return uuidRegex.test(uuid);
}

/**
 * Проверка валидности email
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * Проверка валидности ТНВЭД кода (10 цифр)
 */
export function isValidTnvedCode(code: string): boolean {
  return /^\d{10}$/.test(code);
}

/**
 * Проверка валидности номера банковского счета (20 цифр для КР)
 */
export function isValidBankAccountNumber(number: string): boolean {
  return /^\d{20}$/.test(number);
}

/**
 * Проверка валидности БИК (9 цифр)
 */
export function isValidBankCode(code: string): boolean {
  return /^\d{9}$/.test(code);
}

/**
 * Проверка валидности country code (ISO 3166-1 alpha-2)
 */
export function isValidCountryCode(code: string): boolean {
  return /^[A-Z]{2}$/.test(code);
}

/**
 * Форматирование суммы с валютой
 */
export function formatAmount(amount: number, currency: string): string {
  const curr = CURRENCIES.find((c) => c.code === currency);
  const symbol = curr?.symbol || currency;
  return `${amount.toLocaleString("ru-KG", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })} ${symbol}`;
}

/**
 * Форматирование даты из timestamp
 */
export function formatTimestamp(timestamp: number): string {
  return new Date(timestamp * 1000).toLocaleString("ru-KG", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

/**
 * Форматирование даты из timestamp (короткий формат)
 */
export function formatDate(timestamp: number): string {
  return new Date(timestamp * 1000).toLocaleDateString("ru-KG", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  });
}

// ============================================================================
// TYPE GUARDS
// ============================================================================

export function isUser(obj: any): obj is User {
  return (
    typeof obj === "object" &&
    obj !== null &&
    typeof obj.user_id === "string" &&
    typeof obj.email === "string"
  );
}

export function isCompany(obj: any): obj is Company {
  return (
    typeof obj === "object" &&
    obj !== null &&
    typeof obj.id === "string" &&
    typeof obj.name === "string" &&
    typeof obj.owner_id === "string"
  );
}

export function isCatalogItem(obj: any): obj is CatalogItem {
  return (
    typeof obj === "object" &&
    obj !== null &&
    typeof obj.id === "string" &&
    typeof obj.name === "string" &&
    typeof obj.tnved_code === "string"
  );
}

export function isBankAccount(obj: any): obj is BankAccount {
  return (
    typeof obj === "object" &&
    obj !== null &&
    typeof obj.id === "string" &&
    typeof obj.account_number === "string"
  );
}

export function isDocument(obj: any): obj is Document {
  return (
    typeof obj === "object" &&
    obj !== null &&
    typeof obj.id === "string" &&
    typeof obj.title === "string" &&
    typeof obj.organization_id === "string"
  );
}

export function isForeignCompany(obj: any): obj is ForeignCompany {
  return (
    typeof obj === "object" &&
    obj !== null &&
    typeof obj.id === "number" &&
    typeof obj.pin === "string" &&
    typeof obj.country_code === "string"
  );
}
