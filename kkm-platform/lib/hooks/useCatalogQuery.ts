import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { showToast } from '@/lib/toast';

export interface CatalogItem {
  id: string;
  name: string;
  category: string;
  price: number;
  cost: number;
  stock: number;
  sku: string;
  status: 'active' | 'inactive';
  description?: string;
}

export function useCatalogQuery() {
  return useQuery({
    queryKey: ['catalog'],
    queryFn: async () => {
      const response = await apiClient.get<CatalogItem[]>('/api/v1/catalog');
      return response.data;
    },
  });
}

export function useCatalogItemQuery(itemId: string) {
  return useQuery({
    queryKey: ['catalog', itemId],
    queryFn: async () => {
      const response = await apiClient.get<CatalogItem>(`/api/v1/catalog/${itemId}`);
      return response.data;
    },
    enabled: !!itemId,
  });
}

export function useCreateCatalogItemMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: Partial<CatalogItem>) => {
      const response = await apiClient.post<CatalogItem>('/api/v1/catalog', data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['catalog'] });
      showToast.success('Товар успешно добавлен');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка добавления товара');
    },
  });
}

export function useUpdateCatalogItemMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...data }: { id: string } & Partial<CatalogItem>) => {
      const response = await apiClient.put<CatalogItem>(`/api/v1/catalog/${id}`, data);
      return response.data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['catalog'] });
      queryClient.invalidateQueries({ queryKey: ['catalog', variables.id] });
      showToast.success('Товар успешно обновлен');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка обновления товара');
    },
  });
}

export function useDeleteCatalogItemMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (itemId: string) => {
      await apiClient.delete(`/api/v1/catalog/${itemId}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['catalog'] });
      showToast.success('Товар успешно удален');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка удаления товара');
    },
  });
}
