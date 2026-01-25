import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { showToast } from '@/lib/toast';

export interface Invoice {
  id: string;
  invoiceNumber: string;
  companyId: string;
  amount: number;
  totalAmount: number;
  status: 'draft' | 'sent' | 'paid' | 'overdue' | 'cancelled';
  issueDate: string;
  dueDate: string;
  paidDate?: string;
}

export function useInvoicesQuery() {
  return useQuery({
    queryKey: ['invoices'],
    queryFn: async () => {
      const response = await apiClient.get<Invoice[]>('/api/v1/invoices');
      return response.data;
    },
  });
}

export function useInvoiceQuery(invoiceId: string) {
  return useQuery({
    queryKey: ['invoices', invoiceId],
    queryFn: async () => {
      const response = await apiClient.get<Invoice>(`/api/v1/invoices/${invoiceId}`);
      return response.data;
    },
    enabled: !!invoiceId,
  });
}

export function useCreateInvoiceMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: Partial<Invoice>) => {
      const response = await apiClient.post<Invoice>('/api/v1/invoices', data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invoices'] });
      showToast.success('Счет успешно создан');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка создания счета');
    },
  });
}

export function useUpdateInvoiceMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...data }: { id: string } & Partial<Invoice>) => {
      const response = await apiClient.put<Invoice>(`/api/v1/invoices/${id}`, data);
      return response.data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['invoices'] });
      queryClient.invalidateQueries({ queryKey: ['invoices', variables.id] });
      showToast.success('Счет успешно обновлен');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка обновления счета');
    },
  });
}

export function useDeleteInvoiceMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (invoiceId: string) => {
      await apiClient.delete(`/api/v1/invoices/${invoiceId}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invoices'] });
      showToast.success('Счет успешно удален');
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка удаления счета');
    },
  });
}

export function useInvoiceStatusMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, status }: { id: string; status: Invoice['status'] }) => {
      const response = await apiClient.patch<Invoice>(`/api/v1/invoices/${id}/status`, { status });
      return response.data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['invoices'] });
      queryClient.invalidateQueries({ queryKey: ['invoices', variables.id] });
      showToast.success(`Статус счета изменен на ${variables.status}`);
    },
    onError: (error) => {
      showToast.error(error instanceof Error ? error.message : 'Ошибка изменения статуса');
    },
  });
}
