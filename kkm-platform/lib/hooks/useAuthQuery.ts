import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { showToast } from '@/lib/toast';

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'user' | 'manager';
}

export function useAuthQuery() {
  return useQuery({
    queryKey: ['auth', 'user'],
    queryFn: async () => {
      const response = await apiClient.get<AuthUser>('/api/v1/auth/me');
      return response.data;
    },
    retry: false,
    staleTime: 1000 * 60 * 30, // 30 минут
  });
}

export function useLoginMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ email, password }: { email: string; password: string }) => {
      const response = await apiClient.post<{
        user: AuthUser;
        accessToken: string;
        refreshToken: string;
      }>('/api/v1/auth/login', { email, password });
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['auth', 'user'] });
      showToast.success(`Добро пожаловать, ${data.user.name}!`);
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка при входе');
    },
  });
}

export function useLogoutMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      await apiClient.post('/api/v1/auth/logout');
    },
    onSuccess: () => {
      queryClient.clear();
      showToast.success('Вы вышли из системы');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка при выходе');
    },
  });
}

export function useRefreshTokenMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const response = await apiClient.post<{
        accessToken: string;
      }>('/api/v1/auth/refresh');
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth', 'user'] });
    },
  });
}
