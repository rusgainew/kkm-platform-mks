// Централизованные типы сущностей (новые)
export * from "./entities";

// Existing types (для обратной совместимости)
export * from "./product";
export * from "./cart";
export * from "./dashboard";
// export * from "./api-response"; // Конфликт с entities (APIResponse) - импортируйте напрямую при необходимости

// Специфичные типы из auth (только те, которых нет в entities)
export type { Permission, ROLE_CONFIGS, AuthTokens } from "./auth";

// Enum типы (избегая конфликтов)
export { UserStatus, ApiErrorCode, ESFDocumentStatus } from "./enums";

// ESF API типы - Приоритет 1
export * from "./reference";
export * from "./party";
export * from "./invoice"; // Для ESF-специфичных типов (CreateInvoiceRequest, UpdateInvoiceRequest и т.д.)
// export * from "./catalog"; // Removed - types exported from entities
