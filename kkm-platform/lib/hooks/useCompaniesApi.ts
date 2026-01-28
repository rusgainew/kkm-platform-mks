/**
 * React Query хуки для API компаний
 * Управляют получением, созданием, обновлением и удалением компаний
 */

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import {
  createCompany,
  deleteCompany,
  getCompany,
  listCompanies,
  updateCompany,
} from "@/lib/api/companies";
import type {
  CreateCompanyRequest,
  UpdateCompanyRequest,
  Company,
} from "@/types/entities";
import { useAuthStore } from "@/store/authStore";
import { useRefreshTokensMutation } from "./useAuthApi";

/**
 * Получить список всех компаний
 * - Автоматически обновляет список при изменении токена
 * - При 401 пытается обновить токен и повторить запрос
 * - Ждет гидрации перед выполнением запроса
 */
export function useListCompaniesQuery() {
  const { tokens, isAuthenticated } = useAuthStore();
  const token = tokens?.accessToken;
  const [isHydrated, setIsHydrated] = useState(false);
  const refreshTokensMutation = useRefreshTokensMutation();

  // Ждем гидрации localStorage перед использованием токена
  useEffect(() => {
    console.log("[useListCompaniesQuery] Ожидание гидрации...");
    console.log("[useListCompaniesQuery] Аутентифицирован:", isAuthenticated);
    console.log("[useListCompaniesQuery] Токен доступен:", !!token);
    setIsHydrated(true);
  }, [token, isAuthenticated]);

  return useQuery({
    queryKey: ["companies"],
    queryFn: async () => {
      console.log("[useListCompaniesQuery] Выполняется запрос");
      console.log("[useListCompaniesQuery] Аутентификация:", isAuthenticated);

      if (!isAuthenticated) {
        console.warn(
          "[useListCompaniesQuery] Не аутентифицирован, возвращаем пустой массив",
        );
        return [];
      }

      try {
        console.log("[useListCompaniesQuery] Запрашиваем список компаний...");
        const response = await listCompanies();

        // Extract data from response
        if (response && response.items && Array.isArray(response.items)) {
          return response.items;
        }
        return [];
      } catch (error) {
        // Если получили 401, пытаемся обновить токен и повторить запрос
        if (error instanceof Error && error.message.includes("Unauthorized")) {
          console.warn(
            "[useListCompaniesQuery] Получена ошибка 401, пытаемся обновить токен...",
          );
          try {
            if (tokens?.refreshToken) {
              await refreshTokensMutation.mutateAsync(tokens.refreshToken);
              console.log(
                "[useListCompaniesQuery] Токен обновлен, повторяем запрос...",
              );
              // Повторяем с новым токеном
              const response = await listCompanies();
              if (response && response.items && Array.isArray(response.items)) {
                return response.items;
              }
              return [];
            }
          } catch (refreshError) {
            console.error(
              "[useListCompaniesQuery] Ошибка обновления токена:",
              refreshError,
            );
          }
        }
        throw error;
      }
    },
    enabled: !!token && isHydrated && isAuthenticated,
    retry: (failureCount, error) => {
      // Не повторяем 401 (уже попытались обновить токен)
      if (error instanceof Error && error.message.includes("Unauthorized")) {
        return false;
      }
      return failureCount < 2;
    },
  });
}

export function useGetCompanyQuery(id: string | null) {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken;
  const [isHydrated, setIsHydrated] = useState(false);

  useEffect(() => {
    setIsHydrated(true);
  }, []);

  return useQuery({
    queryKey: ["company", id],
    queryFn: async () => {
      if (!id) return null;
      return await getCompany(id);
    },
    enabled: !!id && isHydrated,
  });
}

export function useCreateCompanyMutation() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken || "";
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateCompanyRequest) => createCompany(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["companies"] });
    },
  });
}

export function useUpdateCompanyMutation() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken || "";
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: UpdateCompanyRequest;
    }) => updateCompany(id, payload),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["companies"] });
      queryClient.invalidateQueries({ queryKey: ["company", id] });
    },
  });
}

export function useDeleteCompanyMutation() {
  const { tokens } = useAuthStore();
  const token = tokens?.accessToken || "";
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteCompany(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["companies"] });
    },
  });
}
