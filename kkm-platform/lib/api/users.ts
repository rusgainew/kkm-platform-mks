import { getApiBase } from "./client";

export type UserStatus = "active" | "disabled" | string;
export type UserRoleApi = "admin" | "manager" | "cashier" | string;

export interface ApiUser {
  user_id: string; // API возвращает user_id, не id
  id?: string; // Алиас для обратной совместимости
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  role: UserRoleApi;
  status: UserStatus;
  is_active?: boolean;
  last_login_at?: string;
  created_at: string | number;
  updated_at: string | number;
}

// Users listing response - supports both query and command service formats
export interface ListUsersResponse {
  users: ApiUser[];
  // Query service format
  total_count?: number;
  page?: number;
  per_page?: number;
  // Command service format (matches Go PageInfo struct)
  page_info?: {
    page: number;
    size: number;
    total_count: number;
  };
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export interface UpdateProfileRequest {
  user_id: string;
  first_name?: string;
  last_name?: string;
}

export interface AssignRoleRequest {
  target_user_id?: string;
  role: "user" | "manager";
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

export interface ChangePasswordRequest {
  user_id?: string;
  old_password: string;
  new_password: string;
}

// Helper to get auth token from storage (fallback for SSR)
function getAuthTokenFromStorage(): string | null {
  if (typeof window === "undefined") {
    console.log("[API] SSR environment detected, cannot access localStorage");
    return null;
  }

  try {
    console.log("[API] Attempting to get token from localStorage...");

    const readTokenFromPersistedState = (raw: string | null) => {
      if (!raw) return null;
      const parsed = JSON.parse(raw);
      return (
        parsed.state?.accessToken ||
        parsed.state?.tokens?.access_token ||
        parsed.state?.tokens?.accessToken ||
        parsed.state?.token ||
        parsed.token ||
        null
      );
    };

    const authStorage = localStorage.getItem("auth-storage");
    const authStore = localStorage.getItem("auth-store");

    const token =
      readTokenFromPersistedState(authStorage) ||
      readTokenFromPersistedState(authStore) ||
      localStorage.getItem("accessToken") ||
      localStorage.getItem("authToken") ||
      localStorage.getItem("token");

    if (token) {
      console.log("[API] Found token in localStorage");
      return token;
    }

    console.warn("[API] Token not found in localStorage");
    return null;
  } catch (err) {
    console.error("[API] Failed to parse auth token from storage:", err);
    return null;
  }
}

async function request<T>(
  path: string,
  options?: RequestInit & { token?: string },
): Promise<T> {
  const API_BASE = getApiBase();

  // Prefer token passed as option, fallback to localStorage
  const token = options?.token || getAuthTokenFromStorage();

  const headersObj: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (token) {
    headersObj["Authorization"] = `Bearer ${token}`;
    console.log(
      "[API] Adding Authorization header, token length:",
      token.length,
    );
  } else if (!path.includes("/login") && !path.includes("/register")) {
    // Only warn for endpoints that require auth
    console.warn("[API] No token available, request will likely fail with 401");
  }

  const headers: HeadersInit = {
    ...headersObj,
    ...(options?.headers || {}),
  };

  console.log("[API] Making request to:", `${API_BASE}${path}`, {
    method: options?.method || "GET",
    hasToken: !!token,
    headersKeys: Object.keys(headers),
  });

  // Pass the rest of options to fetch, excluding token
  const fetchOptions = (() => {
    if (!options) return {};
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { token, ...rest } = options;
    return rest;
  })();

  const res = await fetch(`${API_BASE}${path}`, {
    headers,
    credentials: "include",
    ...fetchOptions,
  });

  console.log("[API] Response status:", res.status, res.statusText);

  if (!res.ok) {
    const text = await res.text();
    console.warn("[API] Response error:", res.status, res.statusText);

    // Попытаемся распарсить JSON ошибку
    let errorMessage = text || res.statusText;
    try {
      const errorJson = JSON.parse(text);
      errorMessage = errorJson.message || errorJson.error || text;
    } catch {
      // Если не JSON, используем текст как есть
    }

    const error = new Error(errorMessage || `HTTP ${res.status}`);
    (error as any).status = res.status;
    (error as any).statusText = res.statusText;
    throw error;
  }

  const data = await res.json();
  console.log("[API] Response data:", data);

  // Debug: Log raw JSON to inspect field names
  if (path.includes("/users")) {
    console.log("[API] Raw users JSON:", JSON.stringify(data, null, 2));
    if (data.users && data.users.length > 0) {
      console.log("[API] First user keys:", Object.keys(data.users[0]));
      console.log("[API] First user data:", data.users[0]);
    }
  }

  return data as Promise<T>;
}

// Users
export const listUsers = () => request<ListUsersResponse>("/users");
export const listUsersQuery = (params?: {
  page?: number;
  per_page?: number;
}) => {
  const queryParams = new URLSearchParams();
  // API uses 0-based page indexing, but we use 1-based in UI
  if (params?.page) queryParams.append("page", (params.page - 1).toString());
  // API uses 'size' parameter instead of 'per_page'
  if (params?.per_page) queryParams.append("size", params.per_page.toString());

  const queryString = queryParams.toString();
  const path = queryString ? `/users?${queryString}` : "/users";

  return request<ListUsersResponse>(path);
};
export const getUserById = (id: string) => request<ApiUser>(`/users/${id}`);

// Auth flows
export const loginUser = (data: LoginRequest) =>
  request<{
    user: ApiUser;
    access_token: string;
    refresh_token: string;
    expires_in?: number;
    refresh_expires_in?: number;
    token_type?: string;
    timestamp?: number;
  }>("/users/login", {
    method: "POST",
    body: JSON.stringify(data),
  });

export const logoutUser = async () => {
  try {
    return await request<void>("/users/logout", {
      method: "POST",
      body: JSON.stringify({}),
    });
  } catch {
    // Игнорируем ошибки при logout (например, истек токен)
    // Локальное состояние очистится в любом случае
    return undefined;
  }
};

export const refreshTokens = (refresh_token: string) =>
  request<AuthTokens>("/users/refresh", {
    method: "POST",
    body: JSON.stringify({ refresh_token }),
  });

export const registerUser = (data: RegisterRequest) =>
  request<{ user: ApiUser; tokens: AuthTokens }>("/users/register", {
    method: "POST",
    body: JSON.stringify(data),
  });

export const createUser = (data: RegisterRequest) =>
  request<{ user: ApiUser }>("/users", {
    method: "POST",
    body: JSON.stringify(data),
  });

export const getCurrentUser = () => request<ApiUser>("/users/me");

// Password flows
export const forgotPassword = (data: ForgotPasswordRequest) =>
  request<void>("/users/forgot-password", {
    method: "POST",
    body: JSON.stringify(data),
  });

export const resetPassword = (data: ResetPasswordRequest) =>
  request<void>("/users/reset-password", {
    method: "POST",
    body: JSON.stringify(data),
  });

export const changePassword = (data: ChangePasswordRequest) =>
  request<void>("/users/password", {
    method: "PUT",
    body: JSON.stringify(data),
  });

// Profile management
export const updateProfile = (data: UpdateProfileRequest) =>
  request<ApiUser>("/users/profile", {
    method: "PUT",
    body: JSON.stringify(data),
  });

export const updateUser = (
  userId: string,
  data: Partial<Omit<RegisterRequest, "email" | "password">> & {
    role?: string;
  },
) => {
  console.log("[API] updateUser request:", { userId, data });
  return request<{ user: ApiUser }>(`/users/${userId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
};

export const deleteUser = (userId: string) => {
  console.log("[API] deleteUser request:", { userId });
  return request<void>(`/users/${userId}`, {
    method: "DELETE",
  });
};

export const assignUserRole = (userId: string, data: AssignRoleRequest) => {
  const body = {
    target_user_id: userId,
    role: data.role,
  };
  console.log("[API] assignUserRole request:", body);
  return request<ApiUser>(`/users/${userId}/role`, {
    method: "PUT",
    body: JSON.stringify(body),
  });
};
