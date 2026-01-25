import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { showToast } from '@/lib/toast';

export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'user' | 'manager';
  status: 'active' | 'inactive' | 'pending';
  createdAt: string;
  lastLogin?: string;
}

export function useUsersQuery() {
  return useQuery({
    queryKey: ['users'],
    queryFn: async () => {
      const response = await apiClient.get<User[]>('/api/v1/users');
      return response.data;
    },
  });
}

export function useUserQuery(userId: string) {
  return useQuery({
    queryKey: ['users', userId],
    queryFn: async () => {
      const response = await apiClient.get<User>(`/api/v1/users/${userId}`);
      return response.data;
    },
    enabled: !!userId,
  });
}

export function useCreateUserMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: Partial<User>) => {
      const response = await apiClient.post<User>('/api/v1/users', data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      showToast.success('Пользователь успешно создан');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка создания пользователя');
    },
  });
}

export function useUpdateUserMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...data }: { id: string } & Partial<User>) => {
      const response = await apiClient.put<User>(`/api/v1/users/${id}`, data);
      return response.data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      queryClient.invalidateQueries({ queryKey: ['users', variables.id] });
      showToast.success('Пользователь успешно обновлен');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка обновления пользователя');
    },
  });
}

export function useDeleteUserMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (userId: string) => {
      await apiClient.delete(`/api/v1/users/${userId}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      showToast.success('Пользователь успешно удален');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка удаления пользователя');
    },
  });
}
