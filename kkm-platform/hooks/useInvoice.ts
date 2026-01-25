/**
 * React Hooks для работы со счетами-фактурами
 * Интеграция с React Query (TanStack Query) для управления состоянием
 */

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import {
  invoiceAPI,
  type Invoice,
  type CreateInvoiceRequest,
  type InvoiceFilters,
} from "@/lib/api/invoice";
import { useAuthStore } from "@/store/authStore";

/**
 * Ключи для React Query кэша
 */
export const invoiceQueryKeys = {
  all: ["invoices"] as const,
  lists: () => [...invoiceQueryKeys.all, "list"] as const,
  list: (filters?: InvoiceFilters) =>
    [...invoiceQueryKeys.lists(), filters] as const,
  details: () => [...invoiceQueryKeys.all, "detail"] as const,
  detail: (id: string) => [...invoiceQueryKeys.details(), id] as const,
  statusHistory: (id: string) =>
    [...invoiceQueryKeys.all, "statusHistory", id] as const,
};

/**
 * Hook для получения списка счетов-фактур
 *
 * @param filters - Фильтры поиска
 * @param enabled - Включен ли запрос
 * @returns Query результат с данными и функциями управления
 */
export function useInvoiceList(filters?: InvoiceFilters, enabled = true) {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;

  return useQuery({
    queryKey: invoiceQueryKeys.list(filters),
    queryFn: async () => {
      if (!token) throw new Error("No authentication token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.listInvoices(filters);
    },
    enabled: enabled && !!token,
    staleTime: 5 * 60 * 1000, // 5 минут
    retry: 2,
  });
}

/**
 * Hook для получения одного счета-фактуры
 *
 * @param invoiceUuid - UUID счета
 * @param enabled - Включен ли запрос
 * @returns Query результат с данными счета
 */
export function useInvoiceDetail(invoiceUuid: string | null, enabled = true) {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;

  return useQuery({
    queryKey: invoiceQueryKeys.detail(invoiceUuid || ""),
    queryFn: async () => {
      if (!invoiceUuid || !token)
        throw new Error("Missing invoiceUuid or token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.getInvoice(invoiceUuid);
    },
    enabled: enabled && !!invoiceUuid && !!token,
    staleTime: 2 * 60 * 1000, // 2 минуты
    retry: 2,
  });
}

/**
 * Hook для получения истории статусов счета
 *
 * @param invoiceUuid - UUID счета
 * @returns Query результат с историей статусов
 */
export function useInvoiceStatusHistory(invoiceUuid: string | null) {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;

  return useQuery({
    queryKey: invoiceQueryKeys.statusHistory(invoiceUuid || ""),
    queryFn: async () => {
      if (!invoiceUuid || !token)
        throw new Error("Missing invoiceUuid or token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.getInvoiceStatusHistory(invoiceUuid);
    },
    enabled: !!invoiceUuid && !!token,
    staleTime: 60 * 1000, // 1 минута
    retry: 2,
  });
}

/**
 * Hook для создания нового счета-фактуры
 *
 * @returns Mutation для создания счета с callbacks
 */
export function useCreateInvoice() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (request: CreateInvoiceRequest) => {
      if (!token) throw new Error("No authentication token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.createInvoice(request);
    },
    onSuccess: (newInvoice) => {
      // Инвалидировать список счетов
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.lists(),
      });
      // Добавить в кэш новый счет
      queryClient.setQueryData(
        invoiceQueryKeys.detail(newInvoice.documentUuid),
        newInvoice,
      );
    },
    onError: (error) => {
      console.error("Failed to create invoice:", error);
    },
  });
}

/**
 * Hook для обновления счета-фактуры
 *
 * @returns Mutation для обновления счета
 */
export function useUpdateInvoice() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (
      request: Parameters<typeof invoiceAPI.updateInvoice>[0],
    ) => {
      if (!token) throw new Error("No authentication token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.updateInvoice(request);
    },
    onSuccess: (updatedInvoice) => {
      // Обновить в кэше
      queryClient.setQueryData(
        invoiceQueryKeys.detail(updatedInvoice.documentUuid),
        updatedInvoice,
      );
      // Инвалидировать список
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.lists(),
      });
    },
  });
}

/**
 * Hook для принятия счета-фактуры
 *
 * @returns Mutation для принятия счета
 */
export function useAcceptInvoice() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      invoiceUuid,
      comment,
    }: {
      invoiceUuid: string;
      comment?: string;
    }) => {
      if (!token) throw new Error("No authentication token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.acceptInvoice(invoiceUuid, comment);
    },
    onSuccess: (_, variables) => {
      // Инвалидировать данные счета
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.detail(variables.invoiceUuid),
      });
      // Инвалидировать список
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.lists(),
      });
    },
  });
}

/**
 * Hook для отклонения счета-фактуры
 *
 * @returns Mutation для отклонения счета
 */
export function useRejectInvoice() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      invoiceUuid,
      reason,
      comment,
    }: {
      invoiceUuid: string;
      reason: string;
      comment?: string;
    }) => {
      if (!token) throw new Error("No authentication token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.rejectInvoice(invoiceUuid, reason, comment);
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.detail(variables.invoiceUuid),
      });
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.lists(),
      });
    },
  });
}

/**
 * Hook для подписания счета-фактуры
 *
 * @returns Mutation для подписания счета
 */
export function useSignInvoice() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      invoiceUuid,
      signature,
      certificateData,
    }: {
      invoiceUuid: string;
      signature: string;
      certificateData?: string;
    }) => {
      if (!token) throw new Error("No authentication token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.signInvoice(
        invoiceUuid,
        signature,
        certificateData,
      );
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.detail(variables.invoiceUuid),
      });
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.lists(),
      });
    },
  });
}

/**
 * Hook для отзыва счета-фактуры
 *
 * @returns Mutation для отзыва счета
 */
export function useRevokeInvoice() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      invoiceUuid,
      reason,
      comment,
    }: {
      invoiceUuid: string;
      reason: string;
      comment?: string;
    }) => {
      if (!token) throw new Error("No authentication token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.revokeInvoice(invoiceUuid, reason, comment);
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.detail(variables.invoiceUuid),
      });
      queryClient.invalidateQueries({
        queryKey: invoiceQueryKeys.lists(),
      });
    },
  });
}

/**
 * Hook для скачивания счета в PDF
 *
 * @returns Mutation для скачивания PDF
 */
export function useDownloadInvoicePDF() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;

  return useMutation({
    mutationFn: async (invoiceUuid: string) => {
      if (!token) throw new Error("No authentication token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.downloadInvoicePDF(invoiceUuid);
    },
    onSuccess: (blob, invoiceUuid) => {
      // Создать URL и скачать файл
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `invoice-${invoiceUuid}.pdf`;
      link.click();
      window.URL.revokeObjectURL(url);
    },
  });
}

/**
 * Hook для скачивания счета в XML
 *
 * @returns Mutation для скачивания XML
 */
export function useDownloadInvoiceXML() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;

  return useMutation({
    mutationFn: async (invoiceUuid: string) => {
      if (!token) throw new Error("No authentication token");
      invoiceAPI.setToken(token);
      return await invoiceAPI.downloadInvoiceXML(invoiceUuid);
    },
    onSuccess: (blob, invoiceUuid) => {
      // Создать URL и скачать файл
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `invoice-${invoiceUuid}.xml`;
      link.click();
      window.URL.revokeObjectURL(url);
    },
  });
}
