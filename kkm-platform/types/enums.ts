/**
 * Enumeration types for KKM Platform
 * Provides type-safe status, role, and other enum definitions
 */

/**
 * User Status Enumeration
 * Defines the different states a user account can have
 */
export enum UserStatus {
  ACTIVE = "active",
  INACTIVE = "inactive",
  SUSPENDED = "suspended",
}

/**
 * User Role Enumeration
 * Defines the different permission levels in the system
 */
export enum UserRole {
  ADMIN = "admin",
  USER = "user",
  NETWORK_ADMIN = "network_admin",
}

/**
 * Document Status Enumeration
 * Defines the lifecycle states of documents
 */
export enum DocumentStatus {
  DRAFT = "draft",
  PENDING = "pending",
  APPROVED = "approved",
  REJECTED = "rejected",
  ARCHIVED = "archived",
}

/**
 * Document Type Enumeration
 * Defines the different types of documents in the system
 */
export enum DocumentType {
  INVOICE = "invoice",
  WAYBILL = "waybill",
  ACT = "act",
}

/**
 * Invoice Status Enumeration
 * Defines the different states an invoice can have
 */
export enum InvoiceStatus {
  DRAFT = "draft",
  SENT = "sent",
  APPROVED = "approved",
  REJECTED = "rejected",
}

/**
 * Invoice Action Status
 * Defines actions that can be performed on an invoice
 */
export enum InvoiceAction {
  SIGN = "sign",
  ACCEPT = "accept",
  REJECT = "reject",
  REVOKE = "revoke",
}

/**
 * Currency Code Enumeration (ISO 4217)
 * Defines supported currencies in the system
 */
export enum CurrencyCode {
  KGS = "KGS", // Киргизский сом (основная валюта)
  RUB = "RUB",
  USD = "USD",
  EUR = "EUR",
  GBP = "GBP",
  JPY = "JPY",
  CNY = "CNY",
  INR = "INR",
  AED = "AED",
  SGD = "SGD",
  HKD = "HKD",
  KZT = "KZT", // Казахстанский тенге
  TRY = "TRY", // Турецкая лира
}

/**
 * Company Status Enumeration
 * Defines the different states a company can have
 */
export enum CompanyStatus {
  ACTIVE = "active",
  INACTIVE = "inactive",
  ARCHIVED = "archived",
}

/**
 * HTTP Method Enumeration
 * Defines the standard HTTP methods
 */
export enum HttpMethod {
  GET = "GET",
  POST = "POST",
  PUT = "PUT",
  PATCH = "PATCH",
  DELETE = "DELETE",
  HEAD = "HEAD",
  OPTIONS = "OPTIONS",
}

/**
 * API Error Code Enumeration
 * Defines standard error codes returned by API
 */
export enum ApiErrorCode {
  BAD_REQUEST = "BAD_REQUEST",
  UNAUTHORIZED = "UNAUTHORIZED",
  FORBIDDEN = "FORBIDDEN",
  NOT_FOUND = "NOT_FOUND",
  CONFLICT = "CONFLICT",
  UNPROCESSABLE_ENTITY = "UNPROCESSABLE_ENTITY",
  INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR",
  SERVICE_UNAVAILABLE = "SERVICE_UNAVAILABLE",
}

/**
 * Sort Order Enumeration
 * Defines the sort direction for list operations
 */
export enum SortOrder {
  ASC = "asc",
  DESC = "desc",
}

/**
 * Pagination Direction Enumeration
 * Defines the direction for pagination
 */
export enum PaginationDirection {
  NEXT = "next",
  PREVIOUS = "previous",
}

/**
 * Type definitions for enum values
 * Use these for strict type checking
 */

export type UserStatusType = `${UserStatus}`;
export type UserRoleType = `${UserRole}`;
export type DocumentStatusType = `${DocumentStatus}`;
export type DocumentTypeType = `${DocumentType}`;
export type InvoiceStatusType = `${InvoiceStatus}`;
export type InvoiceActionType = `${InvoiceAction}`;
export type CurrencyCodeType = `${CurrencyCode}`;
export type CompanyStatusType = `${CompanyStatus}`;
export type ApiErrorCodeType = `${ApiErrorCode}`;
export type SortOrderType = `${SortOrder}`;

/**
 * Union type for all status fields
 */
export type StatusType =
  | UserStatusType
  | DocumentStatusType
  | InvoiceStatusType
  | CompanyStatusType;

