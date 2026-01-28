import { create } from "zustand";
import { persist } from "zustand/middleware";
import { User, Permission, ROLE_CONFIGS, AuthTokens } from "@/types/entities";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  tokens: AuthTokens | null;
  login: (payload: { user: User; tokens?: AuthTokens }) => void;
  logout: () => void;
  hasPermission: (permission: Permission) => boolean;
  switchRole: (role: User["role"]) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      tokens: null,

      login: ({ user, tokens }) => {
        set({ user, tokens: tokens || null, isAuthenticated: true });
      },

      logout: () => {
        set({ user: null, tokens: null, isAuthenticated: false });
      },

      hasPermission: (permission) => {
        const { user } = get();
        if (!user || !user.permissions) return false;
        return user.permissions.includes(permission);
      },

      switchRole: (role) => {
        const { user } = get();
        if (!user) return;

        const roleConfig = ROLE_CONFIGS[role];
        set({
          user: {
            ...user,
            role,
            permissions: roleConfig.permissions,
          },
        });
      },
    }),
    {
      name: "auth-storage",
    },
  ),
);
