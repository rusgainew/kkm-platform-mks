import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  ApiUser,
  AuthTokens,
  LoginRequest,
  RegisterRequest,
  ForgotPasswordRequest,
  ResetPasswordRequest,
  ChangePasswordRequest,
  UpdateProfileRequest,
  loginUser,
  logoutUser,
  registerUser,
  createUser,
  refreshTokens,
  getCurrentUser,
  updateProfile,
  forgotPassword,
  resetPassword,
  changePassword,
  listUsers,
  listUsersQuery,
  getUserById,
  assignUserRole,
} from "@/lib/api/users";
import { useAuthStore } from "@/store/authStore";
import {
  ROLE_CONFIGS,
  User,
  UserRole,
  AuthTokens as StoreTokens,
} from "@/types/entities";
import { UserRole as UserRoleEnum } from "@/types/enums";

function mapApiRoleToAppRole(role: string): UserRole {
  // API roles: admin, user, network_admin
  // App roles: admin, manager, cashier, employee, store_manager
  if (role === "admin") return "admin";
  if (role === "network_admin") return "admin"; // Map network_admin to admin
  if (role === "manager") return "manager";
  if (role === "cashier") return "cashier";
  if (role === "employee") return "employee";
  if (role === "store_manager") return "store_manager";
  // Default to employee for unknown roles
  return "employee";
}

// Extract user ID from JWT token
function extractUserIdFromToken(token?: string): string | null {
  if (!token) return null;
  try {
    const parts = token.split(".");
    if (parts.length !== 3) return null;

    // Decode JWT payload (second part)
    const payload = JSON.parse(atob(parts[1]));
    return payload.sub || payload.user_id || payload.id || null;
  } catch (e) {
    console.error("[mapApiUserToUser] Failed to extract ID from token:", e);
    return null;
  }
}

function mapApiUserToUser(apiUser: ApiUser, accessToken?: string): User {
  const role = mapApiRoleToAppRole(apiUser.role);
  const roleConfig = ROLE_CONFIGS[role];

  // Prefer explicit user_id, fallback to id, then JWT
  let userId = apiUser.user_id || apiUser.id;
  if (!userId || userId === "undefined") {
    console.warn(
      "[mapApiUserToUser] apiUser.id is missing, attempting to extract from JWT token",
    );
    userId = extractUserIdFromToken(accessToken) || apiUser.email;
  }

  return {
    // Backend required fields
    user_id: userId,
    email: apiUser.email,
    first_name: apiUser.first_name,
    last_name: apiUser.last_name,
    role,
    is_active: true, // Default to true from login
    created_at: Date.now(),
    updated_at: Date.now(),
    status: "active",

    // Frontend convenience fields
    id: userId,
    name: `${apiUser.first_name} ${apiUser.last_name}`.trim(),
    firstName: apiUser.first_name,
    lastName: apiUser.last_name,
    permissions: roleConfig.permissions,
  };
}

function mapTokens(tokens?: AuthTokens): StoreTokens | undefined {
  if (!tokens) return undefined;
  return {
    accessToken: tokens.access_token,
    refreshToken: tokens.refresh_token,
  };
}

export function useCurrentUserQuery() {
  const login = useAuthStore((s) => s.login);
  const accessToken = useAuthStore((s) => s.tokens?.accessToken);
  const query = useQuery({
    queryKey: ["me"],
    queryFn: getCurrentUser,
    retry: 1,
    staleTime: 1000,
    enabled: !!accessToken,
  });

  if (query.data) {
    const tokens = useAuthStore.getState().tokens;
    const user = mapApiUserToUser(query.data, tokens?.accessToken);
    if (tokens?.accessToken) {
      login({ user, tokens });
    }
  }

  return query;
}

export function useLoginMutation() {
  const login = useAuthStore((s) => s.login);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: LoginRequest) => loginUser(payload),
    onSuccess: (res) => {
      console.log("[useLoginMutation] Login successful, response:", res);
      const user = mapApiUserToUser(res.user, res.access_token);

      // API returns access_token and refresh_token directly, not in a nested object
      const tokens: StoreTokens = {
        accessToken: res.access_token,
        refreshToken: res.refresh_token,
        expiresIn: res.expires_in || 900, // Default 15 minutes if not provided
      };

      console.log("[useLoginMutation] Mapped tokens:", tokens);
      login({ user, tokens });
      console.log("[useLoginMutation] Logged in user:", user);
      queryClient.setQueryData(["me"], res.user);
    },
  });
}

export function useRegisterMutation() {
  const login = useAuthStore((s) => s.login);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: RegisterRequest) => registerUser(payload),
    onSuccess: (res) => {
      const tokens = res.tokens ? mapTokens(res.tokens) : undefined;
      const user = mapApiUserToUser(res.user, tokens?.accessToken);
      login({ user, tokens });
      queryClient.setQueryData(["me"], res.user);
    },
  });
}

export function useCreateUserMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: RegisterRequest) => createUser(payload),
    onSuccess: () => {
      // Refresh users list after creating a new user
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
}

export function useLogoutMutation() {
  const logout = useAuthStore((s) => s.logout);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => logoutUser(),
    onSuccess: () => {
      logout();
      queryClient.removeQueries({ queryKey: ["me"] });
    },
    onError: () => {
      // Даже если есть ошибка (401, истек токен) - разлогимся
      logout();
      queryClient.removeQueries({ queryKey: ["me"] });
    },
  });
}

export function useRefreshTokensMutation() {
  const login = useAuthStore((s) => s.login);
  return useMutation({
    mutationFn: (refreshToken: string) => refreshTokens(refreshToken),
    onSuccess: (tokens) => {
      const stateUser = useAuthStore.getState().user;
      if (stateUser && tokens) {
        const mappedTokens = mapTokens(tokens);
        // API может вернуть expires_in при обновлении
        const storeTokens: StoreTokens = {
          accessToken: mappedTokens?.accessToken || "",
          refreshToken: mappedTokens?.refreshToken || "",
          expiresIn: (tokens as any).expires_in || 900,
        };
        login({ user: stateUser, tokens: storeTokens });
        console.log("[useRefreshTokensMutation] Tokens refreshed successfully");
      }
    },
  });
}

export function useUpdateProfileMutation() {
  const login = useAuthStore((s) => s.login);
  return useMutation({
    mutationFn: (payload: UpdateProfileRequest) => updateProfile(payload),
    onSuccess: (user) => {
      const tokens = useAuthStore.getState().tokens || undefined;
      const mapped = mapApiUserToUser(user, tokens?.accessToken);
      login({ user: mapped, tokens });
    },
  });
}

export function useChangePasswordMutation() {
  return useMutation({
    mutationFn: (payload: ChangePasswordRequest) => changePassword(payload),
  });
}

export function useForgotPasswordMutation() {
  return useMutation({
    mutationFn: (payload: ForgotPasswordRequest) => forgotPassword(payload),
  });
}

export function useResetPasswordMutation() {
  return useMutation({
    mutationFn: (payload: ResetPasswordRequest) => resetPassword(payload),
  });
}

// Admin utilities
export function useListUsersQuery(enabled = true) {
  return useQuery({
    queryKey: ["users"],
    queryFn: () => listUsersQuery(),
    enabled,
  });
}

export function useUserByIdQuery(id?: string) {
  return useQuery({
    queryKey: ["users", id],
    queryFn: () => getUserById(id as string),
    enabled: !!id,
  });
}

export function useAssignUserRoleMutation() {
  return useMutation({
    mutationFn: ({ id, role }: { id: string; role: string }) =>
      assignUserRole(id, { role: role as "user" | "manager" }),
  });
}