/**
 * Helper function to get all enum values as array
 * Usage: getAllValues(UserStatus)
 */
export function getAllEnumValues<T extends Record<string, unknown>>(
  enumObj: T,
): string[] {
  return Object.values(enumObj).filter(
    (value) => typeof value === "string",
  ) as string[];
}

/**
 * Helper function to check if a value is a valid enum value
 * Usage: isValidStatus(value, UserStatus)
 */
export function isValidEnumValue<T extends Record<string, unknown>>(
  value: unknown,
  enumObj: T,
): value is T[keyof T] {
  return getAllEnumValues(enumObj).includes(value as string);
}

/**
 * Helper function to convert string to enum
 * Usage: toEnum('active', UserStatus)
 */
export function toEnum<T extends Record<string, unknown>>(
  value: string,
  enumObj: T,
): T[keyof T] | null {
  const enumValue = Object.values(enumObj).find(
    (v) => v === value.toLowerCase() || v === value.toUpperCase(),
  );
  return enumValue ? (enumValue as T[keyof T]) : null;
}

/**
 * Export all enum values as constants for easier access
 */
export const USER_STATUSES = Object.values(UserStatus);
export const USER_ROLES = Object.values(UserRole);
export const DOCUMENT_STATUSES = Object.values(DocumentStatus);
export const DOCUMENT_TYPES = Object.values(DocumentType);
export const INVOICE_STATUSES = Object.values(InvoiceStatus);
export const INVOICE_ACTIONS = Object.values(InvoiceAction);
export const CURRENCY_CODES = Object.values(CurrencyCode);
export const COMPANY_STATUSES = Object.values(CompanyStatus);
export const API_ERROR_CODES = Object.values(ApiErrorCode);
export const SORT_ORDERS = Object.values(SortOrder);

/**
 * Default values for enums
 */
export const ENUM_DEFAULTS = {
  userStatus: UserStatus.ACTIVE,
  userRole: UserRole.USER,
  documentStatus: DocumentStatus.DRAFT,
  documentType: DocumentType.INVOICE,
  invoiceStatus: InvoiceStatus.DRAFT,
  currency: CurrencyCode.RUB,
  companyStatus: CompanyStatus.ACTIVE,
  sortOrder: SortOrder.DESC,
} as const;

// ============================================================================
// ESF API Специфичные Перечисления (Приоритет 1)
// ============================================================================

/**
 * ESF API Document Status Enumeration
 * Используется числовые коды вместо строк в соответствии с ESF API
 */
export enum ESFDocumentStatus {
  DRAFT = "10", // Новый
  REVOKED = "20", // Отозван
  SENT = "30", // Отправлен
  ACCEPTED = "40", // Принят
  REJECTED = "50", // Отклонен
  DELETED = "70", // Удален
  ALL = "90", // Все (для фильтров)
}

/**
 * ESF Document Type Enumeration
 * Типы электронных счетов-фактур
 */
export enum ESFDocumentType {
  SALES = "1", // СФ-РЕ (Счет-фактура на реализацию)
  PURCHASE = "2", // СФ-ПР (Счет-фактура на покупку)
  CORRECTION = "3", // КСФ (Корректировочный счет-фактура)
}

/**
 * ESF Operation Type Enumeration
 * Тип операции согласно ESF API
 */
export enum ESFOperationType {
  SALES = "10", // Реализация
  PURCHASE_SERVICES = "20", // Покупка/услуги
  OTHER = "30", // Прочее
}

/**
 * ESF Delivery Type Enumeration
 * Способ доставки товаров
 */
export enum ESFDeliveryType {
  DIRECT = "1", // Прямая доставка
  WAREHOUSE = "2", // Через склад
  SERVICE = "3", // Услуга доставки
  COURIER = "4", // Курьер
  PICKUP = "5", // Самовывоз
  POST = "6", // Почта
}

/**
 * ESF Payment Type Enumeration
 * Тип платежа
 */
export enum ESFPaymentType {
  CASH = "1", // Наличные
  BANK_TRANSFER = "2", // Банковский перевод
  CHECK = "3", // Чек
  CARD = "4", // Карточка
  CREDIT = "5", // В кредит
  ELECTRONIC = "6", // Электронный платеж
}

/**
 * ESF Tax Type Enumeration
 * Тип налога/сбора
 */
