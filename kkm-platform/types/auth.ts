import type { UserRole } from "./enums";

// Re-export UserRole for convenience
export type { UserRole };

export interface User {
  id: string;
  name: string;
  firstName?: string;
  lastName?: string;
  email: string;
  role: UserRole | string;
  storeId?: string; // Для store_manager - ID конкретного магазина
  permissions: Permission[];
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresIn?: number; // Время жизни токена в секундах
}

export type Permission =
  | "view_all_stores"
  | "view_own_store"
  | "manage_terminals"
  | "view_reports"
  | "manage_employees"
  | "manage_inventory"
  | "manage_users"
  | "view_users"
  | "view_fiscal_reports"
  | "view_companies"
  | "manage_companies";

export interface RoleConfig {
  role: UserRole | string;
  label: string;
  permissions: Permission[];
}

export const ROLE_CONFIGS: Record<string, RoleConfig> = {
  admin: {
    role: "admin",
    label: "Администратор",
    permissions: [
      "view_all_stores",
      "manage_terminals",
      "view_reports",
      "manage_employees",
      "manage_inventory",
      "manage_users",
      "view_users",
      "view_fiscal_reports",
      "view_companies",
      "manage_companies",
    ],
  },
  network_admin: {
    role: "network_admin",
    label: "Администратор сети",
    permissions: [
      "view_all_stores",
      "manage_terminals",
      "view_reports",
      "manage_employees",
      "manage_inventory",
      "manage_users",
      "view_users",
      "view_fiscal_reports",
      "view_companies",
      "manage_companies",
    ],
  },
  user: {
    role: "user",
    label: "Пользователь",
    permissions: [
      "view_own_store",
      "view_reports",
      "manage_inventory",
      "view_users",
      "view_companies",
    ],
  },
};
