/**
 * Hook for loading documents with error handling
 */

import { useEffect, useState } from "react";
import { listDocuments, Document } from "@/lib/api/documents";

interface UseDocumentsOptions {
  page?: number;
  pageSize?: number;
  status?: string;
  refreshTrigger?: number;
}

export interface PaginationData {
  currentPage: number;
  pageSize: number;
  totalCount: number;
  totalPages: number;
}

export const useDocuments = (options: UseDocumentsOptions = {}) => {
  const { page = 1, pageSize = 10, status, refreshTrigger = 0 } = options;
  const [documents, setDocuments] = useState<Document[]>([]);
  const [pagination, setPagination] = useState<PaginationData>({
    currentPage: page,
    pageSize,
    totalCount: 0,
    totalPages: 0,
  });
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadDocuments = async () => {
      try {
        setIsLoading(true);
        setError(null);

        const response = await listDocuments(
          page,
          pageSize,
          status || undefined
        );

        let documentsData: Document[] = [];

        console.log(
          "[useDocuments] Полный ответ сервера:",
          JSON.stringify(response, null, 2)
        );
        console.log(
          "[useDocuments] Все свойства ответа:",
          Object.keys(response || {})
        );
        console.log("[useDocuments] response.data:", (response as any)?.data);
        console.log("[useDocuments] Тип ответа:", typeof response);
        console.log("[useDocuments] Это массив?:", Array.isArray(response));
        console.log(
          "[useDocuments] Есть свойство documents?:",
          response && typeof response === "object" && "documents" in response
        );

        // Try different response formats
        if (
          response &&
          typeof response === "object" &&
          "data" in response &&
          Array.isArray((response as any).data)
        ) {
          // Format: { data: Document[], ... }
          documentsData = (response as any).data;
          console.log(
            "[useDocuments] Извлечено из свойства data:",
            documentsData
          );
        } else if (
          response &&
          typeof response === "object" &&
          "documents" in response
        ) {
          // Format: { documents: Document[], ... }
          const typedResponse = response as { documents: Document[] };
          documentsData = typedResponse.documents;
          console.log(
            "[useDocuments] Извлечено из свойства documents:",
            documentsData
          );
        } else if (Array.isArray(response)) {
          // Format: Document[]
          documentsData = response;
          console.log("[useDocuments] Ответ является массивом:", documentsData);
        } else {
          console.warn(
            "[useDocuments] Неизвестный формат ответа, ответ пуст или некорректен"
          );
        }

        console.log("[useDocuments] Загружены документы:", documentsData);
        setDocuments(documentsData);

        // Extract pagination data from response
        if (response && typeof response === "object") {
          const totalCountValue = Number(
            (response as any).total_count || (response as any).totalCount || 0
          );
          const totalPages = Math.ceil(totalCountValue / pageSize) || 1;
          setPagination({
            currentPage: page,
            pageSize,
            totalCount: totalCountValue,
            totalPages,
          });
          console.log("[useDocuments] Пагинация:", {
            currentPage: page,
            pageSize,
            totalCount: totalCountValue,
            totalPages,
          });
        }
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Ошибка загрузки документов";
        console.error("[useDocuments] Ошибка:", errorMessage);
        setError(errorMessage);
      } finally {
        setIsLoading(false);
      }
    };

    loadDocuments();
  }, [page, pageSize, status, refreshTrigger]);

  return { documents, pagination, isLoading, error, setDocuments };
};