export enum ESFTaxType {
  VAT = "НДС", // НДС (НДС)
  SPA = "НСП", // НСП (Налог специального назначения)
  EXCISE = "АЦ", // Акциз
  CUSTOMS = "ТП", // Таможенные платежи
}

/**
 * ESF Invoice Action Enumeration
 * Действия которые можно выполнить над счетом-фактурой
 */
export enum ESFInvoiceAction {
  SIGN = "sign", // Подписать
  ACCEPT = "accept", // Принять
  REJECT = "reject", // Отклонить
  REVOKE = "revoke", // Отозвать
  CORRECT = "correct", // Исправить (создать КСФ)
}

/**
 * ESF Unit Classification Enumeration
 * Распространённые единицы измерения в ESF
 */
export enum ESFUnitClassification {
  PIECE = "796", // Штука
  KILOGRAM = "110", // Килограмм
  GRAM = "163", // Грамм
  METER = "004", // Метр
  SQUARE_METER = "055", // Квадратный метр
  CUBIC_METER = "112", // Кубический метр (также Литр)
  HOUR = "255", // Час
  DAY = "259", // День
  KILOMETER = "037", // Километр
  PACKAGE = "303", // Упаковка
}

/**
 * Маппинг ESF статусов на русские названия
 */
export const ESF_STATUS_LABELS: Record<ESFDocumentStatus, string> = {
  [ESFDocumentStatus.DRAFT]: "Новый",
  [ESFDocumentStatus.REVOKED]: "Отозван",
  [ESFDocumentStatus.SENT]: "Отправлен",
  [ESFDocumentStatus.ACCEPTED]: "Принят",
  [ESFDocumentStatus.REJECTED]: "Отклонен",
  [ESFDocumentStatus.DELETED]: "Удален",
  [ESFDocumentStatus.ALL]: "Все",
};

/**
 * Маппинг типов операций на названия
 */
export const OPERATION_TYPE_LABELS: Record<ESFOperationType, string> = {
  [ESFOperationType.SALES]: "Реализация",
  [ESFOperationType.PURCHASE_SERVICES]: "Покупка/Услуги",
  [ESFOperationType.OTHER]: "Прочее",
};

/**
 * Маппинг типов документов ESF на названия
 */
export const ESF_DOCUMENT_TYPE_LABELS: Record<ESFDocumentType, string> = {
  [ESFDocumentType.SALES]: "СФ-РЕ (Реализация)",
  [ESFDocumentType.PURCHASE]: "СФ-ПР (Покупка)",
  [ESFDocumentType.CORRECTION]: "КСФ (Коррекция)",
};

/**
 * Маппинг доставки на названия
 */
export const DELIVERY_TYPE_LABELS: Record<ESFDeliveryType, string> = {
  [ESFDeliveryType.DIRECT]: "Прямая доставка",
  [ESFDeliveryType.WAREHOUSE]: "Через склад",
  [ESFDeliveryType.SERVICE]: "Услуга доставки",
  [ESFDeliveryType.COURIER]: "Курьер",
  [ESFDeliveryType.PICKUP]: "Самовывоз",
  [ESFDeliveryType.POST]: "Почта",
};

/**
 * Маппинг типов платежей на названия
 */
export const PAYMENT_TYPE_LABELS: Record<ESFPaymentType, string> = {
  [ESFPaymentType.CASH]: "Наличные",
  [ESFPaymentType.BANK_TRANSFER]: "Банковский перевод",
  [ESFPaymentType.CHECK]: "Чек",
  [ESFPaymentType.CARD]: "Карточка",
  [ESFPaymentType.CREDIT]: "В кредит",
  [ESFPaymentType.ELECTRONIC]: "Электронный платеж",
};

/**
 * Маппинг единиц измерения на названия
 */
export const UNIT_CLASSIFICATION_LABELS: Record<ESFUnitClassification, string> =
  {
    [ESFUnitClassification.PIECE]: "Штука",
    [ESFUnitClassification.KILOGRAM]: "Килограмм",
    [ESFUnitClassification.GRAM]: "Грамм",
    [ESFUnitClassification.METER]: "Метр",
    [ESFUnitClassification.SQUARE_METER]: "Квадратный метр",
    [ESFUnitClassification.CUBIC_METER]: "Кубический метр",
    [ESFUnitClassification.HOUR]: "Час",
    [ESFUnitClassification.DAY]: "День",
    [ESFUnitClassification.KILOMETER]: "Километр",
    [ESFUnitClassification.PACKAGE]: "Упаковка",
  };
