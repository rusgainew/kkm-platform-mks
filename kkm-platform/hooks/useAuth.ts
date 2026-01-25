import { useAuthStore } from '@/store';
import { useCallback } from 'react';

/**
 * useAuth Hook - управление аутентификацией
 * 
 * Примеры использования:
 * const { user, login, logout, isLoading } = useAuth();
 * 
 * // Login
 * await login('user@example.com', 'password');
 * 
 * // Logout
 * logout();
 * 
 * // Check if authenticated
 * if (isAuthenticated) { ... }
 */
export const useAuth = () => {
  const {
    user,
    loading,
    error,
    isAuthenticated,
    login,
    logout,
    setError,
    clearError,
    setUser,
    refreshAccessToken,
    checkAuth,
  } = useAuthStore();

  const handleLogin = useCallback(
    async (email: string, password: string) => {
      try {
        await login(email, password);
        return true;
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Login failed';
        return false;
      }
    },
    [login]
  );

  const handleLogout = useCallback(() => {
    logout();
  }, [logout]);

  const handleRefreshToken = useCallback(
    async () => {
      try {
        await refreshAccessToken();
        return true;
      } catch (error) {
        return false;
      }
    },
    [refreshAccessToken]
  );

  return {
    user,
    isLoading: loading,
    error,
    isAuthenticated,
    login: handleLogin,
    logout: handleLogout,
    refreshToken: handleRefreshToken,
    checkAuth,
    setError,
    clearError,
    setUser,
  };
};

export default useAuth;
