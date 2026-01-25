'use client';

import { useAuthStore } from '@/store/authStore';

export function useApiToken(): string | null {
  // Select only the access token to avoid creating new objects on every render
  // This prevents infinite loops in Zustand selectors
  const accessToken = useAuthStore((state) => state.tokens?.accessToken ?? null);
  
  console.log('[useApiToken] Token:', accessToken ? '***' + accessToken.slice(-10) : null);
  
  return accessToken;
}
