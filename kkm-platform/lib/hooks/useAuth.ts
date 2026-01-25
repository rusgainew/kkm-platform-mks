'use client';

import { useAuthStore } from '@/store/authStore';

/**
 * Hook для удобного доступа к состоянию авторизации
 * Использует для проверки наличия пользователя, токена и прав доступа
 */
export function useAuth() {
  const user = useAuthStore((state) => state.user);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const tokens = useAuthStore((state) => state.tokens);
  const hasPermission = useAuthStore((state) => state.hasPermission);
  const logout = useAuthStore((state) => state.logout);

  return {
    user,
    isAuthenticated,
    tokens,
    hasToken: !!tokens?.accessToken,
    hasRefreshToken: !!tokens?.refreshToken,
    tokenExpiresIn: tokens?.expiresIn,
    hasPermission,
    logout,
  };
}
